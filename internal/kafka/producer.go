package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/dw-account-service/configs"
	"github.com/dw-account-service/internal/utilities"
	"github.com/dw-account-service/internal/utilities/str"

	//"github.com/dw-account-service/pkg/tools"
	"strings"
)

var Producer sarama.SyncProducer

func initProducer() error {
	splitBrokers := strings.Split(configs.MainConfig.Kafka.Brokers, ",")

	conf := configs.NewSaramaConfig()
	conf.Producer.Retry.Max = configs.MainConfig.Kafka.Producer.RetryMax
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	//conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Idempotent = configs.MainConfig.Kafka.Producer.Idempotent

	syncProducer, err := sarama.NewSyncProducer(splitBrokers, conf)
	if err != nil {
		return errors.New(fmt.Sprintf("| failed to create producer: %s", err.Error()))
	}

	Producer = syncProducer
	utilities.Log.Println("| producer >> created")

	return nil
}

func Initialize() error {
	switch configs.MainConfig.Kafka.Mode {
	case "producer":
		if err := initProducer(); err != nil {
			return err
		}
	case "consumer":
		return nil
	}

	return nil
}

func ProduceMsg(topic string, payload []byte) error {
	utilities.Log.SetPrefix("[PRODUCER] ")
	_, _, err := Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(str.GetUnixTime()),
		Value: sarama.StringEncoder(payload),
	})
	if err != nil {
		utilities.Log.Println("| failed to send message to ", topic, err)
		return err
	}

	//utilities.Log.Printf("| message successfully wrote at partition: %d, offset: %d\n", partition, offset)
	return nil
}
