package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/dw-account-service/configs"
	"log"
	"os"
	"strings"
)

var (
	Producer     sarama.SyncProducer
	SaramaLogger = log.New(os.Stdout, "[PRODUCER] ", log.LstdFlags)
)

func StartMessageProducer() error {
	splitBrokers := strings.Split(configs.MainConfig.Kafka.Brokers, ",")

	conf := configs.NewSaramaConfig()
	conf.Producer.Retry.Max = configs.MainConfig.Kafka.Producer.RetryMax
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	//conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Idempotent = configs.MainConfig.Kafka.Producer.Idempotent

	syncProducer, err := sarama.NewSyncProducer(splitBrokers, conf)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create producer: %s", err.Error()))
	}

	Producer = syncProducer
	log.Println("[INIT] kafka producer >> created")

	return nil
}
