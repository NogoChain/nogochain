package blockchain

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/state"
	"nogochain/core/types"
	"nogochain/metrics"
)

// Blockchain represents the blockchain structure
// Blockchain 区块链结构
type Blockchain struct {
	blocks      map[common.Hash]*types.Block
	blockNumber map[uint64]common.Hash
	stateDB     state.StateDB
	genesis     *types.Block
	currentHead *types.Block
	mu          sync.RWMutex
}

// NewBlockchain creates a new blockchain instance
// NewBlockchain 创建新的区块链
func NewBlockchain(genesis *types.Block) *Blockchain {
	blocks := make(map[common.Hash]*types.Block)
	blockNumber := make(map[uint64]common.Hash)

	if genesis == nil {
		genesis = createGenesisBlock()
	}

	blocks[genesis.Hash()] = genesis
	blockNumber[genesis.NumberU64()] = genesis.Hash()

	return &Blockchain{
		blocks:      blocks,
		blockNumber: blockNumber,
		stateDB:     state.NewMemoryStateDB(),
		genesis:     genesis,
		currentHead: genesis,
	}
}

// createGenesisBlock creates the genesis block
// createGenesisBlock 创建创世区块
func createGenesisBlock() *types.Block {
	genesisTime := uint64(1700000000)
	genesisDifficulty := big.NewInt(1000000)
	genesisNumber := big.NewInt(0)
	genesisGasLimit := uint64(10000000)

	return types.NewBlock(
		common.Hash{},
		common.Address{},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		genesisDifficulty,
		genesisNumber,
		genesisGasLimit,
		0,
		genesisTime,
		[]byte("NogoChain Genesis Block"),
		common.Hash{},
		0,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)
}

// Genesis returns the genesis block
// Genesis 获取创世区块
func (bc *Blockchain) Genesis() *types.Block {
	return bc.genesis
}

// CurrentHead returns the current head block
// CurrentHead 获取当前头部区块
func (bc *Blockchain) CurrentHead() *types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.currentHead
}

// GetBlock retrieves a block by its hash
// GetBlock 通过哈希获取区块
func (bc *Blockchain) GetBlock(hash common.Hash) *types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.blocks[hash]
}

// GetBlockByNumber retrieves a block by its number
// GetBlockByNumber 通过区块号获取区块
func (bc *Blockchain) GetBlockByNumber(number uint64) *types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	if hash, exists := bc.blockNumber[number]; exists {
		return bc.blocks[hash]
	}
	return nil
}

// AddBlock adds a new block to the blockchain
// AddBlock 添加区块
func (bc *Blockchain) AddBlock(block *types.Block) error {
	startTime := time.Now()
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Check if block already exists
	// 检查区块是否已存在
	if _, exists := bc.blocks[block.Hash()]; exists {
		return nil
	}

	// Check if parent block exists
	// 检查父区块是否存在
	parent := bc.blocks[block.ParentHash()]
	if parent == nil {
		return nil
	}

	// Check if block number is correct
	// 检查区块号是否正确
	if block.NumberU64() != parent.NumberU64()+1 {
		return nil
	}

	// Add block to storage
	// 添加区块到存储
	bc.blocks[block.Hash()] = block
	bc.blockNumber[block.NumberU64()] = block.Hash()

	// Update current head
	// 更新当前头部
	if block.NumberU64() > bc.currentHead.NumberU64() {
		bc.currentHead = block
		// Update block height metric
		// 更新区块高度指标
		metrics.BlockHeight.Set(float64(block.NumberU64()))
		// Update block size metric (using transaction count as approximation)
		// 更新区块大小指标（使用交易数量作为近似值）
		metrics.BlockSize.Set(float64(len(block.Transactions)))
		// Update transaction count metric
		// 更新交易计数指标
		metrics.TransactionCount.Add(float64(len(block.Transactions)))
	}

	// Record block processing time
	// 记录区块处理时间
	metrics.BlockProcessingTime.Observe(time.Since(startTime).Seconds())

	return nil
}

// Length returns the length of the blockchain
// Length 获取链长度
func (bc *Blockchain) Length() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.currentHead.NumberU64() + 1
}

// StateDB returns the state database
// StateDB 获取状态数据库
func (bc *Blockchain) StateDB() state.StateDB {
	return bc.stateDB
}

// TransactionPool represents the transaction pool
// TransactionPool 交易池
type TransactionPool struct {
	txs map[common.Hash]*types.Transaction
	mu  sync.RWMutex
}

// NewTransactionPool creates a new transaction pool
// NewTransactionPool 创建新的交易池
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		txs: make(map[common.Hash]*types.Transaction),
	}
}

// AddTransaction adds a transaction to the pool
// AddTransaction 添加交易到交易池
func (tp *TransactionPool) AddTransaction(tx *types.Transaction) error {
	startTime := time.Now()
	tp.mu.Lock()
	defer tp.mu.Unlock()

	txHash := tx.Hash()
	if _, exists := tp.txs[txHash]; exists {
		return nil
	}

	tp.txs[txHash] = tx

	// Record transaction processing time
	// 记录交易处理时间
	metrics.TransactionProcessingTime.Observe(time.Since(startTime).Seconds())

	return nil
}

// GetTransaction retrieves a transaction by its hash
// GetTransaction 通过哈希获取交易
func (tp *TransactionPool) GetTransaction(hash common.Hash) *types.Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	return tp.txs[hash]
}

// GetTransactions retrieves all transactions in the pool
// GetTransactions 获取交易池中的所有交易
func (tp *TransactionPool) GetTransactions() []*types.Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	txs := make([]*types.Transaction, 0, len(tp.txs))
	for _, tx := range tp.txs {
		txs = append(txs, tx)
	}
	return txs
}

// RemoveTransaction removes a transaction from the pool
// RemoveTransaction 从交易池移除交易
func (tp *TransactionPool) RemoveTransaction(hash common.Hash) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	delete(tp.txs, hash)
}

// RemoveTransactions removes multiple transactions from the pool
// RemoveTransactions 从交易池移除多个交易
func (tp *TransactionPool) RemoveTransactions(hashes []common.Hash) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	for _, hash := range hashes {
		delete(tp.txs, hash)
	}
}

// Size returns the size of the transaction pool
// Size 获取交易池大小
func (tp *TransactionPool) Size() int {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	return len(tp.txs)
}

// ValidateTransaction validates a transaction
// ValidateTransaction 验证交易
func (tp *TransactionPool) ValidateTransaction(tx *types.Transaction) error {
	return tx.Validate()
}
