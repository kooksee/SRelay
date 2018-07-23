package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/inconshreveable/log15"
	"github.com/patrickmn/go-cache"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Debug     bool
	Nat       bool
	Whitelist string
	Host      string
	Port      int64

	Cache *cache.Cache

	l  log15.Logger
	wl map[string]interface{}
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

// 地址白名单
func (t *Config) InitWhitelist() {
	if t.Whitelist == "" {
		t.wl = nil
		return
	}

	dt, err := ioutil.ReadFile(t.Whitelist)
	if err != nil {
		panic(err.Error())
	}

	t.wl = make(map[string]interface{})
	if err := json.Unmarshal(dt, &t.wl); err != nil {
		panic(err.Error())
	}
}

// IsWhitelist 是否属于白名单
func (t *Config) IsWhitelist(k string) bool {
	if t.wl == nil {
		return true
	} else {
		_, ok := t.wl[k]
		return ok
	}
}

// CheckAddress 签名是否正确以及是否在白名单中
func (t *Config) CheckAddress(id string, sign []byte) bool {
	pk, err := crypto.Ecrecover(common.Hex2Bytes(id), sign)
	if err != nil {
		Log().Error(err.Error())
		return false
	}
	addr := common.BytesToAddress(pk).Hex()
	return t.IsWhitelist(addr)
}

func Log() log15.Logger {
	if GetCfg().l == nil {
		panic("please init srelay log")
	}
	return GetCfg().l
}
