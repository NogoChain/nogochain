package nogopow

import (
	"math/big"
	"testing"
)

func TestCalculateReward(t *testing.T) {
	testCases := []struct {
		blockNumber    uint64
		expectedReward *big.Int
		description    string
	}{
		{
			blockNumber:    0,
			expectedReward: big.NewInt(8).Mul(big.NewInt(8), big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)),
			description:    "创世区块奖励",
		},
		{
			blockNumber:    5200000 - 1,
			expectedReward: big.NewInt(8).Mul(big.NewInt(8), big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)),
			description:    "第一次减产前的最后一个区块",
		},
		{
			blockNumber:    5200000,
			expectedReward: big.NewInt(64).Mul(big.NewInt(64), big.NewInt(10).Exp(big.NewInt(10), big.NewInt(17), nil)),
			description:    "第一次减产后的第一个区块",
		},
		{
			blockNumber:    10400000,
			expectedReward: big.NewInt(512).Mul(big.NewInt(512), big.NewInt(10).Exp(big.NewInt(10), big.NewInt(16), nil)),
			description:    "第二次减产后的第一个区块",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualReward := CalculateReward(tc.blockNumber)
			if actualReward.Cmp(tc.expectedReward) != 0 {
				t.Errorf("区块高度 %d: 期望奖励 %s, 实际奖励 %s", tc.blockNumber, tc.expectedReward, actualReward)
			}
		})
	}
}

func TestGetRewardForBlock(t *testing.T) {
	blockNumber := uint64(12345)
	reward1 := CalculateReward(blockNumber)
	reward2 := GetRewardForBlock(blockNumber)

	if reward1.Cmp(reward2) != 0 {
		t.Errorf("GetRewardForBlock 与 CalculateReward 结果不一致")
	}
}
