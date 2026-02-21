package synchronizer

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/core/validator"
)

// SyncMode 同步模式
type SyncMode uint8

const (
	FullSync SyncMode = iota
	FastSync
	LightSync
)

// SyncState 同步状态
type SyncState struct {
	CurrentBlock  uint64
	HighestBlock  uint64
	StartingBlock uint64
	PulledStates  uint64
	KnownStates   uint64
}

// Peer 对等节点结构
type Peer struct {
	ID          string
	Head        common.Hash
	Td          uint64
	BlockNumber uint64
}

// Synchronizer 链同步器
type Synchronizer struct {
	blockchain *blockchain.Blockchain
	validator  *validator.Validator
	peers      map[string]*Peer
	mode       SyncMode
	state      SyncState
	mu         sync.RWMutex
}

// NewSynchronizer 创建新的同步器
func NewSynchronizer(blockchain *blockchain.Blockchain, mode SyncMode) *Synchronizer {
	return &Synchronizer{
		blockchain: blockchain,
		validator:  validator.NewValidator(),
		peers:      make(map[string]*Peer),
		mode:       mode,
		state: SyncState{
			CurrentBlock:  blockchain.CurrentHead().NumberU64(),
			HighestBlock:  blockchain.CurrentHead().NumberU64(),
			StartingBlock: blockchain.CurrentHead().NumberU64(),
			PulledStates:  0,
			KnownStates:   0,
		},
	}
}

// Start 开始同步
func (s *Synchronizer) Start() {
	go s.syncLoop()
}

// Stop 停止同步
func (s *Synchronizer) Stop() {
}

// syncLoop 同步循环
func (s *Synchronizer) syncLoop() {
	for {
		s.sync()
		time.Sleep(2 * time.Second)
	}
}

// sync 执行同步
func (s *Synchronizer) sync() {
	// 获取最佳对等节点
	bestPeer := s.getBestPeer()
	if bestPeer == nil {
		return
	}

	// 检查是否需要同步
	currentHead := s.blockchain.CurrentHead()
	if bestPeer.BlockNumber <= currentHead.NumberU64() {
		return
	}

	// 更新同步状态
	s.mu.Lock()
	s.state.HighestBlock = bestPeer.BlockNumber
	s.state.CurrentBlock = currentHead.NumberU64()
	s.mu.Unlock()

	// 同步区块
	s.syncBlocks(bestPeer)
}

// syncBlocks 同步区块
func (s *Synchronizer) syncBlocks(peer *Peer) {
	currentHead := s.blockchain.CurrentHead()
	startNumber := currentHead.NumberU64() + 1
	endNumber := peer.BlockNumber

	// 批量同步区块
	for i := startNumber; i <= endNumber; i++ {
		// 从对等节点获取区块
		block := s.fetchBlockFromPeer(peer, i)
		if block == nil {
			continue
		}

		// 验证区块
		parent := s.blockchain.GetBlock(block.ParentHash())
		if parent == nil {
			continue
		}

		// 验证区块
		if err := s.validator.ValidateBlock(block, parent, s.blockchain.StateDB()); err != nil {
			continue
		}

		// 添加区块到区块链
		if err := s.blockchain.AddBlock(block); err != nil {
			continue
		}

		// 更新同步状态
		s.mu.Lock()
		s.state.CurrentBlock = block.NumberU64()
		s.mu.Unlock()
	}
}

// fetchBlockFromPeer 从对等节点获取区块
func (s *Synchronizer) fetchBlockFromPeer(peer *Peer, number uint64) *types.Block {
	// 模拟从对等节点获取区块
	// 实际实现中应该通过P2P网络获取
	return nil
}

// getBestPeer 获取最佳对等节点
func (s *Synchronizer) getBestPeer() *Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bestPeer *Peer
	for _, peer := range s.peers {
		if bestPeer == nil || peer.BlockNumber > bestPeer.BlockNumber {
			bestPeer = peer
		}
	}
	return bestPeer
}

// AddPeer 添加对等节点
func (s *Synchronizer) AddPeer(peer *Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.peers[peer.ID] = peer
}

// RemovePeer 移除对等节点
func (s *Synchronizer) RemovePeer(peerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.peers, peerID)
}

// GetPeers 获取所有对等节点
func (s *Synchronizer) GetPeers() []*Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peers := make([]*Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	return peers
}

// GetSyncState 获取同步状态
func (s *Synchronizer) GetSyncState() SyncState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SetSyncMode 设置同步模式
func (s *Synchronizer) SetSyncMode(mode SyncMode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mode = mode
}

// GetSyncMode 获取同步模式
func (s *Synchronizer) GetSyncMode() SyncMode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mode
}
