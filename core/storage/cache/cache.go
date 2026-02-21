package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// CacheItem 缓存项
 type CacheItem struct {
	Value      interface{}
	Expiration int64
	AccessTime int64
}

// MemoryCache 内存缓存
 type MemoryCache struct {
	data      map[string]*CacheItem
	mutex     sync.RWMutex
	capacity  int
	eviction  EvictionPolicy
	hits      int64
	misses    int64
	evictions int64
}

// EvictionPolicy 缓存淘汰策略
 type EvictionPolicy interface {
	SelectVictim(cache *MemoryCache) string
}

// LRUCache LRU缓存淘汰策略
 type LRUCache struct {}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(capacity int, eviction EvictionPolicy) *MemoryCache {
	if eviction == nil {
		eviction = &LRUCache{}
	}

	cache := &MemoryCache{
		data:     make(map[string]*CacheItem),
		capacity: capacity,
		eviction: eviction,
	}

	// 启动过期清理协程
	go cache.cleanExpired()

	return cache
}

// Get 获取缓存项
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		c.misses++
		return nil, false
	}

	// 检查是否过期
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		c.misses++
		return nil, false
	}

	// 更新访问时间
	item.AccessTime = time.Now().UnixNano()
	c.hits++
	return item.Value, true
}

// Set 设置缓存项
func (c *MemoryCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查容量
	if len(c.data) >= c.capacity {
		// 执行淘汰
		victim := c.eviction.SelectVictim(c)
		if victim != "" {
			delete(c.data, victim)
			c.evictions++
		}
	}

	// 设置缓存项
	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	c.data[key] = &CacheItem{
		Value:      value,
		Expiration: exp,
		AccessTime: time.Now().UnixNano(),
	}
}

// Delete 删除缓存项
func (c *MemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear 清空缓存
func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*CacheItem)
	c.hits = 0
	c.misses = 0
	c.evictions = 0
}

// GetStats 获取缓存统计信息
func (c *MemoryCache) GetStats() map[string]int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return map[string]int64{
		"hits":      c.hits,
		"misses":    c.misses,
		"evictions": c.evictions,
		"size":      int64(len(c.data)),
		"capacity":  int64(c.capacity),
	}
}

// GetHitRate 获取缓存命中率
func (c *MemoryCache) GetHitRate() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	total := c.hits + c.misses
	if total == 0 {
		return 0
	}

	return float64(c.hits) / float64(total)
}

// cleanExpired 清理过期缓存项
func (c *MemoryCache) cleanExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now().UnixNano()

		for key, item := range c.data {
			if item.Expiration > 0 && item.Expiration < now {
				delete(c.data, key)
			}
		}

		c.mutex.Unlock()
	}
}

// SelectVictim 选择淘汰的缓存项
func (lru *LRUCache) SelectVictim(cache *MemoryCache) string {
	var victimKey string
	var oldestAccessTime int64 = time.Now().UnixNano()

	for key, item := range cache.data {
		if item.AccessTime < oldestAccessTime {
			oldestAccessTime = item.AccessTime
			victimKey = key
		}
	}

	return victimKey
}

// DiskCache 磁盘缓存
 type DiskCache struct {
	path     string
	capacity int64
	mutex    sync.RWMutex
	size     int64
}

// NewDiskCache 创建磁盘缓存
func NewDiskCache(path string, capacity int64) *DiskCache {
	// 确保路径存在
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}

	return &DiskCache{
		path:     path,
		capacity: capacity,
		size:     0,
	}
}

// Get 获取缓存项
func (c *DiskCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// 构建文件路径
	filePath := filepath.Join(c.path, key)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, false
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	// 反序列化数据
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, false
	}

	return value, true
}

// Set 设置缓存项
func (c *DiskCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 序列化数据
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	// 检查容量
	filePath := filepath.Join(c.path, key)
	fileInfo, err := os.Stat(filePath)
	var fileSize int64
	if err == nil {
		fileSize = fileInfo.Size()
	}

	// 检查是否需要清理空间
	if c.size-fileSize+int64(len(data)) > c.capacity {
		c.evict()
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return
	}

	// 更新大小
	if err == nil {
		c.size = c.size - fileSize + int64(len(data))
	} else {
		c.size += int64(len(data))
	}
}

// Delete 删除缓存项
func (c *DiskCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 构建文件路径
	filePath := filepath.Join(c.path, key)

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return
	}

	// 更新大小
	c.size -= fileInfo.Size()
}

// Clear 清空缓存
func (c *DiskCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 读取目录内容
	files, err := os.ReadDir(c.path)
	if err != nil {
		return
	}

	// 删除所有文件
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(c.path, file.Name())
			if err := os.Remove(filePath); err != nil {
				continue
			}
		}
	}

	// 重置大小
	c.size = 0
}

// evict 淘汰缓存项
func (c *DiskCache) evict() {
	// 读取目录内容
	files, err := os.ReadDir(c.path)
	if err != nil {
		return
	}

	// 按修改时间排序
	type fileInfo struct {
		name       string
		modTime    time.Time
		size       int64
	}

	var fileInfos []fileInfo
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(c.path, file.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				continue
			}
			fileInfos = append(fileInfos, fileInfo{
				name:    file.Name(),
				modTime: info.ModTime(),
				size:    info.Size(),
			})
		}
	}

	// 按修改时间排序（最旧的在前）
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].modTime.Before(fileInfos[j].modTime)
	})

	// 淘汰最旧的文件，直到有足够的空间
	targetSize := c.capacity * 7 / 10 // 保留70%的容量
	for c.size > targetSize && len(fileInfos) > 0 {
		// 删除最旧的文件
		oldest := fileInfos[0]
		fileInfos = fileInfos[1:]

		filePath := filepath.Join(c.path, oldest.name)
		if err := os.Remove(filePath); err != nil {
			continue
		}

		// 更新大小
		c.size -= oldest.size
	}
}

// MultiLevelCache 多级缓存
 type MultiLevelCache struct {
	memoryCache *MemoryCache
	diskCache   *DiskCache
	mutex       sync.RWMutex
}

// NewMultiLevelCache 创建多级缓存
func NewMultiLevelCache(memoryCapacity int, diskPath string, diskCapacity int64) *MultiLevelCache {
	return &MultiLevelCache{
		memoryCache: NewMemoryCache(memoryCapacity, &LRUCache{}),
		diskCache:   NewDiskCache(diskPath, diskCapacity),
	}
}

// Get 获取缓存项
func (c *MultiLevelCache) Get(key string) (interface{}, bool) {
	// 先从内存缓存获取
	if value, exists := c.memoryCache.Get(key); exists {
		return value, true
	}

	// 再从磁盘缓存获取
	if value, exists := c.diskCache.Get(key); exists {
		// 提升到内存缓存
		c.memoryCache.Set(key, value, 10*time.Minute)
		return value, true
	}

	return nil, false
}

// Set 设置缓存项
func (c *MultiLevelCache) Set(key string, value interface{}, expiration time.Duration) {
	// 设置到内存缓存
	c.memoryCache.Set(key, value, expiration)

	// 设置到磁盘缓存
	c.diskCache.Set(key, value, expiration)
}

// Delete 删除缓存项
func (c *MultiLevelCache) Delete(key string) {
	// 从内存缓存删除
	c.memoryCache.Delete(key)

	// 从磁盘缓存删除
	c.diskCache.Delete(key)
}

// Clear 清空缓存
func (c *MultiLevelCache) Clear() {
	// 清空内存缓存
	c.memoryCache.Clear()

	// 清空磁盘缓存
	c.diskCache.Clear()
}

// GetStats 获取缓存统计信息
func (c *MultiLevelCache) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"memory": c.memoryCache.GetStats(),
		"hit_rate": c.memoryCache.GetHitRate(),
	}
}
