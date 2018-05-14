package tserver

import (
	"bytes"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"

	kts "SRelay/types"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("tcp client conneted", conn.RemoteAddr().String())

	go func(conn *websocket.Conn) {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Error(err.Error())
				break
			}

			if messageType != websocket.TextMessage {
				conn.WriteMessage(websocket.TextMessage, kts.ResultError("数据类型错误"))
				continue
			}

			p = bytes.Trim(p, "\n")

			// 解析请求数据
			msg := &kts.KMsg{}
			if err := jsoniter.Unmarshal(p, msg); err != nil {
				conn.WriteMessage(websocket.TextMessage, kts.ResultError(err.Error()))
				return
			}

			switch msg.Event {
			}
		}
	}(conn)
}