package types

import "fmt"

func DecodeClient(data []byte) (KMsg, error) {
	c := KMsg{}
	return c, json.Unmarshal(data, &c)
}

type KMsg struct {
	ID   string `json:"id,omitempty"`
	Addr string `json:"addr,omitempty"`
	Data []byte `json:"data,omitempty"`
}

func (c KMsg) Bytes() []byte {
	d, _ := json.Marshal(c)
	return append(d, "\n"...)
}

type ErrCode struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func (e ErrCode) Bytes() []byte {
	d, _ := json.Marshal(e)
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
		Code: 10001,
		Msg:  fmt.Sprintf("peer write error,%s", err.Error()),
	}.Bytes()
}
