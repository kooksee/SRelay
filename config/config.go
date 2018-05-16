package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func GetCfg() *Config {
	once.Do(func() {
		instance = &Config{
			Crypt:    "aes-128",
			Key:      "hello",
			Salt:     "hello",
			Host:     "0.0.0.0",
			KcpPort:  8080,
			HttpPort: 8081,
			Debug:    true,
			LogLevel: "info",
			Cache:    cache.New(time.Minute, 5*time.Minute),
		}
	})
	return instance
}
