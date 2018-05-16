package main

import (
	"encoding/json"
	"flag"
	"runtime"

	"github.com/kooksee/log"
	"github.com/kooksee/srelay/config"
	"github.com/kooksee/srelay/server"
)

const Version = "1.0"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()

	cfg := config.GetCfg()
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "debug mode")
	flag.StringVar(&cfg.Crypt, "crypt", cfg.Crypt, "crypt")
	flag.StringVar(&cfg.Key, "key", cfg.Key, "key")
	flag.StringVar(&cfg.Salt, "salt", cfg.Salt, "salt")
	flag.StringVar(&cfg.LogLevel, "level", cfg.LogLevel, "log level")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "app host")
	flag.IntVar(&cfg.KcpPort, "kcp", cfg.KcpPort, "kcp port")
	flag.IntVar(&cfg.HttpPort, "http", cfg.HttpPort, "http port")
	flag.Parse()

	d, _ := json.Marshal(cfg)
	log.Info(string(d))

	server.SetCfg(cfg)
	server.GetKcpServer().Listen()
	server.RunHttpServer()

	select {}

}
