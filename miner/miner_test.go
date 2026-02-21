package miner

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"nogochain/consensus/nogopow"
	"nogochain/core/types"
)

// 测试NewMiner函数
func TestNewMiner(t *testing.T) {
	// 创建挖矿配置
	config := &Config{
		Enabled:          true,
		Coinbase:         common.Address{0x01},
		ExtraData:        []byte("NogoChain Miner"),
		MinGasPrice:      big.NewInt(1),
		MaxGasLimit:      10000000,
		RecommitInterval: 3 * time.Second,
		NumThreads:       2,
	}

	// 创建共识引擎
	engine := nogopow.NewNogoPow()

	// 创建挖矿实例
	miner := NewMiner(config, engine)
	if miner == nil {
		t.Errorf("NewMiner returned nil")
	}

	if miner.config != config {
		t.Errorf("Miner config mismatch")
	}

	if miner.engine != engine {
		t.Errorf("Miner engine mismatch")
	}

	if miner.config.NumThreads != 2 {
		t.Errorf("Miner num threads should be 2, got %d", miner.config.NumThreads)
	}

	// 测试默认线程数
	configWithNoThreads := &Config{
		Enabled:          true,
		Coinbase:         common.Address{0x01},
		ExtraData:        []byte("NogoChain Miner"),
		MinGasPrice:      big.NewInt(1),
		MaxGasLimit:      10000000,
		RecommitInterval: 3 * time.Second,
		NumThreads:       0, // 0表示使用默认值
	}

	minerWithDefaultThreads := NewMiner(configWithNoThreads, engine)
	if minerWithDefaultThreads.config.NumThreads <= 0 {
		t.Errorf("Miner num threads should be set to default value, got %d", minerWithDefaultThreads.config.NumThreads)
	}
}

// 测试SetChain方法
func TestSetChain(t *testing.T) {
	// 创建挖矿实例
	config := &Config{
		Enabled:          true,
		Coinbase:         common.Address{0x01},
		ExtraData:        []byte("NogoChain Miner"),
		MinGasPrice:      big.NewInt(1),
		MaxGasLimit:      10000000,
		RecommitInterval: 3 * time.Second,
		NumThreads:       2,
	}

	engine := nogopow.NewNogoPow()
	miner := NewMiner(config, engine)

	// 创建模拟区块链接口
	mockChain := struct{}{}

	// 设置区块链
	miner.SetChain(mockChain)
	if miner.chain != mockChain {
		t.Errorf("Miner chain mismatch")
	}
}

// 测试Start/Stop/IsRunning方法
func TestStartStop(t *testing.T) {
	// 创建挖矿实例
	config := &Config{
		Enabled:          true,
		Coinbase:         common.Address{0x01},
		ExtraData:        []byte("NogoChain Miner"),
		MinGasPrice:      big.NewInt(1),
		MaxGasLimit:      10000000,
		RecommitInterval: 3 * time.Second,
		NumThreads:       2,
	}

	engine := nogopow.NewNogoPow()
	miner := NewMiner(config, engine)

	// 测试初始状态
	if miner.IsRunning() {
		t.Errorf("Miner should not be running initially")
	}

	// 测试启动挖矿
	err := miner.Start()
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}

	if !miner.IsRunning() {
		t.Errorf("Miner should be running after Start")
	}

	// 测试重复启动（应该无错误）
	err = miner.Start()
	if err != nil {
		t.Errorf("Start should not return error when already running: %v", err)
	}

	if !miner.IsRunning() {
		t.Errorf("Miner should still be running after repeated Start")
	}

	// 测试停止挖矿
	err = miner.Stop()
	if err != nil {
		t.Errorf("Stop returned error: %v", err)
	}

	if miner.IsRunning() {
		t.Errorf("Miner should not be running after Stop")
	}

	// 测试重复停止（应该无错误）
	err = miner.Stop()
	if err != nil {
		t.Errorf("Stop should not return error when not running: %v", err)
	}

	if miner.IsRunning() {
		t.Errorf("Miner should still not be running after repeated Stop")
	}
}

// 测试SerializeHeader函数
func TestSerializeHeader(t *testing.T) {
	// 创建测试区块头
	header := &types.BlockHeader{
		ParentHash:  common.HexToHash("0x01"),
		Coinbase:    common.Address{0x01},
		Root:        common.HexToHash("0x02"),
		TxHash:      common.HexToHash("0x03"),
		ReceiptHash: common.HexToHash("0x04"),
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(1),
		GasLimit:    10000000,
		GasUsed:     0,
		Time:        1700000000,
		Extra:       []byte("test extra"),
		MixDigest:   common.HexToHash("0x05"),
		Nonce:       12345,
	}

	// 序列化区块头
	serialized := SerializeHeader(header)
	if len(serialized) == 0 {
		t.Errorf("SerializeHeader returned empty bytes")
	}

	// 验证序列化结果长度（至少应该大于0）
	if len(serialized) <= 0 {
		t.Errorf("Serialized header should not be empty")
	}
}

// 测试CalculateSeed函数
func TestCalculateSeed(t *testing.T) {
	// 测试不同区块号的种子计算
	testCases := []struct {
		number      *big.Int
		expectedLen int
	}{
		{big.NewInt(0), 32},
		{big.NewInt(1), 32},
		{big.NewInt(1000), 32},
		{big.NewInt(999999), 32},
	}

	for _, tc := range testCases {
		seed := CalculateSeed(tc.number)
		if len(seed) != tc.expectedLen {
			t.Errorf("CalculateSeed should return %d bytes, got %d", tc.expectedLen, len(seed))
		}
	}
}

// 测试CalculateTarget函数
func TestCalculateTarget(t *testing.T) {
	// 测试不同难度的目标值计算
	testCases := []struct {
		difficulty  *big.Int
		expectedLen int
	}{
		{big.NewInt(1), 32},
		{big.NewInt(1000), 32},
		{big.NewInt(1000000), 32},
	}

	for _, tc := range testCases {
		target := CalculateTarget(tc.difficulty)
		if len(target) != tc.expectedLen {
			t.Errorf("CalculateTarget should return %d bytes, got %d", tc.expectedLen, len(target))
		}
	}
}

// 测试EstimateHashRate函数
func TestEstimateHashRate(t *testing.T) {
	// 测试估算哈希率
	testCases := []struct {
		elapsed  time.Duration
		attempts uint64
		expected float64
	}{
		{0, 1000, 0},                         // 零时间应该返回0
		{1 * time.Second, 1000, 1000},        // 1秒1000次尝试
		{2 * time.Second, 1000, 500},         // 2秒1000次尝试
		{500 * time.Millisecond, 1000, 2000}, // 500毫秒1000次尝试
	}

	for _, tc := range testCases {
		hashRate := EstimateHashRate(tc.elapsed, tc.attempts)
		if hashRate != tc.expected {
			t.Errorf("EstimateHashRate should return %f, got %f", tc.expected, hashRate)
		}
	}
}

// 测试ValidateShare函数
func TestValidateShare(t *testing.T) {
	// 创建测试区块头
	header := &types.BlockHeader{
		ParentHash:  common.HexToHash("0x01"),
		Coinbase:    common.Address{0x01},
		Root:        common.HexToHash("0x02"),
		TxHash:      common.HexToHash("0x03"),
		ReceiptHash: common.HexToHash("0x04"),
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(1),
		GasLimit:    10000000,
		GasUsed:     0,
		Time:        1700000000,
		Extra:       []byte("test extra"),
		MixDigest:   common.HexToHash("0x05"),
		Nonce:       12345,
	}

	// 测试验证份额（当前实现返回true）
	result := ValidateShare(header, 12345, common.HexToHash("0x05"), big.NewInt(1000000))
	if !result {
		t.Errorf("ValidateShare should return true (implementation returns true)")
	}
}

// 测试SetExtra和SetCoinbase方法
func TestSetExtraAndCoinbase(t *testing.T) {
	// 创建挖矿实例
	config := &Config{
		Enabled:          true,
		Coinbase:         common.Address{0x01},
		ExtraData:        []byte("NogoChain Miner"),
		MinGasPrice:      big.NewInt(1),
		MaxGasLimit:      10000000,
		RecommitInterval: 3 * time.Second,
		NumThreads:       2,
	}

	engine := nogopow.NewNogoPow()
	miner := NewMiner(config, engine)

	// 测试设置额外数据
	newExtra := []byte("Updated Extra Data")
	miner.SetExtra(newExtra)
	if string(miner.config.ExtraData) != string(newExtra) {
		t.Errorf("ExtraData should be updated")
	}

	// 测试设置矿工地址
	newCoinbase := common.Address{0x02}
	miner.SetCoinbase(newCoinbase)
	if miner.config.Coinbase != newCoinbase {
		t.Errorf("Coinbase should be updated")
	}
}

// 测试GetHashRate方法
func TestGetHashRate(t *testing.T) {
	// 创建挖矿实例
	config := &Config{
		Enabled:          true,
		Coinbase:         common.Address{0x01},
		ExtraData:        []byte("NogoChain Miner"),
		MinGasPrice:      big.NewInt(1),
		MaxGasLimit:      10000000,
		RecommitInterval: 3 * time.Second,
		NumThreads:       2,
	}

	engine := nogopow.NewNogoPow()
	miner := NewMiner(config, engine)

	// 测试获取哈希率（当前实现返回0）
	hashRate := miner.GetHashRate()
	if hashRate != 0 {
		t.Errorf("GetHashRate should return 0 (implementation returns 0)")
	}
}

// 集成测试：测试完整的挖矿流程
func TestMiningIntegration(t *testing.T) {
	// 创建挖矿实例
	config := &Config{
		Enabled:          true,
		Coinbase:         common.Address{0x01},
		ExtraData:        []byte("NogoChain Miner"),
		MinGasPrice:      big.NewInt(1),
		MaxGasLimit:      10000000,
		RecommitInterval: 3 * time.Second,
		NumThreads:       2,
	}

	engine := nogopow.NewNogoPow()
	miner := NewMiner(config, engine)

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

	// 测试Seal方法
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := make(chan *types.Block, 1)
	stop := make(chan struct{})
	defer close(stop)

	err := miner.Seal(ctx, testBlock, results, stop)
	if err != nil && err != context.DeadlineExceeded {
		// 允许超时错误，因为挖矿可能需要时间
		t.Errorf("Seal returned unexpected error: %v", err)
	}

	// 测试启动和停止挖矿
	err = miner.Start()
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}

	// 等待一小段时间
	time.Sleep(500 * time.Millisecond)

	err = miner.Stop()
	if err != nil {
		t.Errorf("Stop returned error: %v", err)
	}
}

// 测试挖矿配置
func TestMiningConfig(t *testing.T) {
	// 测试不同的挖矿配置
	testConfigs := []struct {
		name               string
		config             *Config
		expectedNumThreads int
	}{
		{
			name: "positive threads",
			config: &Config{
				Enabled:          true,
				Coinbase:         common.Address{0x01},
				ExtraData:        []byte("NogoChain Miner"),
				MinGasPrice:      big.NewInt(1),
				MaxGasLimit:      10000000,
				RecommitInterval: 3 * time.Second,
				NumThreads:       4,
			},
			expectedNumThreads: 4,
		},
		{
			name: "zero threads",
			config: &Config{
				Enabled:          true,
				Coinbase:         common.Address{0x01},
				ExtraData:        []byte("NogoChain Miner"),
				MinGasPrice:      big.NewInt(1),
				MaxGasLimit:      10000000,
				RecommitInterval: 3 * time.Second,
				NumThreads:       0,
			},
			expectedNumThreads: 0, // 会被设置为默认值
		},
		{
			name: "negative threads",
			config: &Config{
				Enabled:          true,
				Coinbase:         common.Address{0x01},
				ExtraData:        []byte("NogoChain Miner"),
				MinGasPrice:      big.NewInt(1),
				MaxGasLimit:      10000000,
				RecommitInterval: 3 * time.Second,
				NumThreads:       -1,
			},
			expectedNumThreads: 0, // 会被设置为默认值
		},
	}

	for _, tc := range testConfigs {
		engine := nogopow.NewNogoPow()
		miner := NewMiner(tc.config, engine)

		if tc.config.NumThreads > 0 {
			if miner.config.NumThreads != tc.expectedNumThreads {
				t.Errorf("%s: NumThreads should be %d, got %d", tc.name, tc.expectedNumThreads, miner.config.NumThreads)
			}
		} else {
			if miner.config.NumThreads <= 0 {
				t.Errorf("%s: NumThreads should be set to default value, got %d", tc.name, miner.config.NumThreads)
			}
		}
	}
}
