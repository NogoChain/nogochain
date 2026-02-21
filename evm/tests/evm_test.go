package tests

import (
	"math/big"
	"testing"

	"nogochain/evm/core/vm"
	"nogochain/evm/core/vm/gas"
)

// MockStateDB 模拟状态数据库
type MockStateDB struct {
	balances map[string]*big.Int
	codes    map[string][]byte
	nonces   map[string]uint64
	states   map[string]map[string][]byte
	accounts map[string]bool
}

// NewMockStateDB 创建新的模拟状态数据库
func NewMockStateDB() *MockStateDB {
	return &MockStateDB{
		balances: make(map[string]*big.Int),
		codes:    make(map[string][]byte),
		nonces:   make(map[string]uint64),
		states:   make(map[string]map[string][]byte),
		accounts: make(map[string]bool),
	}
}

func (m *MockStateDB) GetBalance(addr []byte) *big.Int {
	key := string(addr)
	if balance, exists := m.balances[key]; exists {
		return balance
	}
	return big.NewInt(0)
}

func (m *MockStateDB) GetCode(addr []byte) []byte {
	key := string(addr)
	if code, exists := m.codes[key]; exists {
		return code
	}
	return nil
}

func (m *MockStateDB) GetNonce(addr []byte) uint64 {
	key := string(addr)
	if nonce, exists := m.nonces[key]; exists {
		return nonce
	}
	return 0
}

func (m *MockStateDB) SetNonce(addr []byte, nonce uint64) {
	key := string(addr)
	m.nonces[key] = nonce
}

func (m *MockStateDB) GetState(addr, key []byte) []byte {
	addrKey := string(addr)
	keyStr := string(key)
	if state, exists := m.states[addrKey]; exists {
		if value, exists := state[keyStr]; exists {
			return value
		}
	}
	return make([]byte, 32)
}

func (m *MockStateDB) SetState(addr, key, value []byte) {
	addrKey := string(addr)
	keyStr := string(key)
	if _, exists := m.states[addrKey]; !exists {
		m.states[addrKey] = make(map[string][]byte)
	}
	m.states[addrKey][keyStr] = value
}

func (m *MockStateDB) SetCode(addr []byte, code []byte) {
	key := string(addr)
	m.codes[key] = code
}

func (m *MockStateDB) AddBalance(addr []byte, amount *big.Int) {
	key := string(addr)
	if _, exists := m.balances[key]; !exists {
		m.balances[key] = big.NewInt(0)
	}
	m.balances[key] = m.balances[key].Add(m.balances[key], amount)
}

func (m *MockStateDB) SubBalance(addr []byte, amount *big.Int) {
	key := string(addr)
	if _, exists := m.balances[key]; !exists {
		m.balances[key] = big.NewInt(0)
	}
	m.balances[key] = m.balances[key].Sub(m.balances[key], amount)
}

func (m *MockStateDB) CreateAccount(addr []byte) {
	key := string(addr)
	m.accounts[key] = true
	if _, exists := m.balances[key]; !exists {
		m.balances[key] = big.NewInt(0)
	}
	if _, exists := m.nonces[key]; !exists {
		m.nonces[key] = 0
	}
}

func (m *MockStateDB) Exist(addr []byte) bool {
	key := string(addr)
	return m.accounts[key]
}

// TestEVMExecute 测试EVM执行
func TestEVMExecute(t *testing.T) {
	// 创建测试环境
	stateDB := NewMockStateDB()
	context := vm.Context{
		Caller:      []byte{0x01}, // 测试地址
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1),
		Timestamp:   big.NewInt(1000),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(1),
	}

	header := &vm.BlockHeader{
		Coinbase:   []byte{0x00},
		GasLimit:   1000000,
		Number:     big.NewInt(1),
		Timestamp:  big.NewInt(1000),
		BaseFee:    big.NewInt(1),
		Difficulty: big.NewInt(1000),
	}

	// 创建EVM
	evm := vm.NewEVM(context, stateDB, header)

	// 测试简单的加法指令
	testCode := []byte{
		0x60, 0x01, // PUSH1 1
		0x60, 0x02, // PUSH1 2
		0x01,       // ADD
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0x52,       // MSTORE (存储结果到内存)
		0x60, 0x20, // PUSH1 32 (大小)
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0xf3, // RETURN
	}

	// 执行代码
	returnData, err := evm.Run(testCode)
	if err != nil {
		t.Errorf("EVM execution failed: %v", err)
	}

	// 检查结果
	if len(returnData) == 0 {
		t.Error("Expected return data, got empty")
	}
}

// TestGasCalculation 测试Gas计算
func TestGasCalculation(t *testing.T) {
	// 测试基础Gas成本
	testCases := []struct {
		opcode   byte
		expected uint64
	}{
		{0x00, 0},     // STOP
		{0x01, 3},     // ADD
		{0x02, 5},     // MUL
		{0x03, 3},     // SUB
		{0x04, 5},     // DIV
		{0x54, 800},   // SLOAD
		{0x55, 20000}, // SSTORE
	}

	for _, tc := range testCases {
		actual := gas.CalculateBaseGas(tc.opcode)
		if actual != tc.expected {
			t.Errorf("Gas cost for opcode 0x%02x: expected %d, got %d", tc.opcode, tc.expected, actual)
		}
	}
}

// TestStackOperations 测试栈操作
func TestStackOperations(t *testing.T) {
	// 测试栈的基本操作
	stateDB := NewMockStateDB()
	context := vm.Context{
		Caller:      []byte{0x01},
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1),
		Timestamp:   big.NewInt(1000),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(1),
	}

	header := &vm.BlockHeader{
		Coinbase:   []byte{0x00},
		GasLimit:   1000000,
		Number:     big.NewInt(1),
		Timestamp:  big.NewInt(1000),
		BaseFee:    big.NewInt(1),
		Difficulty: big.NewInt(1000),
	}

	evm := vm.NewEVM(context, stateDB, header)

	// 测试栈压入和弹出
	testValue := big.NewInt(42)
	err := evm.Stack.Push(testValue)
	if err != nil {
		t.Errorf("Failed to push to stack: %v", err)
	}

	popped, err := evm.Stack.Pop()
	if err != nil {
		t.Errorf("Failed to pop from stack: %v", err)
	}

	if popped.Cmp(testValue) != 0 {
		t.Errorf("Stack pop: expected %v, got %v", testValue, popped)
	}
}

// TestMemoryOperations 测试内存操作
func TestMemoryOperations(t *testing.T) {
	stateDB := NewMockStateDB()
	context := vm.Context{
		Caller:      []byte{0x01},
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1),
		Timestamp:   big.NewInt(1000),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(1),
	}

	header := &vm.BlockHeader{
		Coinbase:   []byte{0x00},
		GasLimit:   1000000,
		Number:     big.NewInt(1),
		Timestamp:  big.NewInt(1000),
		BaseFee:    big.NewInt(1),
		Difficulty: big.NewInt(1000),
	}

	evm := vm.NewEVM(context, stateDB, header)

	// 测试内存写入和读取
	testData := []byte{0x01, 0x02, 0x03, 0x04}
	evm.Memory.Set(0, testData)

	readData := evm.Memory.Get(0, 4)
	for i, b := range readData {
		if b != testData[i] {
			t.Errorf("Memory read at offset %d: expected %02x, got %02x", i, testData[i], b)
		}
	}
}

// TestStorageOperations 测试存储操作
func TestStorageOperations(t *testing.T) {
	stateDB := NewMockStateDB()
	context := vm.Context{
		Caller:      []byte{0x01},
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1),
		Timestamp:   big.NewInt(1000),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(1),
	}

	header := &vm.BlockHeader{
		Coinbase:   []byte{0x00},
		GasLimit:   1000000,
		Number:     big.NewInt(1),
		Timestamp:  big.NewInt(1000),
		BaseFee:    big.NewInt(1),
		Difficulty: big.NewInt(1000),
	}

	evm := vm.NewEVM(context, stateDB, header)

	// 测试存储写入和读取
	testKey := []byte{0x00}
	testValue := []byte{0x01, 0x02, 0x03, 0x04}
	evm.Storage.Set(testKey, testValue)

	readValue := evm.Storage.Get(testKey)
	for i, b := range readValue {
		if i < len(testValue) && b != testValue[i] {
			t.Errorf("Storage read at key %v: expected %02x, got %02x", testKey, testValue[i], b)
		}
	}
}
