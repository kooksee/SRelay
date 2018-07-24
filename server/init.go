package server

import (
	"github.com/inconshreveable/log15"
	"github.com/kooksee/srelay/config"
)

var (
	cfg    *config.Config
	logger log15.Logger
)

func Init() {
	cfg = config.GetCfg()
	clientsCache = cfg.Cache
	logger = config.Log().New("package", "server")
}
