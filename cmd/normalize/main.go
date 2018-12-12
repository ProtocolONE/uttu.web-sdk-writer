package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"time"
	"uttu-web-sdk-writer/pkg/types"
)

var (
	config struct {
		BrokersIn  []string `env:"BROKERS-IN" envDefault:"localhost:9092" envSeparator:","`
		BrokersOut []string `env:"BROKERS-OUT" envDefault:"localhost:9092" envSeparator:","`
		Group      string   `env:"GROUP" envDefault:"normalize"`
		Frequency  int64    `env:"FREQUENCY" envDefault:"500"`
		ReadTopics []string `env:"READ-TOPIC" envDefault:"WebSDK-123-1" envSeparator:","`
		Offset     int64    `env:"OFFSET" envDefault:"-1"` //-2:oldest -1:newest
	}

	clientID string
	err      error

	consumer *cluster.Consumer
	producer sarama.AsyncProducer

	part cluster.PartitionConsumer
	ok   bool
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

func partitionConsumer(pc cluster.PartitionConsumer) {
	for msg := range pc.Messages() {
		log.Println(string(msg.Value))
		lineMsg := new(types.Msg)
		err := json.Unmarshal(msg.Value, lineMsg)
		if err != nil {
			log.Println(err.Error())
			sendToKafka(msg.Value, "1", "-error-unmarshall")
		} else {
			//log.Printf("%+v\n", *lineMsg)
			msgNorm := new(types.MsgNorm)
			msgNorm.ClientId, err = strconv.ParseInt(clientID, 10, 64)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			msgNorm.IsCookie = lineMsg.Cookie != "0"
			msgNorm.FingerPrint = lineMsg.FingerPrint
			msgNorm.FirstContentfulPaint = lineMsg.FirstContentfulPaint

			msgNorm.BodyFramesCount = lineMsg.Body.BodyFramesCount
			msgNorm.BodyIsIframe = lineMsg.Body.BodyIframe != 0
			msgNorm.BodyIsInstantArticle = lineMsg.Body.BodyInstantArticle != "0"
			msgNorm.BodyNotificationPermission = lineMsg.Body.BodyNotificationPermission
			msgNorm.BodyIsSameOriginTopFrame = lineMsg.Body.BodySameOriginTopFrame != 0

			msgNorm.NetDownlink = lineMsg.Net.NetDownlink
			msgNorm.NetDownlinkMax = lineMsg.Net.NetDownlinkMax
			msgNorm.NetEffectiveType = lineMsg.Net.NetEffectiveType
			msgNorm.NetRTT = lineMsg.Net.NetRTT
			msgNorm.NetType = lineMsg.Net.NetType

			msgNorm.PageCharset = lineMsg.Page.PageCharset
			msgNorm.PageLanguage = lineMsg.Page.PageLanguage
			msgNorm.PageRef = lineMsg.Page.PageRef
			msgNorm.PageUrl = lineMsg.Page.PageUrl

			msgNorm.PerformanceConnect = lineMsg.Performance.PerformanceConnect
			msgNorm.PerformanceDomainLookup = lineMsg.Performance.PerformanceDomainLookup
			msgNorm.PerformanceDomCompleteFromNav = lineMsg.Performance.PerformanceDomCompleteFromNav
			msgNorm.PerformanceDomContentLoadedEvent = lineMsg.Performance.PerformanceDomContentLoadedEvent
			msgNorm.PerformanceDomContentLoadedEventFromNav = lineMsg.Performance.PerformanceDomContentLoadedEventFromNav
			msgNorm.PerformanceDomInteractive = lineMsg.Performance.PerformanceDomInteractive
			msgNorm.PerformanceFetchFromNav = lineMsg.Performance.PerformanceFetchFromNav
			msgNorm.PerformanceLoadEvent = lineMsg.Performance.PerformanceLoadEvent
			msgNorm.PerformanceLoadEventStartFromNav = lineMsg.Performance.PerformanceLoadEventStartFromNav
			msgNorm.PerformanceRedirect = lineMsg.Performance.PerformanceRedirect
			msgNorm.PerformanceRedirectCount = lineMsg.Performance.PerformanceRedirectCount
			msgNorm.PerformanceRequest = lineMsg.Performance.PerformanceRequest
			msgNorm.PerformanceResponse = lineMsg.Performance.PerformanceResponse

			msgNorm.PlatformFlash = lineMsg.Platform.PlatformFlash
			msgNorm.PlatformIsJava = lineMsg.Platform.PlatformJava != "0"
			msgNorm.PlatformSystem = lineMsg.Platform.PlatformSystem

			msgNorm.ScreenMetricsAvailableHeight = lineMsg.Screen.ScreenMetricsAvailableHeight
			msgNorm.ScreenMetricsAvailableWight = lineMsg.Screen.ScreenMetricsAvailableWight
			msgNorm.ScreenMetricsBodyHeight = lineMsg.Screen.ScreenMetricsBodyHeight
			msgNorm.ScreenMetricsBodyWidth = lineMsg.Screen.ScreenMetricsBodyWidth
			msgNorm.ScreenMetricsDevicePixelRatio = lineMsg.Screen.ScreenMetricsDevicePixelRatio
			msgNorm.ScreenMetricsHeight = lineMsg.Screen.ScreenMetricsHeight
			msgNorm.ScreenMetricsOrientAngle = lineMsg.Screen.ScreenMetricsOrientAngle
			msgNorm.ScreenMetricsOrientType = lineMsg.Screen.ScreenMetricsOrientType
			msgNorm.ScreenMetricsWidth = lineMsg.Screen.ScreenMetricsWidth

			msgNorm.TimeStamp = lineMsg.Time.TimeStamp
			msgNorm.TimeZone = lineMsg.Time.TimeZone
			msgNorm.TimeZoneOffset = lineMsg.Time.TimeZoneOffset

			msgNorm.TrackerEventNum = lineMsg.Tracker[0].TrackerEventNum
			msgNorm.TrackerHitTime = lineMsg.Tracker[0].TrackerHitTime
			msgNorm.TrackerId = lineMsg.Tracker[0].TrackerId
			msgNorm.TrackerInactivity = lineMsg.Tracker[0].TrackerInactivity
			msgNorm.TrackerIsExpired = lineMsg.Tracker[0].TrackerIsExpired
			msgNorm.TrackerNum = lineMsg.Tracker[0].TrackerNum

			log.Printf("%+v\n", *msgNorm)

			msgNormByte, errMarshall := json.Marshal(msgNorm)
			if errMarshall != nil {
				sendToKafka(msg.Value, "1", "-error-marshall")
			} else {
				sendToKafka(msgNormByte, "1", "")
			}

		}
		consumer.MarkOffset(msg, "")
	}
}

func sendToKafka(msg []byte, versionOut string, postfix string) {
	//log.Println("WebSDK_" + clientId + "_" + version)
	//log.Println(string(msg))
	producer.Input() <- &sarama.ProducerMessage{
		Topic: "WebSDKNormalize-" + versionOut + postfix,
		Value: sarama.ByteEncoder(msg),
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	// go func() {
	//      log.Println(http.ListenAndServe("0.0.0.0:8888", nil))
	// }()

	loadEnvs()
	log.Println(config)

	//brokersString := "localhost:9092"
	clientID = "123"
	//broker := sarama.NewBroker(brokersString)

	//configBroker := sarama.NewConfig()
	//configBroker.Version = sarama.V2_0_0_0

	configProducer := sarama.NewConfig()
	configProducer.Version = sarama.V2_0_0_0
	configProducer.Producer.RequiredAcks = sarama.WaitForLocal
	configProducer.Producer.Compression = sarama.CompressionLZ4
	configProducer.Producer.Flush.Frequency = time.Duration(500) * time.Millisecond

	configConsumer := cluster.NewConfig()
	configConsumer.Version = sarama.V2_0_0_0
	//configConsumer.ChannelBufferSize = 10000
	configConsumer.Group.Mode = cluster.ConsumerModePartitions
	configConsumer.Group.Return.Notifications = true
	configConsumer.Consumer.Return.Errors = true
	//configConsumer.Consumer.Offsets.Initial = offset
	configConsumer.Consumer.Offsets.Initial = config.Offset

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

	//err = broker.Open(configBroker)
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
	//
	//loadKafkaTopics(broker)

	producer, err = sarama.NewAsyncProducer(config.BrokersOut, configProducer)
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		log.Printf("New kafka producer created\n")
	}
	defer producer.AsyncClose()

	for {
		select {
		case part, ok = <-consumer.Partitions():
			if !ok {
				return
			}
			go partitionConsumer(part)
		case err = <-producer.Errors():
			log.Println("Failed to produce message", err)
		}
	}

}
