package types

import (
	"github.com/json-iterator/go"
)

type Resp struct {
	Code string      `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ResultOk() string {

	d, _ := jsoniter.MarshalToString(Resp{
		Code: "ok",
	})
	return d
}

func ResultOkWithData(data interface{}) string {

	d, _ := jsoniter.MarshalToString(Resp{
		Code: "ok",
		Data: data,
	})
	return d
}

func ResultError(msg error) string {
	d, _ := jsoniter.MarshalToString(Resp{
		Code: "error",
		Msg:  msg.Error(),
	})
	return d
}
