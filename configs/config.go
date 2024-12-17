package configs

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	ENV            = "BILLING_SERVICE_ENV"
	LOCAL_ENV      = "local"
	STAGING_ENV    = "staging"
	PRODUCTION_ENV = "production"
)

type Config struct {
	Http     HttpConfig
	DB       DBConfig
	RabbitMQ RabbitMQConfig
}

type HttpConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ApiPrefix    string
	BaseURL      string
}

type DBConfig struct {
	DriverName      string
	DBName          string
	Source          string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
}

type RabbitMQConfig struct {
	Dsn       string
	QueueName string
}

func InitConfig() Config {
	var config Config
	env := os.Getenv(ENV)
	if env == LOCAL_ENV {
		viper.SetConfigName("env-local")
	} else if env == STAGING_ENV {
		viper.SetConfigName("env")
	} else {
		viper.SetConfigName("env-local")
	}

	viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/appname/") // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath("/opt/")
	viper.AddConfigPath("./configs/") // optionally look for config in the working directory
	err := viper.ReadInConfig()       // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		log.Fatalf("Error read config file: %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshall config struct: %s", err)
	}

	return config
}
