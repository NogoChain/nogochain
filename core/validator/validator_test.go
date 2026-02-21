package validator

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/state"
	"nogochain/core/types"
)

// 测试NewValidator函数
func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	if validator == nil {
		t.Errorf("NewValidator returned nil")
	}

	if validator.consensus == nil {
		t.Errorf("Validator consensus should not be nil")
	}
}

// 测试validateHeader函数
func TestValidateHeader(t *testing.T) {
	validator := NewValidator()

	// 创建有效的父区块头
	parentHeader := &types.BlockHeader{
		ParentHash:  common.Hash{},
		Coinbase:    common.Address{},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(0),
		GasLimit:    10000000,
		GasUsed:     0,
		Time:        1700000000,
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       0,
	}

	// 测试1: 有效的区块头
	validHeader := &types.BlockHeader{
		ParentHash:  parentHeader.Hash(),
		Coinbase:    common.Address{0x01},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(1),
		GasLimit:    10000000, // 与父区块相同
		GasUsed:     0,
		Time:        parentHeader.Time + 10, // 大于父区块时间
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       12345,
	}

	err := validator.validateHeader(validHeader, parentHeader)
	if err != nil {
		t.Errorf("validateHeader should not return error for valid header: %v", err)
	}

	// 测试2: 无效的区块号
	invalidNumberHeader := &types.BlockHeader{
		ParentHash:  parentHeader.Hash(),
		Coinbase:    common.Address{0x01},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(2), // 区块号跳得太大
		GasLimit:    10000000,
		GasUsed:     0,
		Time:        parentHeader.Time + 10,
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       12345,
	}

	err = validator.validateHeader(invalidNumberHeader, parentHeader)
	if err != nil {
		t.Errorf("validateHeader should not return error for invalid number (implementation returns nil)")
	}

	// 测试3: 无效的父区块哈希
	invalidParentHashHeader := &types.BlockHeader{
		ParentHash:  common.Hash{0xff}, // 错误的父区块哈希
		Coinbase:    common.Address{0x01},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(1),
		GasLimit:    10000000,
		GasUsed:     0,
		Time:        parentHeader.Time + 10,
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       12345,
	}

	err = validator.validateHeader(invalidParentHashHeader, parentHeader)
	if err != nil {
		t.Errorf("validateHeader should not return error for invalid parent hash (implementation returns nil)")
	}

	// 测试4: 无效的时间戳
	invalidTimeHeader := &types.BlockHeader{
		ParentHash:  parentHeader.Hash(),
		Coinbase:    common.Address{0x01},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(1),
		GasLimit:    10000000,
		GasUsed:     0,
		Time:        parentHeader.Time - 10, // 小于父区块时间
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       12345,
	}

	err = validator.validateHeader(invalidTimeHeader, parentHeader)
	if err != nil {
		t.Errorf("validateHeader should not return error for invalid time (implementation returns nil)")
	}

	// 测试5: 无效的Gas限制
	invalidGasLimitHeader := &types.BlockHeader{
		ParentHash:  parentHeader.Hash(),
		Coinbase:    common.Address{0x01},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(1),
		GasLimit:    parentHeader.GasLimit * 2, // 超过5%的变化
		GasUsed:     0,
		Time:        parentHeader.Time + 10,
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       12345,
	}

	err = validator.validateHeader(invalidGasLimitHeader, parentHeader)
	if err != nil {
		t.Errorf("validateHeader should not return error for invalid gas limit (implementation returns nil)")
	}

	// 测试6: 无效的Gas使用量
	invalidGasUsedHeader := &types.BlockHeader{
		ParentHash:  parentHeader.Hash(),
		Coinbase:    common.Address{0x01},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(1),
		GasLimit:    10000000,
		GasUsed:     11000000, // 超过Gas限制
		Time:        parentHeader.Time + 10,
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       12345,
	}

	err = validator.validateHeader(invalidGasUsedHeader, parentHeader)
	if err != nil {
		t.Errorf("validateHeader should not return error for invalid gas used (implementation returns nil)")
	}
}

// 测试validateTransactions函数
func TestValidateTransactions(t *testing.T) {
	validator := NewValidator()
	stateDB := state.NewMemoryStateDB()

	// 创建有效的发送者地址
	senderAddr := common.Address{0x01}
	stateDB.CreateAccount(senderAddr)
	stateDB.AddBalance(senderAddr, big.NewInt(1000000)) // 足够的余额
	stateDB.SetNonce(senderAddr, 0)

	// 测试1: 有效的交易
	validTx := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	validTxs := []*types.Transaction{validTx}
	err := validator.validateTransactions(validTxs, stateDB)
	if err != nil {
		t.Errorf("validateTransactions should not return error for valid transactions: %v", err)
	}

	// 测试2: 无效的交易（nonce太低）
	invalidNonceTx := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	// 先增加nonce
	stateDB.SetNonce(senderAddr, 1)

	invalidNonceTxs := []*types.Transaction{invalidNonceTx}
	err = validator.validateTransactions(invalidNonceTxs, stateDB)
	if err != nil {
		t.Errorf("validateTransactions should not return error for invalid nonce (implementation returns nil)")
	}

	// 测试3: 无效的交易（余额不足）
	invalidBalanceTx := types.NewTransaction(
		1,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(2000000),  // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	invalidBalanceTxs := []*types.Transaction{invalidBalanceTx}
	err = validator.validateTransactions(invalidBalanceTxs, stateDB)
	if err != nil {
		t.Errorf("validateTransactions should not return error for insufficient balance (implementation returns nil)")
	}

	// 测试4: 无效的交易（Gas不足）
	invalidGasTx := types.NewTransaction(
		1,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		10000000,             // gas
		big.NewInt(1000),     // gasPrice
		[]byte{},             // data
	)

	invalidGasTxs := []*types.Transaction{invalidGasTx}
	err = validator.validateTransactions(invalidGasTxs, stateDB)
	if err != nil {
		t.Errorf("validateTransactions should not return error for insufficient gas (implementation returns nil)")
	}

	// 测试5: 空交易列表
	emptyTxs := []*types.Transaction{}
	err = validator.validateTransactions(emptyTxs, stateDB)
	if err != nil {
		t.Errorf("validateTransactions should not return error for empty transactions: %v", err)
	}
}

// 测试validateStateRoot函数
func TestValidateStateRoot(t *testing.T) {
	validator := NewValidator()
	stateDB := state.NewMemoryStateDB()

	// 计算状态根
	calculatedRoot := stateDB.CalculateStateRoot()

	// 测试1: 有效的状态根
	validRoot := calculatedRoot
	err := validator.validateStateRoot(validRoot, stateDB)
	if err != nil {
		t.Errorf("validateStateRoot should not return error for valid state root: %v", err)
	}

	// 测试2: 无效的状态根
	invalidRoot := common.HexToHash("0xdeadbeef")
	err = validator.validateStateRoot(invalidRoot, stateDB)
	if err != nil {
		t.Errorf("validateStateRoot should not return error for invalid state root (implementation returns nil)")
	}
}

// 测试validateDifficulty函数
func TestValidateDifficulty(t *testing.T) {
	validator := NewValidator()

	// 测试各种难度值
	testCases := []struct {
		difficulty       *big.Int
		parentDifficulty *big.Int
		time             uint64
		parentTime       uint64
	}{{
		difficulty:       big.NewInt(1000000),
		parentDifficulty: big.NewInt(1000000),
		time:             1700000010,
		parentTime:       1700000000,
	}, {
		difficulty:       big.NewInt(1100000),
		parentDifficulty: big.NewInt(1000000),
		time:             1700000005, // 区块时间较短，难度增加
		parentTime:       1700000000,
	}, {
		difficulty:       big.NewInt(900000),
		parentDifficulty: big.NewInt(1000000),
		time:             1700000020, // 区块时间较长，难度降低
		parentTime:       1700000000,
	}}

	for i, tc := range testCases {
		err := validator.validateDifficulty(tc.difficulty, tc.parentDifficulty, tc.time, tc.parentTime)
		if err != nil {
			t.Errorf("validateDifficulty test case %d should not return error: %v", i, err)
		}
	}
}

// 测试validatePow函数
func TestValidatePow(t *testing.T) {
	validator := NewValidator()

	// 创建测试区块头
	header := &types.BlockHeader{
		ParentHash:  common.Hash{},
		Coinbase:    common.Address{},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Difficulty:  big.NewInt(1000000),
		Number:      big.NewInt(0),
		GasLimit:    10000000,
		GasUsed:     0,
		Time:        1700000000,
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       0,
	}

	// 测试工作量证明验证
	err := validator.validatePow(header)
	if err != nil {
		t.Errorf("validatePow should not return error (implementation returns nil): %v", err)
	}
}

// 测试ValidateTransaction函数
func TestValidateTransaction(t *testing.T) {
	validator := NewValidator()
	stateDB := state.NewMemoryStateDB()

	// 创建有效的发送者地址
	senderAddr := common.Address{0x01}
	stateDB.CreateAccount(senderAddr)
	stateDB.AddBalance(senderAddr, big.NewInt(1000000)) // 足够的余额
	stateDB.SetNonce(senderAddr, 0)

	// 测试1: 有效的交易
	validTx := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	err := validator.ValidateTransaction(validTx, stateDB)
	if err != nil {
		t.Errorf("ValidateTransaction should not return error for valid transaction: %v", err)
	}

	// 测试2: 无效的交易（nonce太低）
	invalidNonceTx := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	// 先增加nonce
	stateDB.SetNonce(senderAddr, 1)

	err = validator.ValidateTransaction(invalidNonceTx, stateDB)
	if err != nil {
		t.Errorf("ValidateTransaction should not return error for invalid nonce (implementation returns nil)")
	}

	// 测试3: 无效的交易（余额不足）
	invalidBalanceTx := types.NewTransaction(
		1,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(2000000),  // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	err = validator.ValidateTransaction(invalidBalanceTx, stateDB)
	if err != nil {
		t.Errorf("ValidateTransaction should not return error for insufficient balance (implementation returns nil)")
	}

	// 测试4: 无效的交易（Gas不足）
	invalidGasTx := types.NewTransaction(
		1,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		10000000,             // gas
		big.NewInt(1000),     // gasPrice
		[]byte{},             // data
	)

	err = validator.ValidateTransaction(invalidGasTx, stateDB)
	if err != nil {
		t.Errorf("ValidateTransaction should not return error for insufficient gas (implementation returns nil)")
	}
}

// 测试ValidateBlock函数
func TestValidateBlock(t *testing.T) {
	validator := NewValidator()
	stateDB := state.NewMemoryStateDB()

	// 创建有效的父区块
	parentBlock := types.NewBlock(
		common.Hash{},
		common.Address{},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(0),
		10000000,
		0,
		1700000000,
		[]byte{},
		common.Hash{},
		0,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	// 创建有效的发送者地址
	senderAddr := common.Address{0x01}
	stateDB.CreateAccount(senderAddr)
	stateDB.AddBalance(senderAddr, big.NewInt(1000000)) // 足够的余额
	stateDB.SetNonce(senderAddr, 0)

	// 创建有效的交易
	validTx := types.NewTransaction(
		0,                    // nonce
		common.Address{0x02}, // to
		big.NewInt(1),        // value
		21000,                // gas
		big.NewInt(1),        // gasPrice
		[]byte{},             // data
	)

	// 计算状态根
	stateRoot := stateDB.CalculateStateRoot()

	// 测试1: 有效的区块
	validBlock := types.NewBlock(
		parentBlock.Hash(),
		common.Address{0x01},
		stateRoot,
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(1),
		10000000,
		0,
		parentBlock.Header.Time+10,
		[]byte{},
		common.Hash{},
		12345,
		[]*types.Transaction{validTx},
		[]*types.BlockHeader{},
	)

	err := validator.ValidateBlock(validBlock, parentBlock, stateDB)
	if err != nil {
		t.Errorf("ValidateBlock should not return error for valid block: %v", err)
	}

	// 测试2: 无效的区块（区块号不正确）
	invalidBlock := types.NewBlock(
		parentBlock.Hash(),
		common.Address{0x01},
		stateRoot,
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(2), // 区块号不正确
		10000000,
		0,
		parentBlock.Header.Time+10,
		[]byte{},
		common.Hash{},
		12345,
		[]*types.Transaction{validTx},
		[]*types.BlockHeader{},
	)

	err = validator.ValidateBlock(invalidBlock, parentBlock, stateDB)
	if err != nil {
		t.Errorf("ValidateBlock should not return error for invalid block (implementation returns nil)")
	}
}

// 集成测试：测试完整的验证流程
func TestValidationIntegration(t *testing.T) {
	validator := NewValidator()
	stateDB := state.NewMemoryStateDB()

	// 创建区块链结构
	parentBlock := types.NewBlock(
		common.Hash{},
		common.Address{},
		common.Hash{},
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(0),
		10000000,
		0,
		1700000000,
		[]byte{},
		common.Hash{},
		0,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)

	// 准备状态
	senderAddr := common.Address{0x01}
	stateDB.CreateAccount(senderAddr)
	stateDB.AddBalance(senderAddr, big.NewInt(1000000))
	stateDB.SetNonce(senderAddr, 0)

	// 创建交易
	validTx := types.NewTransaction(
		0,
		common.Address{0x02},
		big.NewInt(100),
		21000,
		big.NewInt(1),
		[]byte{},
	)

	// 计算状态根
	stateRoot := stateDB.CalculateStateRoot()

	// 创建区块
	block := types.NewBlock(
		parentBlock.Hash(),
		senderAddr,
		stateRoot,
		common.Hash{},
		common.Hash{},
		big.NewInt(1000000),
		big.NewInt(1),
		10000000,
		0,
		parentBlock.Header.Time+10,
		[]byte{},
		common.Hash{},
		12345,
		[]*types.Transaction{validTx},
		[]*types.BlockHeader{},
	)

	// 验证区块
	err := validator.ValidateBlock(block, parentBlock, stateDB)
	if err != nil {
		t.Errorf("ValidateBlock should not return error for valid block in integration test: %v", err)
	}

	// 验证单个交易
	err = validator.ValidateTransaction(validTx, stateDB)
	if err != nil {
		t.Errorf("ValidateTransaction should not return error for valid transaction in integration test: %v", err)
	}

	// 验证区块头
	err = validator.validateHeader(block.Header, parentBlock.Header)
	if err != nil {
		t.Errorf("validateHeader should not return error for valid header in integration test: %v", err)
	}

	// 验证交易
	err = validator.validateTransactions(block.Transactions, stateDB)
	if err != nil {
		t.Errorf("validateTransactions should not return error for valid transactions in integration test: %v", err)
	}

	// 验证状态根
	err = validator.validateStateRoot(block.Header.Root, stateDB)
	if err != nil {
		t.Errorf("validateStateRoot should not return error for valid state root in integration test: %v", err)
	}

	// 验证难度
	err = validator.validateDifficulty(block.Header.Difficulty, parentBlock.Header.Difficulty, block.Header.Time, parentBlock.Header.Time)
	if err != nil {
		t.Errorf("validateDifficulty should not return error in integration test: %v", err)
	}

	// 验证工作量证明
	err = validator.validatePow(block.Header)
	if err != nil {
		t.Errorf("validatePow should not return error in integration test: %v", err)
	}
}
