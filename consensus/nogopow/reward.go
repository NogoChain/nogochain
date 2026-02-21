package nogopow

import (
	"math/big"

	"nogochain/params"
)

// CalculateReward 计算区块奖励
// blockNumber: 区块高度
// 返回: 区块奖励（单位：wei）
func CalculateReward(blockNumber uint64) *big.Int {
	// 计算减产次数
	reductions := blockNumber / params.HalvingInterval

	// 初始奖励：8 NOGO
	initialReward := big.NewInt(int64(params.BlockReward))
	// 转换为wei
	initialRewardWei := new(big.Int).Mul(
		initialReward,
		new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(params.Decimals)), nil),
	)

	// 计算当前奖励
	// 每次减少20%，相当于乘以4/5
	reward := new(big.Int).Set(initialRewardWei)
	for i := uint64(0); i < reductions; i++ {
		// 乘以4，再除以5
		reward = new(big.Int).Div(
			new(big.Int).Mul(reward, big.NewInt(4)),
			big.NewInt(5),
		)
	}

	// 检查是否低于最低奖励
	minRewardWei := calculateMinRewardWei()
	if reward.Cmp(minRewardWei) < 0 {
		return minRewardWei
	}

	return reward
}

// calculateMinRewardWei 计算最低奖励（单位：wei）
func calculateMinRewardWei() *big.Int {
	// 最低奖励：0.1 NOGO
	minRewardFloat := new(big.Float).Mul(
		new(big.Float).SetFloat64(params.MinBlockReward),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(params.Decimals)), nil)),
	)

	// 转换为big.Int
	minReward := new(big.Int)
	minRewardFloat.Int(minReward)

	return minReward
}

// GetRewardForBlock 获取指定区块的奖励
// blockNumber: 区块高度
// 返回: 区块奖励（单位：wei）
func GetRewardForBlock(blockNumber uint64) *big.Int {
	return CalculateReward(blockNumber)
}
