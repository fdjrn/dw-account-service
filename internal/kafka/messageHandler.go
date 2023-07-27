package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/handlers/consumer"
	"github.com/dw-account-service/internal/kafka/topic"
	"github.com/dw-account-service/internal/utilities"
	"time"
)

// HandleMessages contain function/logic that will be executed depends on topic name
func HandleMessages(message *sarama.ConsumerMessage) {
	var (
		handler              = consumer.NewTransactionHandler()
		trx                  = new(entity.BalanceTransaction)
		err                  error
		resultTopicMsg, pMsg string
	)

	utilities.Log.SetPrefix("[CONSUMER] ")

	switch message.Topic {
	case topic.TopUpRequest:
		resultTopicMsg = topic.TopUpResult
		pMsg = "topup"
	case topic.DeductRequest:
		resultTopicMsg = topic.DeductResult
		pMsg = "payment"
	case topic.DistributionRequest:
		resultTopicMsg = topic.DistributionResult
		pMsg = "balance distribution"
	default:
		utilities.Log.Println("| unknown topic message")
		return
	}

	trx, err = handler.DoHandleTransactionRequest(message)
	if err != nil {
		utilities.Log.Printf("| failed to process consumed message for topic: %s, with err: %s\n",
			message.Topic,
			err.Error())
	} else {
		utilities.Log.Printf("| %s with RefNo: %s, has been successfully processed with receipt number: %s\n",
			pMsg,
			trx.ReferenceNo,
			trx.ReceiptNumber,
		)
	}

	payload, _ := json.Marshal(trx)
	err = ProduceMsg(resultTopicMsg, payload)
	if err != nil {
		utilities.Log.Println("| cannot produce message for topic: ", message.Topic, ", with err: ", err.Error())
	}

	// Do Balance Distribution among members
	if trx.TransType == utilities.TransTypeDistribution && trx.Status == utilities.TrxStatusSuccess {
		start := time.Now()
		utilities.Log.Println("| starting merchant balance distribution ... ")
		err = DoBalanceDistribution(trx)
		if err != nil {
			utilities.Log.Println("| error occurred: ", err.Error())
		}
		utilities.Log.Println("| balance distribution finished in", time.Since(start).Seconds(), "seconds")
	}

}
