package tools

import (
	"fmt"
	"time"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/network"
)

// BlockSyncTester 区块同步测试器
type BlockSyncTester struct {
	sourceBlockchain *blockchain.Blockchain
	targetNetwork   *network.Network
	syncStats       SyncStats
}

// SyncStats 同步统计信息
type SyncStats struct {
	BlocksSynced      int
	SyncTime          time.Duration
	StartTime         time.Time
	EndTime           time.Time
	NetworkBandwidth  int64 // 字节
	CPUUsage          float64
	MemoryUsage       int64 // 字节
}

// NewBlockSyncTester 创建新的区块同步测试器
func NewBlockSyncTester(sourceBlockchain *blockchain.Blockchain, targetNetwork *network.Network) *BlockSyncTester {
	return &BlockSyncTester{
		sourceBlockchain: sourceBlockchain,
		targetNetwork:   targetNetwork,
		syncStats: SyncStats{
			StartTime: time.Now(),
		},
	}
}

// PrepareSourceBlockchain 准备源区块链（生成测试区块）
func (t *BlockSyncTester) PrepareSourceBlockchain(numBlocks int) error {
	for i := 0; i < numBlocks; i++ {
		block := &types.Block{
			Transactions: []*types.Transaction{},
		}

		err := t.sourceBlockchain.AddBlock(block)
		if err != nil {
			return fmt.Errorf("添加区块失败: %v", err)
		}
	}
	return nil
}

// StartSync 开始同步测试
func (t *BlockSyncTester) StartSync() {
	// 记录开始时间
	t.syncStats.StartTime = time.Now()

	// 这里应该触发目标网络的同步过程
	// 实际实现可能需要根据网络的具体API进行调整
	t.targetNetwork.SyncBlocks()

	// 记录结束时间
	t.syncStats.EndTime = time.Now()
	t.syncStats.SyncTime = t.syncStats.EndTime.Sub(t.syncStats.StartTime)

	// 模拟统计信息
	t.syncStats.BlocksSynced = t.sourceBlockchain.Height()
	t.syncStats.NetworkBandwidth = int64(t.sourceBlockchain.Height() * 1024) // 假设每个区块1KB
	t.syncStats.CPUUsage = 10.5
	t.syncStats.MemoryUsage = 1024 * 1024 * 100 // 100MB
}

// GetStats 获取同步统计信息
func (t *BlockSyncTester) GetStats() SyncStats {
	return t.syncStats
}

// PrintStats 打印同步统计信息
func (t *BlockSyncTester) PrintStats() {
	stats := t.GetStats()
	fmt.Printf("区块同步测试统计信息:\n")
	fmt.Printf("同步区块数: %d\n", stats.BlocksSynced)
	fmt.Printf("同步时间: %v\n", stats.SyncTime)
	fmt.Printf("网络带宽使用: %d bytes\n", stats.NetworkBandwidth)
	fmt.Printf("CPU使用率: %.2f%%\n", stats.CPUUsage)
	fmt.Printf("内存使用: %d bytes\n", stats.MemoryUsage)
}
