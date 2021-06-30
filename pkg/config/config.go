package config

import (
	"github.com/jinzhu/configor"
	"log"
)

func LoadConfig(path string) *Config {
	var cfg Config
	err := configor.New(&configor.Config{}).Load(&cfg, path)
	if err != nil {
		log.Panicf("Couldn't read config file")
	}
	return &cfg
}
