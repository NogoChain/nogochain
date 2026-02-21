package config

import (
	"math/big"
)

// 生产环境基础链参数
const (
	// MainChainID - 生产环境链ID
	MainChainID uint64 = 318

	// MainSymbol - 生产环境代币符号
	MainSymbol string = "NOGO"

	// MainDecimals - 生产环境代币小数位数
	MainDecimals uint8 = 18
)

// 生产环境区块奖励参数
const (
	// MainBlockReward - 生产环境基础区块奖励，8 NOGO
	MainBlockReward uint64 = 8

	// MainHalvingInterval - 生产环境减半间隔，5,000,000个区块
	MainHalvingInterval uint64 = 5000000

	// MainMinBlockReward - 生产环境最低区块奖励，0.1 NOGO
	MainMinBlockReward float64 = 0.1

	// MainRewardReductionRate - 生产环境每次减半的奖励减少率，20%
	MainRewardReductionRate float64 = 0.2
)

// 生产环境Gas参数
const (
	// MainGenesisGasLimit - 生产环境创世区块Gas限制
	MainGenesisGasLimit uint64 = 40000000

	// MainMaxGasLimit - 生产环境最大Gas限制
	MainMaxGasLimit uint64 = 120000000

	// MainMinGasPrice - 生产环境最低Gas价格
	MainMinGasPrice uint64 = 1

	// MainBaseFeeChangeDenominator - 生产环境基础费用变化分母（EIP-1559）
	MainBaseFeeChangeDenominator uint64 = 8

	// MainElasticityMultiplier - 生产环境弹性乘数（EIP-1559）
	MainElasticityMultiplier uint64 = 2

	// MainInitialBaseFee - 生产环境初始基础费用（EIP-1559）
	MainInitialBaseFee uint64 = 500000000
)

// 生产环境共识参数
const (
	// MainTargetBlockTime - 生产环境目标区块时间，20秒
	MainTargetBlockTime uint64 = 20

	// MainDifficultyAdjustmentInterval - 生产环境难度调整间隔，10个区块
	MainDifficultyAdjustmentInterval uint64 = 10

	// MainMaxDifficultyAdjustment - 生产环境最大难度调整幅度，50%
	MainMaxDifficultyAdjustment float64 = 0.5

	// MainInitialDifficulty - 生产环境初始难度
	MainInitialDifficulty uint64 = 10000
)

// 生产环境网络参数
const (
	// MainProtocolVersion - 生产环境协议版本
	MainProtocolVersion uint = 1

	// MainNetworkID - 生产环境网络ID，与MainChainID相同
	MainNetworkID uint64 = MainChainID
)

// CalculateMainBlockReward - 计算生产环境区块奖励
// blockNumber: 区块高度
// 返回: 区块奖励（单位：NOGO）
func CalculateMainBlockReward(blockNumber uint64) float64 {
	// 计算减半次数
	 halvings := blockNumber / MainHalvingInterval

	// 初始奖励
	reward := float64(MainBlockReward)

	// 应用减半
	for i := uint64(0); i < halvings; i++ {
		// 每次减少20%
		reward *= (1 - MainRewardReductionRate)

		// 确保不低于最低奖励
		if reward < MainMinBlockReward {
			return MainMinBlockReward
		}
	}

	// 确保不低于最低奖励
	if reward < MainMinBlockReward {
		return MainMinBlockReward
	}

	return reward
}

// GetMainBlockRewardBigInt - 获取生产环境区块奖励的big.Int表示（单位：wei）
// blockNumber: 区块高度
// 返回: 区块奖励（单位：wei）
func GetMainBlockRewardBigInt(blockNumber uint64) *big.Int {
	reward := CalculateMainBlockReward(blockNumber)

	// 转换为wei
	wei := new(big.Float).Mul(
		new(big.Float).SetFloat64(reward),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(MainDecimals)), nil)),
	)

	// 转换为big.Int
	result := new(big.Int)
	wei.Int(result)

	return result
}
