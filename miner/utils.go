package miner

import (
	"encoding/binary"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"nogochain/core/types"
)

// SerializeHeader 序列化区块头
func SerializeHeader(header *types.BlockHeader) []byte {
	// 简化实现，实际需要更详细的序列化
	var data []byte

	// ParentHash
	data = append(data, header.ParentHash.Bytes()...)

	// UncleHash
	data = append(data, header.UncleHash.Bytes()...)

	// Coinbase
	data = append(data, header.Coinbase.Bytes()...)

	// Root
	data = append(data, header.Root.Bytes()...)

	// TxHash
	data = append(data, header.TxHash.Bytes()...)

	// ReceiptHash
	data = append(data, header.ReceiptHash.Bytes()...)

	// Bloom
	data = append(data, header.Bloom...)

	// Difficulty
	diffBytes := header.Difficulty.Bytes()
	data = append(data, make([]byte, 32-len(diffBytes))...)
	data = append(data, diffBytes...)

	// Number
	numBytes := header.Number.Bytes()
	data = append(data, make([]byte, 32-len(numBytes))...)
	data = append(data, numBytes...)

	// GasLimit
	gasLimitBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(gasLimitBytes, header.GasLimit)
	data = append(data, gasLimitBytes...)

	// GasUsed
	gasUsedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(gasUsedBytes, header.GasUsed)
	data = append(data, gasUsedBytes...)

	// Time
	timeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeBytes, header.Time)
	data = append(data, timeBytes...)

	// Extra
	data = append(data, header.Extra...)

	return data
}

// CalculateSeed 计算种子
func CalculateSeed(number *big.Int) []byte {
	// 简化实现，实际需要基于区块号计算种子
	seed := make([]byte, 32)
	numBytes := number.Bytes()
	copy(seed[32-len(numBytes):], numBytes)
	return seed
}

// CalculateTarget 计算目标值
func CalculateTarget(difficulty *big.Int) []byte {
	// 计算目标值：2^256 / difficulty
	target := new(big.Int)
	target.Exp(big.NewInt(2), big.NewInt(256), nil)
	target.Div(target, difficulty)

	targetBytes := target.Bytes()
	result := make([]byte, 32)
	copy(result[32-len(targetBytes):], targetBytes)

	return result
}

// ValidateShare 验证份额
func ValidateShare(header *types.BlockHeader, nonce uint64, mixDigest common.Hash, difficulty *big.Int) bool {
	// TODO: 实现份额验证
	return true
}

// EstimateHashRate 估算哈希率
func EstimateHashRate(elapsed time.Duration, attempts uint64) float64 {
	if elapsed == 0 {
		return 0
	}

	return float64(attempts) / elapsed.Seconds()
}
