package configs

import (
	"github.com/dw-account-service/internal/utilities"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type MongoConfig struct {
	Uri    string `mapstructure:"uri"`
	DBName string `mapstructure:"dbName"`
}

type DBConfig struct {
	Mongo MongoConfig `mapstructure:"mongodb"`
}

type KafkaSASLConfig struct {
	Enable       bool   `mapstructure:"enable"`
	Algorithm    string `mapstructure:"algorithm"`
	SASLUserName string `mapstructure:"user"`
	SASLPassword string `mapstructure:"password"`
}

type KafkaTlsConfig struct {
	Enable             bool `mapstructure:"enable"`
	InsecureSkipVerify bool `mapstructure:"enable"`
}

type KafkaProducerConfig struct {
	Idempotent bool `mapstructure:"idempotent"`
	RetryMax   int  `mapstructure:"retryMax"`
}

type KafkaConsumerConfig struct {
	Assignor          string `mapstructure:"assignor"`
	Oldest            bool   `mapstructure:"oldest"`
	Verbose           int    `mapstructure:"verbose"`
	ConsumerGroupName string `mapstructure:"consumerGroupName"`
	ConsumerTopics    string `mapstructure:"topics"`
}

type KafkaConfig struct {
	// mode: producer|consumer|both
	Mode string `mapstructure:"mode"`
	// brokers: comma separated list
	Brokers  string              `mapstructure:"brokers"`
	SASL     KafkaSASLConfig     `mapstructure:"sasl"`
	TLS      KafkaTlsConfig      `mapstructure:"tls"`
	Producer KafkaProducerConfig `mapstructure:"producer"`
	Consumer KafkaConsumerConfig `mapstructure:"consumer"`
}

type AppConfig struct {
	AppName   string `mapstructure:"appName"`
	DebugMode bool   `mapstructure:"debugMode"`
	// os | file
	LogOutput          string       `mapstructure:"logOutput"`
	LogPath            string       `mapstructure:"logPath"`
	VerboseAPIResponse bool         `mapstructure:"verboseApiResponse"`
	APIServer          ServerConfig `mapstructure:"server"`
	Database           DBConfig     `mapstructure:"database"`
	Kafka              KafkaConfig  `mapstructure:"kafka"`
}

var MainConfig AppConfig

func Initialize() error {

	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&MainConfig)
	if err != nil {
		return err
	}

	// --- override config from os-env if its defined ---
	if os.Getenv("DATABASE_MONGODB_URI") != "" {
		MainConfig.Database.Mongo.Uri = os.Getenv("DATABASE_MONGODB_URI")
	}

	if os.Getenv("DATABASE_MONGODB_DB_NAME") != "" {
		MainConfig.Database.Mongo.DBName = os.Getenv("DATABASE_MONGODB_DB_NAME")
	}

	if os.Getenv("KAFKA_BROKERS") != "" {
		MainConfig.Kafka.Brokers = os.Getenv("KAFKA_BROKERS")
	}

	if os.Getenv("KAFKA_SASL_USER") != "" {
		MainConfig.Kafka.SASL.SASLUserName = os.Getenv("KAFKA_SASL_USER")
	}

	if os.Getenv("KAFKA_SASL_PASSWORD") != "" {
		MainConfig.Kafka.SASL.SASLPassword = os.Getenv("KAFKA_SASL_PASSWORD")
	}
	// --- end config overrides ---

	err = (&utilities.AppLogger{
		LogPath:     MainConfig.LogPath,
		CompressLog: true,
		DailyRotate: true,
	}).SetAppLogger()

	if err != nil {
		log.Fatalln("Logger error: ", err.Error())
		return err
	}

	utilities.Log.SetPrefix("[INIT-APP] ")
	utilities.Log.Println(strings.Repeat("-", 40))
	utilities.Log.Println("| configuration >> loaded")
	return nil
}
