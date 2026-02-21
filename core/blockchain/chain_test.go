package blockchain

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/types"
)

// 测试NewBlockchain函数
func TestNewBlockchain(t *testing.T) {
	// 测试创建默认区块链（带创世区块）
	bc := NewBlockchain(nil)
	if bc == nil {
		t.Errorf("NewBlockchain returned nil")
	}

	// 检查创世区块
	genesis := bc.Genesis()
	if genesis == nil {
		t.Errorf("Genesis block is nil")
	}

	if genesis.NumberU64() != 0 {
		t.Errorf("Genesis block number should be 0, got %d", genesis.NumberU64())
	}

	// 检查当前头部
	currentHead := bc.CurrentHead()
	if currentHead == nil {
		t.Errorf("Current head is nil")
	}

	if currentHead.Hash() != genesis.Hash() {
		t.Errorf("Current head should be genesis block")
	}

	// 检查链长度
	length := bc.Length()
	if length != 1 {
		t.Errorf("Chain length should be 1, got %d", length)
	}
}

// 测试AddBlock函数
func TestAddBlock(t *testing.T) {
	bc := NewBlockchain(nil)
	genesis := bc.Genesis()

	// 创建一个有效的新区块
	newBlock := types.NewBlock(
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
		[]byte("Test Block 1"),
		common.Hash{},
		12345,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	// 添加区块
	err := bc.AddBlock(newBlock)
	if err != nil {
		t.Errorf("AddBlock returned error: %v", err)
	}

	// 检查区块是否添加成功
	addedBlock := bc.GetBlock(newBlock.Hash())
	if addedBlock == nil {
		t.Errorf("Block not found after AddBlock")
	}

	if addedBlock.Hash() != newBlock.Hash() {
		t.Errorf("Block hash mismatch")
	}

	// 检查链长度
	length := bc.Length()
	if length != 2 {
		t.Errorf("Chain length should be 2, got %d", length)
	}

	// 检查当前头部
	currentHead := bc.CurrentHead()
	if currentHead.Hash() != newBlock.Hash() {
		t.Errorf("Current head should be the new block")
	}

	// 测试添加已存在的区块
	err = bc.AddBlock(newBlock)
	if err != nil {
		t.Errorf("AddBlock should not return error for existing block")
	}

	// 测试添加无效的区块（父区块不存在）
	invalidBlock := types.NewBlock(
		common.Hash{0xff}, // 不存在的父区块哈希
		common.Address{0x01},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(2),
		10000000,
		0,
		newBlock.Header.Time+10,
		[]byte("Invalid Block"),
		common.Hash{},
		12346,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	err = bc.AddBlock(invalidBlock)
	if err != nil {
		t.Errorf("AddBlock should not return error for block with non-existent parent")
	}

	// 检查无效区块是否未被添加
	invalidBlockAdded := bc.GetBlock(invalidBlock.Hash())
	if invalidBlockAdded != nil {
		t.Errorf("Invalid block should not be added")
	}

	// 测试添加无效的区块（区块号不正确）
	invalidBlockNumber := types.NewBlock(
		newBlock.Hash(),
		common.Address{0x01},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(10), // 区块号跳得太大
		10000000,
		0,
		newBlock.Header.Time+10,
		[]byte("Invalid Block Number"),
		common.Hash{},
		12346,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	err = bc.AddBlock(invalidBlockNumber)
	if err != nil {
		t.Errorf("AddBlock should not return error for block with invalid number")
	}

	// 检查无效区块是否未被添加
	invalidBlockNumberAdded := bc.GetBlock(invalidBlockNumber.Hash())
	if invalidBlockNumberAdded != nil {
		t.Errorf("Invalid block (wrong number) should not be added")
	}
}

// 测试GetBlock和GetBlockByNumber函数
func TestGetBlock(t *testing.T) {
	bc := NewBlockchain(nil)
	genesis := bc.Genesis()

	// 创建两个区块
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
		[]byte("Test Block 1"),
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
		[]byte("Test Block 2"),
		common.Hash{},
		12346,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	// 添加区块
	bc.AddBlock(block1)
	bc.AddBlock(block2)

	// 测试GetBlock
	genesisByHash := bc.GetBlock(genesis.Hash())
	if genesisByHash == nil {
		t.Errorf("GetBlock failed for genesis block")
	}

	block1ByHash := bc.GetBlock(block1.Hash())
	if block1ByHash == nil {
		t.Errorf("GetBlock failed for block 1")
	}

	block2ByHash := bc.GetBlock(block2.Hash())
	if block2ByHash == nil {
		t.Errorf("GetBlock failed for block 2")
	}

	// 测试GetBlockByNumber
	genesisByNumber := bc.GetBlockByNumber(0)
	if genesisByNumber == nil {
		t.Errorf("GetBlockByNumber failed for genesis block")
	}

	block1ByNumber := bc.GetBlockByNumber(1)
	if block1ByNumber == nil {
		t.Errorf("GetBlockByNumber failed for block 1")
	}

	block2ByNumber := bc.GetBlockByNumber(2)
	if block2ByNumber == nil {
		t.Errorf("GetBlockByNumber failed for block 2")
	}

	// 测试获取不存在的区块
	nonExistentBlock := bc.GetBlock(common.Hash{0xff})
	if nonExistentBlock != nil {
		t.Errorf("GetBlock should return nil for non-existent block")
	}

	nonExistentBlockByNumber := bc.GetBlockByNumber(999)
	if nonExistentBlockByNumber != nil {
		t.Errorf("GetBlockByNumber should return nil for non-existent block number")
	}
}

// 测试Length函数
func TestLength(t *testing.T) {
	bc := NewBlockchain(nil)

	// 初始长度应该是1（只有创世区块）
	length := bc.Length()
	if length != 1 {
		t.Errorf("Initial chain length should be 1, got %d", length)
	}

	// 添加一个区块
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
		[]byte("Test Block 1"),
		common.Hash{},
		12345,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)
	bc.AddBlock(block1)

	// 长度应该是2
	length = bc.Length()
	if length != 2 {
		t.Errorf("Chain length should be 2 after adding one block, got %d", length)
	}

	// 添加另一个区块
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
		[]byte("Test Block 2"),
		common.Hash{},
		12346,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)
	bc.AddBlock(block2)

	// 长度应该是3
	length = bc.Length()
	if length != 3 {
		t.Errorf("Chain length should be 3 after adding two blocks, got %d", length)
	}
}

// 测试TransactionPool
func TestTransactionPool(t *testing.T) {
	tp := NewTransactionPool()

	// 测试初始大小
	size := tp.Size()
	if size != 0 {
		t.Errorf("Initial transaction pool size should be 0, got %d", size)
	}

	// 创建测试交易
	tx1 := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	tx2 := types.NewTransaction(
		0,
		common.Address{0x03},
		big.NewInt(1),
		21000,
		big.NewInt(2000),
		[]byte{},
	)

	// 测试AddTransaction
	err := tp.AddTransaction(tx1)
	if err != nil {
		t.Errorf("AddTransaction returned error: %v", err)
	}

	size = tp.Size()
	if size != 1 {
		t.Errorf("Transaction pool size should be 1 after adding one transaction, got %d", size)
	}

	// 测试添加重复交易
	err = tp.AddTransaction(tx1)
	if err != nil {
		t.Errorf("AddTransaction should not return error for duplicate transaction")
	}

	size = tp.Size()
	if size != 1 {
		t.Errorf("Transaction pool size should still be 1 after adding duplicate transaction, got %d", size)
	}

	// 测试GetTransaction
	tx1Retrieved := tp.GetTransaction(tx1.Hash())
	if tx1Retrieved == nil {
		t.Errorf("GetTransaction failed for existing transaction")
	}

	if tx1Retrieved.Hash() != tx1.Hash() {
		t.Errorf("Transaction hash mismatch")
	}

	// 测试获取不存在的交易
	nonExistentTx := tp.GetTransaction(common.Hash{0xff})
	if nonExistentTx != nil {
		t.Errorf("GetTransaction should return nil for non-existent transaction")
	}

	// 测试GetTransactions
	txs := tp.GetTransactions()
	if len(txs) != 1 {
		t.Errorf("GetTransactions should return 1 transaction, got %d", len(txs))
	}

	// 添加第二个交易
	err = tp.AddTransaction(tx2)
	if err != nil {
		t.Errorf("AddTransaction returned error: %v", err)
	}

	size = tp.Size()
	if size != 2 {
		t.Errorf("Transaction pool size should be 2 after adding two transactions, got %d", size)
	}

	txs = tp.GetTransactions()
	if len(txs) != 2 {
		t.Errorf("GetTransactions should return 2 transactions, got %d", len(txs))
	}

	// 测试RemoveTransaction
	tp.RemoveTransaction(tx1.Hash())
	size = tp.Size()
	if size != 1 {
		t.Errorf("Transaction pool size should be 1 after removing one transaction, got %d", size)
	}

	tx1Retrieved = tp.GetTransaction(tx1.Hash())
	if tx1Retrieved != nil {
		t.Errorf("GetTransaction should return nil for removed transaction")
	}

	// 测试RemoveTransactions
	tp.AddTransaction(tx1) // 重新添加
	size = tp.Size()
	if size != 2 {
		t.Errorf("Transaction pool size should be 2 after re adding transaction, got %d", size)
	}

	txsToRemove := []common.Hash{tx1.Hash(), tx2.Hash()}
	tp.RemoveTransactions(txsToRemove)
	size = tp.Size()
	if size != 0 {
		t.Errorf("Transaction pool size should be 0 after removing all transactions, got %d", size)
	}
}

// 测试ValidateTransaction函数
func TestValidateTransaction(t *testing.T) {
	tp := NewTransactionPool()

	// 创建有效的交易
	validTx := types.NewTransaction(
		0,
		common.Address{0x02},
		big.NewInt(1),
		21000,
		big.NewInt(1000),
		[]byte{},
	)

	// 测试验证有效交易
	err := tp.ValidateTransaction(validTx)
	if err != nil {
		t.Errorf("ValidateTransaction should not return error for valid transaction: %v", err)
	}

	// 创建无效的交易（负的 gasPrice）
	invalidTx := types.NewTransaction(
		0,
		common.Address{0x02},
		big.NewInt(1),
		21000,
		big.NewInt(-1), // 负的 gasPrice
		[]byte{},
	)

	// 测试验证无效交易
	err = tp.ValidateTransaction(invalidTx)
	if err == nil {
		t.Errorf("ValidateTransaction should return error for invalid transaction")
	}
}

// 集成测试：测试区块链和交易池的交互
func TestBlockchainWithTransactions(t *testing.T) {
	// 创建区块链
	bc := NewBlockchain(nil)
	genesis := bc.Genesis()

	// 创建交易池
	tp := NewTransactionPool()

	// 创建测试交易
	tx1 := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	tx2 := types.NewTransaction(
		0,
		common.Address{0x03},
		big.NewInt(1),
		21000,
		big.NewInt(2000),
		[]byte{},
	)

	// 添加交易到交易池
	tp.AddTransaction(tx1)
	tp.AddTransaction(tx2)

	// 创建包含这些交易的区块
	blockWithTxs := types.NewBlock(
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
		[]byte("Block with Transactions"),
		common.Hash{},
		12345,
		[]*types.Transaction{tx1, tx2},
		[]*types.BlockHeader{},
	)

	// 添加区块到区块链
	err := bc.AddBlock(blockWithTxs)
	if err != nil {
		t.Errorf("AddBlock returned error: %v", err)
	}

	// 检查区块是否添加成功
	addedBlock := bc.GetBlock(blockWithTxs.Hash())
	if addedBlock == nil {
		t.Errorf("Block with transactions not found after AddBlock")
	}

	// 检查区块中的交易
	if len(addedBlock.Transactions) != 2 {
		t.Errorf("Block should contain 2 transactions, got %d", len(addedBlock.Transactions))
	}

	// 从交易池移除已包含在区块中的交易
	txsToRemove := []common.Hash{tx1.Hash(), tx2.Hash()}
	tp.RemoveTransactions(txsToRemove)

	// 检查交易池是否为空
	size := tp.Size()
	if size != 0 {
		t.Errorf("Transaction pool should be empty after removing transactions in block, got %d", size)
	}
}
