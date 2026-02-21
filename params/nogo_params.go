package params

import (
	"math/big"
)

// 基础链参数
const (
	// ChainID - 链ID，用于标识NogoChain
	ChainID uint64 = 318

	// Symbol - 代币符号
	Symbol string = "NOGO"

	// Decimals - 代币小数位数
	Decimals uint8 = 18
)

// 区块奖励参数
const (
	// BlockReward - 基础区块奖励，8 NOGO
	BlockReward uint64 = 8

	// HalvingInterval - 减半间隔，5,200,000个区块
	HalvingInterval uint64 = 5200000

	// MinBlockReward - 最低区块奖励，0.1 NOGO
	MinBlockReward float64 = 0.1

	// RewardReductionRate - 每次减半的奖励减少率，20%
	RewardReductionRate float64 = 0.2
)

// Gas参数（兼容以太坊London硬分叉）
const (
	// GenesisGasLimit - 创世区块Gas限制
	GenesisGasLimit uint64 = 30000000

	// MaxGasLimit - 最大Gas限制
	MaxGasLimit uint64 = 100000000

	// MinGasPrice - 最低Gas价格
	MinGasPrice uint64 = 1

	// BaseFeeChangeDenominator - 基础费用变化分母（EIP-1559）
	BaseFeeChangeDenominator uint64 = 8

	// ElasticityMultiplier - 弹性乘数（EIP-1559）
	ElasticityMultiplier uint64 = 2

	// InitialBaseFee - 初始基础费用（EIP-1559）
	InitialBaseFee uint64 = 1000000000
)

// 共识参数
const (
	// TargetBlockTime - 目标区块时间，20秒
	TargetBlockTime uint64 = 20

	// DifficultyAdjustmentInterval - 难度调整间隔，10个区块
	DifficultyAdjustmentInterval uint64 = 10

	// MaxDifficultyAdjustment - 最大难度调整幅度，50%
	MaxDifficultyAdjustment float64 = 0.5
)

// 网络参数
const (
	// ProtocolVersion - 协议版本
	ProtocolVersion uint = 1

	// NetworkID - 网络ID，与ChainID相同
	NetworkID uint64 = ChainID
)

// 计算区块奖励
// blockNumber: 区块高度
// 返回: 区块奖励（单位：NOGO）
func CalculateBlockReward(blockNumber uint64) float64 {
	// 计算减半次数
	halvings := blockNumber / HalvingInterval

	// 初始奖励
	reward := float64(BlockReward)

	// 应用减半
	for i := uint64(0); i < halvings; i++ {
		// 每次减少20%
		reward *= (1 - RewardReductionRate)

		// 确保不低于最低奖励
		if reward < MinBlockReward {
			return MinBlockReward
		}
	}

	// 确保不低于最低奖励
	if reward < MinBlockReward {
		return MinBlockReward
	}

	return reward
}

// GetBlockRewardBigInt - 获取区块奖励的big.Int表示（单位：wei）
// blockNumber: 区块高度
// 返回: 区块奖励（单位：wei）
func GetBlockRewardBigInt(blockNumber uint64) *big.Int {
	reward := CalculateBlockReward(blockNumber)

	// 转换为wei
	wei := new(big.Float).Mul(
		new(big.Float).SetFloat64(reward),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(Decimals)), nil)),
	)

	// 转换为big.Int
	result := new(big.Int)
	wei.Int(result)

	return result
}
