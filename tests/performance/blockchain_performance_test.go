package performance

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
)

// BenchmarkBlockchain_AddBlock 测试添加区块的性能
func BenchmarkBlockchain_AddBlock(b *testing.B) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 获取创世区块
	genesis := bc.Genesis()

	// 重置计时器
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 创建区块
		block := types.NewBlock(
			genesis.Hash(),
			common.Address{},
			common.Hash{},
			common.Hash{},
			common.Hash{},
			big.NewInt(1000000),
			big.NewInt(int64(i+1)),
			10000000,
			0,
			uint64(1700000000+i),
			[]byte("Test Block"),
			common.Hash{},
			0,
			[]*types.Transaction{},
			[]*types.BlockHeader{},
		)

		// 添加区块
		err := bc.AddBlock(block)
		if err != nil {
			b.Fatalf("添加区块失败: %v", err)
		}
	}
}

// BenchmarkBlockchain_GetBlock 测试获取区块的性能
func BenchmarkBlockchain_GetBlock(b *testing.B) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 获取创世区块
	genesis := bc.Genesis()

	// 添加100个区块
	for i := 0; i < 100; i++ {
		block := types.NewBlock(
			genesis.Hash(),
			common.Address{},
			common.Hash{},
			common.Hash{},
			common.Hash{},
			big.NewInt(1000000),
			big.NewInt(int64(i+1)),
			10000000,
			0,
			uint64(1700000000+i),
			[]byte("Test Block"),
			common.Hash{},
			0,
			[]*types.Transaction{},
			[]*types.BlockHeader{},
		)

		err := bc.AddBlock(block)
		if err != nil {
			b.Fatalf("添加区块失败: %v", err)
		}
	}

	// 重置计时器
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 获取区块
		_ = bc.GetBlockByNumber(uint64(i%100 + 1))
	}
}

// BenchmarkBlockchain_GetLatestBlock 测试获取最新区块的性能
func BenchmarkBlockchain_GetLatestBlock(b *testing.B) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 获取创世区块
	genesis := bc.Genesis()

	// 添加100个区块
	for i := 0; i < 100; i++ {
		block := types.NewBlock(
			genesis.Hash(),
			common.Address{},
			common.Hash{},
			common.Hash{},
			common.Hash{},
			big.NewInt(1000000),
			big.NewInt(int64(i+1)),
			10000000,
			0,
			uint64(1700000000+i),
			[]byte("Test Block"),
			common.Hash{},
			0,
			[]*types.Transaction{},
			[]*types.BlockHeader{},
		)

		err := bc.AddBlock(block)
		if err != nil {
			b.Fatalf("添加区块失败: %v", err)
		}
	}

	// 重置计时器
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 获取最新区块
		_ = bc.CurrentHead()
	}
}
