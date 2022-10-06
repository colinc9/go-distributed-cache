package config

import (
	"github.com/spf13/viper"
	"sync"
)

var ins *Config
var once sync.Once

func GetDefaultInsCfg() *Config {
	once.Do(func() {
		SetUpViperDefault()
		ins = &Config{
			AppAddress:    viper.GetString("app.cache.address"),
		}
	})
	return ins
}

type Config struct {
	AppAddress    string
}

func SetUpViperDefault() {
	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		print(err)
	}
}