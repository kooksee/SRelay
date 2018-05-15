package handlers

import (
	"github.com/json-iterator/go"
	"github.com/kooksee/srelay/config"
	"github.com/kooksee/srelay/sp2p"
	"github.com/kooksee/srelay/types"
)

type IHandler func(*sp2p.SP2p, *types.KMsg)

var (
	cfg  *config.Config
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func SetCfg(cfg1 *config.Config) {
	cfg = cfg1
}
