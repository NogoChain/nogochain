package state

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Account represents an account structure
// Account 账户结构
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash
	CodeHash []byte
}

// StateDB defines the state database interface
// StateDB 状态数据库接口
type StateDB interface {
	CreateAccount(common.Address)
	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int
	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)
	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int
	AddRefund(uint64)
	GetRefund() uint64
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)
	Suicide(common.Address) bool
	HasSuicided(common.Address) bool
	Empty(common.Address) bool
	RevertToSnapshot(int)
	Snapshot() int
	AddLog(Log)
	GetLogs() []Log
	AddPreimage(common.Hash, []byte)
	GetPreimage(common.Hash) []byte
	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)
}

// Log represents a log structure
// Log 日志结构
type Log struct {
	Address     common.Address
	Topics      []common.Hash
	Data        []byte
	BlockNumber uint64
	TxHash      common.Hash
	TxIndex     uint
	BlockHash   common.Hash
	Index       uint
}

// MemoryStateDB implements an in-memory state database
// MemoryStateDB 内存状态数据库实现
type MemoryStateDB struct {
	accounts  map[common.Address]*Account
	storage   map[common.Address]map[common.Hash]common.Hash
	code      map[common.Address][]byte
	logs      []Log
	refund    uint64
	preimages map[common.Hash][]byte
	// 快照相关
	snapshots []snapshot
	// 缓存相关
	accountCache   sync.Map
	storageCache   sync.Map
	codeCache      sync.Map
	stateRootCache common.Hash
	rootCalculated bool
}

// snapshot represents a state snapshot
// snapshot 状态快照
type snapshot struct {
	accounts map[common.Address]*Account
	storage  map[common.Address]map[common.Hash]common.Hash
}

// NewMemoryStateDB creates a new in-memory state database
// NewMemoryStateDB 创建内存状态数据库
func NewMemoryStateDB() *MemoryStateDB {
	return &MemoryStateDB{
		accounts:  make(map[common.Address]*Account),
		storage:   make(map[common.Address]map[common.Hash]common.Hash),
		code:      make(map[common.Address][]byte),
		logs:      make([]Log, 0),
		refund:    0,
		preimages: make(map[common.Hash][]byte),
		snapshots: make([]snapshot, 0),
	}
}

// CreateAccount creates a new account
// CreateAccount 创建账户
func (s *MemoryStateDB) CreateAccount(addr common.Address) {
	if _, exists := s.accounts[addr]; !exists {
		s.accounts[addr] = &Account{
			Nonce:    0,
			Balance:  big.NewInt(0),
			Root:     common.Hash{},
			CodeHash: []byte{},
		}
		s.storage[addr] = make(map[common.Hash]common.Hash)
		// 更新缓存
		s.accountCache.Store(addr, s.accounts[addr])
		// 标记状态根需要重新计算
		s.rootCalculated = false
	}
}

// SubBalance subtracts balance from an account
// SubBalance 减少余额
func (s *MemoryStateDB) SubBalance(addr common.Address, amount *big.Int) {
	s.CreateAccount(addr)
	s.accounts[addr].Balance.Sub(s.accounts[addr].Balance, amount)
	// 更新缓存
	s.accountCache.Store(addr, s.accounts[addr])
	// 标记状态根需要重新计算
	s.rootCalculated = false
}

// AddBalance adds balance to an account
// AddBalance 增加余额
func (s *MemoryStateDB) AddBalance(addr common.Address, amount *big.Int) {
	s.CreateAccount(addr)
	s.accounts[addr].Balance.Add(s.accounts[addr].Balance, amount)
	// 更新缓存
	s.accountCache.Store(addr, s.accounts[addr])
	// 标记状态根需要重新计算
	s.rootCalculated = false
}

// GetBalance retrieves the balance of an account
// GetBalance 获取余额
func (s *MemoryStateDB) GetBalance(addr common.Address) *big.Int {
	// 先从缓存获取
	if acc, ok := s.accountCache.Load(addr); ok {
		return acc.(*Account).Balance
	}
	// 再从内存获取
	if acc, exists := s.accounts[addr]; exists {
		s.accountCache.Store(addr, acc)
		return acc.Balance
	}
	return big.NewInt(0)
}

// GetNonce retrieves the nonce of an account
// GetNonce 获取Nonce
func (s *MemoryStateDB) GetNonce(addr common.Address) uint64 {
	// 先从缓存获取
	if acc, ok := s.accountCache.Load(addr); ok {
		return acc.(*Account).Nonce
	}
	// 再从内存获取
	if acc, exists := s.accounts[addr]; exists {
		s.accountCache.Store(addr, acc)
		return acc.Nonce
	}
	return 0
}

// SetNonce sets the nonce of an account
// SetNonce 设置Nonce
func (s *MemoryStateDB) SetNonce(addr common.Address, nonce uint64) {
	s.CreateAccount(addr)
	s.accounts[addr].Nonce = nonce
	// 更新缓存
	s.accountCache.Store(addr, s.accounts[addr])
	// 标记状态根需要重新计算
	s.rootCalculated = false
}

// GetCodeHash retrieves the code hash of an account
// GetCodeHash 获取代码哈希
func (s *MemoryStateDB) GetCodeHash(addr common.Address) common.Hash {
	// 先从缓存获取
	if acc, ok := s.accountCache.Load(addr); ok {
		return crypto.Keccak256Hash(acc.(*Account).CodeHash)
	}
	// 再从内存获取
	if acc, exists := s.accounts[addr]; exists {
		s.accountCache.Store(addr, acc)
		return crypto.Keccak256Hash(acc.CodeHash)
	}
	return common.Hash{}
}

// GetCode retrieves the code of an account
// GetCode 获取代码
func (s *MemoryStateDB) GetCode(addr common.Address) []byte {
	// 先从缓存获取
	if code, ok := s.codeCache.Load(addr); ok {
		return code.([]byte)
	}
	// 再从内存获取
	if code, exists := s.code[addr]; exists {
		s.codeCache.Store(addr, code)
		return code
	}
	return nil
}

// SetCode sets the code of an account
// SetCode 设置代码
func (s *MemoryStateDB) SetCode(addr common.Address, code []byte) {
	s.CreateAccount(addr)
	s.code[addr] = code
	s.accounts[addr].CodeHash = crypto.Keccak256(code)
	// 更新缓存
	s.codeCache.Store(addr, code)
	s.accountCache.Store(addr, s.accounts[addr])
	// 标记状态根需要重新计算
	s.rootCalculated = false
}

// GetCodeSize retrieves the code size of an account
// GetCodeSize 获取代码大小
func (s *MemoryStateDB) GetCodeSize(addr common.Address) int {
	if code, exists := s.code[addr]; exists {
		return len(code)
	}
	return 0
}

// AddRefund adds gas to the refund counter
// AddRefund 增加退款
func (s *MemoryStateDB) AddRefund(gas uint64) {
	s.refund += gas
}

// GetRefund retrieves the refund counter
// GetRefund 获取退款
func (s *MemoryStateDB) GetRefund() uint64 {
	return s.refund
}

// GetState retrieves the storage state of an account
// GetState 获取存储状态
func (s *MemoryStateDB) GetState(addr common.Address, key common.Hash) common.Hash {
	// 构建缓存键
	cacheKey := common.BytesToHash(append(addr.Bytes(), key.Bytes()...))
	// 先从缓存获取
	if value, ok := s.storageCache.Load(cacheKey); ok {
		return value.(common.Hash)
	}
	// 再从内存获取
	if storage, exists := s.storage[addr]; exists {
		if value, exists := storage[key]; exists {
			s.storageCache.Store(cacheKey, value)
			return value
		}
	}
	return common.Hash{}
}

// SetState sets the storage state of an account
// SetState 设置存储状态
func (s *MemoryStateDB) SetState(addr common.Address, key common.Hash, value common.Hash) {
	s.CreateAccount(addr)
	if _, exists := s.storage[addr]; !exists {
		s.storage[addr] = make(map[common.Hash]common.Hash)
	}
	s.storage[addr][key] = value
	// 更新缓存
	cacheKey := common.BytesToHash(append(addr.Bytes(), key.Bytes()...))
	s.storageCache.Store(cacheKey, value)
	// 标记状态根需要重新计算
	s.rootCalculated = false
}

// Suicide marks an account as suicided
// Suicide 标记账户为自杀
func (s *MemoryStateDB) Suicide(addr common.Address) bool {
	s.CreateAccount(addr)
	return true
}

// HasSuicided checks if an account has suicided
// HasSuicided 检查账户是否已自杀
func (s *MemoryStateDB) HasSuicided(addr common.Address) bool {
	return false
}

// Empty checks if an account is empty
// Empty 检查账户是否为空
func (s *MemoryStateDB) Empty(addr common.Address) bool {
	if acc, exists := s.accounts[addr]; exists {
		return acc.Nonce == 0 && acc.Balance.Sign() == 0 && len(acc.CodeHash) == 0
	}
	return true
}

// RevertToSnapshot reverts the state to a snapshot
// RevertToSnapshot 回滚到快照
func (s *MemoryStateDB) RevertToSnapshot(idx int) {
	if idx < 0 || idx >= len(s.snapshots) {
		return
	}
	snap := s.snapshots[idx]
	s.accounts = snap.accounts
	s.storage = snap.storage
	// 清空缓存
	s.accountCache = sync.Map{}
	s.storageCache = sync.Map{}
	// 标记状态根需要重新计算
	s.rootCalculated = false
}

// Snapshot creates a state snapshot
// Snapshot 创建快照
func (s *MemoryStateDB) Snapshot() int {
	snap := snapshot{
		accounts: make(map[common.Address]*Account),
		storage:  make(map[common.Address]map[common.Hash]common.Hash),
	}
	for addr, acc := range s.accounts {
		accCopy := *acc
		accCopy.Balance = new(big.Int).Set(acc.Balance)
		snap.accounts[addr] = &accCopy
	}
	for addr, storage := range s.storage {
		snap.storage[addr] = make(map[common.Hash]common.Hash)
		for key, value := range storage {
			snap.storage[addr][key] = value
		}
	}
	s.snapshots = append(s.snapshots, snap)
	return len(s.snapshots) - 1
}

// AddLog adds a log to the state
// AddLog 添加日志
func (s *MemoryStateDB) AddLog(log Log) {
	s.logs = append(s.logs, log)
}

// GetLogs retrieves all logs
// GetLogs 获取日志
func (s *MemoryStateDB) GetLogs() []Log {
	return s.logs
}

// AddPreimage adds a preimage to the state
// AddPreimage 添加预映像
func (s *MemoryStateDB) AddPreimage(hash common.Hash, preimage []byte) {
	s.preimages[hash] = preimage
}

// GetPreimage retrieves a preimage from the state
// GetPreimage 获取预映像
func (s *MemoryStateDB) GetPreimage(hash common.Hash) []byte {
	return s.preimages[hash]
}

// ForEachStorage iterates over the storage of an account
// ForEachStorage 遍历存储
func (s *MemoryStateDB) ForEachStorage(addr common.Address, cb func(common.Hash, common.Hash) bool) {
	if storage, exists := s.storage[addr]; exists {
		for key, value := range storage {
			if !cb(key, value) {
				break
			}
		}
	}
}

// CalculateStateRoot calculates the state root
// CalculateStateRoot 计算状态根
func (s *MemoryStateDB) CalculateStateRoot() common.Hash {
	// 检查缓存
	if s.rootCalculated {
		return s.stateRootCache
	}

	// 简化实现，实际应该使用Merkle Patricia Trie
	data := make([]byte, 0)
	for addr, acc := range s.accounts {
		data = append(data, addr.Bytes()...)
		data = append(data, acc.Balance.Bytes()...)
	}

	root := crypto.Keccak256Hash(data)
	// 缓存结果
	s.stateRootCache = root
	s.rootCalculated = true

	return root
}
