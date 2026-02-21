package performance

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/network"
	"nogochain/network/config"
)

// BenchmarkNetwork_BroadcastBlock 测试广播区块的性能
func BenchmarkNetwork_BroadcastBlock(b *testing.B) {
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
		b.Fatalf("启动网络失败: %v", err)
	}

	// 创建区块
	block := &types.Block{
		Transactions: []*types.Transaction{},
	}

	// 重置计时器
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 广播区块
		n.BroadcastBlock(block)
	}
}

// BenchmarkNetwork_BroadcastTransaction 测试广播交易的性能
func BenchmarkNetwork_BroadcastTransaction(b *testing.B) {
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
		b.Fatalf("启动网络失败: %v", err)
	}

	// 创建交易
	tx := &types.Transaction{
		Nonce:    0,
		GasPrice: big.NewInt(1),
		Gas:      21000,
		To:       &common.Address{},
		Value:    big.NewInt(1000),
		Data:     nil,
	}

	// 重置计时器
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 广播交易
		n.BroadcastTransaction(tx)
	}
}

// BenchmarkNetwork_GetPeers 测试获取对等节点的性能
func BenchmarkNetwork_GetPeers(b *testing.B) {
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
		b.Fatalf("启动网络失败: %v", err)
	}

	// 重置计时器
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 获取对等节点
		_ = n.GetPeers()
	}
}
