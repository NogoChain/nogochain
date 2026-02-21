package rpc

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NetService represents the Net RPC service
type NetService struct{}

// NewNetService creates a new Net service
func NewNetService() *NetService {
	return &NetService{}
}

// Version returns the current network protocol version
func (s *NetService) Version() string {
	return "318"
}

// Listening returns whether the node is listening for connections
func (s *NetService) Listening() bool {
	return false
}

// PeerCount returns the number of peers
func (s *NetService) PeerCount() hexutil.Uint {
	return hexutil.Uint(0)
}
