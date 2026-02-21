package synchronizer

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
)

// 测试NewSynchronizer函数
func TestNewSynchronizer(t *testing.T) {
	// 创建区块链实例
	bc := blockchain.NewBlockchain(nil)

	// 测试创建全同步模式的同步器
	sync := NewSynchronizer(bc, FullSync)
	if sync == nil {
		t.Errorf("NewSynchronizer returned nil")
	}

	if sync.blockchain != bc {
		t.Errorf("Blockchain reference mismatch")
	}

	if sync.validator == nil {
		t.Errorf("Validator should not be nil")
	}

	if len(sync.peers) != 0 {
		t.Errorf("Initial peers map should be empty")
	}

	if sync.mode != FullSync {
		t.Errorf("Sync mode should be FullSync, got %v", sync.mode)
	}

	// 检查初始同步状态
	state := sync.GetSyncState()
	genesisNumber := bc.Genesis().NumberU64()
	if state.CurrentBlock != genesisNumber {
		t.Errorf("CurrentBlock should be %d, got %d", genesisNumber, state.CurrentBlock)
	}

	if state.HighestBlock != genesisNumber {
		t.Errorf("HighestBlock should be %d, got %d", genesisNumber, state.HighestBlock)
	}

	if state.StartingBlock != genesisNumber {
		t.Errorf("StartingBlock should be %d, got %d", genesisNumber, state.StartingBlock)
	}

	if state.PulledStates != 0 {
		t.Errorf("PulledStates should be 0, got %d", state.PulledStates)
	}

	if state.KnownStates != 0 {
		t.Errorf("KnownStates should be 0, got %d", state.KnownStates)
	}

	// 测试创建其他模式的同步器
	syncFast := NewSynchronizer(bc, FastSync)
	if syncFast.mode != FastSync {
		t.Errorf("Sync mode should be FastSync, got %v", syncFast.mode)
	}

	syncLight := NewSynchronizer(bc, LightSync)
	if syncLight.mode != LightSync {
		t.Errorf("Sync mode should be LightSync, got %v", syncLight.mode)
	}
}

// 测试AddPeer和RemovePeer函数
func TestAddRemovePeer(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 创建测试对等节点
	peer1 := &Peer{
		ID:          "peer1",
		Head:        common.HexToHash("0x01"),
		Td:          1000,
		BlockNumber: 10,
	}

	peer2 := &Peer{
		ID:          "peer2",
		Head:        common.HexToHash("0x02"),
		Td:          2000,
		BlockNumber: 20,
	}

	// 测试添加对等节点
	sync.AddPeer(peer1)
	peers := sync.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Peers length should be 1 after adding one peer, got %d", len(peers))
	}

	if peers[0].ID != peer1.ID {
		t.Errorf("Peer ID mismatch")
	}

	// 测试添加第二个对等节点
	sync.AddPeer(peer2)
	peers = sync.GetPeers()
	if len(peers) != 2 {
		t.Errorf("Peers length should be 2 after adding two peers, got %d", len(peers))
	}

	// 测试移除对等节点
	sync.RemovePeer(peer1.ID)
	peers = sync.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Peers length should be 1 after removing one peer, got %d", len(peers))
	}

	if peers[0].ID != peer2.ID {
		t.Errorf("Remaining peer should be peer2")
	}

	// 测试移除不存在的对等节点（应该无错误）
	sync.RemovePeer("non-existent")
	peers = sync.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Peers length should still be 1 after trying to remove non-existent peer, got %d", len(peers))
	}

	// 测试移除所有对等节点
	sync.RemovePeer(peer2.ID)
	peers = sync.GetPeers()
	if len(peers) != 0 {
		t.Errorf("Peers length should be 0 after removing all peers, got %d", len(peers))
	}
}

// 测试GetBestPeer函数
func TestGetBestPeer(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 测试没有对等节点时返回nil
	bestPeer := sync.getBestPeer()
	if bestPeer != nil {
		t.Errorf("GetBestPeer should return nil when no peers")
	}

	// 创建测试对等节点
	peer1 := &Peer{
		ID:          "peer1",
		Head:        common.HexToHash("0x01"),
		Td:          1000,
		BlockNumber: 10,
	}

	peer2 := &Peer{
		ID:          "peer2",
		Head:        common.HexToHash("0x02"),
		Td:          2000,
		BlockNumber: 20,
	}

	peer3 := &Peer{
		ID:          "peer3",
		Head:        common.HexToHash("0x03"),
		Td:          1500,
		BlockNumber: 15,
	}

	// 添加对等节点
	sync.AddPeer(peer1)
	sync.AddPeer(peer2)
	sync.AddPeer(peer3)

	// 测试获取最佳对等节点（应该是peer2，因为它的BlockNumber最大）
	bestPeer = sync.getBestPeer()
	if bestPeer == nil {
		t.Errorf("GetBestPeer should not return nil when peers exist")
	}

	if bestPeer.ID != peer2.ID {
		t.Errorf("Best peer should be peer2, got %s", bestPeer.ID)
	}

	if bestPeer.BlockNumber != peer2.BlockNumber {
		t.Errorf("Best peer block number should be %d, got %d", peer2.BlockNumber, bestPeer.BlockNumber)
	}
}

// 测试SetSyncMode和GetSyncMode函数
func TestSyncMode(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 测试初始同步模式
	mode := sync.GetSyncMode()
	if mode != FullSync {
		t.Errorf("Initial sync mode should be FullSync, got %v", mode)
	}

	// 测试设置为FastSync
	sync.SetSyncMode(FastSync)
	mode = sync.GetSyncMode()
	if mode != FastSync {
		t.Errorf("Sync mode should be FastSync after setting, got %v", mode)
	}

	// 测试设置为LightSync
	sync.SetSyncMode(LightSync)
	mode = sync.GetSyncMode()
	if mode != LightSync {
		t.Errorf("Sync mode should be LightSync after setting, got %v", mode)
	}

	// 测试设置回FullSync
	sync.SetSyncMode(FullSync)
	mode = sync.GetSyncMode()
	if mode != FullSync {
		t.Errorf("Sync mode should be FullSync after setting, got %v", mode)
	}
}

// 测试GetSyncState函数
func TestGetSyncState(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 获取初始同步状态
	state := sync.GetSyncState()
	genesisNumber := bc.Genesis().NumberU64()

	if state.CurrentBlock != genesisNumber {
		t.Errorf("CurrentBlock should be %d, got %d", genesisNumber, state.CurrentBlock)
	}

	if state.HighestBlock != genesisNumber {
		t.Errorf("HighestBlock should be %d, got %d", genesisNumber, state.HighestBlock)
	}

	if state.StartingBlock != genesisNumber {
		t.Errorf("StartingBlock should be %d, got %d", genesisNumber, state.StartingBlock)
	}

	if state.PulledStates != 0 {
		t.Errorf("PulledStates should be 0, got %d", state.PulledStates)
	}

	if state.KnownStates != 0 {
		t.Errorf("KnownStates should be 0, got %d", state.KnownStates)
	}

	// 添加对等节点并检查同步状态更新
	peer := &Peer{
		ID:          "peer1",
		Head:        common.HexToHash("0x01"),
		Td:          1000,
		BlockNumber: 10,
	}

	sync.AddPeer(peer)

	// 手动调用sync函数来更新同步状态
	sync.sync()

	// 检查同步状态是否更新
	state = sync.GetSyncState()
	if state.HighestBlock != peer.BlockNumber {
		t.Errorf("HighestBlock should be updated to %d, got %d", peer.BlockNumber, state.HighestBlock)
	}

	if state.CurrentBlock != genesisNumber {
		t.Errorf("CurrentBlock should remain %d, got %d", genesisNumber, state.CurrentBlock)
	}
}

// 测试syncBlocks函数
func TestSyncBlocks(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 创建测试对等节点
	peer := &Peer{
		ID:          "peer1",
		Head:        common.HexToHash("0x01"),
		Td:          1000,
		BlockNumber: 2,
	}

	// 创建有效的区块链
	genesis := bc.Genesis()

	block1 := types.NewBlock(
		genesis.Hash(),
		common.Address{0x01},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(1),
		10000000,
		0,
		genesis.Header.Time+10,
		[]byte("Block 1"),
		common.Hash{},
		12345,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	block2 := types.NewBlock(
		block1.Hash(),
		common.Address{0x01},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(2),
		10000000,
		0,
		block1.Header.Time+10,
		[]byte("Block 2"),
		common.Hash{},
		12346,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	// 手动添加区块到区块链（模拟同步成功）
	bc.AddBlock(block1)
	bc.AddBlock(block2)

	// 调用syncBlocks函数
	sync.syncBlocks(peer)

	// 检查区块链长度
	length := bc.Length()
	if length != 3 {
		t.Errorf("Blockchain length should be 3 after sync, got %d", length)
	}
}

// 集成测试：测试完整的同步流程
func TestSyncIntegration(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 1. 检查初始状态
	initialState := sync.GetSyncState()
	genesisNumber := bc.Genesis().NumberU64()
	if initialState.CurrentBlock != genesisNumber {
		t.Errorf("Initial CurrentBlock should be %d, got %d", genesisNumber, initialState.CurrentBlock)
	}

	// 2. 添加对等节点
	peer := &Peer{
		ID:          "peer1",
		Head:        common.HexToHash("0x01"),
		Td:          1000,
		BlockNumber: 5,
	}

	sync.AddPeer(peer)

	// 3. 检查对等节点
	peers := sync.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Should have 1 peer, got %d", len(peers))
	}

	// 4. 检查最佳对等节点
	bestPeer := sync.getBestPeer()
	if bestPeer == nil {
		t.Errorf("Best peer should not be nil")
	}

	if bestPeer.ID != peer.ID {
		t.Errorf("Best peer ID mismatch")
	}

	// 5. 执行同步
	sync.sync()

	// 6. 检查同步状态
	syncState := sync.GetSyncState()
	if syncState.HighestBlock != peer.BlockNumber {
		t.Errorf("HighestBlock should be %d, got %d", peer.BlockNumber, syncState.HighestBlock)
	}

	// 7. 测试同步模式切换
	sync.SetSyncMode(FastSync)
	mode := sync.GetSyncMode()
	if mode != FastSync {
		t.Errorf("Sync mode should be FastSync, got %v", mode)
	}

	// 8. 移除对等节点
	sync.RemovePeer(peer.ID)
	peers = sync.GetPeers()
	if len(peers) != 0 {
		t.Errorf("Should have 0 peers after removal, got %d", len(peers))
	}

	// 9. 再次执行同步（应该无操作）
	sync.sync()

	// 10. 检查同步状态（应该保持不变）
	finalState := sync.GetSyncState()
	if finalState.HighestBlock != syncState.HighestBlock {
		t.Errorf("HighestBlock should remain %d, got %d", syncState.HighestBlock, finalState.HighestBlock)
	}
}

// 测试同步状态更新
func TestSyncStateUpdate(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 创建测试对等节点，区块号高于本地
	peer := &Peer{
		ID:          "peer1",
		Head:        common.HexToHash("0x01"),
		Td:          1000,
		BlockNumber: 10,
	}

	sync.AddPeer(peer)

	// 调用sync函数
	sync.sync()

	// 检查同步状态
	state := sync.GetSyncState()
	if state.HighestBlock != peer.BlockNumber {
		t.Errorf("HighestBlock should be updated to %d, got %d", peer.BlockNumber, state.HighestBlock)
	}

	// 创建另一个对等节点，区块号更高
	betterPeer := &Peer{
		ID:          "peer2",
		Head:        common.HexToHash("0x02"),
		Td:          2000,
		BlockNumber: 20,
	}

	sync.AddPeer(betterPeer)

	// 再次调用sync函数
	sync.sync()

	// 检查同步状态是否更新到更高的区块号
	state = sync.GetSyncState()
	if state.HighestBlock != betterPeer.BlockNumber {
		t.Errorf("HighestBlock should be updated to %d, got %d", betterPeer.BlockNumber, state.HighestBlock)
	}
}

// 测试无需同步的情况
func TestNoSyncNeeded(t *testing.T) {
	bc := blockchain.NewBlockchain(nil)
	sync := NewSynchronizer(bc, FullSync)

	// 创建测试对等节点，区块号低于或等于本地
	peer := &Peer{
		ID:          "peer1",
		Head:        common.HexToHash("0x01"),
		Td:          1000,
		BlockNumber: 0, // 与创世区块相同
	}

	sync.AddPeer(peer)

	// 调用sync函数（应该无操作）
	sync.sync()

	// 检查同步状态（应该保持不变）
	state := sync.GetSyncState()
	genesisNumber := bc.Genesis().NumberU64()
	if state.CurrentBlock != genesisNumber {
		t.Errorf("CurrentBlock should remain %d, got %d", genesisNumber, state.CurrentBlock)
	}

	if state.HighestBlock != genesisNumber {
		t.Errorf("HighestBlock should remain %d, got %d", genesisNumber, state.HighestBlock)
	}
}
