package vm

import (
	"fmt"
	"math/big"

	"nogochain/evm/core/vm/gas"
)

// Instruction 指令接口
type Instruction interface {
	Execute(evm *EVM) error
	GasCost() uint64
}

// 指令映射
var instructions = make(map[byte]Instruction)

// RegisterInstruction 注册指令
func RegisterInstruction(opcode byte, instr Instruction) {
	instructions[opcode] = instr
}

// GetInstruction 获取指令
func GetInstruction(opcode byte) Instruction {
	return instructions[opcode]
}

// ADD 加法指令
type ADD struct{}

func (ADD) Execute(evm *EVM) error {
	// 消耗Gas
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x01)); err != nil {
		return err
	}

	// 弹出两个操作数
	a, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	b, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	// 执行加法
	result := new(big.Int).Add(a, b)

	// 压入结果
	return evm.Stack.Push(result)
}

func (ADD) GasCost() uint64 {
	return gas.CalculateBaseGas(0x01)
}

// STOP 停止指令
type STOP struct{}

func (STOP) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x00)); err != nil {
		return err
	}

	evm.Stop()
	return nil
}

func (STOP) GasCost() uint64 {
	return gas.CalculateBaseGas(0x00)
}

// PUSH 推送指令
type PUSH struct {
	Size int
}

func (p PUSH) Execute(evm *EVM) error {
	// 消耗Gas
	if err := evm.ConsumeGas(gas.CalculateBaseGas(byte(0x60 + p.Size - 1))); err != nil {
		return err
	}

	// 从代码中读取数据
	if evm.ProgramCounter+1+p.Size > len(evm.Context.Code) {
		return fmt.Errorf("insufficient code for PUSH%d", p.Size)
	}

	data := evm.Context.Code[evm.ProgramCounter+1 : evm.ProgramCounter+1+p.Size]

	// 转换为big.Int
	value := new(big.Int).SetBytes(data)

	// 压入栈
	return evm.Stack.Push(value)
}

func (p PUSH) GasCost() uint64 {
	return gas.CalculateBaseGas(byte(0x60 + p.Size - 1))
}

// RETURN 返回指令
type RETURN struct{}

func (RETURN) Execute(evm *EVM) error {
	// 消耗Gas
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0xf3)); err != nil {
		return err
	}

	// 弹出两个操作数
	offset, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	size, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	// 从内存读取数据
	data := evm.Memory.Get(int(offset.Uint64()), int(size.Uint64()))

	// 设置返回数据
	evm.SetReturnData(data)

	// 停止执行
	evm.Stop()

	return nil
}

func (RETURN) GasCost() uint64 {
	return gas.CalculateBaseGas(0xf3)
}

// MUL 乘法指令
type MUL struct{}

func (MUL) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x02)); err != nil {
		return err
	}

	a, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	b, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	result := new(big.Int).Mul(a, b)
	return evm.Stack.Push(result)
}

func (MUL) GasCost() uint64 {
	return gas.CalculateBaseGas(0x02)
}

// SUB 减法指令
type SUB struct{}

func (SUB) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x03)); err != nil {
		return err
	}

	a, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	b, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	result := new(big.Int).Sub(b, a)
	return evm.Stack.Push(result)
}

func (SUB) GasCost() uint64 {
	return gas.CalculateBaseGas(0x03)
}

// DIV 除法指令
type DIV struct{}

func (DIV) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x04)); err != nil {
		return err
	}

	a, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	b, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	if a.Sign() == 0 {
		return evm.Stack.Push(big.NewInt(0))
	}

	result := new(big.Int).Div(b, a)
	return evm.Stack.Push(result)
}

func (DIV) GasCost() uint64 {
	return gas.CalculateBaseGas(0x04)
}

// POP 弹出指令
type POP struct{}

func (POP) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x50)); err != nil {
		return err
	}

	_, err := evm.Stack.Pop()
	return err
}

func (POP) GasCost() uint64 {
	return gas.CalculateBaseGas(0x50)
}

// MLOAD 内存加载指令
type MLOAD struct{}

func (MLOAD) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x51)); err != nil {
		return err
	}

	offset, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	offsetVal := int(offset.Uint64())

	// 计算内存扩展的Gas成本
	memorySize := offsetVal + 32
	memoryGas := evm.Memory.CalculateGasCost(memorySize)
	if err := evm.ConsumeGas(memoryGas); err != nil {
		return err
	}

	data := evm.Memory.Get(offsetVal, 32)
	value := new(big.Int).SetBytes(data)
	return evm.Stack.Push(value)
}

func (MLOAD) GasCost() uint64 {
	return gas.CalculateBaseGas(0x51)
}

// MSTORE 内存存储指令
type MSTORE struct{}

func (MSTORE) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x52)); err != nil {
		return err
	}

	offset, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	value, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	data := value.Bytes()
	offsetVal := int(offset.Uint64())

	// 计算内存扩展的Gas成本
	memorySize := offsetVal + len(data)
	memoryGas := evm.Memory.CalculateGasCost(memorySize)
	if err := evm.ConsumeGas(memoryGas); err != nil {
		return err
	}

	evm.Memory.Set(offsetVal, data)
	return nil
}

func (MSTORE) GasCost() uint64 {
	return gas.CalculateBaseGas(0x52)
}

// MSTORE8 内存存储单字节指令
type MSTORE8 struct{}

func (MSTORE8) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x53)); err != nil {
		return err
	}

	offset, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	value, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	offsetVal := int(offset.Uint64())

	// 计算内存扩展的Gas成本
	memorySize := offsetVal + 1
	memoryGas := evm.Memory.CalculateGasCost(memorySize)
	if err := evm.ConsumeGas(memoryGas); err != nil {
		return err
	}

	evm.Memory.SetByte(offsetVal, byte(value.Uint64()))
	return nil
}

func (MSTORE8) GasCost() uint64 {
	return gas.CalculateBaseGas(0x53)
}

// SLOAD 存储加载指令
type SLOAD struct{}

func (SLOAD) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x54)); err != nil {
		return err
	}

	key, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	// 这里应该使用合约地址，暂时使用空地址
	value := evm.Storage.Get(key.Bytes())
	result := new(big.Int).SetBytes(value)
	return evm.Stack.Push(result)
}

func (SLOAD) GasCost() uint64 {
	return gas.CalculateBaseGas(0x54)
}

// SSTORE 存储存储指令
type SSTORE struct{}

func (SSTORE) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x55)); err != nil {
		return err
	}

	key, err := evm.Stack.Pop()
	if err != nil {
		return err
	}
	value, err := evm.Stack.Pop()
	if err != nil {
		return err
	}

	// 这里应该使用合约地址，暂时使用空地址
	evm.Storage.Set(key.Bytes(), value.Bytes())
	return nil
}

func (SSTORE) GasCost() uint64 {
	return gas.CalculateBaseGas(0x55)
}

// PC 程序计数器指令
type PC struct{}

func (PC) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x58)); err != nil {
		return err
	}

	return evm.Stack.Push(big.NewInt(int64(evm.ProgramCounter)))
}

func (PC) GasCost() uint64 {
	return gas.CalculateBaseGas(0x58)
}

// MSIZE 内存大小指令
type MSIZE struct{}

func (MSIZE) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x59)); err != nil {
		return err
	}

	return evm.Stack.Push(big.NewInt(int64(evm.Memory.Size())))
}

func (MSIZE) GasCost() uint64 {
	return gas.CalculateBaseGas(0x59)
}

// GAS 剩余Gas指令
type GAS struct{}

func (GAS) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(0x5a)); err != nil {
		return err
	}

	return evm.Stack.Push(big.NewInt(int64(evm.GetGasLeft())))
}

func (GAS) GasCost() uint64 {
	return gas.CalculateBaseGas(0x5a)
}

// DUP 复制指令
type DUP struct {
	Index int
}

func (d DUP) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(byte(0x80 + d.Index - 1))); err != nil {
		return err
	}

	value, err := evm.Stack.PeekN(d.Index - 1)
	if err != nil {
		return err
	}

	return evm.Stack.Push(value)
}

func (d DUP) GasCost() uint64 {
	return gas.CalculateBaseGas(byte(0x80 + d.Index - 1))
}

// SWAP 交换指令
type SWAP struct {
	Index int
}

func (s SWAP) Execute(evm *EVM) error {
	if err := evm.ConsumeGas(gas.CalculateBaseGas(byte(0x90 + s.Index - 1))); err != nil {
		return err
	}

	return evm.Stack.Swap(s.Index)
}

func (s SWAP) GasCost() uint64 {
	return gas.CalculateBaseGas(byte(0x90 + s.Index - 1))
}

// 初始化指令
func init() {
	// 注册基础指令
	RegisterInstruction(0x00, STOP{})
	RegisterInstruction(0x01, ADD{})
	RegisterInstruction(0x02, MUL{})
	RegisterInstruction(0x03, SUB{})
	RegisterInstruction(0x04, DIV{})
	RegisterInstruction(0x50, POP{})
	RegisterInstruction(0x51, MLOAD{})
	RegisterInstruction(0x52, MSTORE{})
	RegisterInstruction(0x53, MSTORE8{})
	RegisterInstruction(0x54, SLOAD{})
	RegisterInstruction(0x55, SSTORE{})
	RegisterInstruction(0x58, PC{})
	RegisterInstruction(0x59, MSIZE{})
	RegisterInstruction(0x5a, GAS{})
	RegisterInstruction(0xf3, RETURN{})

	// 注册PUSH指令
	for i := 1; i <= 32; i++ {
		RegisterInstruction(byte(0x60+i-1), PUSH{Size: i})
	}

	// 注册DUP指令
	for i := 1; i <= 16; i++ {
		RegisterInstruction(byte(0x80+i-1), DUP{Index: i})
	}

	// 注册SWAP指令
	for i := 1; i <= 16; i++ {
		RegisterInstruction(byte(0x90+i-1), SWAP{Index: i})
	}

	// 注册其他指令...
}
