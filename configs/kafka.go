package configs

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/dw-account-service/pkg/tools"
	"github.com/google/uuid"
)

func NewSaramaConfig() *sarama.Config {
	uClientID, _ := uuid.NewUUID()
	cfg := sarama.NewConfig()

	cfg.Version = sarama.V0_11_0_2
	cfg.ClientID = fmt.Sprintf("mdw-client-%s", uClientID)
	cfg.Metadata.Full = true

	cfg.Net.SASL.Enable = MainConfig.Kafka.SASL.Enable
	cfg.Net.SASL.User = MainConfig.Kafka.SASL.SASLUserName
	cfg.Net.SASL.Password = MainConfig.Kafka.SASL.SASLPassword
	cfg.Net.SASL.Handshake = true
	cfg.Net.MaxOpenRequests = 1

	cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &tools.XDGSCRAMClient{HashGeneratorFcn: tools.SHA256} }

	if MainConfig.Kafka.SASL.Algorithm == "sha512" {
		cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &tools.XDGSCRAMClient{HashGeneratorFcn: tools.SHA512} }
	}

	cfg.Net.TLS.Enable = MainConfig.Kafka.TLS.Enable
	cfg.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: MainConfig.Kafka.TLS.Enable,
	}

	return cfg
}
