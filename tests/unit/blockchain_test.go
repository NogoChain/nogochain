package unit

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"nogochain/core/blockchain"
	"nogochain/core/types"
)

// TestBlockchain_CreateBlockchain 测试创建区块链
func TestBlockchain_CreateBlockchain(t *testing.T) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 验证区块链不为空
	if bc == nil {
		t.Fatal("区块链为nil")
	}
}

// TestBlockchain_AddBlock 测试添加区块
func TestBlockchain_AddBlock(t *testing.T) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 创建区块
	genesis := bc.Genesis()
	block := types.NewBlock(
		genesis.Hash(),
		common.Address{},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000),
		big.NewInt(1),
		1000000,
		0,
		uint64(1001),
		[]byte{},
		common.Hash{},
		0,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	// 添加区块
	err := bc.AddBlock(block)
	if err != nil {
		t.Fatalf("添加区块失败: %v", err)
	}
}

// TestBlockchain_GetBlock 测试获取区块
func TestBlockchain_GetBlock(t *testing.T) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 创建区块
	genesis := bc.Genesis()
	block := types.NewBlock(
		genesis.Hash(),
		common.Address{},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000),
		big.NewInt(1),
		1000000,
		0,
		uint64(1001),
		[]byte{},
		common.Hash{},
		0,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	// 添加区块
	err := bc.AddBlock(block)
	if err != nil {
		t.Fatalf("添加区块失败: %v", err)
	}

	// 获取区块（通过哈希）
	retrievedBlock := bc.GetBlock(block.Hash())
	if retrievedBlock == nil {
		t.Fatal("获取的区块为nil")
	}
}

// TestBlockchain_GetLatestBlock 测试获取最新区块
func TestBlockchain_GetLatestBlock(t *testing.T) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 获取最新区块（创世区块）
	latestBlock := bc.CurrentHead()
	if latestBlock == nil {
		t.Fatal("获取的最新区块为nil")
	}

	if latestBlock.NumberU64() != 0 {
		t.Errorf("最新区块编号错误，期望0，实际%v", latestBlock.NumberU64())
	}
}
