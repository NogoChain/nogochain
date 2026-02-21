package rpc

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NogoService represents the Nogo RPC service
type NogoService struct{}

// NewNogoService creates a new Nogo service
func NewNogoService() *NogoService {
	return &NogoService{}
}

// GetDifficulty returns the current difficulty
func (s *NogoService) GetDifficulty() hexutil.Uint64 {
	return hexutil.Uint64(1000000)
}

// GetReward returns the current block reward
func (s *NogoService) GetReward() hexutil.Big {
	return hexutil.Big{}
}

// GetChainInfo returns the chain information
func (s *NogoService) GetChainInfo() map[string]interface{} {
	return map[string]interface{}{
		"chainId":    318,
		"symbol":     "NOGO",
		"decimals":   18,
		"consensus":  "NogoPow",
		"difficulty": "1000000",
	}
}

// GetMiningInfo returns the mining information
func (s *NogoService) GetMiningInfo() map[string]interface{} {
	return map[string]interface{}{
		"difficulty":      "1000000",
		"hashrate":        "0",
		"miner":           "0x0000000000000000000000000000000000000000",
		"networkHashrate": "0",
	}
}

// GetWork returns the mining work for NogoPow
func (s *NogoService) GetWork() []string {
	return []string{
		"0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
	}
}

// SubmitWork submits the mining work for NogoPow
func (s *NogoService) SubmitWork(nonce string, headerHash string, mixDigest string) bool {
	return false
}

// SubmitHashrate submits the hashrate for NogoPow
func (s *NogoService) SubmitHashrate(hashrate hexutil.Uint64, id string) bool {
	return false
}
