package state

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// 测试NewMemoryStateDB函数
func TestNewMemoryStateDB(t *testing.T) {
	sdb := NewMemoryStateDB()
	if sdb == nil {
		t.Errorf("NewMemoryStateDB returned nil")
	}

	// 检查初始状态
	if len(sdb.accounts) != 0 {
		t.Errorf("Initial accounts map should be empty")
	}

	if len(sdb.storage) != 0 {
		t.Errorf("Initial storage map should be empty")
	}

	if len(sdb.code) != 0 {
		t.Errorf("Initial code map should be empty")
	}

	if len(sdb.logs) != 0 {
		t.Errorf("Initial logs should be empty")
	}

	if sdb.refund != 0 {
		t.Errorf("Initial refund should be 0")
	}

	if len(sdb.preimages) != 0 {
		t.Errorf("Initial preimages map should be empty")
	}

	if len(sdb.snapshots) != 0 {
		t.Errorf("Initial snapshots should be empty")
	}
}

// 测试账户创建和余额操作
func TestAccountBalanceOperations(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr := common.Address{0x01}

	// 测试创建账户
	sdb.CreateAccount(addr)
	if _, exists := sdb.accounts[addr]; !exists {
		t.Errorf("Account should be created")
	}

	// 测试初始余额
	balance := sdb.GetBalance(addr)
	if balance.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Initial balance should be 0, got %v", balance)
	}

	// 测试增加余额
	amount := big.NewInt(1000)
	sdb.AddBalance(addr, amount)
	balance = sdb.GetBalance(addr)
	if balance.Cmp(amount) != 0 {
		t.Errorf("Balance should be %v, got %v", amount, balance)
	}

	// 测试减少余额
	subAmount := big.NewInt(500)
	sdb.SubBalance(addr, subAmount)
	balance = sdb.GetBalance(addr)
	expectedBalance := big.NewInt(500)
	if balance.Cmp(expectedBalance) != 0 {
		t.Errorf("Balance should be %v, got %v", expectedBalance, balance)
	}

	// 测试对不存在账户的操作（应该自动创建）
	newAddr := common.Address{0x02}
	sdb.AddBalance(newAddr, big.NewInt(2000))
	balance = sdb.GetBalance(newAddr)
	if balance.Cmp(big.NewInt(2000)) != 0 {
		t.Errorf("Balance for new account should be 2000, got %v", balance)
	}

	if _, exists := sdb.accounts[newAddr]; !exists {
		t.Errorf("Account should be automatically created")
	}
}

// 测试Nonce操作
func TestNonceOperations(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr := common.Address{0x01}

	// 测试初始Nonce（账户不存在）
	nonce := sdb.GetNonce(addr)
	if nonce != 0 {
		t.Errorf("Initial nonce for non-existent account should be 0, got %d", nonce)
	}

	// 创建账户并测试Nonce
	sdb.CreateAccount(addr)
	nonce = sdb.GetNonce(addr)
	if nonce != 0 {
		t.Errorf("Initial nonce for new account should be 0, got %d", nonce)
	}

	// 测试设置Nonce
	newNonce := uint64(5)
	sdb.SetNonce(addr, newNonce)
	nonce = sdb.GetNonce(addr)
	if nonce != newNonce {
		t.Errorf("Nonce should be %d, got %d", newNonce, nonce)
	}

	// 测试对不存在账户设置Nonce（应该自动创建）
	newAddr := common.Address{0x02}
	sdb.SetNonce(newAddr, 10)
	nonce = sdb.GetNonce(newAddr)
	if nonce != 10 {
		t.Errorf("Nonce for new account should be 10, got %d", nonce)
	}

	if _, exists := sdb.accounts[newAddr]; !exists {
		t.Errorf("Account should be automatically created")
	}
}

// 测试代码操作
func TestCodeOperations(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr := common.Address{0x01}

	// 测试获取不存在账户的代码
	code := sdb.GetCode(addr)
	if code != nil {
		t.Errorf("Code for non-existent account should be nil")
	}

	codeSize := sdb.GetCodeSize(addr)
	if codeSize != 0 {
		t.Errorf("Code size for non-existent account should be 0, got %d", codeSize)
	}

	codeHash := sdb.GetCodeHash(addr)
	if codeHash != (common.Hash{}) {
		t.Errorf("Code hash for non-existent account should be zero hash")
	}

	// 创建账户并设置代码
	testCode := []byte{0x60, 0x01, 0x60, 0x00, 0xf3} // PUSH1 0x01 PUSH1 0x00 RETURN
	sdb.CreateAccount(addr)
	sdb.SetCode(addr, testCode)

	// 测试获取代码
	code = sdb.GetCode(addr)
	if len(code) != len(testCode) {
		t.Errorf("Code length should be %d, got %d", len(testCode), len(code))
	}

	codeSize = sdb.GetCodeSize(addr)
	if codeSize != len(testCode) {
		t.Errorf("Code size should be %d, got %d", len(testCode), codeSize)
	}

	codeHash = sdb.GetCodeHash(addr)
	if codeHash == (common.Hash{}) {
		t.Errorf("Code hash should not be zero hash")
	}
}

// 测试存储操作
func TestStorageOperations(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr := common.Address{0x01}
	key := common.HexToHash("0x01")
	value := common.HexToHash("0x1234")

	// 测试获取不存在账户的存储
	storedValue := sdb.GetState(addr, key)
	if storedValue != (common.Hash{}) {
		t.Errorf("Storage value for non-existent account should be zero hash")
	}

	// 创建账户并设置存储
	sdb.CreateAccount(addr)
	sdb.SetState(addr, key, value)

	// 测试获取存储
	storedValue = sdb.GetState(addr, key)
	if storedValue != value {
		t.Errorf("Storage value should be %v, got %v", value, storedValue)
	}

	// 测试遍历存储
	count := 0
	sdb.ForEachStorage(addr, func(k, v common.Hash) bool {
		count++
		if k != key {
			t.Errorf("Storage key should be %v, got %v", key, k)
		}
		if v != value {
			t.Errorf("Storage value should be %v, got %v", value, v)
		}
		return true
	})

	if count != 1 {
		t.Errorf("ForEachStorage should iterate over 1 item, got %d", count)
	}

	// 测试对不存在账户设置存储（应该自动创建）
	newAddr := common.Address{0x02}
	newKey := common.HexToHash("0x02")
	newValue := common.HexToHash("0x5678")
	sdb.SetState(newAddr, newKey, newValue)

	storedValue = sdb.GetState(newAddr, newKey)
	if storedValue != newValue {
		t.Errorf("Storage value for new account should be %v, got %v", newValue, storedValue)
	}

	if _, exists := sdb.accounts[newAddr]; !exists {
		t.Errorf("Account should be automatically created")
	}
}

// 测试快照和回滚功能
func TestSnapshotAndRevert(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr := common.Address{0x01}

	// 创建账户并设置初始状态
	sdb.CreateAccount(addr)
	sdb.AddBalance(addr, big.NewInt(1000))
	sdb.SetNonce(addr, 1)
	key := common.HexToHash("0x01")
	value := common.HexToHash("0x1234")
	sdb.SetState(addr, key, value)

	// 创建快照
	snapshotIdx := sdb.Snapshot()

	// 修改状态
	sdb.AddBalance(addr, big.NewInt(500))
	sdb.SetNonce(addr, 2)
	newValue := common.HexToHash("0x5678")
	sdb.SetState(addr, key, newValue)

	// 检查修改后的状态
	balance := sdb.GetBalance(addr)
	if balance.Cmp(big.NewInt(1500)) != 0 {
		t.Errorf("Balance should be 1500 after modification, got %v", balance)
	}

	nonce := sdb.GetNonce(addr)
	if nonce != 2 {
		t.Errorf("Nonce should be 2 after modification, got %d", nonce)
	}

	storedValue := sdb.GetState(addr, key)
	if storedValue != newValue {
		t.Errorf("Storage value should be %v after modification, got %v", newValue, storedValue)
	}

	// 回滚到快照
	sdb.RevertToSnapshot(snapshotIdx)

	// 检查回滚后的状态
	balance = sdb.GetBalance(addr)
	if balance.Cmp(big.NewInt(1000)) != 0 {
		t.Errorf("Balance should be 1000 after revert, got %v", balance)
	}

	nonce = sdb.GetNonce(addr)
	if nonce != 1 {
		t.Errorf("Nonce should be 1 after revert, got %d", nonce)
	}

	storedValue = sdb.GetState(addr, key)
	if storedValue != value {
		t.Errorf("Storage value should be %v after revert, got %v", value, storedValue)
	}

	// 测试回滚到无效索引
	sdb.RevertToSnapshot(999) // 无效索引，应该无操作
	balance = sdb.GetBalance(addr)
	if balance.Cmp(big.NewInt(1000)) != 0 {
		t.Errorf("Balance should still be 1000 after invalid revert, got %v", balance)
	}
}

// 测试日志操作
func TestLogOperations(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr := common.Address{0x01}

	// 测试初始日志
	logs := sdb.GetLogs()
	if len(logs) != 0 {
		t.Errorf("Initial logs should be empty")
	}

	// 创建并添加日志
	testLog := Log{
		Address:     addr,
		Topics:      []common.Hash{common.HexToHash("0x01"), common.HexToHash("0x02")},
		Data:        []byte("test data"),
		BlockNumber: 1,
		TxHash:      common.HexToHash("0x1234"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0x5678"),
		Index:       0,
	}

	sdb.AddLog(testLog)

	// 测试获取日志
	logs = sdb.GetLogs()
	if len(logs) != 1 {
		t.Errorf("Logs length should be 1, got %d", len(logs))
	}

	if logs[0].Address != testLog.Address {
		t.Errorf("Log address mismatch")
	}

	if len(logs[0].Topics) != len(testLog.Topics) {
		t.Errorf("Log topics length mismatch")
	}

	if string(logs[0].Data) != string(testLog.Data) {
		t.Errorf("Log data mismatch")
	}

	// 添加更多日志
	anotherLog := Log{
		Address:     addr,
		Topics:      []common.Hash{common.HexToHash("0x03")},
		Data:        []byte("more test data"),
		BlockNumber: 1,
		TxHash:      common.HexToHash("0x1234"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0x5678"),
		Index:       1,
	}

	sdb.AddLog(anotherLog)

	logs = sdb.GetLogs()
	if len(logs) != 2 {
		t.Errorf("Logs length should be 2, got %d", len(logs))
	}
}

// 测试预映像操作
func TestPreimageOperations(t *testing.T) {
	sdb := NewMemoryStateDB()
	hash := common.HexToHash("0x1234")
	preimage := []byte("test preimage")

	// 测试获取不存在的预映像
	retrievedPreimage := sdb.GetPreimage(hash)
	if retrievedPreimage != nil {
		t.Errorf("Preimage for non-existent hash should be nil")
	}

	// 添加预映像
	sdb.AddPreimage(hash, preimage)

	// 测试获取预映像
	retrievedPreimage = sdb.GetPreimage(hash)
	if string(retrievedPreimage) != string(preimage) {
		t.Errorf("Retrieved preimage should match original")
	}

	// 添加多个预映像
	anotherHash := common.HexToHash("0x5678")
	anotherPreimage := []byte("another test preimage")
	sdb.AddPreimage(anotherHash, anotherPreimage)

	retrievedPreimage = sdb.GetPreimage(anotherHash)
	if string(retrievedPreimage) != string(anotherPreimage) {
		t.Errorf("Retrieved preimage should match original")
	}

	// 确保第一个预映像仍然存在
	retrievedPreimage = sdb.GetPreimage(hash)
	if string(retrievedPreimage) != string(preimage) {
		t.Errorf("First preimage should still exist")
	}
}

// 测试其他辅助功能
func TestHelperFunctions(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr := common.Address{0x01}

	// 测试Empty函数（账户不存在）
	empty := sdb.Empty(addr)
	if !empty {
		t.Errorf("Non-existent account should be empty")
	}

	// 创建空账户
	sdb.CreateAccount(addr)
	empty = sdb.Empty(addr)
	if !empty {
		t.Errorf("Newly created account should be empty")
	}

	// 修改账户使其非空
	sdb.AddBalance(addr, big.NewInt(1000))
	empty = sdb.Empty(addr)
	if empty {
		t.Errorf("Account with balance should not be empty")
	}

	// 测试Suicide和HasSuicided函数
	suicided := sdb.Suicide(addr)
	if !suicided {
		t.Errorf("Suicide should return true")
	}

	hasSuicided := sdb.HasSuicided(addr)
	if hasSuicided {
		t.Errorf("HasSuicided should return false (not implemented)")
	}

	// 测试AddRefund和GetRefund函数
	sdb.AddRefund(1000)
	refund := sdb.GetRefund()
	if refund != 1000 {
		t.Errorf("Refund should be 1000, got %d", refund)
	}

	sdb.AddRefund(500)
	refund = sdb.GetRefund()
	if refund != 1500 {
		t.Errorf("Refund should be 1500, got %d", refund)
	}

	// 测试CalculateStateRoot函数
	root := sdb.CalculateStateRoot()
	if root == (common.Hash{}) {
		t.Errorf("State root should not be zero hash")
	}
}

// 集成测试：测试完整的状态操作流程
func TestStateIntegration(t *testing.T) {
	sdb := NewMemoryStateDB()
	addr1 := common.Address{0x01}
	addr2 := common.Address{0x02}

	// 1. 创建账户并设置初始状态
	sdb.CreateAccount(addr1)
	sdb.AddBalance(addr1, big.NewInt(10000))
	sdb.SetNonce(addr1, 0)

	// 2. 创建快照
	snapshotIdx := sdb.Snapshot()

	// 3. 执行一些操作
	// 转账（从addr1到addr2）
	sdb.SubBalance(addr1, big.NewInt(2000))
	sdb.AddBalance(addr2, big.NewInt(2000))

	// 增加nonce
	sdb.SetNonce(addr1, 1)

	// 设置存储
	key := common.HexToHash("0x01")
	value := common.HexToHash("0x1234")
	sdb.SetState(addr1, key, value)

	// 添加日志
	log := Log{
		Address:     addr1,
		Topics:      []common.Hash{common.HexToHash("0x01")},
		Data:        []byte("transfer"),
		BlockNumber: 1,
		TxHash:      common.HexToHash("0x1234"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0x5678"),
		Index:       0,
	}
	sdb.AddLog(log)

	// 4. 检查状态
	balance1 := sdb.GetBalance(addr1)
	if balance1.Cmp(big.NewInt(8000)) != 0 {
		t.Errorf("Balance of addr1 should be 8000, got %v", balance1)
	}

	balance2 := sdb.GetBalance(addr2)
	if balance2.Cmp(big.NewInt(2000)) != 0 {
		t.Errorf("Balance of addr2 should be 2000, got %v", balance2)
	}

	nonce := sdb.GetNonce(addr1)
	if nonce != 1 {
		t.Errorf("Nonce of addr1 should be 1, got %d", nonce)
	}

	storedValue := sdb.GetState(addr1, key)
	if storedValue != value {
		t.Errorf("Storage value should be %v, got %v", value, storedValue)
	}

	logs := sdb.GetLogs()
	if len(logs) != 1 {
		t.Errorf("Should have 1 log, got %d", len(logs))
	}

	// 5. 回滚到快照
	sdb.RevertToSnapshot(snapshotIdx)

	// 6. 检查回滚后的状态
	balance1 = sdb.GetBalance(addr1)
	if balance1.Cmp(big.NewInt(10000)) != 0 {
		t.Errorf("Balance of addr1 should be 10000 after revert, got %v", balance1)
	}

	balance2 = sdb.GetBalance(addr2)
	if balance2.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Balance of addr2 should be 0 after revert, got %v", balance2)
	}

	nonce = sdb.GetNonce(addr1)
	if nonce != 0 {
		t.Errorf("Nonce of addr1 should be 0 after revert, got %d", nonce)
	}

	storedValue = sdb.GetState(addr1, key)
	if storedValue != (common.Hash{}) {
		t.Errorf("Storage value should be zero hash after revert, got %v", storedValue)
	}

	logs = sdb.GetLogs()
	if len(logs) != 1 {
		t.Errorf("Logs should still have 1 entry (not reverted), got %d", len(logs))
	}
}
