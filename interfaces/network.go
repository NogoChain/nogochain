package interfaces

import (
	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/types"
	"nogochain/network/config"
)

// NodeID 节点ID
type NodeID [32]byte

// NetworkInterface 网络接口
type NetworkInterface interface {
	AddPeer(peer *Peer)
	RemovePeer(id string)
	GetPeers() []*Peer
	GetPeer(id string) *Peer
	BroadcastBlock(block *types.Block)
	BroadcastTransaction(tx *types.Transaction)
	SyncBlocks() error
	IsStarted() bool
	GetNode() interface{}
	GetConfig() *config.Config
}

// Peer 对等节点
type Peer struct {
	ID       NodeID
	Node     interface{}
	Conn     interface{}
	Head     common.Hash
	Td       uint64
	BlockNum uint64
	LastSeen int64
}
