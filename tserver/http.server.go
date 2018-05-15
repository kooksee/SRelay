package tserver

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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
	router := httprouter.New()
	router.POST("/:sid", indexPost)
	router.GET("/:sid", indexGet)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort), router); err != nil {
		panic(err)
	}
}
