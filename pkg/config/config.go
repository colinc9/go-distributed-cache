package config

import (
	"github.com/spf13/viper"
	"sync"
)

var ins *Config
var once sync.Once

func GetDefaultInsCfg() *Config {
	once.Do(func() {
		ins = &Config{
			AppAddress:    viper.GetString("app.cache.port"),
		}
	})
	return ins
}

type Config struct {
	AppAddress    string
}