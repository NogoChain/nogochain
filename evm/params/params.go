package params

import (
	"math/big"
)

// EVM 相关参数

// 链ID
const (
	ChainID = 318
)

// 代币参数
const (
	TokenSymbol   = "NOGO"
	TokenDecimals = 18
)

// EVM 指令集参数
const (
	// 栈深度限制
	StackLimit = 1024

	// 内存限制（按32字节字计算）
	MemoryLimit = 1024 * 1024 // 1MB

	// 代码大小限制
	CodeSizeLimit = 24576 // 24KB

	// 交易数据Gas成本
	TxDataZeroGas           = 4
	TxDataNonZeroGas        = 16
	TxDataNonZeroGasEIP2028 = 16

	// 交易基础Gas
	TxGas                 = 21000
	TxGasContractCreation = 53000

	// 区块Gas限制相关
	MinGasLimit          = 5000
	MaxGasLimit          = 10000000
	GasLimitBoundDivisor = 1024

	// 存储操作Gas成本
	SloadGas        = 800
	SstoreSetGas    = 20000
	SstoreResetGas  = 5000
	SstoreClearGas  = 5000
	SstoreRefundGas = 15000

	// 调用操作Gas成本
	CallGas              = 700
	CallStipend          = 2300
	CallValueTransferGas = 9000
	CallNewAccountGas    = 25000

	// 日志操作Gas成本
	LogGas      = 375
	LogDataGas  = 8
	LogTopicGas = 375

	// SHA3操作Gas成本
	Sha3Gas     = 30
	Sha3WordGas = 6

	// 内存操作Gas成本
	MemoryGas    = 3
	QuadCoeffDiv = 512

	// 合约创建Gas成本
	CreateGas  = 32000
	Create2Gas = 32000

	// 自毁操作Gas成本
	SuicideGas       = 5000
	SuicideRefundGas = 24000

	// 其他操作Gas成本
	ExpGas         = 10
	ExpByteGas     = 10
	SigExtendGas   = 5
	ExtCodeSizeGas = 700
	ExtCodeCopyGas = 700
	ExtCodeHashGas = 700
	BalanceGas     = 400
	BlockHashGas   = 20
)

// 硬分叉激活区块
var (
	// 基础费相关硬分叉
	EIP1559Block = big.NewInt(0)

	// 伦敦硬分叉
	LondonBlock = big.NewInt(0)

	// 柏林硬分叉
	BerlinBlock = big.NewInt(0)

	// 伊斯坦布尔硬分叉
	IstanbulBlock = big.NewInt(0)

	// 君士坦丁堡硬分叉
	ConstantinopleBlock = big.NewInt(0)

	// 彼得斯堡硬分叉
	PetersburgBlock = big.NewInt(0)

	//  Byzantium硬分叉
	ByzantiumBlock = big.NewInt(0)

	// Spurious Dragon硬分叉
	SpuriousDragonBlock = big.NewInt(0)

	// Tangerine Whistle硬分叉
	TangerineWhistleBlock = big.NewInt(0)

	// Homestead硬分叉
	HomesteadBlock = big.NewInt(0)
)

// 硬分叉配置
func IsEIP1559(blockNumber *big.Int) bool {
	return blockNumber.Cmp(EIP1559Block) >= 0
}

func IsLondon(blockNumber *big.Int) bool {
	return blockNumber.Cmp(LondonBlock) >= 0
}

func IsBerlin(blockNumber *big.Int) bool {
	return blockNumber.Cmp(BerlinBlock) >= 0
}

func IsIstanbul(blockNumber *big.Int) bool {
	return blockNumber.Cmp(IstanbulBlock) >= 0
}

func IsConstantinople(blockNumber *big.Int) bool {
	return blockNumber.Cmp(ConstantinopleBlock) >= 0
}

func IsPetersburg(blockNumber *big.Int) bool {
	return blockNumber.Cmp(PetersburgBlock) >= 0
}

func IsByzantium(blockNumber *big.Int) bool {
	return blockNumber.Cmp(ByzantiumBlock) >= 0
}

func IsSpuriousDragon(blockNumber *big.Int) bool {
	return blockNumber.Cmp(SpuriousDragonBlock) >= 0
}

func IsTangerineWhistle(blockNumber *big.Int) bool {
	return blockNumber.Cmp(TangerineWhistleBlock) >= 0
}

func IsHomestead(blockNumber *big.Int) bool {
	return blockNumber.Cmp(HomesteadBlock) >= 0
}

// Gas 计算辅助函数
func CalculateDataGas(data []byte, isEIP2028 bool) uint64 {
	var gas uint64
	for _, b := range data {
		if b == 0 {
			gas += TxDataZeroGas
		} else {
			if isEIP2028 {
				gas += TxDataNonZeroGasEIP2028
			} else {
				gas += TxDataNonZeroGas
			}
		}
	}
	return gas
}

func CalculateMemoryGas(size uint64) uint64 {
	if size == 0 {
		return 0
	}
	// 按32字节对齐
	words := (size + 31) / 32
	return MemoryGas*words + words*words/QuadCoeffDiv
}

func CalculateLogGas(topics int, dataSize int) uint64 {
	return LogGas + uint64(topics)*LogTopicGas + uint64(dataSize)*LogDataGas
}

func CalculateSha3Gas(dataSize int) uint64 {
	words := (dataSize + 31) / 32
	return Sha3Gas + Sha3WordGas*uint64(words)
}
