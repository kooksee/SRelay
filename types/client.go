package types

import "fmt"

func DecodeKClient(data []byte) (KClient, error) {
	c := KClient{}
	return c, JsonUnmarshal(data, &c)
}

type KClient struct {
	Sign []byte `json:"sign,omitempty"`
}

func (c KClient) Bytes() []byte {
	d, _ := JsonMarshal(c)
	return append(d, "\n"...)
}

func DecodeKMsg(data []byte) (KMsg, error) {
	c := KMsg{}
	return c, JsonUnmarshal(data, &c)
}

type KMsg struct {
	TID string `json:"tid,omitempty"`
}

func (c KMsg) Bytes() []byte {
	d, _ := JsonMarshal(c)
	return append(d, "\n"...)
}

type ErrCode struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func (e ErrCode) Bytes() []byte {
	d, _ := JsonMarshal(e)
	return append(d, "\n"...)
}

func ErrJsonParse(err error) []byte {
	return ErrCode{
		Code: 10000,
		Msg:  fmt.Sprintf("json parse error,%s", err.Error()),
	}.Bytes()
}

func ErrPeerNotFound(err error) []byte {
	return ErrCode{
		Code: 10001,
		Msg:  fmt.Sprintf("the peer is nonexistent,%s", err.Error()),
	}.Bytes()
}

func ErrPeerWrite(err error) []byte {
	return ErrCode{
		Code: 10002,
		Msg:  fmt.Sprintf("peer write error,%s", err.Error()),
	}.Bytes()
}

func ErrNotWhitelist(err error) []byte {
	return ErrCode{
		Code: 10003,
		Msg:  fmt.Sprintf("not in white list,%s", err.Error()),
	}.Bytes()
}

func ErrSignError(err error) []byte {
	return ErrCode{
		Code: 10004,
		Msg:  fmt.Sprintf("sign error,%s", err.Error()),
	}.Bytes()
}
