package util

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	ListenerAddress string   `mapstructure:"listenerAddress"`
	Database        Database `mapstructure:"database"`
}

type Database struct {
	Dsn  string `mapstructure:"dsn"`
	Type string `mapstructure:"type"`
}

var lock = &sync.Mutex{}

var configInstance *Config

func GetInstance() *Config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
			log.Print("Loading config")
			config, err := LoadConfig(".")
			if err != nil {
				log.Fatal(err)
			}
			configInstance = &config
		}
	}
	log.Print("Reusing config")
	return configInstance
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
