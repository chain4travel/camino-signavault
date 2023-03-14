package util

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	ListenerAddress string   `mapstructure:"listenerAddress"`
	Database        Database `mapstructure:"database"`
	CaminoNode      string   `mapstructure:"caminoNode"`
}

type Database struct {
	Dsn  string `mapstructure:"dsn"`
	Type string `mapstructure:"type"`
}

var lock = &sync.Mutex{}

var configInstance *Config

func GetInstance() *Config {
	var configName = "config.yml" // default config name
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
			config, err := loadConfig(".", configName)
			if err != nil {
				log.Fatal(err)
			}
			configInstance = &config
		}
	}
	return configInstance
}

func loadConfig(path string, name string) (config Config, err error) {
	log.Printf("Loading config %s/%s", path, name)
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
