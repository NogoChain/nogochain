package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestTransactionCreation(t *testing.T) {
	// 创建测试交易
	nonce := uint64(1)
	to := common.Address{0x01}
	value := big.NewInt(1000000000000000000) // 1 NOGO
	gas := uint64(21000)
	gasPrice := big.NewInt(1000000000) // 1 Gwei
	data := []byte("test")

	tx := NewTransaction(nonce, to, value, gas, gasPrice, data)

	if tx.Nonce != nonce {
		t.Errorf("Nonce mismatch: expected %v, got %v", nonce, tx.Nonce)
	}

	if tx.To == nil {
		t.Errorf("To is nil")
	} else if *tx.To != to {
		t.Errorf("To mismatch: expected %v, got %v", to, *tx.To)
	}

	if tx.Value.Cmp(value) != 0 {
		t.Errorf("Value mismatch: expected %v, got %v", value, tx.Value)
	}

	if tx.Gas != gas {
		t.Errorf("Gas mismatch: expected %v, got %v", gas, tx.Gas)
	}

	if tx.GasPrice.Cmp(gasPrice) != 0 {
		t.Errorf("GasPrice mismatch: expected %v, got %v", gasPrice, tx.GasPrice)
	}

	if string(tx.Data) != string(data) {
		t.Errorf("Data mismatch: expected %v, got %v", data, tx.Data)
	}
}

func TestContractCreation(t *testing.T) {
	// 创建合约创建交易
	nonce := uint64(1)
	value := big.NewInt(0)
	gas := uint64(1000000)
	gasPrice := big.NewInt(1000000000)
	data := []byte("contract code")

	tx := NewContractCreation(nonce, value, gas, gasPrice, data)

	if tx.Nonce != nonce {
		t.Errorf("Nonce mismatch: expected %v, got %v", nonce, tx.Nonce)
	}

	if tx.To != nil {
		t.Errorf("To should be nil for contract creation")
	}

	if tx.Value.Cmp(value) != 0 {
		t.Errorf("Value mismatch: expected %v, got %v", value, tx.Value)
	}

	if tx.Gas != gas {
		t.Errorf("Gas mismatch: expected %v, got %v", gas, tx.Gas)
	}

	if tx.GasPrice.Cmp(gasPrice) != 0 {
		t.Errorf("GasPrice mismatch: expected %v, got %v", gasPrice, tx.GasPrice)
	}

	if string(tx.Data) != string(data) {
		t.Errorf("Data mismatch: expected %v, got %v", data, tx.Data)
	}

	if !tx.IsContractCreation() {
		t.Errorf("IsContractCreation should return true")
	}
}

func TestTransactionHash(t *testing.T) {
	// 创建测试交易
	tx := NewTransaction(
		1,
		common.Address{0x01},
		big.NewInt(1000000000000000000),
		21000,
		big.NewInt(1000000000),
		[]byte("test"),
	)

	hash := tx.Hash()
	if hash == (common.Hash{}) {
		t.Errorf("Transaction hash is zero")
	}
}

func TestTransactionMethods(t *testing.T) {
	// 创建测试交易
	tx := NewTransaction(
		1,
		common.Address{0x01},
		big.NewInt(1000000000000000000),
		21000,
		big.NewInt(1000000000),
		[]byte("test"),
	)

	if tx.GasPriceU64() != 1000000000 {
		t.Errorf("GasPriceU64 mismatch: expected 1000000000, got %v", tx.GasPriceU64())
	}

	if tx.ValueU64() != 1000000000000000000 {
		t.Errorf("ValueU64 mismatch: expected 1000000000000000000, got %v", tx.ValueU64())
	}

	if tx.DataLength() != 4 {
		t.Errorf("DataLength mismatch: expected 4, got %v", tx.DataLength())
	}

	if tx.IsContractCreation() {
		t.Errorf("IsContractCreation should return false")
	}
}

func TestTransactionCopy(t *testing.T) {
	// 创建测试交易
	tx := NewTransaction(
		1,
		common.Address{0x01},
		big.NewInt(1000000000000000000),
		21000,
		big.NewInt(1000000000),
		[]byte("test"),
	)

	copyTx := tx.Copy()

	if copyTx.Nonce != tx.Nonce {
		t.Errorf("Nonce mismatch after copy")
	}

	if copyTx.To == nil || tx.To == nil || *copyTx.To != *tx.To {
		t.Errorf("To mismatch after copy")
	}

	if copyTx.Value.Cmp(tx.Value) != 0 {
		t.Errorf("Value mismatch after copy")
	}

	if copyTx.Gas != tx.Gas {
		t.Errorf("Gas mismatch after copy")
	}

	if copyTx.GasPrice.Cmp(tx.GasPrice) != 0 {
		t.Errorf("GasPrice mismatch after copy")
	}

	if string(copyTx.Data) != string(tx.Data) {
		t.Errorf("Data mismatch after copy")
	}
}
