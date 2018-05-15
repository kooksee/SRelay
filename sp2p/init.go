package sp2p

import (
	"time"

	"github.com/json-iterator/go"
	"github.com/kooksee/log"
	"github.com/kooksee/srelay/handlers"
)

const (
	connReadTimeout = 10 * time.Second
)

var (
	cfg    *KConfig
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger = log.New("package", "sp2p")
	hm     = handlers.GetHManager()
)

func SetCfg(cfg1 *KConfig) {
	cfg = cfg1
}

func init() {
}
