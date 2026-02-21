package gas

import (
	"errors"
	"math/big"
)

// GasMeter 以太坊虚拟机Gas计量器实现
type GasMeter struct {
	GasLimit  uint64
	GasUsed   uint64
	GasRefund uint64
}

// 错误定义
var (
	ErrOutOfGas = errors.New("out of gas")
)

// NewGasMeter 创建新的Gas计量器
func NewGasMeter(limit uint64) *GasMeter {
	return &GasMeter{
		GasLimit:  limit,
		GasUsed:   0,
		GasRefund: 0,
	}
}

// GetGasLeft 获取剩余Gas
func (g *GasMeter) GetGasLeft() uint64 {
	return g.GasLimit - g.GasUsed
}

// GetGasUsed 获取已使用的Gas
func (g *GasMeter) GetGasUsed() uint64 {
	return g.GasUsed
}

// GetGasRefund 获取Gas退款
func (g *GasMeter) GetGasRefund() uint64 {
	return g.GasRefund
}

// ConsumeGas 消耗Gas
func (g *GasMeter) ConsumeGas(amount uint64) error {
	if g.GasUsed+amount > g.GasLimit {
		return ErrOutOfGas
	}
	g.GasUsed += amount
	return nil
}

// RefundGas 退款Gas
func (g *GasMeter) RefundGas(amount uint64) {
	g.GasRefund += amount
}

// SetGasLimit 设置Gas限制
func (g *GasMeter) SetGasLimit(limit uint64) {
	g.GasLimit = limit
}

// Reset 重置Gas计量器
func (g *GasMeter) Reset() {
	g.GasUsed = 0
	g.GasRefund = 0
}

// ApplyRefund 应用Gas退款
func (g *GasMeter) ApplyRefund() {
	// 退款不能超过已使用Gas的一半
	maxRefund := g.GasUsed / 2
	if g.GasRefund > maxRefund {
		g.GasRefund = maxRefund
	}
}

// CalculateBaseGas 计算基础Gas成本
func CalculateBaseGas(opcode byte) uint64 {
	// 基础Gas成本表 - 与以太坊完全一致
	switch opcode {
	case 0x00: // STOP
		return 0
	case 0x01: // ADD
		return 3
	case 0x02: // MUL
		return 5
	case 0x03: // SUB
		return 3
	case 0x04: // DIV
		return 5
	case 0x05: // SDIV
		return 5
	case 0x06: // MOD
		return 5
	case 0x07: // SMOD
		return 5
	case 0x08: // ADDMOD
		return 8
	case 0x09: // MULMOD
		return 8
	case 0x0a: // EXP
		return 10
	case 0x0b: // SIGNEXTEND
		return 5
	case 0x10: // LT
		return 3
	case 0x11: // GT
		return 3
	case 0x12: // SLT
		return 3
	case 0x13: // SGT
		return 3
	case 0x14: // EQ
		return 3
	case 0x15: // ISZERO
		return 3
	case 0x16: // AND
		return 3
	case 0x17: // OR
		return 3
	case 0x18: // XOR
		return 3
	case 0x19: // NOT
		return 3
	case 0x1a: // BYTE
		return 3
	case 0x20: // SHA3
		return 30
	case 0x30: // ADDRESS
		return 2
	case 0x31: // BALANCE
		return 400
	case 0x32: // ORIGIN
		return 2
	case 0x33: // CALLER
		return 2
	case 0x34: // CALLVALUE
		return 2
	case 0x35: // CALLDATALOAD
		return 3
	case 0x36: // CALLDATASIZE
		return 2
	case 0x37: // CALLDATACOPY
		return 3
	case 0x38: // CODESIZE
		return 2
	case 0x39: // CODECOPY
		return 3
	case 0x3a: // GASPRICE
		return 2
	case 0x3b: // EXTCODESIZE
		return 700
	case 0x3c: // EXTCODECOPY
		return 700
	case 0x3d: // RETURNDATASIZE
		return 2
	case 0x3e: // RETURNDATACOPY
		return 3
	case 0x3f: // EXTCODEHASH
		return 700
	case 0x40: // BLOCKHASH
		return 20
	case 0x41: // COINBASE
		return 2
	case 0x42: // TIMESTAMP
		return 2
	case 0x43: // NUMBER
		return 2
	case 0x44: // DIFFICULTY
		return 2
	case 0x45: // GASLIMIT
		return 2
	case 0x46: // CHAINID
		return 2
	case 0x47: // SELFBALANCE
		return 5
	case 0x48: // BASEFEE
		return 2
	case 0x50: // POP
		return 2
	case 0x51: // MLOAD
		return 3
	case 0x52: // MSTORE
		return 3
	case 0x53: // MSTORE8
		return 3
	case 0x54: // SLOAD
		return 800
	case 0x55: // SSTORE
		return 20000
	case 0x56: // JUMP
		return 8
	case 0x57: // JUMPI
		return 10
	case 0x58: // PC
		return 2
	case 0x59: // MSIZE
		return 2
	case 0x5a: // GAS
		return 2
	case 0x5b: // JUMPDEST
		return 1
	case 0xa0: // LOG0
		return 375
	case 0xa1: // LOG1
		return 750
	case 0xa2: // LOG2
		return 1125
	case 0xa3: // LOG3
		return 1500
	case 0xa4: // LOG4
		return 1875
	case 0xf0: // CREATE
		return 32000
	case 0xf1: // CALL
		return 700
	case 0xf2: // CALLCODE
		return 700
	case 0xf3: // RETURN
		return 0
	case 0xf4: // DELEGATECALL
		return 700
	case 0xf5: // CREATE2
		return 32000
	case 0xf6: // STATICCALL
		return 700
	case 0xf7: // REVERT
		return 0
	case 0xf8: // INVALID
		return 0
	case 0xf9: // SELFDESTRUCT
		return 5000
	default:
		// PUSH*, DUP*, SWAP* 等指令
		if opcode >= 0x60 && opcode <= 0x7f {
			return 3
		} else if opcode >= 0x80 && opcode <= 0x8f {
			return 3
		} else if opcode >= 0x90 && opcode <= 0x9f {
			return 3
		}
		return 0
	}
}

// CalculateSstoreGas 计算SSTORE操作的Gas成本
func CalculateSstoreGas(current, new []byte) uint64 {
	// 检查是否为零值
	isCurrentZero := isZeroValue(current)
	isNewZero := isZeroValue(new)

	switch {
	case isCurrentZero && !isNewZero:
		// 从0设置为非0
		return 20000
	case !isCurrentZero && isNewZero:
		// 从非0设置为0（退款15000）
		return 5000
	case !isCurrentZero && !isNewZero:
		// 从非0设置为另一个非0
		return 5000
	default:
		// 从0设置为0
		return 20000
	}
}

// CalculateCallGasWithValue 计算带值的调用操作的Gas成本
func CalculateCallGasWithValue(gasLimit, gasPrice uint64, value *big.Int, isNewAccount bool) uint64 {
	baseGas := uint64(700)

	// 如果有值转移
	if value != nil && value.Sign() > 0 {
		baseGas += 9000
	}

	// 如果创建新账户
	if isNewAccount {
		baseGas += 25000
	}

	// 添加调用的Gas限制
	if gasLimit > 0 {
		return baseGas + gasLimit
	}

	return baseGas
}

// CalculateCreateGas 计算创建操作的Gas成本
func CalculateCreateGas() uint64 {
	return 32000
}

// CalculateLogGas 计算日志操作的Gas成本
func CalculateLogGas(topics int, dataSize int) uint64 {
	baseGas := uint64(375)
	topicGas := uint64(topics) * 375
	dataGas := uint64(dataSize) * 8

	return baseGas + topicGas + dataGas
}

// CalculateSha3Gas 计算SHA3操作的Gas成本
func CalculateSha3Gas(dataSize int) uint64 {
	baseGas := uint64(30)
	wordGas := uint64((dataSize+31)/32) * 6

	return baseGas + wordGas
}

// CalculateExpGas 计算EXP操作的Gas成本
func CalculateExpGas(exponent *big.Int) uint64 {
	baseGas := uint64(10)
	byteGas := uint64(exponent.BitLen()/8) * 10

	return baseGas + byteGas
}

// 辅助函数
func isZeroValue(value []byte) bool {
	for _, b := range value {
		if b != 0 {
			return false
		}
	}
	return true
}

// CalculateMemoryGas 计算内存扩展的Gas成本
func CalculateMemoryGas(size int) uint64 {
	if size <= 0 {
		return 0
	}

	// 按32字节对齐
	words := (size + 31) / 32
	// 内存Gas计算：3 * words^2 / 512
	return uint64(3 * words * words / 512)
}

// CalculateCallGas 计算调用的Gas成本
func CalculateCallGas(gasLimit, gasPrice uint64) uint64 {
	// 基础调用Gas
	baseGas := uint64(700)

	// 计算实际使用的Gas
	if gasLimit > 0 {
		return baseGas + gasLimit
	}

	return baseGas
}
