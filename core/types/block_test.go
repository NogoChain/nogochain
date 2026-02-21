package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestBlockCreation(t *testing.T) {
	// 创建测试区块
	parentHash := common.Hash{0x01}
	coinbase := common.Address{0x02}
	root := common.Hash{0x03}
	txHash := common.Hash{0x04}
	receiptHash := common.Hash{0x05}
	difficulty := big.NewInt(1000)
	number := big.NewInt(1)
	gasLimit := uint64(10000000)
	gasUsed := uint64(5000000)
	time := uint64(1700000001)
	extra := []byte("test")
	mixDigest := common.Hash{0x06}
	nonce := uint64(12345)
	transactions := []*Transaction{}
	uncles := []*BlockHeader{}

	block := NewBlock(parentHash, coinbase, root, txHash, receiptHash, difficulty, number, gasLimit, gasUsed, time, extra, mixDigest, nonce, transactions, uncles)

	if block.Header.ParentHash != parentHash {
		t.Errorf("ParentHash mismatch: expected %v, got %v", parentHash, block.Header.ParentHash)
	}

	if block.Header.Coinbase != coinbase {
		t.Errorf("Coinbase mismatch: expected %v, got %v", coinbase, block.Header.Coinbase)
	}

	if block.Header.Root != root {
		t.Errorf("Root mismatch: expected %v, got %v", root, block.Header.Root)
	}

	if block.Header.TxHash != txHash {
		t.Errorf("TxHash mismatch: expected %v, got %v", txHash, block.Header.TxHash)
	}

	if block.Header.ReceiptHash != receiptHash {
		t.Errorf("ReceiptHash mismatch: expected %v, got %v", receiptHash, block.Header.ReceiptHash)
	}

	if block.Header.Difficulty.Cmp(difficulty) != 0 {
		t.Errorf("Difficulty mismatch: expected %v, got %v", difficulty, block.Header.Difficulty)
	}

	if block.Header.Number.Cmp(number) != 0 {
		t.Errorf("Number mismatch: expected %v, got %v", number, block.Header.Number)
	}

	if block.Header.GasLimit != gasLimit {
		t.Errorf("GasLimit mismatch: expected %v, got %v", gasLimit, block.Header.GasLimit)
	}

	if block.Header.GasUsed != gasUsed {
		t.Errorf("GasUsed mismatch: expected %v, got %v", gasUsed, block.Header.GasUsed)
	}

	if block.Header.Time != time {
		t.Errorf("Time mismatch: expected %v, got %v", time, block.Header.Time)
	}

	if string(block.Header.Extra) != string(extra) {
		t.Errorf("Extra mismatch: expected %v, got %v", extra, block.Header.Extra)
	}

	if block.Header.MixDigest != mixDigest {
		t.Errorf("MixDigest mismatch: expected %v, got %v", mixDigest, block.Header.MixDigest)
	}

	if block.Header.Nonce != nonce {
		t.Errorf("Nonce mismatch: expected %v, got %v", nonce, block.Header.Nonce)
	}
}

func TestBlockHash(t *testing.T) {
	// 创建测试区块
	block := NewBlock(
		common.Hash{0x01},
		common.Address{0x02},
		common.Hash{0x03},
		common.Hash{0x04},
		common.Hash{0x05},
		big.NewInt(1000),
		big.NewInt(1),
		10000000,
		5000000,
		1700000001,
		[]byte("test"),
		common.Hash{0x06},
		12345,
		[]*Transaction{},
		[]*BlockHeader{},
	)

	hash := block.Hash()
	if hash == (common.Hash{}) {
		t.Errorf("Block hash is zero")
	}

	headerHash := block.Header.Hash()
	if headerHash == (common.Hash{}) {
		t.Errorf("Header hash is zero")
	}
}

func TestBlockMethods(t *testing.T) {
	// 创建测试区块
	block := NewBlock(
		common.Hash{0x01},
		common.Address{0x02},
		common.Hash{0x03},
		common.Hash{0x04},
		common.Hash{0x05},
		big.NewInt(1000),
		big.NewInt(5),
		10000000,
		5000000,
		1700000001,
		[]byte("test"),
		common.Hash{0x06},
		12345,
		[]*Transaction{},
		[]*BlockHeader{},
	)

	if block.NumberU64() != 5 {
		t.Errorf("NumberU64 mismatch: expected 5, got %v", block.NumberU64())
	}

	if block.DifficultyU64() != 1000 {
		t.Errorf("DifficultyU64 mismatch: expected 1000, got %v", block.DifficultyU64())
	}

	if block.GasLimit() != 10000000 {
		t.Errorf("GasLimit mismatch: expected 10000000, got %v", block.GasLimit())
	}

	if block.GasUsed() != 5000000 {
		t.Errorf("GasUsed mismatch: expected 5000000, got %v", block.GasUsed())
	}

	if block.Coinbase() != common.HexToAddress("0x02") {
		t.Errorf("Coinbase mismatch: expected %v, got %v", common.HexToAddress("0x02"), block.Coinbase())
	}

	if block.ParentHash() != common.HexToHash("0x01") {
		t.Errorf("ParentHash mismatch: expected %v, got %v", common.HexToHash("0x01"), block.ParentHash())
	}

	if block.TxCount() != 0 {
		t.Errorf("TxCount mismatch: expected 0, got %v", block.TxCount())
	}

	if block.UncleCount() != 0 {
		t.Errorf("UncleCount mismatch: expected 0, got %v", block.UncleCount())
	}
}
