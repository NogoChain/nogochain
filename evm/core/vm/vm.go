package vm

import (
	"fmt"
	"math/big"

	gasMeter "nogochain/evm/core/vm/gas"
	"nogochain/evm/core/vm/memory"
	"nogochain/evm/core/vm/stack"
	"nogochain/evm/core/vm/storage"
	"nogochain/evm/params"
)

// EVM 以太坊虚拟机实现
type EVM struct {
	// 执行环境
	Context     Context
	StateDB     StateDB
	BlockHeader *BlockHeader

	// 虚拟机组件
	Stack    *stack.Stack
	Memory   *memory.Memory
	Storage  *storage.Storage
	GasMeter *gasMeter.GasMeter

	// 执行状态
	ProgramCounter int
	ReturnData     []byte
	Running        bool
	Err            error
}

// Context 执行上下文
type Context struct {
	Caller      []byte
	GasPrice    *big.Int
	Origin      []byte
	BlockNumber *big.Int
	Timestamp   *big.Int
	GasLimit    uint64
	BaseFee     *big.Int
	Code        []byte
}

// StateDB 状态数据库接口
type StateDB interface {
	GetBalance(addr []byte) *big.Int
	GetCode(addr []byte) []byte
	GetNonce(addr []byte) uint64
	SetNonce(addr []byte, nonce uint64)
	GetState(addr, key []byte) []byte
	SetState(addr, key, value []byte)
	SetCode(addr []byte, code []byte)
	AddBalance(addr []byte, amount *big.Int)
	SubBalance(addr []byte, amount *big.Int)
	CreateAccount(addr []byte)
	Exist(addr []byte) bool
}

// BlockHeader 区块头
type BlockHeader struct {
	Coinbase   []byte
	GasLimit   uint64
	Number     *big.Int
	Timestamp  *big.Int
	BaseFee    *big.Int
	Difficulty *big.Int
}

// Transaction 交易结构体
type Transaction struct {
	To        []byte
	From      []byte
	Nonce     uint64
	GasPrice  *big.Int
	GasFeeCap *big.Int // EIP-1559
	GasTipCap *big.Int // EIP-1559
	GasLimit  uint64
	Value     *big.Int
	Data      []byte
	V         *big.Int
	R         *big.Int
	S         *big.Int
}

// IsEIP1559Transaction 检查是否为EIP-1559交易
func (tx *Transaction) IsEIP1559Transaction() bool {
	return tx.GasFeeCap != nil && tx.GasTipCap != nil
}

// EffectiveGasPrice 计算有效Gas价格
func (tx *Transaction) EffectiveGasPrice(baseFee *big.Int) *big.Int {
	if tx.IsEIP1559Transaction() {
		// EIP-1559: max(tip, feeCap - baseFee)
		feeCapMinusBase := new(big.Int).Sub(tx.GasFeeCap, baseFee)
		if feeCapMinusBase.Cmp(tx.GasTipCap) < 0 {
			return tx.GasTipCap
		}
		return feeCapMinusBase
	}
	// 传统交易：使用GasPrice
	return tx.GasPrice
}

// NewEVM 创建新的EVM实例
func NewEVM(context Context, stateDB StateDB, header *BlockHeader) *EVM {
	return &EVM{
		Context:        context,
		StateDB:        stateDB,
		BlockHeader:    header,
		Stack:          stack.NewStack(),
		Memory:         memory.NewMemory(),
		Storage:        storage.NewStorage(),
		GasMeter:       gasMeter.NewGasMeter(context.GasLimit),
		ProgramCounter: 0,
		Running:        true,
	}
}

// NewGasMeter 创建新的Gas计量器
func NewGasMeter(limit uint64) *gasMeter.GasMeter {
	return gasMeter.NewGasMeter(limit)
}

// IsHardForkActive 检查硬分叉是否激活
func (evm *EVM) IsHardForkActive(forkName string) bool {
	blockNumber := evm.BlockHeader.Number

	switch forkName {
	case "homestead":
		return blockNumber.Cmp(params.HomesteadBlock) >= 0
	case "tangerineWhistle":
		return blockNumber.Cmp(params.TangerineWhistleBlock) >= 0
	case "spuriousDragon":
		return blockNumber.Cmp(params.SpuriousDragonBlock) >= 0
	case "byzantium":
		return blockNumber.Cmp(params.ByzantiumBlock) >= 0
	case "constantinople":
		return blockNumber.Cmp(params.ConstantinopleBlock) >= 0
	case "petersburg":
		return blockNumber.Cmp(params.PetersburgBlock) >= 0
	case "istanbul":
		return blockNumber.Cmp(params.IstanbulBlock) >= 0
	case "berlin":
		return blockNumber.Cmp(params.BerlinBlock) >= 0
	case "london":
		return blockNumber.Cmp(params.LondonBlock) >= 0
	case "eip1559":
		return blockNumber.Cmp(params.EIP1559Block) >= 0
	default:
		return false
	}
}

// CalculateBaseFee 计算基础费（EIP-1559）
func (evm *EVM) CalculateBaseFee() *big.Int {
	// 检查EIP-1559是否激活
	if !evm.IsHardForkActive("eip1559") {
		return big.NewInt(0)
	}

	// 这里实现EIP-1559的基础费计算逻辑
	// 简化实现：使用区块头中的基础费
	return evm.BlockHeader.BaseFee
}

// ApplyHardForkRules 应用硬分叉规则
func (evm *EVM) ApplyHardForkRules() {
	// 根据当前硬分叉状态应用相应的规则
	if evm.IsHardForkActive("berlin") {
		// Berlin硬分叉规则：调整Gas成本等
	}

	if evm.IsHardForkActive("london") {
		// London硬分叉规则：EIP-1559等
	}

	// 其他硬分叉规则...
}

// ValidateTransaction 验证交易是否符合当前硬分叉规则
func (evm *EVM) ValidateTransaction(tx *Transaction) error {
	// 检查交易是否符合当前硬分叉的规则
	if evm.IsHardForkActive("eip1559") {
		// EIP-1559交易验证
		if tx.IsEIP1559Transaction() {
			// 检查EIP-1559交易格式
			if tx.GasFeeCap == nil || tx.GasTipCap == nil {
				return fmt.Errorf("EIP-1559 transaction missing fee cap or tip cap")
			}
			// 检查GasFeeCap >= GasTipCap
			if tx.GasFeeCap.Cmp(tx.GasTipCap) < 0 {
				return fmt.Errorf("gas fee cap must be greater than or equal to gas tip cap")
			}
		}
	}

	// 检查Gas限制
	if tx.GasLimit == 0 {
		return fmt.Errorf("gas limit cannot be zero")
	}

	return nil
}

// Run 运行EVM
func (evm *EVM) Run(code []byte) ([]byte, error) {
	// 设置当前执行的代码
	evm.Context.Code = code

	for evm.Running {
		if evm.ProgramCounter >= len(code) {
			evm.Stop()
			break
		}

		opcode := code[evm.ProgramCounter]

		// 执行指令
		err := evm.executeOpcode(opcode)
		if err != nil {
			evm.Err = err
			evm.Stop()
			break
		}
	}

	return evm.ReturnData, evm.Err
}

// Create 创建智能合约
func (evm *EVM) Create(caller []byte, code []byte, value *big.Int, gas uint64) ([]byte, []byte, error) {
	// 生成合约地址（使用CREATE算法）
	address := generateCreateAddress(caller, evm.StateDB.GetNonce(caller))

	// 检查地址是否已存在
	if evm.StateDB.Exist(address) {
		return nil, nil, fmt.Errorf("contract address already exists")
	}

	// 创建新账户
	evm.StateDB.CreateAccount(address)

	// 转移价值
	if value.Sign() > 0 {
		evm.StateDB.SubBalance(caller, value)
		evm.StateDB.AddBalance(address, value)
	}

	// 增加调用者的nonce
	evm.StateDB.SetNonce(caller, evm.StateDB.GetNonce(caller)+1)

	// 创建子EVM执行环境
	childEVM := evm.createChildEVM(address, code, value, gas)

	// 执行初始化代码
	returnData, err := childEVM.Run(code)

	// 如果执行成功且有返回数据，部署合约代码
	if err == nil && len(returnData) > 0 {
		// 检查代码大小限制
		if len(returnData) > 24576 {
			return nil, nil, fmt.Errorf("contract code size exceeds limit")
		}
		evm.StateDB.SetCode(address, returnData)
	}

	return address, returnData, err
}

// Create2 创建智能合约（使用CREATE2算法）
func (evm *EVM) Create2(caller []byte, code []byte, salt []byte, value *big.Int, gas uint64) ([]byte, []byte, error) {
	// 生成合约地址（使用CREATE2算法）
	address := generateCreate2Address(caller, salt, code)

	// 检查地址是否已存在
	if evm.StateDB.Exist(address) {
		return nil, nil, fmt.Errorf("contract address already exists")
	}

	// 创建新账户
	evm.StateDB.CreateAccount(address)

	// 转移价值
	if value.Sign() > 0 {
		evm.StateDB.SubBalance(caller, value)
		evm.StateDB.AddBalance(address, value)
	}

	// 创建子EVM执行环境
	childEVM := evm.createChildEVM(address, code, value, gas)

	// 执行初始化代码
	returnData, err := childEVM.Run(code)

	// 如果执行成功且有返回数据，部署合约代码
	if err == nil && len(returnData) > 0 {
		// 检查代码大小限制
		if len(returnData) > 24576 {
			return nil, nil, fmt.Errorf("contract code size exceeds limit")
		}
		evm.StateDB.SetCode(address, returnData)
	}

	return address, returnData, err
}

// Call 调用智能合约
func (evm *EVM) Call(caller []byte, to []byte, input []byte, value *big.Int, gas uint64) ([]byte, error) {
	// 检查目标地址是否存在
	if !evm.StateDB.Exist(to) {
		return nil, fmt.Errorf("contract address does not exist")
	}

	// 转移价值
	if value.Sign() > 0 {
		evm.StateDB.SubBalance(caller, value)
		evm.StateDB.AddBalance(to, value)
	}

	// 获取合约代码
	code := evm.StateDB.GetCode(to)
	if len(code) == 0 {
		return nil, fmt.Errorf("contract has no code")
	}

	// 创建子EVM执行环境
	childEVM := evm.createChildEVM(to, code, value, gas)

	// 执行合约代码
	return childEVM.Run(code)
}

// StaticCall 静态调用智能合约
func (evm *EVM) StaticCall(caller []byte, to []byte, input []byte, gas uint64) ([]byte, error) {
	// 检查目标地址是否存在
	if !evm.StateDB.Exist(to) {
		return nil, fmt.Errorf("contract address does not exist")
	}

	// 获取合约代码
	code := evm.StateDB.GetCode(to)
	if len(code) == 0 {
		return nil, fmt.Errorf("contract has no code")
	}

	// 创建静态子EVM执行环境
	childEVM := evm.createChildEVM(to, code, big.NewInt(0), gas)

	// 执行合约代码（禁止状态修改）
	return childEVM.Run(code)
}

// DelegateCall 委托调用智能合约
func (evm *EVM) DelegateCall(caller []byte, to []byte, input []byte, gas uint64) ([]byte, error) {
	// 检查目标地址是否存在
	if !evm.StateDB.Exist(to) {
		return nil, fmt.Errorf("contract address does not exist")
	}

	// 获取合约代码
	code := evm.StateDB.GetCode(to)
	if len(code) == 0 {
		return nil, fmt.Errorf("contract has no code")
	}

	// 创建委托调用EVM执行环境（使用调用者的上下文）
	childEVM := evm.createChildEVM(caller, code, big.NewInt(0), gas)

	// 执行合约代码
	return childEVM.Run(code)
}

// CallCode 代码调用智能合约
func (evm *EVM) CallCode(caller []byte, to []byte, input []byte, value *big.Int, gas uint64) ([]byte, error) {
	// 检查目标地址是否存在
	if !evm.StateDB.Exist(to) {
		return nil, fmt.Errorf("contract address does not exist")
	}

	// 转移价值
	if value.Sign() > 0 {
		evm.StateDB.SubBalance(caller, value)
		evm.StateDB.AddBalance(caller, value) // 注意：CallCode使用调用者的地址
	}

	// 获取合约代码
	code := evm.StateDB.GetCode(to)
	if len(code) == 0 {
		return nil, fmt.Errorf("contract has no code")
	}

	// 创建代码调用EVM执行环境
	childEVM := evm.createChildEVM(caller, code, value, gas)

	// 执行合约代码
	return childEVM.Run(code)
}

// 创建子EVM执行环境
func (evm *EVM) createChildEVM(address []byte, code []byte, value *big.Int, gasLimit uint64) *EVM {
	// 创建新的执行上下文
	context := evm.Context
	context.Caller = address
	context.Code = code

	// 创建新的EVM实例
	childEVM := NewEVM(context, evm.StateDB, evm.BlockHeader)
	childEVM.GasMeter = NewGasMeter(gasLimit)

	return childEVM
}

// 生成CREATE地址
func generateCreateAddress(caller []byte, nonce uint64) []byte {
	// 简单实现：实际应该使用keccak256(rlp([caller, nonce]))
	// 这里返回一个简化的地址
	address := make([]byte, 20)
	copy(address, caller)
	// 在实际实现中，这里应该使用正确的地址生成算法
	return address
}

// 生成CREATE2地址
func generateCreate2Address(caller []byte, salt []byte, initCode []byte) []byte {
	// 简单实现：实际应该使用keccak256(0xff + caller + salt + keccak256(initCode))
	// 这里返回一个简化的地址
	address := make([]byte, 20)
	copy(address, caller)
	// 在实际实现中，这里应该使用正确的地址生成算法
	return address
}

// executeOpcode 执行操作码
func (evm *EVM) executeOpcode(opcode byte) error {
	// 获取指令
	instr := GetInstruction(opcode)
	if instr == nil {
		// 处理未知指令
		return fmt.Errorf("unknown opcode: 0x%02x", opcode)
	}

	// 执行指令
	err := instr.Execute(evm)

	// 对于大多数指令，增加ProgramCounter
	// 对于PUSH指令，需要根据数据大小增加更多
	if opcode >= 0x60 && opcode <= 0x7f {
		// PUSH1-PUSH32
		dataSize := int(opcode - 0x60 + 1)
		evm.ProgramCounter += 1 + dataSize
	} else {
		evm.ProgramCounter++
	}

	return err
}

// Stop 停止EVM执行
func (evm *EVM) Stop() {
	evm.Running = false
}

// SetReturnData 设置返回数据
func (evm *EVM) SetReturnData(data []byte) {
	evm.ReturnData = data
}

// GetGasLeft 获取剩余Gas
func (evm *EVM) GetGasLeft() uint64 {
	return evm.GasMeter.GetGasLeft()
}

// ConsumeGas 消耗Gas
func (evm *EVM) ConsumeGas(amount uint64) error {
	return evm.GasMeter.ConsumeGas(amount)
}

// Revert 回滚状态
func (evm *EVM) Revert() {
	// 实现状态回滚逻辑
}

// Snapshot 创建状态快照
func (evm *EVM) Snapshot() int {
	// 实现快照创建逻辑
	return 0
}

// RevertToSnapshot 回滚到指定快照
func (evm *EVM) RevertToSnapshot(id int) {
	// 实现快照回滚逻辑
}
