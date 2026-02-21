package unit

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/state"
)

// TestState_NewState 测试创建新状态
func TestState_NewState(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 验证状态不为空
	if s == nil {
		t.Fatal("状态为nil")
	}
}

// TestState_CreateAccount 测试创建账户
func TestState_CreateAccount(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 验证账户存在（通过检查余额是否为0，因为创建账户后余额会被初始化为0）
	balance := s.GetBalance(addr)
	if balance == nil {
		t.Fatal("账户不存在")
	}
}

// TestState_GetBalance 测试获取余额
func TestState_GetBalance(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 验证初始余额为0
	balance := s.GetBalance(addr)
	if balance.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("初始余额错误，期望0，实际%v", balance)
	}
}

// TestState_AddBalance 测试增加余额
func TestState_AddBalance(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 增加余额
	amount := big.NewInt(1000)
	s.AddBalance(addr, amount)

	// 验证余额
	balance := s.GetBalance(addr)
	if balance.Cmp(amount) != 0 {
		t.Errorf("余额错误，期望%v，实际%v", amount, balance)
	}
}

// TestState_SubBalance 测试减少余额
func TestState_SubBalance(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 增加余额
	initialAmount := big.NewInt(1000)
	s.AddBalance(addr, initialAmount)

	// 减少余额
	subAmount := big.NewInt(500)
	s.SubBalance(addr, subAmount)

	// 验证余额
	expectedBalance := big.NewInt(500)
	balance := s.GetBalance(addr)
	if balance.Cmp(expectedBalance) != 0 {
		t.Errorf("余额错误，期望%v，实际%v", expectedBalance, balance)
	}
}

// TestState_GetNonce 测试获取nonce
func TestState_GetNonce(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 验证初始nonce为0
	nonce := s.GetNonce(addr)
	if nonce != 0 {
		t.Errorf("初始nonce错误，期望0，实际%v", nonce)
	}
}

// TestState_SetNonce 测试设置nonce
func TestState_SetNonce(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 设置nonce
	expectedNonce := uint64(5)
	s.SetNonce(addr, expectedNonce)

	// 验证nonce
	nonce := s.GetNonce(addr)
	if nonce != expectedNonce {
		t.Errorf("nonce错误，期望%v，实际%v", expectedNonce, nonce)
	}
}

// TestState_GetState 测试获取状态
func TestState_GetState(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 设置状态
	key := common.Hash{0x02}
	value := common.Hash{0x03}
	s.SetState(addr, key, value)

	// 获取状态
	retrievedValue := s.GetState(addr, key)
	if retrievedValue != value {
		t.Errorf("状态值错误，期望%v，实际%v", value, retrievedValue)
	}
}

// TestState_SetState 测试设置状态
func TestState_SetState(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 设置状态
	key := common.Hash{0x02}
	value := common.Hash{0x03}
	s.SetState(addr, key, value)

	// 验证状态存在
	retrievedValue := s.GetState(addr, key)
	if retrievedValue == (common.Hash{}) {
		t.Fatal("状态不存在")
	}
}

// TestState_GetCode 测试获取代码
func TestState_GetCode(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 设置代码
	code := []byte{0x60, 0x01, 0x60, 0x00, 0xf3} // PUSH1 0x01 PUSH1 0x00 RETURN
	s.SetCode(addr, code)

	// 获取代码
	retrievedCode := s.GetCode(addr)
	if len(retrievedCode) != len(code) {
		t.Errorf("代码长度错误，期望%v，实际%v", len(code), len(retrievedCode))
	}
}

// TestState_SetCode 测试设置代码
func TestState_SetCode(t *testing.T) {
	// 创建新状态
	s := state.NewMemoryStateDB()

	// 创建账户
	addr := common.Address{0x01}
	s.CreateAccount(addr)

	// 设置代码
	code := []byte{0x60, 0x01, 0x60, 0x00, 0xf3} // PUSH1 0x01 PUSH1 0x00 RETURN
	s.SetCode(addr, code)

	// 验证代码存在
	retrievedCode := s.GetCode(addr)
	if retrievedCode == nil {
		t.Fatal("代码不存在")
	}
}
