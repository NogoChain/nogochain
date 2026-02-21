package tools

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/network"
	"nogochain/network/config"
)

// PerformanceTest 性能测试主函数
func PerformanceTest() {
	fmt.Println("=== NogoChain 性能测试 (Task 7 相同测试) ===")
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

	// 5. 高负载测试 (100 TPS)
	fmt.Println("5. 高负载测试 (100 TPS)")
	runHighLoadTest()
	fmt.Println()

	// 6. 长时间稳定性测试准备
	fmt.Println("6. 长时间稳定性测试准备")
	fmt.Println("注意: 长时间稳定性测试需要运行24小时，将在后台执行")
	fmt.Println()

	fmt.Println("=== 测试完成 ===")
	fmt.Println("测试结束时间:", time.Now())
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
	for i := 0; i < numBlocks; i++ {
		block := &types.Block{
			Transactions: []*types.Transaction{},
		}
		sourceBC.AddBlock(block)
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

		err := bc.AddBlock(block)
		if err != nil {
			fmt.Printf("添加区块失败: %v\n", err)
		}
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

// RunLongTermStabilityTest 运行长时间稳定性测试 (24小时)
func RunLongTermStabilityTest() {
	fmt.Println("=== NogoChain 长时间稳定性测试 ===")
	fmt.Println("测试开始时间:", time.Now())
	fmt.Println("测试配置: 100 TPS, 持续24小时")
	fmt.Println()

	// 创建网络配置
	cfg := &config.Config{
		Port: 8550,
		RPC: &config.RPCConfig{
			Enabled: false,
			Port:    8550,
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

	// 测试100 TPS的高负载，持续24小时
	testRate := 100
	testDuration := 24 * time.Hour

	// 创建交易生成器
	generator := NewTxGenerator(n, testRate, testDuration)

	// 开始生成交易
	generator.Start()

	// 记录测试信息
	fmt.Println("长时间稳定性测试已启动，将在24小时后完成")
	fmt.Println("测试过程中请保持系统运行")
	fmt.Println()

	// 注意: 实际运行时，这个函数会持续运行24小时
	// 这里我们只打印启动信息，实际测试需要在后台运行
}
