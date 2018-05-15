package types

type FindNodeReq struct {
	N   int
	NID string
}
type FindNodeResp struct {
	Nodes []string
}
