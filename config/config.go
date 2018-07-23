package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func GetCfg() *Config {
	once.Do(func() {
		instance = &Config{
			Host:  "0.0.0.0",
			Port:  8081,
			Debug: false,
			Cache: cache.New(time.Minute*5, 10*time.Minute),
		}
	})
	return instance
}
