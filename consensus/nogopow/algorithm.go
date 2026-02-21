package nogopow

import (
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"runtime"
	"sync"
)

const (
	// DatasetSize Dataset size limit is 1GB (optimization: reduce memory usage)
	// DatasetSize 数据集大小限制为1GB（优化：减少内存使用）
	DatasetSize = 1 * 1024 * 1024 * 1024
	// ItemSize Size of each data item
	// ItemSize 每个数据项大小
	ItemSize = 64
	// CacheSize Cache size
	// CacheSize 缓存大小
	CacheSize = 16 * 1024 * 1024
	// NumItems Number of data items
	// NumItems 数据项数量
	NumItems = DatasetSize / ItemSize
	// CacheItems Number of cache items
	// CacheItems 缓存项数量
	CacheItems = CacheSize / ItemSize
)

// NogoPow NogoPow algorithm structure
// NogoPow 算法结构体
type NogoPow struct {
	cache       []uint64
	dataset     []uint64
	cacheOnce   sync.Once
	datasetOnce sync.Once
}

// Global cache pool for reusing NogoPow instances
// 全局缓存池，用于重用NogoPow实例
var (
	powPool     = make(map[string]*NogoPow)
	powPoolMu   sync.RWMutex
	seedCache   = make(map[string][]uint64)
	seedCacheMu sync.RWMutex
)

// NewNogoPow Create a new NogoPow instance
// NewNogoPow 创建新的NogoPow实例
func NewNogoPow() *NogoPow {
	return &NogoPow{}
}

// GetCachedNogoPow Get NogoPow instance from cache
// GetCachedNogoPow 从缓存获取NogoPow实例
func GetCachedNogoPow(seedHash string) *NogoPow {
	powPoolMu.RLock()
	if pow, exists := powPool[seedHash]; exists {
		powPoolMu.RUnlock()
		return pow
	}
	powPoolMu.RUnlock()

	pow := NewNogoPow()
	powPoolMu.Lock()
	powPool[seedHash] = pow
	powPoolMu.Unlock()

	return pow
}

// Initialize Initialize the algorithm
// Initialize 初始化算法
func (n *NogoPow) Initialize(seed []byte) {
	n.cacheOnce.Do(func() {
		n.buildCache(seed)
	})
	n.datasetOnce.Do(func() {
		n.buildDataset()
	})
}

// buildCache Build cache
// buildCache 构建缓存
func (n *NogoPow) buildCache(seed []byte) {
	seedHash := sha256.Sum256(seed)
	seedHashStr := string(seedHash[:])

	// Check if already in cache
	// 检查缓存中是否已有
	seedCacheMu.RLock()
	if cachedCache, exists := seedCache[seedHashStr]; exists {
		n.cache = cachedCache
		seedCacheMu.RUnlock()
		return
	}
	seedCacheMu.RUnlock()

	n.cache = make([]uint64, CacheItems)
	h := seedHash
	for i := 0; i < 4; i++ {
		n.cache[0] ^= binary.LittleEndian.Uint64(h[i*8:])
	}

	for i := 1; i < CacheItems; i++ {
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, n.cache[i-1])
		h := sha256.Sum256(data)
		val := uint64(0)
		for j := 0; j < 4; j++ {
			val ^= binary.LittleEndian.Uint64(h[j*8:])
		}
		n.cache[i] = val
	}

	// Cache result
	// 缓存结果
	seedCacheMu.Lock()
	seedCache[seedHashStr] = n.cache
	seedCacheMu.Unlock()
}

// buildDataset Build dataset
// buildDataset 构建数据集
func (n *NogoPow) buildDataset() {
	n.dataset = make([]uint64, NumItems)
	threads := runtime.GOMAXPROCS(0)
	chunkSize := NumItems / threads
	if chunkSize < 1 {
		chunkSize = 1
	}

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + chunkSize
			if end > NumItems {
				end = NumItems
			}
			for j := start; j < end; j++ {
				n.dataset[j] = n.calculateDatasetItem(j)
			}
		}(i * chunkSize)
	}
	wg.Wait()
}

// calculateDatasetItem Calculate dataset item
// calculateDatasetItem 计算数据集项
func (n *NogoPow) calculateDatasetItem(index int) uint64 {
	result := n.cache[index%CacheItems]
	for i := 0; i < 2; i++ { // Optimization: reduce iteration count
		pos := (index ^ (i << 11)) % CacheItems
		if pos < 0 {
			pos = -pos
		}
		result ^= n.cache[pos]
	}
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, result)
	h := sha256.Sum256(data)
	val := uint64(0)
	for j := 0; j < 4; j++ {
		val ^= binary.LittleEndian.Uint64(h[j*8:])
	}
	return val
}

// Hashimoto Hashimoto algorithm
// Hashimoto 哈希imoto算法
func (n *NogoPow) Hashimoto(header []byte, nonce uint64) ([]byte, []byte) {
	data := append(header, make([]byte, 8)...)
	binary.LittleEndian.PutUint64(data[len(header):], nonce)

	mix := make([]uint64, 8)
	// Safely initialize mix array
	// 安全初始化mix数组
	for i := 0; i < 8; i++ {
		if i*8+8 <= len(data) {
			mix[i] = binary.LittleEndian.Uint64(data[i*8:])
		} else {
			mix[i] = 0
		}
	}

	for i := 0; i < 32; i++ { // Optimization: reduce iteration count
		pos := int(mix[i%8] % uint64(NumItems))
		for j := 0; j < 8; j++ {
			mix[j] ^= n.dataset[pos+j]
		}
	}

	result := make([]byte, 64)
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint64(result[i*8:], mix[i])
	}

	hash := sha256.Sum256(result)
	mixDigest := sha256.Sum256(hash[:])

	return hash[:], mixDigest[:]
}

// Verify Verify hash
// Verify 验证哈希
func (n *NogoPow) Verify(header []byte, nonce uint64, target *big.Int) bool {
	hash, _ := n.Hashimoto(header, nonce)
	hashInt := new(big.Int).SetBytes(hash)
	return hashInt.Cmp(target) <= 0
}

// Mine Mine
// Mine 挖矿
func (n *NogoPow) Mine(header []byte, target *big.Int, startNonce uint64, iterations uint64) (uint64, []byte, []byte, bool) {
	for nonce := startNonce; nonce < startNonce+iterations; nonce++ {
		hash, mixDigest := n.Hashimoto(header, nonce)
		hashInt := new(big.Int).SetBytes(hash)
		if hashInt.Cmp(target) <= 0 {
			return nonce, hash, mixDigest, true
		}
	}
	return 0, nil, nil, false
}

// MineParallel Parallel mining
// MineParallel 并行挖矿
func (n *NogoPow) MineParallel(header []byte, target *big.Int, iterations uint64) (uint64, []byte, []byte, bool) {
	threads := runtime.GOMAXPROCS(0)
	chunkSize := iterations / uint64(threads)
	if chunkSize < 1 {
		chunkSize = 1
	}

	var (
		wg             sync.WaitGroup
		mu             sync.Mutex
		foundNonce     uint64
		foundHash      []byte
		foundMixDigest []byte
		found          bool
	)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			startNonce := uint64(threadID) * chunkSize
			localNonce, localHash, localMixDigest, localFound := n.Mine(header, target, startNonce, chunkSize)
			if localFound {
				mu.Lock()
				if !found {
					foundNonce = localNonce
					foundHash = localHash
					foundMixDigest = localMixDigest
					found = true
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return foundNonce, foundHash, foundMixDigest, found
}
