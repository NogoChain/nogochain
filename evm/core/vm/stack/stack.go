package stack

import (
	"errors"
	"math/big"
)

// Stack 以太坊虚拟机栈实现
type Stack struct {
	data  []*big.Int
	limit int
	used  int
}

// 错误定义
var (
	ErrStackOverflow  = errors.New("stack overflow")
	ErrStackUnderflow = errors.New("stack underflow")
)

// NewStack 创建新的栈实例
func NewStack() *Stack {
	return &Stack{
		data:  make([]*big.Int, 1024), // 预分配最大深度
		limit: 1024,                   // EVM栈深度限制
		used:  0,
	}
}

// Push 压入栈
func (s *Stack) Push(value *big.Int) error {
	if s.used >= s.limit {
		return ErrStackOverflow
	}

	// 复制值以避免引用问题
	copy := new(big.Int).Set(value)
	s.data[s.used] = copy
	s.used++
	return nil
}

// Pop 弹出栈顶元素
func (s *Stack) Pop() (*big.Int, error) {
	if s.used == 0 {
		return nil, ErrStackUnderflow
	}

	s.used--
	value := s.data[s.used]
	// 清理引用
	s.data[s.used] = nil
	return value, nil
}

// Peek 查看栈顶元素但不弹出
func (s *Stack) Peek() (*big.Int, error) {
	if s.used == 0 {
		return nil, ErrStackUnderflow
	}

	return s.data[s.used-1], nil
}

// PeekN 查看栈中指定位置的元素
func (s *Stack) PeekN(n int) (*big.Int, error) {
	if n < 0 || n >= s.used {
		return nil, ErrStackUnderflow
	}

	return s.data[s.used-n-1], nil
}

// Swap 交换栈顶两个元素
func (s *Stack) Swap(n int) error {
	if n < 0 || n >= s.used {
		return ErrStackUnderflow
	}

	idx1 := s.used - 1
	idx2 := s.used - n - 1
	s.data[idx1], s.data[idx2] = s.data[idx2], s.data[idx1]
	return nil
}

// Dup 复制栈顶元素到指定位置
func (s *Stack) Dup(n int) error {
	if n < 0 || n >= s.used {
		return ErrStackUnderflow
	}

	if s.used >= s.limit {
		return ErrStackOverflow
	}

	value := new(big.Int).Set(s.data[s.used-n-1])
	s.data[s.used] = value
	s.used++
	return nil
}

// Depth 获取当前栈深度
func (s *Stack) Depth() int {
	return s.used
}

// Reset 重置栈
func (s *Stack) Reset() {
	// 清理所有引用
	for i := 0; i < s.used; i++ {
		s.data[i] = nil
	}
	s.used = 0
}

// Back 查看栈中指定深度的元素
func (s *Stack) Back(n int) (*big.Int, error) {
	if n < 0 || n >= s.used {
		return nil, ErrStackUnderflow
	}

	return s.data[n], nil
}
