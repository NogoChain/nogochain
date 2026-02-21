package storage

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"nogochain/core/storage/cache"
	"nogochain/core/storage/compression"
)

// StorageType 存储类型
 type StorageType int

const (
	// HotStorage 热数据存储
	HotStorage StorageType = iota
	// ColdStorage 冷数据存储
	ColdStorage
)

// StorageItem 存储项
 type StorageItem struct {
	Value      interface{}
	AccessTime time.Time
	StorageType StorageType
	Compressed bool
}

// Storage 存储接口
 type Storage interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Clear()
	GetStats() map[string]interface{}
}

// OptimizedStorage 优化的存储实现
 type OptimizedStorage struct {
	hotStorage     *cache.MemoryCache
	coldStorage    *cache.DiskCache
	compressor     compression.Compressor
	mutex          sync.RWMutex
	hotToColdThreshold time.Duration
	coldAccessCount   map[string]int
	coldAccessMutex   sync.Mutex
}

// NewOptimizedStorage 创建优化的存储
func NewOptimizedStorage(dataDir string, hotCapacity int, coldCapacity int64, hotToColdThreshold time.Duration) *OptimizedStorage {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic(err)
	}

	// 创建热数据存储（内存缓存）
	hotStorage := cache.NewMemoryCache(hotCapacity, &cache.LRUCache{})

	// 创建冷数据存储（磁盘缓存）
	coldPath := filepath.Join(dataDir, "cold")
	if err := os.MkdirAll(coldPath, 0755); err != nil {
		panic(err)
	}
	coldStorage := cache.NewDiskCache(coldPath, coldCapacity)

	// 创建压缩器
	compressor := compression.NewGzipCompressor(compression.DefaultCompression)

	return &OptimizedStorage{
		hotStorage:         hotStorage,
		coldStorage:        coldStorage,
		compressor:         compressor,
		hotToColdThreshold: hotToColdThreshold,
		coldAccessCount:    make(map[string]int),
	}
}

// Get 获取存储项
func (s *OptimizedStorage) Get(key string) (interface{}, bool) {
	// 先从热数据存储获取
	if value, exists := s.hotStorage.Get(key); exists {
		return value, true
	}

	// 再从冷数据存储获取
	if value, exists := s.coldStorage.Get(key); exists {
		// 增加访问计数
		s.coldAccessMutex.Lock()
		s.coldAccessCount[key]++
		count := s.coldAccessCount[key]
		s.coldAccessMutex.Unlock()

		// 如果访问次数超过阈值，迁移到热数据存储
		if count > 3 {
			s.migrateToHot(key, value)
			s.coldAccessMutex.Lock()
			delete(s.coldAccessCount, key)
			s.coldAccessMutex.Unlock()
		}

		return value, true
	}

	return nil, false
}

// Set 设置存储项
func (s *OptimizedStorage) Set(key string, value interface{}) {
	// 先存储到热数据存储
	s.hotStorage.Set(key, value, 0)

	// 启动后台任务，检查是否需要迁移到冷数据存储
	go s.checkHotToColdMigration()
}

// Delete 删除存储项
func (s *OptimizedStorage) Delete(key string) {
	// 从热数据存储删除
	s.hotStorage.Delete(key)

	// 从冷数据存储删除
	s.coldStorage.Delete(key)

	// 从冷数据访问计数中删除
	s.coldAccessMutex.Lock()
	delete(s.coldAccessCount, key)
	s.coldAccessMutex.Unlock()
}

// Clear 清空存储
func (s *OptimizedStorage) Clear() {
	// 清空热数据存储
	s.hotStorage.Clear()

	// 清空冷数据存储
	s.coldStorage.Clear()

	// 清空冷数据访问计数
	s.coldAccessMutex.Lock()
	s.coldAccessCount = make(map[string]int)
	s.coldAccessMutex.Unlock()
}

// GetStats 获取存储统计信息
func (s *OptimizedStorage) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"hot_storage":  s.hotStorage.GetStats(),
		"cache_hit_rate": s.hotStorage.GetHitRate(),
		"cold_access_count": len(s.coldAccessCount),
	}
}

// checkHotToColdMigration 检查是否需要将热数据迁移到冷数据存储
func (s *OptimizedStorage) checkHotToColdMigration() {
	// 实现热数据到冷数据的迁移逻辑
	// 这里简化实现，实际应该定期检查热数据的访问时间
}

// migrateToHot 将冷数据迁移到热数据存储
func (s *OptimizedStorage) migrateToHot(key string, value interface{}) {
	// 将数据从冷数据存储迁移到热数据存储
	s.hotStorage.Set(key, value, 0)
	s.coldStorage.Delete(key)
}

// migrateToCold 将热数据迁移到冷数据存储
func (s *OptimizedStorage) migrateToCold(key string, value interface{}) {
	// 压缩数据
	compressed, err := compression.CompressObject(value, s.compressor)
	if err != nil {
		// 压缩失败，直接存储原始数据
		s.coldStorage.Set(key, value, 0)
		return
	}

	// 存储压缩后的数据
	s.coldStorage.Set(key, compressed, 0)
	s.hotStorage.Delete(key)
}

// BlockStorage 区块存储
 type BlockStorage struct {
	storage *OptimizedStorage
	mutex   sync.RWMutex
}

// NewBlockStorage 创建区块存储
func NewBlockStorage(dataDir string) *BlockStorage {
	return &BlockStorage{
		storage: NewOptimizedStorage(
			filepath.Join(dataDir, "blocks"),
			1000, // 热数据容量
			1024*1024*1024, // 冷数据容量（1GB）
			24*time.Hour, // 热到冷的迁移阈值
		),
	}
}

// Get 获取区块
func (s *BlockStorage) Get(key string) (interface{}, bool) {
	return s.storage.Get(key)
}

// Set 设置区块
func (s *BlockStorage) Set(key string, value interface{}) {
	s.storage.Set(key, value)
}

// Delete 删除区块
func (s *BlockStorage) Delete(key string) {
	s.storage.Delete(key)
}

// Clear 清空区块存储
func (s *BlockStorage) Clear() {
	s.storage.Clear()
}

// GetStats 获取区块存储统计信息
func (s *BlockStorage) GetStats() map[string]interface{} {
	return s.storage.GetStats()
}

// StateStorage 状态存储
 type StateStorage struct {
	storage *OptimizedStorage
	mutex   sync.RWMutex
}

// NewStateStorage 创建状态存储
func NewStateStorage(dataDir string) *StateStorage {
	return &StateStorage{
		storage: NewOptimizedStorage(
			filepath.Join(dataDir, "state"),
			5000, // 热数据容量
			2*1024*1024*1024, // 冷数据容量（2GB）
			12*time.Hour, // 热到冷的迁移阈值
		),
	}
}

// Get 获取状态
func (s *StateStorage) Get(key string) (interface{}, bool) {
	return s.storage.Get(key)
}

// Set 设置状态
func (s *StateStorage) Set(key string, value interface{}) {
	s.storage.Set(key, value)
}

// Delete 删除状态
func (s *StateStorage) Delete(key string) {
	s.storage.Delete(key)
}

// Clear 清空状态存储
func (s *StateStorage) Clear() {
	s.storage.Clear()
}

// GetStats 获取状态存储统计信息
func (s *StateStorage) GetStats() map[string]interface{} {
	return s.storage.GetStats()
}
