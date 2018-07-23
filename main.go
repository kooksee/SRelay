package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kooksee/srelay/config"
	"github.com/kooksee/srelay/nat"
	"github.com/kooksee/srelay/server"
)

const Version = "1.0"

func runApp(port int64) {
	us := &server.UdpServer{}
	if err := us.Listen(port); err != nil {
		panic(err.Error())
	}

	ts := &server.TcpServer{}
	if err := ts.Listen(port); err != nil {
		panic(err.Error())
	}

	fmt.Println("listen tcp", "0.0.0.0", port)
	fmt.Println("listen udp", "0.0.0.0", port)
}

func main() {

	cfg := config.GetCfg()
	flag.BoolVar(&cfg.Debug, "d", cfg.Debug, "debug mode")
	flag.Int64Var(&cfg.Port, "p", cfg.Port, "app port")
	flag.BoolVar(&cfg.Nat, "nat", cfg.Nat, "is pnp or pmp")
	flag.StringVar(&cfg.Whitelist, "wl", cfg.Whitelist, "white list file")
	flag.Parse()
	cfg.InitLog()

	var a nat.Interface
	if cfg.Nat {
		a = nat.Any()
		ip, err := a.ExternalIP()
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("ext ip", ip.String())

		if err := a.AddMapping("tcp", int(cfg.Port), int(cfg.Port), "srealy", time.Hour*24*365); err != nil {
			panic(err.Error())
		}
		if err := a.AddMapping("udp", int(cfg.Port), int(cfg.Port), "srealy", time.Hour*24*365); err != nil {
			panic(err.Error())
		}
	}

	server.Init()
	runApp(cfg.Port)

	// 处理程序退出问题
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Printf("captured %v, exiting...\n", sig)

			// 退出程序之后删除端口映射
			if a != nil {
				if err := a.DeleteMapping("udp", int(cfg.Port), int(cfg.Port)); err != nil {
					panic(err.Error())
				}
				if err := a.DeleteMapping("udp", int(cfg.Port), int(cfg.Port)); err != nil {
					panic(err.Error())
				}
			}

			os.Exit(1)
		}
	}()
	select {}
}
