package integration

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/network"
	"nogochain/network/config"
)

// TestNetwork_Integration_BroadcastBlock 测试广播区块的集成功能
func TestNetwork_Integration_BroadcastBlock(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Port: 8545,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8545,
		},
	}

	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 创建网络
	n := network.NewNetwork(cfg, bc)

	// 启动网络
	err := n.Start()
	if err != nil {
		t.Fatalf("启动网络失败: %v", err)
	}

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

	// 添加区块到区块链
	err = bc.AddBlock(block)
	if err != nil {
		t.Fatalf("添加区块失败: %v", err)
	}

	// 广播区块
	n.BroadcastBlock(block)
}

// TestNetwork_Integration_BroadcastTransaction 测试广播交易的集成功能
func TestNetwork_Integration_BroadcastTransaction(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Port: 8545,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8545,
		},
	}

	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 创建网络
	n := network.NewNetwork(cfg, bc)

	// 启动网络
	err := n.Start()
	if err != nil {
		t.Fatalf("启动网络失败: %v", err)
	}

	// 创建交易
	tx := types.NewTransaction(
		0,
		common.Address{},
		big.NewInt(1000),
		21000,
		big.NewInt(1),
		nil,
	)

	// 广播交易
	n.BroadcastTransaction(tx)
}

// TestNetwork_Integration_GetPeers 测试获取对等节点的集成功能
func TestNetwork_Integration_GetPeers(t *testing.T) {
	// 创建配置
	cfg := &config.Config{
		Port: 8545,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8545,
		},
	}

	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 创建网络
	n := network.NewNetwork(cfg, bc)

	// 启动网络
	err := n.Start()
	if err != nil {
		t.Fatalf("启动网络失败: %v", err)
	}

	// 获取对等节点
	peers := n.GetPeers()
	if peers == nil {
		t.Fatal("对等节点为nil")
	}
}
