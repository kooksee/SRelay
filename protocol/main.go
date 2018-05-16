package protocol

type FindNodeReq struct {
	N   int
	NID string
}
type FindNodeResp struct {
	Nodes []string
}

type KV struct {
	K string      `json:"k,omitempty"`
	V interface{} `json:"v,omitempty"`
}

type KVGetReq struct {
	KV
}

type KVGetResp struct {
	KV
}