package sync

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/interfaces"
)

// SyncManager 同步管理器
type SyncManager struct {
	network    interfaces.NetworkInterface
	blockchain *blockchain.Blockchain
	ctx        context.Context
	cancel     context.CancelFunc
	isRunning  bool
	isSyncing  bool
	mu         sync.RWMutex
}

// NewSyncManager 创建新的同步管理器
func NewSyncManager(net interfaces.NetworkInterface, bc *blockchain.Blockchain) *SyncManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &SyncManager{
		network:    net,
		blockchain: bc,
		ctx:        ctx,
		cancel:     cancel,
		isRunning:  false,
		isSyncing:  false,
	}
}

// Start 启动同步管理器
func (sm *SyncManager) Start() error {
	if sm.isRunning {
		return nil
	}

	sm.isRunning = true
	log.Println("Sync manager started")

	// 启动同步循环
	go sm.syncLoop()

	return nil
}

// Stop 停止同步管理器
func (sm *SyncManager) Stop() error {
	if !sm.isRunning {
		return nil
	}

	sm.cancel()
	sm.isRunning = false
	sm.isSyncing = false
	log.Println("Sync manager stopped")

	return nil
}

// StartSync 开始同步
func (sm *SyncManager) StartSync() error {
	if !sm.isRunning {
		return nil
	}

	sm.mu.Lock()
	if sm.isSyncing {
		sm.mu.Unlock()
		return nil
	}
	sm.isSyncing = true
	sm.mu.Unlock()

	// 执行同步
	go sm.sync()

	return nil
}

// syncLoop 同步循环
func (sm *SyncManager) syncLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-ticker.C:
			// 检查是否需要同步
			sm.checkSyncNeeded()
		}
	}
}

// checkSyncNeeded 检查是否需要同步
func (sm *SyncManager) checkSyncNeeded() {
	peers := sm.network.GetPeers()
	if len(peers) == 0 {
		return
	}

	// 获取最佳对等节点
	bestPeer := sm.getBestPeer(peers)
	if bestPeer == nil {
		return
	}

	// 检查区块高度差
	currentHeight := sm.blockchain.CurrentHead().NumberU64()
	if bestPeer.BlockNum > currentHeight {
		// 开始同步
		sm.StartSync()
	}
}

// sync 执行同步
func (sm *SyncManager) sync() {
	defer func() {
		sm.mu.Lock()
		sm.isSyncing = false
		sm.mu.Unlock()
	}()

	peers := sm.network.GetPeers()
	if len(peers) == 0 {
		return
	}

	// 获取最佳对等节点
	bestPeer := sm.getBestPeer(peers)
	if bestPeer == nil {
		return
	}

	currentHeight := sm.blockchain.CurrentHead().NumberU64()
	targetHeight := bestPeer.BlockNum

	log.Printf("Starting sync: current=%d, target=%d", currentHeight, targetHeight)

	// 执行区块同步
	sm.syncBlocks(bestPeer, currentHeight+1, targetHeight)

	log.Printf("Sync completed: current=%d", sm.blockchain.CurrentHead().NumberU64())
}

// syncBlocks 同步区块
func (sm *SyncManager) syncBlocks(peer *interfaces.Peer, startHeight, endHeight uint64) {
	// 批量同步区块
	const batchSize = 128

	for height := startHeight; height <= endHeight; height += batchSize {
		batchEnd := height + batchSize - 1
		if batchEnd > endHeight {
			batchEnd = endHeight
		}

		// 获取区块批次
		blocks := sm.fetchBlocks(peer, height, batchEnd)
		if len(blocks) == 0 {
			continue
		}

		// 处理区块
		for _, block := range blocks {
			// 验证区块
			parent := sm.blockchain.GetBlock(block.ParentHash())
			if parent == nil {
				// 父区块不存在，需要回溯
				height = block.NumberU64() - 1
				break
			}

			// 添加区块到区块链
			err := sm.blockchain.AddBlock(block)
			if err != nil {
				log.Printf("Failed to add block %d: %v", block.NumberU64(), err)
				continue
			}

			// 广播区块
			sm.network.BroadcastBlock(block)
		}
	}
}

// fetchBlocks 从对等节点获取区块
func (sm *SyncManager) fetchBlocks(peer *interfaces.Peer, startHeight, endHeight uint64) []*types.Block {
	// 模拟从对等节点获取区块
	// 实际实现中应该通过P2P协议发送请求
	var blocks []*types.Block

	for height := startHeight; height <= endHeight; height++ {
		// 这里应该发送区块请求并等待响应
		// 暂时创建模拟区块
		block := sm.createMockBlock(height)
		if block != nil {
			blocks = append(blocks, block)
		}
	}

	return blocks
}

// createMockBlock 创建模拟区块
func (sm *SyncManager) createMockBlock(height uint64) *types.Block {
	// 获取前一个区块
	var parentHash common.Hash
	if height > 0 {
		parent := sm.blockchain.GetBlockByNumber(height - 1)
		if parent != nil {
			parentHash = parent.Hash()
		}
	}

	// 创建新区块
	return types.NewBlock(
		parentHash,
		common.Address{},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(int64(height)),
		10000000,
		0,
		uint64(time.Now().Unix()),
		[]byte(fmt.Sprintf("Block %d", height)),
		common.Hash{},
		0,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)
}

// getBestPeer 获取最佳对等节点
func (sm *SyncManager) getBestPeer(peers []*interfaces.Peer) *interfaces.Peer {
	if len(peers) == 0 {
		return nil
	}

	var bestPeer *interfaces.Peer
	for _, peer := range peers {
		if bestPeer == nil || peer.BlockNum > bestPeer.BlockNum {
			bestPeer = peer
		}
	}

	return bestPeer
}

// IsSyncing 检查是否正在同步
func (sm *SyncManager) IsSyncing() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.isSyncing
}

// IsRunning 检查是否运行
func (sm *SyncManager) IsRunning() bool {
	return sm.isRunning
}

// GetSyncStatus 获取同步状态
func (sm *SyncManager) GetSyncStatus() (currentHeight, highestHeight uint64) {
	currentHeight = sm.blockchain.CurrentHead().NumberU64()

	// 获取最高区块高度
	peers := sm.network.GetPeers()
	bestPeer := sm.getBestPeer(peers)
	if bestPeer != nil {
		highestHeight = bestPeer.BlockNum
	} else {
		highestHeight = currentHeight
	}

	return
}
