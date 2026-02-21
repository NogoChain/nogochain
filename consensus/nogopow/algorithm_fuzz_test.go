package nogopow

import (
	"math/big"
	"testing"
)

// Fuzz测试Initialize函数
func FuzzNogoPow_Initialize(f *testing.F) {
	f.Add([]byte("test seed"))
	f.Add([]byte(""))
	f.Add([]byte("long seed with many characters to test different input sizes"))
	f.Add([]byte{0x00})
	f.Add([]byte{0xff, 0x00, 0xff, 0x00})

	f.Fuzz(func(t *testing.T, seed []byte) {
		pow := NewNogoPow()
		pow.Initialize(seed)

		if len(pow.cache) != CacheItems {
			t.Errorf("Cache size mismatch: expected %d, got %d", CacheItems, len(pow.cache))
		}

		if len(pow.dataset) != NumItems {
			t.Errorf("Dataset size mismatch: expected %d, got %d", NumItems, len(pow.dataset))
		}
	})
}

// Fuzz测试Hashimoto函数
func FuzzNogoPow_Hashimoto(f *testing.F) {
	f.Add([]byte("test header"), uint64(12345))
	f.Add([]byte(""), uint64(0))
	f.Add([]byte("long header with many characters"), uint64(999999999))
	f.Add([]byte{0x00}, uint64(1))
	f.Add([]byte{0xff, 0x00, 0xff, 0x00}, uint64(123456789))

	f.Fuzz(func(t *testing.T, header []byte, nonce uint64) {
		pow := NewNogoPow()
		pow.Initialize([]byte("test seed"))

		hash, mixDigest := pow.Hashimoto(header, nonce)

		if len(hash) != 32 {
			t.Errorf("Hash length mismatch: expected 32, got %d", len(hash))
		}

		if len(mixDigest) != 32 {
			t.Errorf("MixDigest length mismatch: expected 32, got %d", len(mixDigest))
		}
	})
}

// Fuzz测试Verify函数
func FuzzNogoPow_Verify(f *testing.F) {
	f.Add([]byte("test header"), uint64(12345))
	f.Add([]byte(""), uint64(0))
	f.Add([]byte("long header"), uint64(999999999))

	f.Fuzz(func(t *testing.T, header []byte, nonce uint64) {
		pow := NewNogoPow()
		pow.Initialize([]byte("test seed"))

		target := big.NewInt(1000000)
		target.Lsh(target, 240)

		result := pow.Verify(header, nonce, target)
		_ = result // 验证函数应该不会panic
	})
}

// Fuzz测试Mine函数
func FuzzNogoPow_Mine(f *testing.F) {
	f.Add([]byte("test header"), uint64(0), uint64(100))
	f.Add([]byte(""), uint64(100), uint64(50))

	f.Fuzz(func(t *testing.T, header []byte, startNonce uint64, iterations uint64) {
		pow := NewNogoPow()
		pow.Initialize([]byte("test seed"))

		target := big.NewInt(1000000)
		target.Lsh(target, 240)

		nonce, hash, mixDigest, found := pow.Mine(header, target, startNonce, iterations)
		_ = nonce
		_ = hash
		_ = mixDigest
		_ = found
	})
}

// Benchmark测试Initialize函数
func BenchmarkNogoPow_Initialize(b *testing.B) {
	pow := NewNogoPow()
	seed := []byte("test seed")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pow.Initialize(seed)
	}
}

// Benchmark测试Hashimoto函数
func BenchmarkNogoPow_Hashimoto(b *testing.B) {
	pow := NewNogoPow()
	pow.Initialize([]byte("test seed"))
	header := []byte("test header")
	nonce := uint64(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pow.Hashimoto(header, nonce)
	}
}

// Benchmark测试Mine函数
func BenchmarkNogoPow_Mine(b *testing.B) {
	pow := NewNogoPow()
	pow.Initialize([]byte("test seed"))
	header := []byte("test header")
	target := big.NewInt(1000000)
	target.Lsh(target, 240)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pow.Mine(header, target, 0, 100)
	}
}

// Benchmark测试MineParallel函数
func BenchmarkNogoPow_MineParallel(b *testing.B) {
	pow := NewNogoPow()
	pow.Initialize([]byte("test seed"))
	header := []byte("test header")
	target := big.NewInt(1000000)
	target.Lsh(target, 240)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pow.MineParallel(header, target, 1000)
	}
}

// Benchmark测试buildCache函数
func BenchmarkNogoPow_buildCache(b *testing.B) {
	pow := NewNogoPow()
	seed := []byte("test seed")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pow.buildCache(seed)
	}
}

// Benchmark测试buildDataset函数
func BenchmarkNogoPow_buildDataset(b *testing.B) {
	pow := NewNogoPow()
	pow.buildCache([]byte("test seed"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pow.buildDataset()
	}
}
