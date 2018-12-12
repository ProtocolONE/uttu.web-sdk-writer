package main

import (
	"github.com/Shopify/sarama"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	config struct {
		BrokersOut []string `env:"BROKERS-OUT" envDefault:"localhost:9092" envSeparator:","`
		Frequency  int64    `env:"FREQUENCY" envDefault:"500"`
	}

	err      error
	producer sarama.AsyncProducer
)

func loadEnvs() {
	if err := godotenv.Load(); err != nil {
		//log.Println("File .env not found, reading configuration from ENV")
	}

	if err := env.Parse(&config); err != nil {
		log.Fatal("Failed to parse ENV")
	}
}

//func loadKafkaTopics(broker *sarama.Broker) {
//	// TO DO: make check topics
//	req := sarama.MetadataRequest{}
//	//for {
//	resp, err := broker.GetMetadata(&req)
//	if err != nil {
//		log.Printf("%#v", &err)
//	}
//	t := resp.Topics
//	for _, val := range t {
//		//log.Printf("Key is %s", key)
//		log.Printf("topic: %s", val.Name)
//	}
//	time.Sleep(5 * time.Second)
//	//}
//}

func collect(c echo.Context) error {
	clientId := c.Param("clientId")
	if _, err := strconv.ParseInt(clientId, 10, 64); err != nil {
		return c.String(http.StatusConflict, "clientId not a number")
	}
	version := c.Param("version")
	if _, err := strconv.ParseInt(version, 10, 64); err != nil {
		return c.String(http.StatusConflict, "version not a number")
	}

	//getBody, err := c.Request().GetBody()
	//if err!=nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(string(getBody))
	//}

	bodyBuffer, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(string(bodyBuffer))
		sendToKafka(bodyBuffer, clientId, version)
	}
	return c.HTMLBlob(http.StatusOK, bodyBuffer)

	//return c.String(http.StatusOK, "")
}

func sendToKafka(msg []byte, clientId string, version string) {
	//log.Println("WebSDK_" + clientId + "_" + version)
	//log.Println(string(msg))
	producer.Input() <- &sarama.ProducerMessage{
		Topic: "WebSDK-" + clientId + "-" + version,
		Value: sarama.ByteEncoder(msg),
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	loadEnvs()

	//brokersString := "localhost:9092"
	//broker := sarama.NewBroker(brokersString)

	configKafka := sarama.NewConfig()
	configKafka.Version = sarama.V2_0_0_0
	configKafka.Producer.RequiredAcks = sarama.WaitForLocal
	configKafka.Producer.Compression = sarama.CompressionLZ4
	configKafka.Producer.Flush.Frequency = time.Duration(config.Frequency) * time.Millisecond

	//err = broker.Open(configKafka)
	//if err != nil {
	//	log.Fatalln("Error:", err)
	//} else {
	//	log.Printf("Open broker\n")
	//}
	//defer func() {
	//	if err := broker.Close(); err != nil {
	//		log.Println(err)
	//	}
	//}()
	//
	//connected, err := broker.Connected()
	//if err != nil {
	//	log.Print(err.Error())
	//}
	//if connected {
	//	log.Print("connected")
	//}

	producer, err = sarama.NewAsyncProducer(config.BrokersOut, configKafka)
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		log.Printf("New kafka producer created\n")
	}
	defer producer.AsyncClose()

	//go loadKafkaTopics(broker)
	e := echo.New()
	e.POST("/collect/:clientId/:version", collect)
	e.Logger.Fatal(e.Start(":80"))
}
