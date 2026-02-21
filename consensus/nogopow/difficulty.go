package nogopow

import (
	"math/big"
	"time"
)

const (
	// TargetBlockTime Target block time (seconds)
	// TargetBlockTime 目标区块时间（秒）
	TargetBlockTime = 20
	// DifficultyAdjustmentInterval Difficulty adjustment interval (blocks)
	// DifficultyAdjustmentInterval 难度调整间隔（区块数）
	DifficultyAdjustmentInterval = 10
	// MaxDifficultyAdjustment Maximum difficulty adjustment (±50%)
	// MaxDifficultyAdjustment 最大难度调整幅度（±50%）
	MaxDifficultyAdjustment = 2
	// MinDifficultyAdjustment Minimum difficulty adjustment (±50%)
	// MinDifficultyAdjustment 最小难度调整幅度（±50%）
	MinDifficultyAdjustment = 1
	// InitialDifficulty Initial difficulty
	// InitialDifficulty 初始难度
	InitialDifficulty = 1000000
)

// CalculateDifficulty Calculate new difficulty
// CalculateDifficulty 计算新难度
// 每10个区块调整一次，目标出块时间为20秒，限制调整幅度在±50%
func CalculateDifficulty(parentTimestamp time.Time, currentTimestamp time.Time, parentDifficulty *big.Int, height uint64) *big.Int {
	// Initial difficulty: First 10 blocks use fixed initial difficulty
	// 初始难度：前10个区块使用固定初始难度
	if height < DifficultyAdjustmentInterval {
		return big.NewInt(InitialDifficulty)
	}

	// Calculate actual block time: Current block time minus parent block time
	// 计算实际区块时间：当前区块时间减去父区块时间
	actualTime := currentTimestamp.Sub(parentTimestamp)
	// Calculate target block time: Adjustment interval * target block time
	// 计算目标区块时间：调整间隔 * 目标出块时间
	targetTime := time.Duration(DifficultyAdjustmentInterval*TargetBlockTime) * time.Second

	// Calculate time ratio: Actual time / target time
	// If actual time > target time, network hashrate is insufficient, need to decrease difficulty
	// If actual time < target time, network hashrate is too high, need to increase difficulty
	// 计算时间差比率：实际时间 / 目标时间
	// 如果实际时间大于目标时间，说明网络算力不足，需要降低难度
	// 如果实际时间小于目标时间，说明网络算力过高，需要提高难度
	timeRatio := float64(actualTime) / float64(targetTime)

	// Limit adjustment to ±50%
	// When timeRatio > 2, actual time is more than twice target time, maximum 50% difficulty decrease
	// When timeRatio < 1, actual time is less than target time, maximum 50% difficulty increase
	// 限制调整幅度在±50%
	// 当timeRatio > 2时，说明实际时间是目标时间的2倍以上，最多降低难度50%
	// 当timeRatio < 1时，说明实际时间不到目标时间，最多提高难度50%
	if timeRatio > MaxDifficultyAdjustment {
		timeRatio = MaxDifficultyAdjustment
	} else if timeRatio < MinDifficultyAdjustment {
		timeRatio = MinDifficultyAdjustment
	}

	// Calculate new difficulty: Parent difficulty * (1 / timeRatio)
	// Use big.Float for high precision calculation to avoid integer division precision loss
	// 计算新难度：父区块难度 * (1 / timeRatio)
	// 使用big.Float进行高精度计算，避免整数除法的精度损失
	parentDiff := new(big.Float).SetInt(parentDifficulty)
	adjustment := new(big.Float).SetFloat64(1 / timeRatio)
	newDiff := new(big.Float).Mul(parentDiff, adjustment)

	// Convert result to integer
	// 将计算结果转换为整数
	result := new(big.Int)
	newDiff.Int(result)

	// Ensure difficulty is at least 1, prevent calculation errors caused by difficulty 0
	// 确保难度至少为1，防止难度为0导致计算错误
	if result.Cmp(big.NewInt(1)) < 0 {
		return big.NewInt(1)
	}

	return result
}

// ToTarget Convert difficulty to target value
// In PoW algorithm, miners need to find a nonce such that block hash is less than target value
// Higher difficulty means smaller target value, making it harder to find valid nonce
// ToTarget 将难度转换为目标值
// 在PoW算法中，矿工需要找到一个nonce，使得区块哈希小于目标值
// 难度越高，目标值越小，找到符合条件的nonce就越困难
func ToTarget(difficulty *big.Int) *big.Int {
	// Handle invalid difficulty value
	// 处理无效难度值
	if difficulty.Cmp(big.NewInt(0)) <= 0 {
		return big.NewInt(0)
	}

	// Max target value (2^256 - 1): The largest possible hash value
	// 最大目标值 (2^256 - 1)：所有可能的哈希值中最大的那个
	maxTarget := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	// Calculate target value: maxTarget / difficulty
	// Higher difficulty means smaller target value
	// 计算目标值: maxTarget / difficulty
	// 难度越高，目标值越小
	target := new(big.Int).Div(maxTarget, difficulty)

	return target
}

// FromTarget Convert target value to difficulty
// Used to derive current difficulty from block header target value
// FromTarget 将目标值转换为难度
// 用于从区块头的目标值反推当前难度
func FromTarget(target *big.Int) *big.Int {
	// Handle invalid target value
	// 处理无效目标值
	if target.Cmp(big.NewInt(0)) <= 0 {
		return big.NewInt(0)
	}

	// Max target value (2^256 - 1)
	// 最大目标值 (2^256 - 1)
	maxTarget := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	// Calculate difficulty: maxTarget / target
	// Smaller target value means higher difficulty
	// 计算难度: maxTarget / target
	// 目标值越小，难度越高
	difficulty := new(big.Int).Div(maxTarget, target)

	return difficulty
}
