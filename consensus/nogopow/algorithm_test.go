package nogopow

import (
	"math/big"
	"testing"
	"time"
)

func TestNogoPow_Initialize(t *testing.T) {
	pow := NewNogoPow()
	seed := []byte("test seed")

	start := time.Now()
	pow.Initialize(seed)
	duration := time.Since(start)

	if duration > 500*time.Millisecond {
		t.Errorf("Initialize took too long: %v", duration)
	}

	if len(pow.cache) != CacheItems {
		t.Errorf("Cache size mismatch: expected %d, got %d", CacheItems, len(pow.cache))
	}

	if len(pow.dataset) != NumItems {
		t.Errorf("Dataset size mismatch: expected %d, got %d", NumItems, len(pow.dataset))
	}
}

func TestNogoPow_Hashimoto(t *testing.T) {
	pow := NewNogoPow()
	seed := []byte("test seed")
	pow.Initialize(seed)

	header := []byte("test header")
	nonce := uint64(12345)

	start := time.Now()
	hash, mixDigest := pow.Hashimoto(header, nonce)
	duration := time.Since(start)

	if duration > 15*time.Millisecond {
		t.Errorf("Hashimoto took too long: %v", duration)
	}

	if len(hash) != 32 {
		t.Errorf("Hash length mismatch: expected 32, got %d", len(hash))
	}

	if len(mixDigest) != 32 {
		t.Errorf("MixDigest length mismatch: expected 32, got %d", len(mixDigest))
	}
}

func TestNogoPow_Verify(t *testing.T) {
	pow := NewNogoPow()
	seed := []byte("test seed")
	pow.Initialize(seed)

	header := []byte("test header")
	nonce := uint64(12345)
	target := big.NewInt(1000000)
	target.Lsh(target, 240) // 增大目标值

	result := pow.Verify(header, nonce, target)
	if !result {
		t.Errorf("Verify failed unexpectedly")
	}
}

func TestNogoPow_MineParallel(t *testing.T) {
	pow := NewNogoPow()
	seed := []byte("test seed")
	pow.Initialize(seed)

	header := []byte("test header")
	target := big.NewInt(1000000)
	target.Lsh(target, 240) // 增大目标值

	start := time.Now()
	_, hash, mixDigest, found := pow.MineParallel(header, target, 1000)
	duration := time.Since(start)

	if duration > 100*time.Millisecond {
		t.Errorf("MineParallel took too long: %v", duration)
	}

	if found {
		if len(hash) != 32 {
			t.Errorf("Hash length mismatch: expected 32, got %d", len(hash))
		}
		if len(mixDigest) != 32 {
			t.Errorf("MixDigest length mismatch: expected 32, got %d", len(mixDigest))
		}
	}
}

func TestCalculateDifficulty(t *testing.T) {
	parentTime := time.Now().Add(-20 * time.Second)
	currentTime := time.Now()
	parentDiff := big.NewInt(1000000)
	height := uint64(10)

	difficulty := CalculateDifficulty(parentTime, currentTime, parentDiff, height)
	if difficulty.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("Invalid difficulty: %v", difficulty)
	}
}

func TestGetBlockReward(t *testing.T) {
	tests := []struct {
		height       uint64
		wantPositive bool
	}{
		{0, true},
		{1000000, true},
		{5000000, true},
		{10000000, true},
	}

	for _, tt := range tests {
		reward := GetBlockReward(tt.height)
		if tt.wantPositive && reward.Cmp(big.NewInt(0)) <= 0 {
			t.Errorf("GetBlockReward(%d) = %v, want positive value", tt.height, reward)
		}
	}
}
