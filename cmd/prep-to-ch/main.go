package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/caarlos0/env"
	"github.com/dgryski/go-jump"
	"github.com/joho/godotenv"
	"github.com/kshvakov/clickhouse"
	"github.com/pierrec/xxHash/xxHash64"
	"github.com/valyala/fastjson"
	"log"
	"time"
	"uttu-web-sdk-writer/pkg/constants"
)

type buffer struct {
	Shard int
	Data  []byte
}

type DataBase struct {
	connect *sql.DB
}

var (
	config struct {
		CHDataSourceName string   `env:"CH-DATASOURCE-NAME" envDefault:"tcp://localhost:9001?debug=true&read_timeout=10&write_timeout=20"`
		BrokersIn        []string `env:"BROKERS-IN" envDefault:"localhost:9092" envSeparator:","`
		BrokersOut       []string `env:"BROKERS-OUT" envDefault:"localhost:9092" envSeparator:","`
		Group            string   `env:"GROUP" envDefault:"PrepToCH"`
		Frequency        int64    `env:"FREQUENCY" envDefault:"500"`
		ReadTopics       []string `env:"READ-TOPIC" envDefault:"WebSDKNormalize-1" envSeparator:","`
		Offset           int64    `env:"OFFSET" envDefault:"-1"` //-2:oldest -1:newest
		WriteTopicPrefix string   `env:"WRITE-TOPIC-PREFIX" envDefault:"WebSDKPrepared-1"`
		TopicNums        int      `env:"TOPIC-NUMS" envDefault:"3"`
		DataFlashBuffer  int      `env:"DATA-FLASH-BUFFER" envDefault:"10000"`
	}

	err             error
	consumer        *cluster.Consumer
	producer        sarama.AsyncProducer
	messagesChannel chan buffer
	part            cluster.PartitionConsumer
	ok              bool
)

func loadEnvs() {
	if err := godotenv.Load(); err != nil {
		//log.Println("File .env not found, reading configuration from ENV")
	}

	if err := env.Parse(&config); err != nil {
		log.Fatal("Failed to parse ENV")
	}
}

func (db *DataBase) InitCH(dsn string) error {
	var err error

	db.connect, err = sql.Open("clickhouse", dsn)
	if err != nil {
		return err
	}

	if err := db.connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Fatalf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return err
	}

	if _, err := db.connect.Exec(constants.WebSDK); err != nil {
		return err
	}
	if _, err := db.connect.Exec(constants.WebSDKKafka); err != nil {
		return err
	}
	if _, err := db.connect.Exec(constants.WebSDKKafkaMV); err != nil {
		return err
	}

	return nil
}

func partitionConsumer(pc cluster.PartitionConsumer) {

	var (
		err           error
		clientShard   int
		hash          uint64
		tickerMessage time.Ticker
	)
	messages := make([]bytes.Buffer, config.TopicNums)
	messagesSize := make([]int, config.TopicNums)
	tickerMessage = *time.NewTicker(time.Second * 5)

	for messageIn := range pc.Messages() {

		hash = xxHash64.Checksum(fastjson.GetBytes(messageIn.Value, "clientID"), 0xC0FE)
		clientShard = int(jump.Hash(hash, config.TopicNums))
		log.Println("Read topic, clientShard:", clientShard)

		messages[clientShard].Write(messageIn.Value)
		messages[clientShard].WriteString("\n")
		messagesSize[clientShard] += messages[clientShard].Len()
		if messagesSize[clientShard] >= config.DataFlashBuffer {
			var m buffer
			m.Shard = clientShard
			m.Data = make([]byte, messages[clientShard].Len())
			_, err = messages[clientShard].Read(m.Data)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}
			log.Printf("Data flash len(%d) in topic %d\n", messagesSize[clientShard], clientShard)
			messagesSize[clientShard] = 0
			messagesChannel <- m
			tickerMessage = *time.NewTicker(time.Second * 5)
		}

		select {
		case <-tickerMessage.C:
			{
				for i := 0; i < config.TopicNums; i++ {
					if messages[i].Len() != 0 {
						var m buffer
						m.Shard = i
						m.Data = make([]byte, messages[i].Len())
						_, err = messages[i].Read(m.Data)
						if err != nil {
							log.Printf("Error: %s\n", err.Error())
						}
						messagesChannel <- m
						log.Printf("Time flash len(%d) in topic %d\n", messagesSize[i], i)
						messagesSize[i] = 0
					}
				}
			}
		default:
		}
		consumer.MarkOffset(messageIn, "")
	}
}

func main() {

	// go func() {
	//      log.Println(http.ListenAndServe("0.0.0.0:8888", nil))
	// }()
	loadEnvs()
	log.Println(config)

	log.Println("Init DB...")
	db := DataBase{}
	if err = db.InitCH(config.CHDataSourceName); err != nil {
		log.Fatalln("Error:", err)
	}
	if err = db.connect.Close(); err != nil {
		log.Fatalln("Error:", err)
	}

	log.Println("Started...")

	configConsumer := cluster.NewConfig()
	configConsumer.Version = sarama.V2_0_0_0
	configConsumer.ChannelBufferSize = 10000
	configConsumer.Group.Mode = cluster.ConsumerModePartitions
	configConsumer.Group.Return.Notifications = true
	configConsumer.Consumer.Return.Errors = true
	configConsumer.Consumer.Offsets.Initial = config.Offset
	// configConsumer.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err = cluster.NewConsumer(config.BrokersIn, config.Group, config.ReadTopics, configConsumer)
	// consumer, err = cluster.NewConsumer(brokersIn, *group+strconv.FormatInt(time.Now().Unix(), 10), readTopics, configConsumer)
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		log.Printf("New kafka consumer created\n")
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()
	go func() {
		for note := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", note)
		}
	}()

	configProducer := sarama.NewConfig()
	configProducer.Version = sarama.V2_0_0_0
	configProducer.Producer.RequiredAcks = sarama.WaitForLocal
	configProducer.Producer.Compression = sarama.CompressionLZ4
	configProducer.Producer.Flush.Frequency = time.Duration(config.Frequency) * time.Millisecond
	//configProducer.Producer.MaxMessageBytes = 10 * 1000 * 1000

	producer, err = sarama.NewAsyncProducer(config.BrokersOut, configProducer)
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		log.Printf("New kafka producer created\n")
	}
	defer producer.AsyncClose()

	messagesChannel = make(chan buffer, 200)

	for {
		select {
		case part, ok = <-consumer.Partitions():
			if !ok {
				return
			}
			go partitionConsumer(part)
		case msg := <-messagesChannel:
			{
				log.Printf("Send len:%d msg bytes in %d topic\n", len(msg.Data), msg.Shard)
				producer.Input() <- &sarama.ProducerMessage{
					Topic: config.WriteTopicPrefix + fmt.Sprintf("-%02d", msg.Shard),
					Value: sarama.ByteEncoder(msg.Data),
				}
			}
		case err = <-producer.Errors():
			log.Println("Failed to produce message", err)
		}
	}
}
