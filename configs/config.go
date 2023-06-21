package configs

import (
	"github.com/dw-account-service/pkg/xlogger"
	"github.com/spf13/viper"
	"log"
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

type KafkaConfig struct {
	// mode: producer|consumer|both
	Mode string `mapstructure:"mode"`
	// brokers: comma separated list
	Brokers  string              `mapstructure:"brokers"`
	SASL     KafkaSASLConfig     `mapstructure:"sasl"`
	TLS      KafkaTlsConfig      `mapstructure:"tls"`
	Producer KafkaProducerConfig `mapstructure:"producer"`
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

	err = (&xlogger.AppLogger{
		LogPath:     MainConfig.LogPath,
		CompressLog: true,
		DailyRotate: true,
	}).SetAppLogger()

	if err != nil {
		log.Fatalln("Logger error: ", err.Error())
		return err
	}

	xlogger.Log.SetPrefix("[INIT-APP] ")
	xlogger.Log.Println(strings.Repeat("-", 40))
	xlogger.Log.Println("| configuration >> loaded")
	return nil
}
