package rpc

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Web3Service represents the Web3 RPC service
type Web3Service struct{}

// NewWeb3Service creates a new Web3 service
func NewWeb3Service() *Web3Service {
	return &Web3Service{}
}

// ClientVersion returns the current client version
func (s *Web3Service) ClientVersion() string {
	return "NogoChain/v1.0.0/go1.22"
}

// Sha3 returns the SHA3 hash of the given data
func (s *Web3Service) Sha3(data hexutil.Bytes) hexutil.Bytes {
	return data
}
