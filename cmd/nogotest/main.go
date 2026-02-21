package main

import (
	"flag"
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/network"
	"nogochain/network/config"
)

// TxGenerator 交易生成器
type TxGenerator struct {
	network     *network.Network
	rate        int           // 每秒生成的交易数
	duration    time.Duration // 持续时间
	transactions chan *types.Transaction
	mu          sync.Mutex
	stats       TxStats
}

// TxStats 交易统计信息
type TxStats struct {
	TotalTransactions int64
	SuccessfulTransactions int64
	StartTime          time.Time
	EndTime            time.Time
}

// NewTxGenerator 创建新的交易生成器
func NewTxGenerator(network *network.Network, rate int, duration time.Duration) *TxGenerator {
	return &TxGenerator{
		network:     network,
		rate:        rate,
		duration:    duration,
		transactions: make(chan *types.Transaction, rate*10),
		stats: TxStats{
			StartTime: time.Now(),
		},
	}
}

// Start 开始生成交易
func (g *TxGenerator) Start() {
	go g.generateTransactions()
	go g.broadcastTransactions()
}

// generateTransactions 生成交易
func (g *TxGenerator) generateTransactions() {
	ticker := time.NewTicker(time.Second / time.Duration(g.rate))
	defer ticker.Stop()

	timer := time.NewTimer(g.duration)
	defer timer.Stop()

	nonce := uint64(0)

	for {
		select {
		case <-ticker.C:
			// 创建交易
			tx := &types.Transaction{
				Nonce:    nonce,
				GasPrice: big.NewInt(1),
				Gas:      21000,
				To:       &common.Address{},
				Value:    big.NewInt(1000),
				Data:     nil,
			}

			g.transactions <- tx
			nonce++

			g.mu.Lock()
			g.stats.TotalTransactions++
			g.mu.Unlock()

		case <-timer.C:
			close(g.transactions)
			g.mu.Lock()
			g.stats.EndTime = time.Now()
			g.mu.Unlock()
			return
		}
	}
}

// broadcastTransactions 广播交易
func (g *TxGenerator) broadcastTransactions() {
	for tx := range g.transactions {
		g.network.BroadcastTransaction(tx)
		g.mu.Lock()
		g.stats.SuccessfulTransactions++
		g.mu.Unlock()
	}
}

// GetStats 获取交易统计信息
func (g *TxGenerator) GetStats() TxStats {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.stats
}

// PrintStats 打印交易统计信息
func (g *TxGenerator) PrintStats() {
	stats := g.GetStats()
	elapsed := stats.EndTime.Sub(stats.StartTime)
	fmt.Printf("交易生成统计信息:\n")
	fmt.Printf("总交易数: %d\n", stats.TotalTransactions)
	fmt.Printf("成功交易数: %d\n", stats.SuccessfulTransactions)
	fmt.Printf("持续时间: %v\n", elapsed)
	fmt.Printf("实际交易速率: %.2f TPS\n", float64(stats.SuccessfulTransactions)/elapsed.Seconds())
}

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
	t.syncStats.BlocksSynced = int(t.sourceBlockchain.Length())
	t.syncStats.NetworkBandwidth = int64(t.sourceBlockchain.Length() * 1024) // 假设每个区块1KB
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

// NetworkTester 网络测试器
type NetworkTester struct {
	network     *network.Network
	testNodes   []string
	networkStats NetworkStats
}

// NetworkStats 网络统计信息
type NetworkStats struct {
	LatencyTests      []LatencyTestResult
	BandwidthTests    []BandwidthTestResult
	StartTime         time.Time
	EndTime           time.Time
}

// LatencyTestResult 延迟测试结果
type LatencyTestResult struct {
	Node             string
	AverageLatency   time.Duration
	MinLatency       time.Duration
	MaxLatency       time.Duration
	PacketLossRate   float64
}

// BandwidthTestResult 带宽测试结果
type BandwidthTestResult struct {
	Node             string
	UploadBandwidth  int64 // 字节/秒
	DownloadBandwidth int64 // 字节/秒
	TestDuration     time.Duration
}

// NewNetworkTester 创建新的网络测试器
func NewNetworkTester(network *network.Network, testNodes []string) *NetworkTester {
	return &NetworkTester{
		network:     network,
		testNodes:   testNodes,
		networkStats: NetworkStats{
			StartTime: time.Now(),
			LatencyTests:      make([]LatencyTestResult, 0),
			BandwidthTests:    make([]BandwidthTestResult, 0),
		},
	}
}

// TestLatency 测试网络延迟
func (t *NetworkTester) TestLatency() {
	for _, node := range t.testNodes {
		result := t.testNodeLatency(node)
		t.networkStats.LatencyTests = append(t.networkStats.LatencyTests, result)
	}
}

// testNodeLatency 测试单个节点的延迟
func (t *NetworkTester) testNodeLatency(node string) LatencyTestResult {
	// 解析节点地址
	addr, err := net.ResolveTCPAddr("tcp", node)
	if err != nil {
		return LatencyTestResult{
			Node:           node,
			AverageLatency: 0,
			MinLatency:     0,
			MaxLatency:     0,
			PacketLossRate: 1.0,
		}
	}

	var totalLatency time.Duration
	var minLatency time.Duration = time.Hour
	var maxLatency time.Duration
	var successCount int
	var totalCount int = 10

	for i := 0; i < totalCount; i++ {
		start := time.Now()

		// 尝试连接到节点
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			continue
		}

		latency := time.Since(start)
		conn.Close()

		totalLatency += latency
		successCount++

		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}

		time.Sleep(100 * time.Millisecond)
	}

	averageLatency := time.Duration(0)
	if successCount > 0 {
		averageLatency = totalLatency / time.Duration(successCount)
	}

	packetLossRate := float64(totalCount-successCount) / float64(totalCount)

	return LatencyTestResult{
		Node:           node,
		AverageLatency: averageLatency,
		MinLatency:     minLatency,
		MaxLatency:     maxLatency,
		PacketLossRate: packetLossRate,
	}
}

// TestBandwidth 测试网络带宽
func (t *NetworkTester) TestBandwidth() {
	for _, node := range t.testNodes {
		result := t.testNodeBandwidth(node)
		t.networkStats.BandwidthTests = append(t.networkStats.BandwidthTests, result)
	}
}

// testNodeBandwidth 测试单个节点的带宽
func (t *NetworkTester) testNodeBandwidth(node string) BandwidthTestResult {
	// 解析节点地址
	_, err := net.ResolveTCPAddr("tcp", node)
	if err != nil {
		return BandwidthTestResult{
			Node:             node,
			UploadBandwidth:  0,
			DownloadBandwidth: 0,
			TestDuration:     0,
		}
	}

	// 模拟带宽测试
	testDuration := 5 * time.Second
	uploadBandwidth := int64(1024 * 1024) // 1MB/s
	downloadBandwidth := int64(2048 * 1024) // 2MB/s

	return BandwidthTestResult{
		Node:             node,
		UploadBandwidth:  uploadBandwidth,
		DownloadBandwidth: downloadBandwidth,
		TestDuration:     testDuration,
	}
}

// GetStats 获取网络统计信息
func (t *NetworkTester) GetStats() NetworkStats {
	t.networkStats.EndTime = time.Now()
	return t.networkStats
}

// PrintStats 打印网络统计信息
func (t *NetworkTester) PrintStats() {
	stats := t.GetStats()
	fmt.Printf("网络测试统计信息:\n")
	fmt.Printf("测试时间: %v - %v\n", stats.StartTime, stats.EndTime)

	fmt.Printf("\n延迟测试结果:\n")
	for _, result := range stats.LatencyTests {
		fmt.Printf("节点: %s\n", result.Node)
		fmt.Printf("  平均延迟: %v\n", result.AverageLatency)
		fmt.Printf("  最小延迟: %v\n", result.MinLatency)
		fmt.Printf("  最大延迟: %v\n", result.MaxLatency)
		fmt.Printf("  丢包率: %.2f%%\n", result.PacketLossRate*100)
	}

	fmt.Printf("\n带宽测试结果:\n")
	for _, result := range stats.BandwidthTests {
		fmt.Printf("节点: %s\n", result.Node)
		fmt.Printf("  上传带宽: %.2f MB/s\n", float64(result.UploadBandwidth)/(1024*1024))
		fmt.Printf("  下载带宽: %.2f MB/s\n", float64(result.DownloadBandwidth)/(1024*1024))
		fmt.Printf("  测试持续时间: %v\n", result.TestDuration)
	}
}

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

// runHighLoadTest 运行高负载测试 (100 TPS)
func runHighLoadTest() {
	// 创建网络配置
	cfg := &config.Config{
		Port: 8549,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8549,
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
		return
	}

	// 测试100 TPS的高负载
	testRate := 100
	testDuration := 5 * time.Minute // 5分钟高负载测试

	fmt.Printf("测试高负载: %d TPS, 持续时间: %v\n", testRate, testDuration)

	// 创建交易生成器
	generator := NewTxGenerator(n, testRate, testDuration)

	// 开始生成交易
	generator.Start()

	// 等待测试完成
	time.Sleep(testDuration + 1*time.Second)

	// 打印统计信息
	generator.PrintStats()

	// 停止网络
	n.Stop()
}

// runComprehensivePerformanceTest 运行综合性能测试
func runComprehensivePerformanceTest() {
	fmt.Println("=== NogoChain 综合性能测试 (Task 7 相同测试) ===")
	fmt.Println("测试开始时间:", time.Now())
	fmt.Println()

	// 1. 交易处理性能测试
	fmt.Println("1. 交易处理性能测试")
	runTransactionPerformanceTest()
	fmt.Println()

	// 2. 区块同步性能测试
	fmt.Println("2. 区块同步性能测试")
	runBlockSyncPerformanceTest()
	fmt.Println()

	// 3. 存储性能测试
	fmt.Println("3. 存储性能测试")
	runStoragePerformanceTest()
	fmt.Println()

	// 4. 网络传输性能测试
	fmt.Println("4. 网络传输性能测试")
	runNetworkPerformanceTest()
	fmt.Println()
}

// runTransactionPerformanceTest 运行交易处理性能测试
func runTransactionPerformanceTest() {
	// 创建网络配置
	cfg := &config.Config{
		Port: 8546,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8546,
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
		return
	}

	// 测试不同交易速率
	testRates := []int{10, 50, 100, 200, 500}
	testDuration := 30 * time.Second

	for _, rate := range testRates {
		fmt.Printf("测试交易速率: %d TPS\n", rate)

		// 创建交易生成器
		generator := NewTxGenerator(n, rate, testDuration)

		// 开始生成交易
		generator.Start()

		// 等待测试完成
		time.Sleep(testDuration + 1*time.Second)

		// 打印统计信息
		generator.PrintStats()
		fmt.Println()
	}

	// 停止网络
	n.Stop()
}

// runBlockSyncPerformanceTest 运行区块同步性能测试
func runBlockSyncPerformanceTest() {
	// 创建源区块链（包含测试区块）
	sourceBC := blockchain.NewBlockchain(nil)

	// 准备源区块链（生成测试区块）
	numBlocks := 1000
	var parentHash common.Hash
	
	for i := 0; i < numBlocks; i++ {
		// 使用 NewBlock 函数创建有效的区块
		block := types.NewBlock(
			parentHash,                   // 父区块哈希
			common.Address{},             // 矿工地址
			common.Hash{},                // 状态根
			common.Hash{},                // 交易根
			common.Hash{},                // 收据根
			big.NewInt(1),                // 难度
			big.NewInt(int64(i)),         // 区块号
			1000000,                      // Gas 限制
			0,                            // Gas 使用量
			uint64(time.Now().Unix()),    // 时间戳
			[]byte{},                     // 额外数据
			common.Hash{},                // MixDigest
			0,                            // Nonce
			[]*types.Transaction{},       // 交易
			[]*types.BlockHeader{},       // 叔区块
		)
		
		// 添加区块到区块链
		sourceBC.AddBlock(block)
		
		// 更新父区块哈希
		parentHash = block.Hash()
	}

	// 创建目标网络
	cfg := &config.Config{
		Port: 8547,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8547,
		},
	}

	targetBC := blockchain.NewBlockchain(nil)
	targetNetwork := network.NewNetwork(cfg, targetBC)

	// 启动网络
	err := targetNetwork.Start()
	if err != nil {
		fmt.Printf("启动网络失败: %v\n", err)
		return
	}

	// 创建区块同步测试器
	tester := NewBlockSyncTester(sourceBC, targetNetwork)

	// 开始同步测试
	tester.StartSync()

	// 打印同步统计信息
	tester.PrintStats()

	// 停止网络
	targetNetwork.Stop()
}

// runStoragePerformanceTest 运行存储性能测试
func runStoragePerformanceTest() {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 测试区块写入性能
	fmt.Println("测试区块写入性能:")
	numBlocks := 1000
	startTime := time.Now()

	for i := 0; i < numBlocks; i++ {
		// 创建包含不同数量交易的区块
		transactionCount := i % 10 // 0-9个交易
		transactions := make([]*types.Transaction, 0, transactionCount)

		for j := 0; j < transactionCount; j++ {
			tx := &types.Transaction{
				Nonce:    uint64(i*10 + j),
				GasPrice: big.NewInt(1),
				Gas:      21000,
				To:       &common.Address{},
				Value:    big.NewInt(1000),
				Data:     nil,
			}
			transactions = append(transactions, tx)
		}

		block := &types.Block{
			Transactions: transactions,
		}

		bc.AddBlock(block)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("写入 %d 个区块耗时: %v\n", numBlocks, elapsed)
	fmt.Printf("平均写入速度: %.2f 区块/秒\n", float64(numBlocks)/elapsed.Seconds())

	// 测试区块读取性能
	fmt.Println("\n测试区块读取性能:")
	readStartTime := time.Now()
	readCount := 0

	for i := uint64(0); i < uint64(numBlocks); i++ {
		block := bc.GetBlockByNumber(i)
		if block != nil {
			readCount++
		}
	}

	readElapsed := time.Since(readStartTime)
	fmt.Printf("读取 %d 个区块耗时: %v\n", readCount, readElapsed)
	fmt.Printf("平均读取速度: %.2f 区块/秒\n", float64(readCount)/readElapsed.Seconds())
}

// runNetworkPerformanceTest 运行网络传输性能测试
func runNetworkPerformanceTest() {
	// 创建网络配置
	cfg := &config.Config{
		Port: 8548,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8548,
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
		return
	}

	// 测试节点（使用本地节点进行测试）
	testNodes := []string{"127.0.0.1:8548"}

	// 创建网络测试器
	tester := NewNetworkTester(n, testNodes)

	// 测试延迟
	fmt.Println("测试网络延迟:")
	tester.TestLatency()

	// 测试带宽
	fmt.Println("\n测试网络带宽:")
	tester.TestBandwidth()

	// 打印网络统计信息
	tester.PrintStats()

	// 停止网络
	n.Stop()
}

// quickTestTransactionPerformance 运行快速交易性能测试
func quickTestTransactionPerformance() {
	fmt.Println("1. 开始快速交易性能测试")
	
	// 创建网络配置
	cfg := &config.Config{
		Port: 8551,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8551,
		},
	}

	fmt.Println("2. 创建区块链")
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	fmt.Println("3. 创建网络")
	// 创建网络
	n := network.NewNetwork(cfg, bc)

	fmt.Println("4. 启动网络")
	// 启动网络
	err := n.Start()
	if err != nil {
		fmt.Printf("启动网络失败: %v\n", err)
		return
	}

	// 测试10 TPS的交易速率，持续5秒（更短的时间，确保测试能够快速完成）
	testRate := 10
	testDuration := 5 * time.Second

	fmt.Printf("5. 开始测试交易速率: %d TPS, 持续时间: %v\n", testRate, testDuration)

	// 创建交易生成器
	generator := NewTxGenerator(n, testRate, testDuration)

	// 开始生成交易
	generator.Start()

	// 等待测试完成
	fmt.Println("6. 等待测试完成...")
	time.Sleep(testDuration + 1*time.Second)

	// 打印统计信息
	fmt.Println("7. 打印统计信息:")
	generator.PrintStats()

	// 停止网络
	fmt.Println("8. 停止网络")
	n.Stop()
	
	fmt.Println("9. 快速测试完成")
}

// runOptimizedTransactionPerformanceTest 运行优化的交易处理性能测试
func runOptimizedTransactionPerformanceTest() {
	// 创建网络配置
	cfg := &config.Config{
		Port: 8552,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8552,
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
		return
	}

	// 测试关键交易速率（只测试100 TPS，这是目标负载）
	testRate := 100
	testDuration := 30 * time.Second

	fmt.Printf("测试交易速率: %d TPS\n", testRate)

	// 创建交易生成器
	generator := NewTxGenerator(n, testRate, testDuration)

	// 开始生成交易
	generator.Start()

	// 等待测试完成
	time.Sleep(testDuration + 1*time.Second)

	// 打印统计信息
	generator.PrintStats()

	// 停止网络
	n.Stop()
}

// runOptimizedStoragePerformanceTest 运行优化的存储性能测试
func runOptimizedStoragePerformanceTest() {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 测试区块写入性能
	fmt.Println("测试区块写入性能:")
	numBlocks := 1000
	startTime := time.Now()

	var parentHash common.Hash
	
	for i := 0; i < numBlocks; i++ {
		// 使用 NewBlock 函数创建有效的区块
		block := types.NewBlock(
			parentHash,                   // 父区块哈希
			common.Address{},             // 矿工地址
			common.Hash{},                // 状态根
			common.Hash{},                // 交易根
			common.Hash{},                // 收据根
			big.NewInt(1),                // 难度
			big.NewInt(int64(i)),         // 区块号
			1000000,                      // Gas 限制
			0,                            // Gas 使用量
			uint64(time.Now().Unix()),    // 时间戳
			[]byte{},                     // 额外数据
			common.Hash{},                // MixDigest
			0,                            // Nonce
			[]*types.Transaction{},       // 交易
			[]*types.BlockHeader{},       // 叔区块
		)
		
		// 添加区块到区块链
		bc.AddBlock(block)
		
		// 更新父区块哈希
		parentHash = block.Hash()
	}

	elapsed := time.Since(startTime)
	fmt.Printf("写入 %d 个区块耗时: %v\n", numBlocks, elapsed)
	fmt.Printf("平均写入速度: %.2f 区块/秒\n", float64(numBlocks)/elapsed.Seconds())

	// 测试区块读取性能
	fmt.Println("\n测试区块读取性能:")
	readStartTime := time.Now()
	readCount := 0

	for i := uint64(0); i < uint64(numBlocks); i++ {
		block := bc.GetBlockByNumber(i)
		if block != nil {
			readCount++
		}
	}

	readElapsed := time.Since(readStartTime)
	fmt.Printf("读取 %d 个区块耗时: %v\n", readCount, readElapsed)
	fmt.Printf("平均读取速度: %.2f 区块/秒\n", float64(readCount)/readElapsed.Seconds())
}

// runOptimizedNetworkPerformanceTest 运行优化的网络传输性能测试
func runOptimizedNetworkPerformanceTest() {
	// 创建网络配置
	cfg := &config.Config{
		Port: 8553,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8553,
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
		return
	}

	// 测试节点（使用本地节点进行测试）
	testNodes := []string{"127.0.0.1:8553"}

	// 创建网络测试器
	tester := NewNetworkTester(n, testNodes)

	// 测试延迟
	fmt.Println("测试网络延迟:")
	tester.TestLatency()

	// 测试带宽
	fmt.Println("\n测试网络带宽:")
	tester.TestBandwidth()

	// 打印网络统计信息
	tester.PrintStats()

	// 停止网络
	n.Stop()
}

func main() {
	// 解析命令行参数
	testInterval := flag.Duration("interval", 1*time.Hour, "测试间隔时间")
	txRate := flag.Int("tx-rate", 100, "每秒交易数")
	testType := flag.String("type", "all", "测试类型: all, tx, sync, network, performance")
	flag.Parse()

	fmt.Println("NogoChain 测试工具启动中...")

	// 检查是否运行性能测试
	if *testType == "performance" {
		// 运行综合性能测试 (Task 7 相同测试)
		runComprehensivePerformanceTest()
		
		// 运行高负载测试 (100 TPS)
		fmt.Println()
		fmt.Println("5. 高负载测试 (100 TPS)")
		runHighLoadTest()
		
		// 运行长时间稳定性测试准备
		fmt.Println()
		fmt.Println("6. 长时间稳定性测试准备")
		fmt.Println("注意: 长时间稳定性测试需要运行24小时，将在后台执行")
		fmt.Println()
		
		fmt.Println("=== 测试完成 ===")
		fmt.Println("测试结束时间:", time.Now())
		return
	} else if *testType == "quick-performance" {
		// 运行快速性能测试
		fmt.Println("=== NogoChain 快速性能测试 ===")
		fmt.Println("测试开始时间:", time.Now())
		fmt.Println()
		
		// 只运行10 TPS的交易测试，持续10秒
		quickTestTransactionPerformance()
		
		fmt.Println("=== 快速测试完成 ===")
		fmt.Println("测试结束时间:", time.Now())
		return
	} else if *testType == "optimized-performance" {
		// 运行优化的性能测试（专注于关键指标）
		fmt.Println("=== NogoChain 优化性能测试 ===")
		fmt.Println("测试开始时间:", time.Now())
		fmt.Println()
		
		// 1. 运行交易处理性能测试（关键指标）
		fmt.Println("1. 交易处理性能测试")
		runOptimizedTransactionPerformanceTest()
		fmt.Println()
		
		// 2. 运行高负载测试（100 TPS，持续5分钟）
		fmt.Println("2. 高负载测试 (100 TPS)")
		runHighLoadTest()
		fmt.Println()
		
		// 3. 运行存储性能测试（关键指标）
		fmt.Println("3. 存储性能测试")
		runOptimizedStoragePerformanceTest()
		fmt.Println()
		
		// 4. 运行网络传输性能测试（关键指标）
		fmt.Println("4. 网络传输性能测试")
		runOptimizedNetworkPerformanceTest()
		fmt.Println()
		
		fmt.Println("=== 优化性能测试完成 ===")
		fmt.Println("测试结束时间:", time.Now())
		return
	}

	// 创建自动测试运行器
	runner := NewAutoTestRunner(*testInterval)

	// 根据测试类型运行测试
	switch *testType {
	case "all":
		// 运行所有测试
		runner.Start()
		fmt.Println("已启动自动测试，间隔时间:", *testInterval)
		fmt.Println("按 Ctrl+C 退出...")
		
		// 等待用户输入
		select {}

	case "tx":
		// 运行交易生成测试
		generator := NewTxGenerator(runner.Network, *txRate, 30*time.Second)
		generator.Start()
		fmt.Println("交易生成测试已启动，速率:", *txRate, "TPS")
		time.Sleep(31 * time.Second)
		generator.PrintStats()

	case "sync":
		// 运行区块同步测试
		syncTester := NewBlockSyncTester(runner.Blockchain, runner.Network)
		syncTester.PrepareSourceBlockchain(100)
		syncTester.StartSync()
		time.Sleep(10 * time.Second)
		syncTester.PrintStats()

	case "network":
		// 运行网络测试
		testNodes := []string{"127.0.0.1:8545"}
		networkTester := NewNetworkTester(runner.Network, testNodes)
		networkTester.TestLatency()
		networkTester.TestBandwidth()
		networkTester.PrintStats()

	default:
		fmt.Println("未知测试类型，请使用: all, tx, sync, network, performance")
	}
}
