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

var IsTest = false

func GetInstance() *Config {
	// todo: remove this hack
	var configName = "config.yml" // default config name
	if IsTest {
		configName = "config-test.yml"
	}
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
			log.Print("Loading config")
			config, err := loadConfig(GetRootPath(), configName)
			if err != nil {
				log.Fatal(err)
			}
			configInstance = &config
		}
	}
	log.Print("Reusing config")
	log.Println(configInstance)
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
