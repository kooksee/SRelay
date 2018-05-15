package config

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Debug    bool
	LogLevel string
	Crypt    string
	Key      string
	Salt     string

	Host     string
	KcpPort  int
	HttpPort int
	UdpPort  int

	Cache *cache.Cache
}
