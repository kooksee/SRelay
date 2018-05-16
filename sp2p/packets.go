package sp2p

import (
	"github.com/kooksee/srelay/types"
	"github.com/kooksee/uspnet/common"
)

func findNode(p *SP2p, msg *types.KMsg) {
	s := msg.Data.(types.FindNodeReq)
	nodes := p.tab.FindMinDisNodes(common.StringToHash(s.NID), s.N)
	ns := make([]string, 0)
	for _, n := range nodes {
		ns = append(ns, n.String())
	}
	m := &types.KMsg{
		TAddr: msg.FAddr,
		Data:  types.FindNodeResp{Nodes: ns},
	}
	if err := p.write(m); err != nil {
		logger.Error(err.Error())
		return
	}

	if node := NodeFromKMsg(msg); node != nil {
		p.tab.UpdateNode(node)
	}
}
func ping(p *SP2p, msg *types.KMsg) {
	if node := NodeFromKMsg(msg); node != nil {
		p.tab.UpdateNode(node)
	}
}

func init() {
	hm.Registry("findNode", findNode)
	hm.Registry("ping", ping)
}
