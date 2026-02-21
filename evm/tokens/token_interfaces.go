package tokens

import (
	"math/big"
)

// ERC20Interface ERC-20代币标准接口
type ERC20Interface interface {
	// 代币总供应量
	TotalSupply() *big.Int

	// 账户余额
	BalanceOf(owner []byte) *big.Int

	// 转账
	Transfer(to []byte, value *big.Int) bool

	// 授权转账
	Allowance(owner, spender []byte) *big.Int
	Approve(spender []byte, value *big.Int) bool
	TransferFrom(from, to []byte, value *big.Int) bool

	// 事件
	TransferEvent(from, to []byte, value *big.Int)
	ApprovalEvent(owner, spender []byte, value *big.Int)
}

// ERC721Interface ERC-721代币标准接口
type ERC721Interface interface {
	// 余额（拥有的NFT数量）
	BalanceOf(owner []byte) *big.Int

	// 所有者
	OwnerOf(tokenId *big.Int) []byte

	// 安全转账
	SafeTransferFrom(from, to []byte, tokenId *big.Int, data []byte)
	TransferFrom(from, to []byte, tokenId *big.Int)

	// 授权
	Approve(to []byte, tokenId *big.Int)
	SetApprovalForAll(operator []byte, approved bool)
	GetApproved(tokenId *big.Int) []byte
	IsApprovedForAll(owner, operator []byte) bool

	// 事件
	TransferEvent(from, to []byte, tokenId *big.Int)
	ApprovalEvent(owner, approved []byte, tokenId *big.Int)
	ApprovalForAllEvent(owner, operator []byte, approved bool)
}

// ERC1155Interface ERC-1155代币标准接口
type ERC1155Interface interface {
	// 余额
	BalanceOf(owner []byte, id *big.Int) *big.Int
	BalanceOfBatch(owners [][]byte, ids []*big.Int) []*big.Int

	// 授权
	SetApprovalForAll(operator []byte, approved bool)
	IsApprovedForAll(owner, operator []byte) bool

	// 安全转账
	SafeTransferFrom(from, to []byte, id *big.Int, value *big.Int, data []byte)
	SafeBatchTransferFrom(from, to []byte, ids []*big.Int, values []*big.Int, data []byte)

	// 事件
	TransferSingleEvent(operator, from, to []byte, id, value *big.Int)
	TransferBatchEvent(operator, from, to []byte, ids, values []*big.Int)
	ApprovalForAllEvent(owner, operator []byte, approved bool)
}
