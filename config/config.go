package config

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/inconshreveable/log15"
	"github.com/patrickmn/go-cache"
)

func Log() log15.Logger {
	if GetCfg().l == nil {
		panic("please init srelay log")
	}
	return GetCfg().l
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
func (t *Config) CheckAddress(sign []byte) (string, bool) {
	pk, err := crypto.Ecrecover([]byte(""), sign)
	if err != nil {
		Log().Error("CheckAddress Ecrecover Error", "err", err)
		return "", false
	}
	addr := common.BytesToAddress(pk).Hex()
	return addr, t.IsWhitelist(addr)
}

// GetCache 获得缓存
func (t *Config) GetCache() *cache.Cache {
	if t.c == nil {
		panic("please init cache")
	}
	return t.c
}

func GetCfg() *Config {
	once.Do(func() {
		instance = &Config{
			Host:      "0.0.0.0",
			Port:      8081,
			Debug:     false,
			Nat:       false,
			Whitelist: "",
			c:         cache.New(time.Minute*5, 10*time.Minute),
		}
	})
	return instance
}
