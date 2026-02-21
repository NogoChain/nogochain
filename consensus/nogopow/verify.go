package nogopow

import (
	"crypto/sha256"
	"math/big"
	"time"
)

// VerifyBlock 验证区块
func VerifyBlock(header []byte, nonce uint64, difficulty *big.Int, timestamp time.Time, parentTimestamp time.Time, height uint64) bool {
	// 验证时间戳
	if timestamp.Before(parentTimestamp) {
		return false
	}

	// 验证难度
	target := ToTarget(difficulty)
	
	// 优化：使用缓存的NogoPow实例
	seedHash := sha256.Sum256(header)
	seedHashStr := string(seedHash[:])
	pow := GetCachedNogoPow(seedHashStr)
	pow.Initialize(header)

	// 验证工作量证明
	return pow.Verify(header, nonce, target)
}

// VerifyDifficulty 验证难度调整
func VerifyDifficulty(parentDifficulty *big.Int, currentDifficulty *big.Int, parentTimestamp time.Time, currentTimestamp time.Time, height uint64) bool {
	expectedDifficulty := CalculateDifficulty(parentTimestamp, currentTimestamp, parentDifficulty, height)
	return expectedDifficulty.Cmp(currentDifficulty) == 0
}

// VerifyNonce 验证nonce
func VerifyNonce(header []byte, nonce uint64, target *big.Int) bool {
	// 优化：使用缓存的NogoPow实例
	seedHash := sha256.Sum256(header)
	seedHashStr := string(seedHash[:])
	pow := GetCachedNogoPow(seedHashStr)
	pow.Initialize(header)
	return pow.Verify(header, nonce, target)
}

// GetBlockReward 获取区块奖励
func GetBlockReward(height uint64) *big.Int {
	// 基础奖励: 8 NOGO
	baseReward := big.NewInt(8)
	baseReward.Mul(baseReward, big.NewInt(1000000000000000000)) // 转换为wei

	// 每500万区块减半20%
	halfLife := uint64(5000000)
	decreaseCount := height / halfLife

	if decreaseCount == 0 {
		return baseReward
	}

	// 计算奖励衰减
	reward := new(big.Int).Set(baseReward)
	for i := uint64(0); i < decreaseCount; i++ {
		// 减少20%
		reward.Mul(reward, big.NewInt(8))
		reward.Div(reward, big.NewInt(10))
	}

	// 最小奖励: 0.1 NOGO
	minReward := big.NewInt(1)
	minReward.Mul(minReward, big.NewInt(100000000000000000)) // 转换为wei

	if reward.Cmp(minReward) < 0 {
		return minReward
	}

	return reward
}
