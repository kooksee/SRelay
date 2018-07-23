package server

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kooksee/srelay/types"
)

func indexPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, types.ResultError(err))
		return
	}
	message = bytes.Trim(message, "\n")
	logger.Debug("message data", "data", string(message))

	tx := &types.KMsg{}
	if err := json.Unmarshal(message, tx); err != nil {
		logger.Error("Unmarshal error", "err", err)
		fmt.Fprint(w, types.ResultError(err))
		return
	}
	ksInstance.Send(tx)
	fmt.Fprint(w, types.ResultOk())
}

func indexGet(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	sid := p.ByName("sid")
	d, _ := cfg.Cache.Get(sid)
	if d != nil {
		fmt.Fprint(w, string(d.([]byte)))
		return
	}
	fmt.Fprint(w, types.ResultError(errors.New("not found")))
}

func RunHttpServer() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		// 得到端口号以及数据，根据端口号把数据写入对应的app
		// 得到一个ID，然后再把相应的请求保存到ID中

		port := strings.Trim(req.URL.RawPath, "/")
		d, _ := ioutil.ReadAll(req.Body)
		//	 找到port对应的kcp client，然后把数据发送过去
	})
}
