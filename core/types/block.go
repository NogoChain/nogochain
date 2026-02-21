package types

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// BlockHeader represents a block header structure
// BlockHeader 区块头结构
type BlockHeader struct {
	ParentHash  common.Hash    `json:"parentHash"`
	UncleHash   common.Hash    `json:"uncleHash"`
	Coinbase    common.Address `json:"coinbase"`
	Root        common.Hash    `json:"root"`
	TxHash      common.Hash    `json:"txHash"`
	ReceiptHash common.Hash    `json:"receiptHash"`
	Bloom       []byte         `json:"bloom"`
	Difficulty  *big.Int       `json:"difficulty"`
	Number      *big.Int       `json:"number"`
	GasLimit    uint64         `json:"gasLimit"`
	GasUsed     uint64         `json:"gasUsed"`
	Time        uint64         `json:"time"`
	Extra       []byte         `json:"extra"`
	MixDigest   common.Hash    `json:"mixDigest"`
	Nonce       uint64         `json:"nonce"`
}

// Block represents a block structure
// Block 区块结构
type Block struct {
	Header       *BlockHeader   `json:"header"`
	Transactions []*Transaction `json:"transactions"`
	Uncles       []*BlockHeader `json:"uncles"`
}

// NewBlock creates a new block
// NewBlock 创建新区块
func NewBlock(parentHash common.Hash, coinbase common.Address, root common.Hash, txHash common.Hash, receiptHash common.Hash, difficulty *big.Int, number *big.Int, gasLimit uint64, gasUsed uint64, time uint64, extra []byte, mixDigest common.Hash, nonce uint64, transactions []*Transaction, uncles []*BlockHeader) *Block {
	return &Block{
		Header: &BlockHeader{
			ParentHash:  parentHash,
			UncleHash:   CalcUncleHash(uncles),
			Coinbase:    coinbase,
			Root:        root,
			TxHash:      txHash,
			ReceiptHash: receiptHash,
			Bloom:       make([]byte, 256),
			Difficulty:  difficulty,
			Number:      number,
			GasLimit:    gasLimit,
			GasUsed:     gasUsed,
			Time:        time,
			Extra:       extra,
			MixDigest:   mixDigest,
			Nonce:       nonce,
		},
		Transactions: transactions,
		Uncles:       uncles,
	}
}

// Hash calculates the block hash
// Hash 计算区块哈希
func (b *Block) Hash() common.Hash {
	data, _ := json.Marshal(b.Header)
	return crypto.Keccak256Hash(data)
}

// Hash calculates the block header hash
// HeaderHash 计算区块头哈希
func (h *BlockHeader) Hash() common.Hash {
	data, _ := json.Marshal(h)
	return crypto.Keccak256Hash(data)
}

// CalcUncleHash calculates the uncle hash
// CalcUncleHash 计算叔区块哈希
func CalcUncleHash(uncles []*BlockHeader) common.Hash {
	if len(uncles) == 0 {
		return common.Hash{}
	}
	data, _ := json.Marshal(uncles)
	return crypto.Keccak256Hash(data)
}

// CalcTxHash calculates the transaction root hash
// CalcTxHash 计算交易根哈希
func CalcTxHash(transactions []*Transaction) common.Hash {
	if len(transactions) == 0 {
		return common.Hash{}
	}
	data, _ := json.Marshal(transactions)
	return crypto.Keccak256Hash(data)
}

// Timestamp gets the block timestamp
// Timestamp 获取区块时间戳
func (b *Block) Timestamp() time.Time {
	return time.Unix(int64(b.Header.Time), 0)
}

// NumberU64 gets the block number as uint64
// NumberU64 获取区块号（uint64）
func (b *Block) NumberU64() uint64 {
	return b.Header.Number.Uint64()
}

// DifficultyU64 gets the difficulty as uint64
// DifficultyU64 获取难度（uint64）
func (b *Block) DifficultyU64() uint64 {
	return b.Header.Difficulty.Uint64()
}

// GasLimit gets the gas limit
// GasLimit 获取Gas限制
func (b *Block) GasLimit() uint64 {
	return b.Header.GasLimit
}

// GasUsed gets the gas used
// GasUsed 获取Gas使用量
func (b *Block) GasUsed() uint64 {
	return b.Header.GasUsed
}

// Coinbase gets the miner address
// Coinbase 获取矿工地址
func (b *Block) Coinbase() common.Address {
	return b.Header.Coinbase
}

// ParentHash gets the parent block hash
// ParentHash 获取父区块哈希
func (b *Block) ParentHash() common.Hash {
	return b.Header.ParentHash
}

// TxCount gets the transaction count
// TxCount 获取交易数量
func (b *Block) TxCount() int {
	return len(b.Transactions)
}

// UncleCount gets the uncle count
// UncleCount 获取叔区块数量
func (b *Block) UncleCount() int {
	return len(b.Uncles)
}
