package utils

import (
	"bytes"
	"math/big"
)

// 接口ID定义
var (
	// ERC-20 接口ID
	ERC20InterfaceID = []byte{0x36, 0x37, 0x2b, 0x07}

	// ERC-721 接口ID
	ERC721InterfaceID = []byte{0x80, 0xac, 0x58, 0xcd}

	// ERC-1155 接口ID
	ERC1155InterfaceID = []byte{0xd9, 0xb6, 0x7a, 0x26}

	// ERC-165 接口ID
	ERC165InterfaceID = []byte{0x01, 0xff, 0xc9, 0xa7}
)

// IsERC20 检查合约是否实现了ERC-20标准
func IsERC20(supportsInterface func(interfaceID []byte) bool) bool {
	return supportsInterface(ERC20InterfaceID)
}

// IsERC721 检查合约是否实现了ERC-721标准
func IsERC721(supportsInterface func(interfaceID []byte) bool) bool {
	return supportsInterface(ERC721InterfaceID)
}

// IsERC1155 检查合约是否实现了ERC-1155标准
func IsERC1155(supportsInterface func(interfaceID []byte) bool) bool {
	return supportsInterface(ERC1155InterfaceID)
}

// IsERC165 检查合约是否实现了ERC-165标准
func IsERC165(supportsInterface func(interfaceID []byte) bool) bool {
	return supportsInterface(ERC165InterfaceID)
}

// CalculateERC20TransferGas 计算ERC-20转账的Gas成本
func CalculateERC20TransferGas() uint64 {
	// 基础转账Gas成本
	return 21000 + 60000 // 包含基础交易Gas和合约执行Gas
}

// CalculateERC721TransferGas 计算ERC-721转账的Gas成本
func CalculateERC721TransferGas() uint64 {
	// NFT转账Gas成本
	return 21000 + 150000 // 包含基础交易Gas和合约执行Gas
}

// CalculateERC1155TransferGas 计算ERC-1155转账的Gas成本
func CalculateERC1155TransferGas(batchSize int) uint64 {
	// 多代币转账Gas成本
	baseGas := uint64(21000 + 100000)
	batchGas := uint64(batchSize) * 10000
	return baseGas + batchGas
}

// ValidateERC20Transfer 验证ERC-20转账
func ValidateERC20Transfer(from, to []byte, value *big.Int, balance *big.Int) error {
	// 检查转账金额是否为正
	if value.Sign() <= 0 {
		return ErrInvalidTransferValue
	}

	// 检查余额是否足够
	if balance.Cmp(value) < 0 {
		return ErrInsufficientBalance
	}

	// 检查接收地址是否有效
	if len(to) != 20 {
		return ErrInvalidAddress
	}

	return nil
}

// ValidateERC721Transfer 验证ERC-721转账
func ValidateERC721Transfer(from, to []byte, tokenId *big.Int, owner []byte) error {
	// 检查发送者是否为所有者
	if !bytes.Equal(from, owner) {
		return ErrNotTokenOwner
	}

	// 检查接收地址是否有效
	if len(to) != 20 {
		return ErrInvalidAddress
	}

	// 检查tokenId是否有效
	if tokenId.Sign() < 0 {
		return ErrInvalidTokenId
	}

	return nil
}

// ValidateERC1155Transfer 验证ERC-1155转账
func ValidateERC1155Transfer(from, to []byte, id, value *big.Int, balance *big.Int) error {
	// 检查转账金额是否为正
	if value.Sign() <= 0 {
		return ErrInvalidTransferValue
	}

	// 检查余额是否足够
	if balance.Cmp(value) < 0 {
		return ErrInsufficientBalance
	}

	// 检查接收地址是否有效
	if len(to) != 20 {
		return ErrInvalidAddress
	}

	// 检查tokenId是否有效
	if id.Sign() < 0 {
		return ErrInvalidTokenId
	}

	return nil
}

// 错误定义
var (
	ErrInvalidTransferValue = NewTokenError("invalid transfer value")
	ErrInsufficientBalance  = NewTokenError("insufficient balance")
	ErrInvalidAddress       = NewTokenError("invalid address")
	ErrNotTokenOwner        = NewTokenError("not token owner")
	ErrInvalidTokenId       = NewTokenError("invalid token id")
	ErrApprovalRequired     = NewTokenError("approval required")
	ErrTokenDoesNotExist    = NewTokenError("token does not exist")
)

// TokenError 代币错误
type TokenError struct {
	message string
}

// NewTokenError 创建新的代币错误
func NewTokenError(message string) *TokenError {
	return &TokenError{message: message}
}

// Error 实现error接口
func (e *TokenError) Error() string {
	return e.message
}
