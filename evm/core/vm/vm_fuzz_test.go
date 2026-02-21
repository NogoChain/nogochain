package vm

import (
	"math/big"
	"testing"
)

// MockStateDB 模拟StateDB实现
type MockStateDB struct {
	balances map[string]*big.Int
	codes    map[string][]byte
	nonces   map[string]uint64
	states   map[string]map[string][]byte
	accounts map[string]bool
}

// NewMockStateDB 创建新的MockStateDB
func NewMockStateDB() *MockStateDB {
	return &MockStateDB{
		balances: make(map[string]*big.Int),
		codes:    make(map[string][]byte),
		nonces:   make(map[string]uint64),
		states:   make(map[string]map[string][]byte),
		accounts: make(map[string]bool),
	}
}

// GetBalance 获取余额
func (m *MockStateDB) GetBalance(addr []byte) *big.Int {
	key := string(addr)
	if bal, ok := m.balances[key]; ok {
		return bal
	}
	return big.NewInt(0)
}

// GetCode 获取代码
func (m *MockStateDB) GetCode(addr []byte) []byte {
	key := string(addr)
	if code, ok := m.codes[key]; ok {
		return code
	}
	return nil
}

// GetNonce 获取nonce
func (m *MockStateDB) GetNonce(addr []byte) uint64 {
	key := string(addr)
	if nonce, ok := m.nonces[key]; ok {
		return nonce
	}
	return 0
}

// SetNonce 设置nonce
func (m *MockStateDB) SetNonce(addr []byte, nonce uint64) {
	key := string(addr)
	m.nonces[key] = nonce
}

// GetState 获取状态
func (m *MockStateDB) GetState(addr, key []byte) []byte {
	addrKey := string(addr)
	keyKey := string(key)
	if stateMap, ok := m.states[addrKey]; ok {
		if val, ok := stateMap[keyKey]; ok {
			return val
		}
	}
	return nil
}

// SetState 设置状态
func (m *MockStateDB) SetState(addr, key, value []byte) {
	addrKey := string(addr)
	keyKey := string(key)
	if _, ok := m.states[addrKey]; !ok {
		m.states[addrKey] = make(map[string][]byte)
	}
	m.states[addrKey][keyKey] = value
}

// SetCode 设置代码
func (m *MockStateDB) SetCode(addr []byte, code []byte) {
	key := string(addr)
	m.codes[key] = code
}

// AddBalance 增加余额
func (m *MockStateDB) AddBalance(addr []byte, amount *big.Int) {
	key := string(addr)
	if _, ok := m.balances[key]; !ok {
		m.balances[key] = big.NewInt(0)
	}
	m.balances[key] = m.balances[key].Add(m.balances[key], amount)
}

// SubBalance 减少余额
func (m *MockStateDB) SubBalance(addr []byte, amount *big.Int) {
	key := string(addr)
	if _, ok := m.balances[key]; !ok {
		m.balances[key] = big.NewInt(0)
	}
	m.balances[key] = m.balances[key].Sub(m.balances[key], amount)
}

// CreateAccount 创建账户
func (m *MockStateDB) CreateAccount(addr []byte) {
	key := string(addr)
	m.accounts[key] = true
	if _, ok := m.balances[key]; !ok {
		m.balances[key] = big.NewInt(0)
	}
	if _, ok := m.nonces[key]; !ok {
		m.nonces[key] = 0
	}
	if _, ok := m.states[key]; !ok {
		m.states[key] = make(map[string][]byte)
	}
}

// Exist 检查账户是否存在
func (m *MockStateDB) Exist(addr []byte) bool {
	key := string(addr)
	_, ok := m.accounts[key]
	return ok
}

// Fuzz测试EVM Run方法
func FuzzEVM_Run(f *testing.F) {
	// 测试用例：简单的RETURN指令
	f.Add([]byte{0x60, 0x01, 0x60, 0x00, 0xf3}) // PUSH1 0x01 PUSH1 0x00 RETURN
	// 测试用例：空代码
	f.Add([]byte{})
	// 测试用例：简单的ADD指令
	f.Add([]byte{0x60, 0x01, 0x60, 0x02, 0x01, 0xf3}) // PUSH1 0x01 PUSH1 0x02 ADD RETURN

	f.Fuzz(func(t *testing.T, code []byte) {
		// 创建执行上下文
		context := Context{
			Caller:      []byte{0x01}, // 模拟调用者地址
			GasPrice:    big.NewInt(1),
			Origin:      []byte{0x01},
			BlockNumber: big.NewInt(1),
			Timestamp:   big.NewInt(1),
			GasLimit:    1000000,
			BaseFee:     big.NewInt(0),
			Code:        code,
		}

		// 创建区块头
		header := &BlockHeader{
			Coinbase:   []byte{0x02},
			GasLimit:   1000000,
			Number:     big.NewInt(1),
			Timestamp:  big.NewInt(1),
			BaseFee:    big.NewInt(0),
			Difficulty: big.NewInt(1000),
		}

		// 创建模拟StateDB
		stateDB := NewMockStateDB()

		// 创建EVM实例
		evm := NewEVM(context, stateDB, header)

		// 运行EVM（应该不会panic）
		_, err := evm.Run(code)
		_ = err
	})
}

// Fuzz测试EVM Create方法
func FuzzEVM_Create(f *testing.F) {
	// 测试用例：简单的合约创建
	f.Add([]byte{0x60, 0x01, 0x60, 0x00, 0xf3}, int64(0), uint64(1000000)) // 代码，价值，gas

	f.Fuzz(func(t *testing.T, code []byte, value int64, gas uint64) {
		// 创建执行上下文
		context := Context{
			Caller:      []byte{0x01}, // 模拟调用者地址
			GasPrice:    big.NewInt(1),
			Origin:      []byte{0x01},
			BlockNumber: big.NewInt(1),
			Timestamp:   big.NewInt(1),
			GasLimit:    1000000,
			BaseFee:     big.NewInt(0),
			Code:        code,
		}

		// 创建区块头
		header := &BlockHeader{
			Coinbase:   []byte{0x02},
			GasLimit:   1000000,
			Number:     big.NewInt(1),
			Timestamp:  big.NewInt(1),
			BaseFee:    big.NewInt(0),
			Difficulty: big.NewInt(1000),
		}

		// 创建模拟StateDB
		stateDB := NewMockStateDB()
		// 为调用者添加余额
		stateDB.CreateAccount([]byte{0x01})
		stateDB.AddBalance([]byte{0x01}, big.NewInt(1000000))

		// 创建EVM实例
		evm := NewEVM(context, stateDB, header)

		// 创建合约（应该不会panic）
		_, _, err := evm.Create([]byte{0x01}, code, big.NewInt(value), gas)
		_ = err
	})
}

// Fuzz测试EVM Call方法
func FuzzEVM_Call(f *testing.F) {
	// 测试用例：简单的合约调用
	f.Add([]byte{0x01}, []byte{0x02}, []byte{}, int64(0), uint64(1000000)) // 调用者，目标，输入，价值，gas

	f.Fuzz(func(t *testing.T, caller, to, input []byte, value int64, gas uint64) {
		// 创建执行上下文
		context := Context{
			Caller:      caller,
			GasPrice:    big.NewInt(1),
			Origin:      caller,
			BlockNumber: big.NewInt(1),
			Timestamp:   big.NewInt(1),
			GasLimit:    1000000,
			BaseFee:     big.NewInt(0),
			Code:        nil,
		}

		// 创建区块头
		header := &BlockHeader{
			Coinbase:   []byte{0x03},
			GasLimit:   1000000,
			Number:     big.NewInt(1),
			Timestamp:  big.NewInt(1),
			BaseFee:    big.NewInt(0),
			Difficulty: big.NewInt(1000),
		}

		// 创建模拟StateDB
		stateDB := NewMockStateDB()
		// 为调用者和目标账户添加余额和代码
		stateDB.CreateAccount(caller)
		stateDB.AddBalance(caller, big.NewInt(1000000))
		stateDB.CreateAccount(to)
		stateDB.AddBalance(to, big.NewInt(1000000))
		stateDB.SetCode(to, []byte{0x60, 0x01, 0x60, 0x00, 0xf3}) // 简单的返回指令

		// 创建EVM实例
		evm := NewEVM(context, stateDB, header)

		// 调用合约（应该不会panic）
		_, err := evm.Call(caller, to, input, big.NewInt(value), gas)
		_ = err
	})
}

// Fuzz测试EVM ValidateTransaction方法
func FuzzEVM_ValidateTransaction(f *testing.F) {
	// 测试用例：传统交易参数
	f.Add([]byte{0x02}, []byte{0x01}, uint64(0), int64(1), uint64(21000), int64(1000), []byte{})

	f.Fuzz(func(t *testing.T, to, from []byte, nonce uint64, gasPrice int64, gasLimit uint64, value int64, data []byte) {
		// 创建交易
		tx := &Transaction{
			To:       to,
			From:     from,
			Nonce:    nonce,
			GasPrice: big.NewInt(gasPrice),
			GasLimit: gasLimit,
			Value:    big.NewInt(value),
			Data:     data,
		}

		// 创建执行上下文
		context := Context{
			Caller:      []byte{0x01},
			GasPrice:    big.NewInt(1),
			Origin:      []byte{0x01},
			BlockNumber: big.NewInt(1559), // 激活EIP-1559
			Timestamp:   big.NewInt(1),
			GasLimit:    1000000,
			BaseFee:     big.NewInt(0),
			Code:        nil,
		}

		// 创建区块头
		header := &BlockHeader{
			Coinbase:   []byte{0x03},
			GasLimit:   1000000,
			Number:     big.NewInt(1559), // 激活EIP-1559
			Timestamp:  big.NewInt(1),
			BaseFee:    big.NewInt(0),
			Difficulty: big.NewInt(1000),
		}

		// 创建模拟StateDB
		stateDB := NewMockStateDB()

		// 创建EVM实例
		evm := NewEVM(context, stateDB, header)

		// 验证交易（应该不会panic）
		err := evm.ValidateTransaction(tx)
		_ = err
	})
}

// Benchmark测试EVM Run方法
func BenchmarkEVM_Run(b *testing.B) {
	// 简单的合约代码：返回1
	code := []byte{0x60, 0x01, 0x60, 0x00, 0xf3} // PUSH1 0x01 PUSH1 0x00 RETURN

	// 创建执行上下文
	context := Context{
		Caller:      []byte{0x01},
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1),
		Timestamp:   big.NewInt(1),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(0),
		Code:        code,
	}

	// 创建区块头
	header := &BlockHeader{
		Coinbase:   []byte{0x02},
		GasLimit:   1000000,
		Number:     big.NewInt(1),
		Timestamp:  big.NewInt(1),
		BaseFee:    big.NewInt(0),
		Difficulty: big.NewInt(1000),
	}

	// 创建模拟StateDB
	stateDB := NewMockStateDB()

	// 创建EVM实例
	evm := NewEVM(context, stateDB, header)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evm.Run(code)
	}
}

// Benchmark测试EVM Create方法
func BenchmarkEVM_Create(b *testing.B) {
	// 简单的合约代码：返回1
	code := []byte{0x60, 0x01, 0x60, 0x00, 0xf3} // PUSH1 0x01 PUSH1 0x00 RETURN

	// 创建执行上下文
	context := Context{
		Caller:      []byte{0x01},
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1),
		Timestamp:   big.NewInt(1),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(0),
		Code:        code,
	}

	// 创建区块头
	header := &BlockHeader{
		Coinbase:   []byte{0x02},
		GasLimit:   1000000,
		Number:     big.NewInt(1),
		Timestamp:  big.NewInt(1),
		BaseFee:    big.NewInt(0),
		Difficulty: big.NewInt(1000),
	}

	// 创建模拟StateDB
	stateDB := NewMockStateDB()
	stateDB.CreateAccount([]byte{0x01})
	stateDB.AddBalance([]byte{0x01}, big.NewInt(1000000))

	// 创建EVM实例
	evm := NewEVM(context, stateDB, header)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evm.Create([]byte{0x01}, code, big.NewInt(0), 1000000)
	}
}

// Benchmark测试EVM Call方法
func BenchmarkEVM_Call(b *testing.B) {
	// 简单的合约代码：返回1
	code := []byte{0x60, 0x01, 0x60, 0x00, 0xf3} // PUSH1 0x01 PUSH1 0x00 RETURN

	// 创建执行上下文
	context := Context{
		Caller:      []byte{0x01},
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1),
		Timestamp:   big.NewInt(1),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(0),
		Code:        nil,
	}

	// 创建区块头
	header := &BlockHeader{
		Coinbase:   []byte{0x03},
		GasLimit:   1000000,
		Number:     big.NewInt(1),
		Timestamp:  big.NewInt(1),
		BaseFee:    big.NewInt(0),
		Difficulty: big.NewInt(1000),
	}

	// 创建模拟StateDB
	stateDB := NewMockStateDB()
	stateDB.CreateAccount([]byte{0x01})
	stateDB.AddBalance([]byte{0x01}, big.NewInt(1000000))
	stateDB.CreateAccount([]byte{0x02})
	stateDB.AddBalance([]byte{0x02}, big.NewInt(1000000))
	stateDB.SetCode([]byte{0x02}, code)

	// 创建EVM实例
	evm := NewEVM(context, stateDB, header)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evm.Call([]byte{0x01}, []byte{0x02}, []byte{}, big.NewInt(0), 1000000)
	}
}

// Benchmark测试EVM ValidateTransaction方法
func BenchmarkEVM_ValidateTransaction(b *testing.B) {
	// 创建传统交易
	tx := &Transaction{
		To:       []byte{0x02},
		From:     []byte{0x01},
		Nonce:    0,
		GasPrice: big.NewInt(1),
		GasLimit: 21000,
		Value:    big.NewInt(1000),
		Data:     nil,
	}

	// 创建执行上下文
	context := Context{
		Caller:      []byte{0x01},
		GasPrice:    big.NewInt(1),
		Origin:      []byte{0x01},
		BlockNumber: big.NewInt(1559), // 激活EIP-1559
		Timestamp:   big.NewInt(1),
		GasLimit:    1000000,
		BaseFee:     big.NewInt(0),
		Code:        nil,
	}

	// 创建区块头
	header := &BlockHeader{
		Coinbase:   []byte{0x03},
		GasLimit:   1000000,
		Number:     big.NewInt(1559), // 激活EIP-1559
		Timestamp:  big.NewInt(1),
		BaseFee:    big.NewInt(0),
		Difficulty: big.NewInt(1000),
	}

	// 创建模拟StateDB
	stateDB := NewMockStateDB()

	// 创建EVM实例
	evm := NewEVM(context, stateDB, header)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evm.ValidateTransaction(tx)
	}
}
