package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Transaction 交易结构
type Transaction struct {
	Nonce    uint64          `json:"nonce"`
	GasPrice *big.Int        `json:"gasPrice"`
	Gas      uint64          `json:"gas"`
	To       *common.Address `json:"to"`
	Value    *big.Int        `json:"value"`
	Data     []byte          `json:"data"`
	V        *big.Int        `json:"v"`
	R        *big.Int        `json:"r"`
	S        *big.Int        `json:"s"`
}

// TxType 交易类型
type TxType uint8

const (
	TxTypeLegacy TxType = iota
	TxTypeEIP1559
)

// NewTransaction 创建新交易
func NewTransaction(nonce uint64, to common.Address, value *big.Int, gas uint64, gasPrice *big.Int, data []byte) *Transaction {
	return &Transaction{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       &to,
		Value:    value,
		Data:     data,
		V:        big.NewInt(0),
		R:        big.NewInt(0),
		S:        big.NewInt(0),
	}
}

// NewContractCreation 创建合约创建交易
func NewContractCreation(nonce uint64, value *big.Int, gas uint64, gasPrice *big.Int, data []byte) *Transaction {
	return &Transaction{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       nil,
		Value:    value,
		Data:     data,
		V:        big.NewInt(0),
		R:        big.NewInt(0),
		S:        big.NewInt(0),
	}
}

// Hash 计算交易哈希
func (tx *Transaction) Hash() common.Hash {
	data, _ := json.Marshal(tx)
	return crypto.Keccak256Hash(data)
}

// IsContractCreation 判断是否为合约创建交易
func (tx *Transaction) IsContractCreation() bool {
	return tx.To == nil
}

// Sign 签名交易
func (tx *Transaction) Sign(privateKey []byte) error {
	// 简化实现，实际应该使用正确的ECDSA签名
	tx.V = big.NewInt(27)
	tx.R = big.NewInt(1)
	tx.S = big.NewInt(1)
	return nil
}

// Sender 获取交易发送者地址
func (tx *Transaction) Sender() (common.Address, error) {
	// 简化实现，实际应该从签名中恢复公钥
	return common.Address{}, nil
}

// Validate 验证交易
func (tx *Transaction) Validate() error {
	if tx.GasPrice.Sign() < 0 {
		return errors.New("invalid gas price")
	}
	if tx.Value.Sign() < 0 {
		return errors.New("invalid value")
	}
	if tx.Gas == 0 {
		return errors.New("invalid gas")
	}
	return nil
}

// GasPriceU64 获取Gas价格（uint64）
func (tx *Transaction) GasPriceU64() uint64 {
	return tx.GasPrice.Uint64()
}

// ValueU64 获取交易金额（uint64）
func (tx *Transaction) ValueU64() uint64 {
	return tx.Value.Uint64()
}

// DataLength 获取数据长度
func (tx *Transaction) DataLength() int {
	return len(tx.Data)
}

// Copy 复制交易
func (tx *Transaction) Copy() *Transaction {
	copyTx := *tx
	if tx.To != nil {
		copyTo := *tx.To
		copyTx.To = &copyTo
	}
	copyTx.GasPrice = new(big.Int).Set(tx.GasPrice)
	copyTx.Value = new(big.Int).Set(tx.Value)
	copyTx.V = new(big.Int).Set(tx.V)
	copyTx.R = new(big.Int).Set(tx.R)
	copyTx.S = new(big.Int).Set(tx.S)
	copyTx.Data = make([]byte, len(tx.Data))
	copy(copyTx.Data, tx.Data)
	return &copyTx
}

// Transactions 交易列表
type Transactions []*Transaction

// Hash 计算交易列表哈希
func (txs Transactions) Hash() common.Hash {
	if len(txs) == 0 {
		return common.Hash{}
	}
	data, _ := json.Marshal(txs)
	return crypto.Keccak256Hash(data)
}

// Len 获取交易数量
func (txs Transactions) Len() int {
	return len(txs)
}

// Get 获取指定索引的交易
func (txs Transactions) Get(i int) *Transaction {
	if i < 0 || i >= len(txs) {
		return nil
	}
	return txs[i]
}
