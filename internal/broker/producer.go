package broker

import (
	"github.com/Shopify/sarama"
	"github.com/dw-account-service/internal/conn/kafka"
	"github.com/dw-account-service/pkg/tools"
)

var (
	DeductTopic = "mdw.transaction.deduct.created"
	TopUpTopic  = "mdw.transaction.topup.created"
)

func ProduceMsg(topic string, payload []byte) error {

	partition, offset, err := kafka.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(tools.GetUnixTime()),
		Value: sarama.StringEncoder(payload),
	})
	if err != nil {
		kafka.SaramaLogger.Println("failed to send message to ", topic, err)
		return err
	}

	kafka.SaramaLogger.Printf("wrote message at partition: %d, offset: %d", partition, offset)
	return nil
}
