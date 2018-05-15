package sp2p

import (
	"github.com/kooksee/uspnet/common"
)

type IKRpc interface {
	// 关闭通信
	Close() error
	// 对一个节点发起通信
	Ping(*Node) error
	// 查找节点
	FindNode(*Node) error
	// 发送消息,单播，多播，广播
	SendTx(*Tx) error
	// 查找信息,主要用于分片查找
	FindTx(*Tx) error
	// 获得路由表
	GetTable() *Table
}

type ITable interface {
	// 路由表大小
	Size() int
	// 获得节点列表,把节点列表转换为[enode://<hex node id>@10.3.58.6:30303?discport=30301]的方式
	GetTableList() []string
	// 添加节点
	AddNode(*Node)
	// 更新节点
	UpdateNode(*Node)
	// 删除节点
	DeleteNode(common.Hash)
	// 随机得到路由表中的n个节点
	FindRandomNodes(int) []*Node
	// 查找距离最近的n个节点
	FindMinDisNodes(common.Hash, int) []*Node
	// 查找相比另一个节点的更近的节点
	FindNodeWithTarget(common.Hash, common.Hash) []*Node
}
