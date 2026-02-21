package performance

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/blockchain"
	"nogochain/core/types"
)

// TestStoragePerformance 存储性能测试
func TestStoragePerformance(t *testing.T) {
	fmt.Println("=== NogoChain 存储性能测试 ===")
	fmt.Println("测试开始时间:", time.Now())
	fmt.Println()

	// 测试区块写入性能
	fmt.Println("1. 区块写入性能测试")
	testBlockWritePerformance()
	fmt.Println()

	// 测试区块读取性能
	fmt.Println("2. 区块读取性能测试")
	testBlockReadPerformance()
	fmt.Println()

	// 测试状态存储性能
	fmt.Println("3. 状态存储性能测试")
	testStateStoragePerformance()
	fmt.Println()

	fmt.Println("=== 存储性能测试完成 ===")
	fmt.Println("测试结束时间:", time.Now())
}

// testBlockWritePerformance 测试区块写入性能
func testBlockWritePerformance() {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 测试区块写入性能
	numBlocks := 1000
	startTime := time.Now()

	// 获取创世区块
	genesis := bc.Genesis()
	parentHash := genesis.Hash()

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

		// 创建区块并设置正确的父区块哈希
		block := types.NewBlock(
			parentHash,
			common.Address{},
			common.Hash{},
			common.Hash{},
			common.Hash{},
			big.NewInt(1000000),
			big.NewInt(int64(i+1)),
			10000000,
			0,
			uint64(time.Now().Unix()),
			[]byte(fmt.Sprintf("Block %d", i+1)),
			common.Hash{},
			0,
			transactions,
			[]*types.BlockHeader{},
		)

		err := bc.AddBlock(block)
		if err != nil {
			fmt.Printf("添加区块失败: %v\n", err)
		}

		// 更新父区块哈希为当前区块的哈希
		parentHash = block.Hash()
	}

	elapsed := time.Since(startTime)
	fmt.Printf("写入 %d 个区块耗时: %v\n", numBlocks, elapsed)
	fmt.Printf("平均写入速度: %.2f 区块/秒\n", float64(numBlocks)/elapsed.Seconds())
}

// testBlockReadPerformance 测试区块读取性能
func testBlockReadPerformance() {
	// 创建区块链并预填充区块
	bc := blockchain.NewBlockchain(nil)

	// 预填充区块
	numBlocks := 1000
	// 获取创世区块
	genesis := bc.Genesis()
	parentHash := genesis.Hash()

	for i := 0; i < numBlocks; i++ {
		transactionCount := i % 10
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

		// 创建区块并设置正确的父区块哈希
		block := types.NewBlock(
			parentHash,
			common.Address{},
			common.Hash{},
			common.Hash{},
			common.Hash{},
			big.NewInt(1000000),
			big.NewInt(int64(i+1)),
			10000000,
			0,
			uint64(time.Now().Unix()),
			[]byte(fmt.Sprintf("Block %d", i+1)),
			common.Hash{},
			0,
			transactions,
			[]*types.BlockHeader{},
		)

		bc.AddBlock(block)

		// 更新父区块哈希为当前区块的哈希
		parentHash = block.Hash()
	}

	// 测试区块读取性能
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

// testStateStoragePerformance 测试状态存储性能
func testStateStoragePerformance() {
	// 创建区块链
	bc := blockchain.NewBlockchain(nil)

	// 获取状态数据库
	stateDB := bc.StateDB()

	// 测试状态写入性能
	fmt.Println("测试状态写入性能:")
	numAccounts := 1000
	startTime := time.Now()

	for i := 0; i < numAccounts; i++ {
		// 创建随机地址
		addr := common.Address{byte(i), byte(i+1), byte(i+2), byte(i+3), byte(i+4), byte(i+5), byte(i+6), byte(i+7), byte(i+8), byte(i+9)}

		// 创建账户
		stateDB.CreateAccount(addr)

		// 设置余额
		stateDB.AddBalance(addr, big.NewInt(int64(i*1000)))

		// 设置Nonce
		stateDB.SetNonce(addr, uint64(i))

		// 设置存储状态
		for j := 0; j < 10; j++ {
			key := common.Hash{byte(j), byte(j+1), byte(j+2), byte(j+3), byte(j+4), byte(j+5), byte(j+6), byte(j+7), byte(j+8), byte(j+9)}
			value := common.Hash{byte(i + j), byte(i + j + 1), byte(i + j + 2), byte(i + j + 3), byte(i + j + 4), byte(i + j + 5), byte(i + j + 6), byte(i + j + 7), byte(i + j + 8), byte(i + j + 9)}
			stateDB.SetState(addr, key, value)
		}
	}

	writeElapsed := time.Since(startTime)
	fmt.Printf("写入 %d 个账户状态耗时: %v\n", numAccounts, writeElapsed)
	fmt.Printf("平均写入速度: %.2f 账户/秒\n", float64(numAccounts)/writeElapsed.Seconds())

	// 测试状态读取性能
	fmt.Println("测试状态读取性能:")
	readStartTime := time.Now()
	readCount := 0

	for i := 0; i < numAccounts; i++ {
		// 创建随机地址
		addr := common.Address{byte(i), byte(i+1), byte(i+2), byte(i+3), byte(i+4), byte(i+5), byte(i+6), byte(i+7), byte(i+8), byte(i+9)}

		// 读取余额
		stateDB.GetBalance(addr)

		// 读取Nonce
		stateDB.GetNonce(addr)

		// 读取存储状态
		for j := 0; j < 10; j++ {
			key := common.Hash{byte(j), byte(j+1), byte(j+2), byte(j+3), byte(j+4), byte(j+5), byte(j+6), byte(j+7), byte(j+8), byte(j+9)}
			stateDB.GetState(addr, key)
		}

		readCount++
	}

	readElapsed := time.Since(readStartTime)
	fmt.Printf("读取 %d 个账户状态耗时: %v\n", readCount, readElapsed)
	fmt.Printf("平均读取速度: %.2f 账户/秒\n", float64(readCount)/readElapsed.Seconds())
}
