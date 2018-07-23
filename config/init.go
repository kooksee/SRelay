package config

import (
	"os"
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/patrickmn/go-cache"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Debug bool
	Nat bool

	Host string
	Port int64

	Cache *cache.Cache

	l log15.Logger
}

func (t *Config) InitLog() {
	l := log15.New("app", "srelay")
	if t.Debug {
		l.SetHandler(log15.LvlFilterHandler(log15.LvlDebug, log15.StreamHandler(os.Stdout, log15.TerminalFormat())))
	} else {
		l.SetHandler(log15.LvlFilterHandler(log15.LvlError, log15.StreamHandler(os.Stderr, log15.LogfmtFormat())))
	}
	t.l = l
}

func Log() log15.Logger {
	if GetCfg().l == nil {
		panic("please init srelay log")
	}
	return GetCfg().l
}
