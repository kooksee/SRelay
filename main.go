package main

import (
	"errors"
	"fmt"

	_ "SRelay/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
)

const Version = "1.0"

type DatabaseCheck struct {
}

func (dc *DatabaseCheck) Check() error {
	return errors.New("can't connect database")
}

func main() {

	// beego.LoadAppConfig("ini", "kdata/config/knode.ini")

	// beego.AppConfig.DefaultString()
	beego.BConfig.Listen.EnableAdmin = true
	beego.BConfig.Log.FileLineNum = true
	beego.BConfig.Log.AccessLogsFormat = "JSON_FORMAT"
	toolbox.AddHealthCheck("database", &DatabaseCheck{})

	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async(1e3)

	if beego.BConfig.RunMode == beego.DEV {
		logs.SetLogger(logs.AdapterConsole, `{"level":1}`)
		beego.SetLevel(beego.LevelDebug)
	}

	if beego.BConfig.RunMode == beego.PROD {
		logs.SetLogger(logs.AdapterFile, `{"filename":"knode.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
		beego.SetLevel(beego.LevelError)
	}

	toolbox.AddTask("tk1", toolbox.NewTask("tk1", "0/30 * * * * *", func() error { fmt.Println("tk1"); return nil }))
	toolbox.StartTask()
	defer toolbox.StopTask()

	beego.Run()
}
