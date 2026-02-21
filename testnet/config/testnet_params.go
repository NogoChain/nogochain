package config

import (
	"math/big"
)

// 测试网络基础链参数
const (
	// TestChainID - 测试网络链ID
	TestChainID uint64 = 31888

	// TestSymbol - 测试网络代币符号
	TestSymbol string = "TESTNOGO"

	// TestDecimals - 测试网络代币小数位数
	TestDecimals uint8 = 18
)

// 测试网络区块奖励参数
const (
	// TestBlockReward - 测试网络基础区块奖励，20 TESTNOGO
	TestBlockReward uint64 = 20

	// TestHalvingInterval - 测试网络减半间隔，10,000个区块
	TestHalvingInterval uint64 = 10000

	// TestMinBlockReward - 测试网络最低区块奖励，0.5 TESTNOGO
	TestMinBlockReward float64 = 0.5

	// TestRewardReductionRate - 测试网络每次减半的奖励减少率，25%
	TestRewardReductionRate float64 = 0.25
)

// 测试网络Gas参数
const (
	// TestGenesisGasLimit - 测试网络创世区块Gas限制
	TestGenesisGasLimit uint64 = 40000000

	// TestMaxGasLimit - 测试网络最大Gas限制
	TestMaxGasLimit uint64 = 120000000

	// TestMinGasPrice - 测试网络最低Gas价格
	TestMinGasPrice uint64 = 1

	// TestBaseFeeChangeDenominator - 测试网络基础费用变化分母（EIP-1559）
	TestBaseFeeChangeDenominator uint64 = 8

	// TestElasticityMultiplier - 测试网络弹性乘数（EIP-1559）
	TestElasticityMultiplier uint64 = 2

	// TestInitialBaseFee - 测试网络初始基础费用（EIP-1559）
	TestInitialBaseFee uint64 = 500000000
)

// 测试网络共识参数
const (
	// TestTargetBlockTime - 测试网络目标区块时间，10秒
	TestTargetBlockTime uint64 = 10

	// TestDifficultyAdjustmentInterval - 测试网络难度调整间隔，5个区块
	TestDifficultyAdjustmentInterval uint64 = 5

	// TestMaxDifficultyAdjustment - 测试网络最大难度调整幅度，75%
	TestMaxDifficultyAdjustment float64 = 0.75

	// TestInitialDifficulty - 测试网络初始难度，设置较低以便快速出块
	TestInitialDifficulty uint64 = 1000
)

// 测试网络网络参数
const (
	// TestProtocolVersion - 测试网络协议版本
	TestProtocolVersion uint = 1

	// TestNetworkID - 测试网络ID，与TestChainID相同
	TestNetworkID uint64 = TestChainID
)

// CalculateTestBlockReward - 计算测试网络区块奖励
// blockNumber: 区块高度
// 返回: 区块奖励（单位：TESTNOGO）
func CalculateTestBlockReward(blockNumber uint64) float64 {
	// 计算减半次数
	halvings := blockNumber / TestHalvingInterval

	// 初始奖励
	reward := float64(TestBlockReward)

	// 应用减半
	for i := uint64(0); i < halvings; i++ {
		// 每次减少25%
		reward *= (1 - TestRewardReductionRate)

		// 确保不低于最低奖励
		if reward < TestMinBlockReward {
			return TestMinBlockReward
		}
	}

	// 确保不低于最低奖励
	if reward < TestMinBlockReward {
		return TestMinBlockReward
	}

	return reward
}

// GetTestBlockRewardBigInt - 获取测试网络区块奖励的big.Int表示（单位：wei）
// blockNumber: 区块高度
// 返回: 区块奖励（单位：wei）
func GetTestBlockRewardBigInt(blockNumber uint64) *big.Int {
	reward := CalculateTestBlockReward(blockNumber)

	// 转换为wei
	wei := new(big.Float).Mul(
		new(big.Float).SetFloat64(reward),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(TestDecimals)), nil)),
	)

	// 转换为big.Int
	result := new(big.Int)
	wei.Int(result)

	return result
}