package network

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/interfaces"
	"nogochain/network/config"
)

// 测试NewNetwork函数
func TestNewNetwork(t *testing.T) {
	// 创建网络配置
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8545,
			Host:    "127.0.0.1",
		},
	}

	// 创建区块链实例
	bc := blockchain.NewBlockchain(nil)

	// 创建网络实例
	network := NewNetwork(cfg, bc)
	if network == nil {
		t.Errorf("NewNetwork returned nil")
	}

	if network.config != cfg {
		t.Errorf("Network config mismatch")
	}

	if network.blockchain != bc {
		t.Errorf("Network blockchain mismatch")
	}

	if len(network.peers) != 0 {
		t.Errorf("Initial peers map should be empty")
	}

	if network.isStarted {
		t.Errorf("Network should not be started initially")
	}

	if network.rpcServer != nil {
		t.Errorf("RPC server should be nil when RPC is disabled")
	}

	// 测试启用RPC的情况
	cfgWithRPC := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: true,
			Port:    8545,
			Host:    "127.0.0.1",
		},
	}

	networkWithRPC := NewNetwork(cfgWithRPC, bc)
	if networkWithRPC.rpcServer == nil {
		t.Errorf("RPC server should not be nil when RPC is enabled")
	}
}

// 测试IsStarted方法
func TestIsStarted(t *testing.T) {
	// 创建网络实例
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 测试初始状态
	if network.IsStarted() {
		t.Errorf("Network should not be started initially")
	}

	// 测试启动后状态
	err := network.Start()
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}

	if !network.IsStarted() {
		t.Errorf("Network should be started after Start")
	}

	// 测试停止后状态
	err = network.Stop()
	if err != nil {
		t.Errorf("Stop returned error: %v", err)
	}

	if network.IsStarted() {
		t.Errorf("Network should not be started after Stop")
	}
}

// 测试GetConfig方法
func TestGetConfig(t *testing.T) {
	// 创建网络配置
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 测试获取配置
	retrievedConfig := network.GetConfig()
	if retrievedConfig != cfg {
		t.Errorf("GetConfig should return the original config")
	}

	if retrievedConfig.Port != cfg.Port {
		t.Errorf("Config port mismatch")
	}
}

// 测试AddPeer、RemovePeer、GetPeers、GetPeer方法
func TestPeerManagement(t *testing.T) {
	// 创建网络实例
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 创建测试对等节点
	peerID := [32]byte{0x01}
	peer := &interfaces.Peer{
		ID:       peerID,
		Head:     common.Hash{},
		Td:       1000,
		BlockNum: 10,
	}

	peerIDStr := "0100000000000000000000000000000000000000000000000000000000000000"

	// 测试添加对等节点
	network.AddPeer(peer)
	peers := network.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Peers length should be 1 after adding one peer, got %d", len(peers))
	}

	// 测试获取对等节点
	retrievedPeer := network.GetPeer(peerIDStr)
	if retrievedPeer == nil {
		t.Errorf("GetPeer should return the added peer")
	}

	if retrievedPeer.ID != peer.ID {
		t.Errorf("Peer ID mismatch")
	}

	// 测试移除对等节点
	network.RemovePeer(peerIDStr)
	peers = network.GetPeers()
	if len(peers) != 0 {
		t.Errorf("Peers length should be 0 after removing the peer, got %d", len(peers))
	}

	// 测试获取不存在的对等节点
	nonExistentPeer := network.GetPeer("non-existent")
	if nonExistentPeer != nil {
		t.Errorf("GetPeer should return nil for non-existent peer")
	}

	// 测试移除不存在的对等节点（应该无错误）
	network.RemovePeer("non-existent")
	peers = network.GetPeers()
	if len(peers) != 0 {
		t.Errorf("Peers length should still be 0 after trying to remove non-existent peer, got %d", len(peers))
	}
}

// 测试BroadcastBlock方法
func TestBroadcastBlock(t *testing.T) {
	// 创建网络实例
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 创建测试区块
	testBlock := &types.Block{
		Header: &types.BlockHeader{
			ParentHash:  common.Hash{},
			Coinbase:    common.Address{0x01},
			Root:        common.Hash{},
			TxHash:      common.Hash{},
			ReceiptHash: common.Hash{},
			Difficulty:  big.NewInt(1000000),
			Number:      big.NewInt(1),
			GasLimit:    10000000,
			GasUsed:     0,
			Time:        1700000000,
			Extra:       []byte("test block"),
			MixDigest:   common.Hash{},
			Nonce:       0,
		},
		Transactions: []*types.Transaction{},
		Uncles:       []*types.BlockHeader{},
	}

	// 测试广播区块（应该无错误）
	network.BroadcastBlock(testBlock)

	// 添加对等节点后测试广播
	peerID := [32]byte{0x01}
	peer := &interfaces.Peer{
		ID:       peerID,
		Head:     common.Hash{},
		Td:       1000,
		BlockNum: 10,
	}

	network.AddPeer(peer)
	network.BroadcastBlock(testBlock)
}

// 测试BroadcastTransaction方法
func TestBroadcastTransaction(t *testing.T) {
	// 创建网络实例
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 创建测试交易
	testTx := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1000),     // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	// 测试广播交易（应该无错误）
	network.BroadcastTransaction(testTx)

	// 添加对等节点后测试广播
	peerID := [32]byte{0x01}
	peer := &interfaces.Peer{
		ID:       peerID,
		Head:     common.Hash{},
		Td:       1000,
		BlockNum: 10,
	}

	network.AddPeer(peer)
	network.BroadcastTransaction(testTx)
}

// 测试SyncBlocks方法
func TestSyncBlocks(t *testing.T) {
	// 创建网络实例
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 测试同步区块（应该无错误）
	err := network.SyncBlocks()
	if err != nil {
		t.Errorf("SyncBlocks returned error: %v", err)
	}
}

// 测试Start和Stop方法
func TestStartStop(t *testing.T) {
	// 创建网络实例
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 测试初始状态
	if network.IsStarted() {
		t.Errorf("Network should not be started initially")
	}

	// 测试启动网络
	err := network.Start()
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}

	if !network.IsStarted() {
		t.Errorf("Network should be started after Start")
	}

	// 测试重复启动（应该返回错误）
	err = network.Start()
	if err == nil {
		t.Errorf("Start should return error when already running")
	}

	// 测试停止网络
	err = network.Stop()
	if err != nil {
		t.Errorf("Stop returned error: %v", err)
	}

	if network.IsStarted() {
		t.Errorf("Network should not be started after Stop")
	}

	// 测试重复停止（应该返回错误）
	err = network.Stop()
	if err == nil {
		t.Errorf("Stop should return error when not running")
	}
}

// 集成测试：测试完整的网络操作流程
func TestNetworkIntegration(t *testing.T) {
	// 创建网络实例
	cfg := &config.Config{
		Port: 30303,
		RPC: &config.RPCConfig{
			Enabled: false,
		},
	}

	bc := blockchain.NewBlockchain(nil)
	network := NewNetwork(cfg, bc)

	// 1. 检查初始状态
	if network.IsStarted() {
		t.Errorf("Network should not be started initially")
	}

	// 2. 启动网络
	err := network.Start()
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}

	if !network.IsStarted() {
		t.Errorf("Network should be started after Start")
	}

	// 3. 管理对等节点
	peerID := [32]byte{0x01}
	peer := &interfaces.Peer{
		ID:       peerID,
		Head:     common.Hash{},
		Td:       1000,
		BlockNum: 10,
	}

	network.AddPeer(peer)
	peers := network.GetPeers()
	if len(peers) != 1 {
		t.Errorf("Should have 1 peer, got %d", len(peers))
	}

	// 4. 广播区块和交易
	testBlock := &types.Block{
		Header: &types.BlockHeader{
			ParentHash:  common.Hash{},
			Coinbase:    common.Address{0x01},
			Root:        common.Hash{},
			TxHash:      common.Hash{},
			ReceiptHash: common.Hash{},
			Difficulty:  big.NewInt(1000000),
			Number:      big.NewInt(1),
			GasLimit:    10000000,
			GasUsed:     0,
			Time:        1700000000,
			Extra:       []byte("test block"),
			MixDigest:   common.Hash{},
			Nonce:       0,
		},
		Transactions: []*types.Transaction{},
		Uncles:       []*types.BlockHeader{},
	}

	testTx := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1000),     // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	network.BroadcastBlock(testBlock)
	network.BroadcastTransaction(testTx)

	// 5. 同步区块
	err = network.SyncBlocks()
	if err != nil {
		t.Errorf("SyncBlocks returned error: %v", err)
	}

	// 6. 移除对等节点
	peerIDStr := "0100000000000000000000000000000000000000000000000000000000000000"
	network.RemovePeer(peerIDStr)
	peers = network.GetPeers()
	if len(peers) != 0 {
		t.Errorf("Should have 0 peers after removal, got %d", len(peers))
	}

	// 7. 停止网络
	err = network.Stop()
	if err != nil {
		t.Errorf("Stop returned error: %v", err)
	}

	if network.IsStarted() {
		t.Errorf("Network should not be started after Stop")
	}
}

// 测试网络配置
func TestNetworkConfig(t *testing.T) {
	// 测试不同的网络配置
	testConfigs := []struct {
		name       string
		config     *config.Config
		rpcEnabled bool
	}{
		{
			name: "RPC disabled",
			config: &config.Config{
				Port: 30303,
				RPC: &config.RPCConfig{
					Enabled: false,
					Port:    8545,
					Host:    "127.0.0.1",
				},
			},
			rpcEnabled: false,
		},
		{
			name: "RPC enabled",
			config: &config.Config{
				Port: 30303,
				RPC: &config.RPCConfig{
					Enabled: true,
					Port:    8545,
					Host:    "127.0.0.1",
				},
			},
			rpcEnabled: true,
		},
	}

	for _, tc := range testConfigs {
		bc := blockchain.NewBlockchain(nil)
		network := NewNetwork(tc.config, bc)

		if tc.rpcEnabled {
			if network.rpcServer == nil {
				t.Errorf("%s: RPC server should be initialized when RPC is enabled", tc.name)
			}
		} else {
			if network.rpcServer != nil {
				t.Errorf("%s: RPC server should not be initialized when RPC is disabled", tc.name)
			}
		}

		if network.config.Port != tc.config.Port {
			t.Errorf("%s: Port mismatch", tc.name)
		}

		if network.config.RPC.Port != tc.config.RPC.Port {
			t.Errorf("%s: RPC port mismatch", tc.name)
		}

		if network.config.RPC.Host != tc.config.RPC.Host {
			t.Errorf("%s: RPC host mismatch", tc.name)
		}
	}
}
