package tools

import (
	"fmt"
	"time"

	"nogochain/core/blockchain"
	"nogochain/network"
	"nogochain/network/config"
)

// AutoTestRunner 自动测试运行器
type AutoTestRunner struct {
	testInterval time.Duration
	Network      *network.Network
	Blockchain   *blockchain.Blockchain
	testResults  []TestResult
}

// TestResult 测试结果
type TestResult struct {
	TestType      string
	Timestamp     time.Time
	Duration      time.Duration
	Success       bool
	Details       string
}

// NewAutoTestRunner 创建新的自动测试运行器
func NewAutoTestRunner(testInterval time.Duration) *AutoTestRunner {
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
		fmt.Printf("启动网络失败: %v\n", err)
	}

	return &AutoTestRunner{
		testInterval: testInterval,
		Network:      n,
		Blockchain:   bc,
		testResults:  make([]TestResult, 0),
	}
}

// Start 开始自动测试
func (r *AutoTestRunner) Start() {
	go func() {
		ticker := time.NewTicker(r.testInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			r.runTests()
		}
	}()
}

// runTests 运行所有测试
func (r *AutoTestRunner) runTests() {
	// 运行交易生成测试
	r.runTransactionTest()

	// 运行区块同步测试
	r.runBlockSyncTest()

	// 运行网络测试
	r.runNetworkTest()
}

// runTransactionTest 运行交易生成测试
func (r *AutoTestRunner) runTransactionTest() {
	start := time.Now()

	// 创建交易生成器
	generator := NewTxGenerator(r.Network, 100, 30*time.Second)

	// 开始生成交易
	generator.Start()

	// 等待测试完成
	time.Sleep(31 * time.Second)

	// 获取统计信息
	stats := generator.GetStats()

	// 记录测试结果
	result := TestResult{
		TestType:  "交易生成测试",
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Success:   stats.SuccessfulTransactions > 0, // 至少有一个成功的交易
		Details:   fmt.Sprintf("生成交易数: %d, 成功: %d, TPS: %.2f", 
			stats.TotalTransactions, stats.SuccessfulTransactions, 
			float64(stats.SuccessfulTransactions)/stats.EndTime.Sub(stats.StartTime).Seconds()),
	}

	r.testResults = append(r.testResults, result)

	// 打印测试结果
	fmt.Printf("%s 测试结果: %v\n", result.TestType, result.Details)
}

// runBlockSyncTest 运行区块同步测试
func (r *AutoTestRunner) runBlockSyncTest() {
	start := time.Now()

	// 创建区块同步测试器
	tester := NewBlockSyncTester(r.Blockchain, r.Network)

	// 准备源区块链
	err := tester.PrepareSourceBlockchain(100)
	if err != nil {
		result := TestResult{
			TestType:  "区块同步测试",
			Timestamp: time.Now(),
			Duration:  time.Since(start),
			Success:   false,
			Details:   fmt.Sprintf("准备区块链失败: %v", err),
		}
		r.testResults = append(r.testResults, result)
		fmt.Printf("%s 测试结果: %v\n", result.TestType, result.Details)
		return
	}

	// 开始同步测试
	tester.StartSync()

	// 获取统计信息
	stats := tester.GetStats()

	// 记录测试结果
	result := TestResult{
		TestType:  "区块同步测试",
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Success:   stats.BlocksSynced == 100,
		Details:   fmt.Sprintf("同步区块数: %d, 同步时间: %v, 带宽: %d bytes", 
			stats.BlocksSynced, stats.SyncTime, stats.NetworkBandwidth),
	}

	r.testResults = append(r.testResults, result)

	// 打印测试结果
	fmt.Printf("%s 测试结果: %v\n", result.TestType, result.Details)
}

// runNetworkTest 运行网络测试
func (r *AutoTestRunner) runNetworkTest() {
	start := time.Now()

	// 创建网络测试器
	testNodes := []string{"127.0.0.1:8545"}
	tester := NewNetworkTester(r.Network, testNodes)

	// 测试延迟
	tester.TestLatency()

	// 测试带宽
	tester.TestBandwidth()

	// 获取统计信息
	stats := tester.GetStats()

	// 记录测试结果
	result := TestResult{
		TestType:  "网络测试",
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Success:   len(stats.LatencyTests) > 0,
		Details:   fmt.Sprintf("测试节点数: %d, 延迟测试: %d, 带宽测试: %d", 
			len(testNodes), len(stats.LatencyTests), len(stats.BandwidthTests)),
	}

	r.testResults = append(r.testResults, result)

	// 打印测试结果
	fmt.Printf("%s 测试结果: %v\n", result.TestType, result.Details)
}

// GetResults 获取所有测试结果
func (r *AutoTestRunner) GetResults() []TestResult {
	return r.testResults
}

// PrintResults 打印所有测试结果
func (r *AutoTestRunner) PrintResults() {
	fmt.Printf("自动测试结果汇总:\n")
	for _, result := range r.testResults {
		fmt.Printf("%s - %v: %v - %s\n", result.TestType, result.Timestamp, result.Success, result.Details)
	}
}
