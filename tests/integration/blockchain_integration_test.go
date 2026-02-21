package integration

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
)

// TestBlockchain_Integration_AddBlock 测试添加区块的集成功能
func TestBlockchain_Integration_AddBlock(t *testing.T) {
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

	// 获取区块
	retrievedBlock := bc.GetBlock(block.Hash())
	if retrievedBlock == nil {
		t.Fatal("获取的区块为nil")
	}
}

// TestBlockchain_Integration_AddBlockWithTransactions 测试添加包含交易的区块
func TestBlockchain_Integration_AddBlockWithTransactions(t *testing.T) {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 创建交易
	tx := types.NewTransaction(
		0,
		common.Address{},
		big.NewInt(1000),
		21000,
		big.NewInt(1),
		nil,
	)

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
		21000,
		uint64(1001),
		[]byte{},
		common.Hash{},
		0,
		[]*types.Transaction{tx},
		[]*types.BlockHeader{},
	)

	// 添加区块
	err := bc.AddBlock(block)
	if err != nil {
		t.Fatalf("添加区块失败: %v", err)
	}

	// 获取区块
	retrievedBlock := bc.GetBlock(block.Hash())
	if retrievedBlock == nil {
		t.Fatal("获取的区块为nil")
	}

	if len(retrievedBlock.Transactions) != 1 {
		t.Errorf("交易数量错误，期望1，实际%v", len(retrievedBlock.Transactions))
	}
}
