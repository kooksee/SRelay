package main

import (
	"encoding/json"
	"flag"
	"runtime"

	"github.com/kooksee/log"
	"github.com/kooksee/srelay/config"
	"github.com/kooksee/srelay/tserver"
)

const Version = "1.0"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()

	cfg := config.GetCfg()
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "debug mode")
	flag.StringVar(&cfg.LogLevel, "level", cfg.LogLevel, "log level")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "app host")
	flag.IntVar(&cfg.KcpPort, "port", cfg.KcpPort, "kcp port")
	flag.IntVar(&cfg.HttpPort, "httpport", cfg.HttpPort, "http port")
	flag.IntVar(&cfg.UdpPort, "udpport", cfg.UdpPort, "udp port")
	flag.Parse()

	d, _ := json.Marshal(cfg)
	log.Info(string(d))

	tserver.SetCfg(cfg)

	ks := tserver.NewKcpServer()
	ks.Listen()
	ks.Start()

	tserver.RunHttpServer()

	select {}

}
