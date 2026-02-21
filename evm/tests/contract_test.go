package tests

import (
	"math/big"
	"testing"

	"nogochain/evm/core/vm"
)

// TestContractCreation 测试智能合约创建
func TestContractCreation(t *testing.T) {
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

	// 测试合约创建
	caller := []byte{0x01}
	contractCode := []byte{
		0x60, 0x04, // PUSH1 4
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0x52,       // MSTORE (存储到内存)
		0x60, 0x20, // PUSH1 32 (大小)
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0xf3, // RETURN
	}
	value := big.NewInt(0)
	gas := uint64(100000)

	// 部署合约
	address, returnData, err := evm.Create(caller, contractCode, value, gas)
	if err != nil {
		t.Errorf("Contract creation failed: %v", err)
	}

	// 检查部署结果
	if len(address) == 0 {
		t.Error("Expected contract address, got empty")
	}

	if len(returnData) == 0 {
		t.Error("Expected return data, got empty")
	}

	// 检查合约是否存在
	if !stateDB.Exist(address) {
		t.Error("Contract address does not exist in state DB")
	}

	// 检查合约代码
	code := stateDB.GetCode(address)
	if len(code) == 0 {
		t.Error("Expected contract code, got empty")
	}
}

// TestContractCall 测试智能合约调用
func TestContractCall(t *testing.T) {
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

	// 先部署一个简单的合约
	caller := []byte{0x01}
	// 直接设置合约代码，绕过初始化代码执行
	contractAddress := []byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
	validContractCode := []byte{
		0x60, 0x04, // PUSH1 4
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0x52,       // MSTORE (存储到内存)
		0x60, 0x20, // PUSH1 32 (大小)
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0xf3, // RETURN
	}
	stateDB.CreateAccount(contractAddress)
	stateDB.SetCode(contractAddress, validContractCode)
	value := big.NewInt(0)
	gas := uint64(100000)

	// 调用合约
	input := []byte{}
	returnData, err := evm.Call(caller, contractAddress, input, value, gas)
	if err != nil {
		t.Errorf("Contract call failed: %v", err)
	}

	// 检查调用结果
	if len(returnData) == 0 {
		t.Error("Expected return data, got empty")
	}
}

// TestStaticCall 测试静态调用
func TestStaticCall(t *testing.T) {
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

	// 直接设置合约代码
	caller := []byte{0x01}
	contractAddress := []byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
	validContractCode := []byte{
		0x60, 0x04, // PUSH1 4
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0x52,       // MSTORE (存储到内存)
		0x60, 0x20, // PUSH1 32 (大小)
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0xf3, // RETURN
	}
	stateDB.CreateAccount(contractAddress)
	stateDB.SetCode(contractAddress, validContractCode)
	gas := uint64(100000)

	// 静态调用合约
	input := []byte{}
	returnData, err := evm.StaticCall(caller, contractAddress, input, gas)
	if err != nil {
		t.Errorf("Static call failed: %v", err)
	}

	// 检查调用结果
	if len(returnData) == 0 {
		t.Error("Expected return data, got empty")
	}
}

// TestDelegateCall 测试委托调用
func TestDelegateCall(t *testing.T) {
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

	// 直接设置合约代码
	caller := []byte{0x01}
	contractAddress := []byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
	validContractCode := []byte{
		0x60, 0x04, // PUSH1 4
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0x52,       // MSTORE (存储到内存)
		0x60, 0x20, // PUSH1 32 (大小)
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0xf3, // RETURN
	}
	stateDB.CreateAccount(contractAddress)
	stateDB.SetCode(contractAddress, validContractCode)
	gas := uint64(100000)

	// 委托调用合约
	input := []byte{}
	returnData, err := evm.DelegateCall(caller, contractAddress, input, gas)
	if err != nil {
		t.Errorf("Delegate call failed: %v", err)
	}

	// 检查调用结果
	if len(returnData) == 0 {
		t.Error("Expected return data, got empty")
	}
}

// TestCallCode 测试代码调用
func TestCallCode(t *testing.T) {
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

	// 直接设置合约代码
	caller := []byte{0x01}
	contractAddress := []byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
	validContractCode := []byte{
		0x60, 0x04, // PUSH1 4
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0x52,       // MSTORE (存储到内存)
		0x60, 0x20, // PUSH1 32 (大小)
		0x60, 0x00, // PUSH1 0 (内存偏移量)
		0xf3, // RETURN
	}
	stateDB.CreateAccount(contractAddress)
	stateDB.SetCode(contractAddress, validContractCode)
	value := big.NewInt(0)
	gas := uint64(100000)

	// 代码调用合约
	input := []byte{}
	returnData, err := evm.CallCode(caller, contractAddress, input, value, gas)
	if err != nil {
		t.Errorf("CallCode failed: %v", err)
	}

	// 检查调用结果
	if len(returnData) == 0 {
		t.Error("Expected return data, got empty")
	}
}
