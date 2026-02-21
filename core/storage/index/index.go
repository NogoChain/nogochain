package index

import (
	"sync"
)

// Index 索引接口
 type Index interface {
	Add(key []byte, value []byte) error
	Get(key []byte) ([]byte, bool)
	Delete(key []byte) error
	Clear() error
	Size() int
}

// MemoryIndex 内存索引
 type MemoryIndex struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

// NewMemoryIndex 创建内存索引
func NewMemoryIndex() *MemoryIndex {
	return &MemoryIndex{
		data: make(map[string][]byte),
	}
}

// Add 添加索引项
func (i *MemoryIndex) Add(key []byte, value []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.data[string(key)] = value
	return nil
}

// Get 获取索引项
func (i *MemoryIndex) Get(key []byte) ([]byte, bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	value, exists := i.data[string(key)]
	return value, exists
}

// Delete 删除索引项
func (i *MemoryIndex) Delete(key []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	delete(i.data, string(key))
	return nil
}

// Clear 清空索引
func (i *MemoryIndex) Clear() error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.data = make(map[string][]byte)
	return nil
}

// Size 获取索引大小
func (i *MemoryIndex) Size() int {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return len(i.data)
}

// BlockIndex 区块索引
 type BlockIndex struct {
	numberIndex *MemoryIndex // 区块号到区块哈希的索引
	hashIndex   *MemoryIndex // 区块哈希到区块数据的索引
	mutex       sync.RWMutex
}

// NewBlockIndex 创建区块索引
func NewBlockIndex() *BlockIndex {
	return &BlockIndex{
		numberIndex: NewMemoryIndex(),
		hashIndex:   NewMemoryIndex(),
	}
}

// Add 添加区块索引
func (i *BlockIndex) Add(number uint64, hash []byte, data []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 添加区块号到哈希的索引
	numberKey := []byte{0} // 0 表示区块号索引
	numberKey = append(numberKey, uint64ToBytes(number)...)
	if err := i.numberIndex.Add(numberKey, hash); err != nil {
		return err
	}

	// 添加哈希到数据的索引
	hashKey := []byte{1} // 1 表示区块哈希索引
	hashKey = append(hashKey, hash...)
	if err := i.hashIndex.Add(hashKey, data); err != nil {
		return err
	}

	return nil
}

// GetByNumber 通过区块号获取区块
func (i *BlockIndex) GetByNumber(number uint64) ([]byte, bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// 获取区块哈希
	numberKey := []byte{0}
	numberKey = append(numberKey, uint64ToBytes(number)...)
	hash, exists := i.numberIndex.Get(numberKey)
	if !exists {
		return nil, false
	}

	// 获取区块数据
	hashKey := []byte{1}
	hashKey = append(hashKey, hash...)
	data, exists := i.hashIndex.Get(hashKey)
	if !exists {
		return nil, false
	}

	return data, true
}

// GetByHash 通过区块哈希获取区块
func (i *BlockIndex) GetByHash(hash []byte) ([]byte, bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// 获取区块数据
	hashKey := []byte{1}
	hashKey = append(hashKey, hash...)
	data, exists := i.hashIndex.Get(hashKey)
	if !exists {
		return nil, false
	}

	return data, true
}

// Delete 删除区块索引
func (i *BlockIndex) Delete(number uint64, hash []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 删除区块号索引
	numberKey := []byte{0}
	numberKey = append(numberKey, uint64ToBytes(number)...)
	if err := i.numberIndex.Delete(numberKey); err != nil {
		return err
	}

	// 删除区块哈希索引
	hashKey := []byte{1}
	hashKey = append(hashKey, hash...)
	if err := i.hashIndex.Delete(hashKey); err != nil {
		return err
	}

	return nil
}

// Clear 清空索引
func (i *BlockIndex) Clear() error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if err := i.numberIndex.Clear(); err != nil {
		return err
	}

	if err := i.hashIndex.Clear(); err != nil {
		return err
	}

	return nil
}

// Size 获取索引大小
func (i *BlockIndex) Size() int {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return i.numberIndex.Size()
}

// StateIndex 状态索引
 type StateIndex struct {
	accountIndex *MemoryIndex // 账户地址到账户数据的索引
	storageIndex *MemoryIndex // 账户存储到数据的索引
	mutex        sync.RWMutex
}

// NewStateIndex 创建状态索引
func NewStateIndex() *StateIndex {
	return &StateIndex{
		accountIndex: NewMemoryIndex(),
		storageIndex: NewMemoryIndex(),
	}
}

// AddAccount 添加账户索引
func (i *StateIndex) AddAccount(address []byte, data []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 添加账户地址到数据的索引
	accountKey := []byte{0} // 0 表示账户索引
	accountKey = append(accountKey, address...)
	return i.accountIndex.Add(accountKey, data)
}

// GetAccount 获取账户数据
func (i *StateIndex) GetAccount(address []byte) ([]byte, bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// 获取账户数据
	accountKey := []byte{0}
	accountKey = append(accountKey, address...)
	return i.accountIndex.Get(accountKey)
}

// AddStorage 添加存储索引
func (i *StateIndex) AddStorage(address []byte, key []byte, data []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 添加存储索引
	storageKey := []byte{1} // 1 表示存储索引
	storageKey = append(storageKey, address...)
	storageKey = append(storageKey, key...)
	return i.storageIndex.Add(storageKey, data)
}

// GetStorage 获取存储数据
func (i *StateIndex) GetStorage(address []byte, key []byte) ([]byte, bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// 获取存储数据
	storageKey := []byte{1}
	storageKey = append(storageKey, address...)
	storageKey = append(storageKey, key...)
	return i.storageIndex.Get(storageKey)
}

// DeleteAccount 删除账户索引
func (i *StateIndex) DeleteAccount(address []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 删除账户索引
	accountKey := []byte{0}
	accountKey = append(accountKey, address...)
	return i.accountIndex.Delete(accountKey)
}

// DeleteStorage 删除存储索引
func (i *StateIndex) DeleteStorage(address []byte, key []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 删除存储索引
	storageKey := []byte{1}
	storageKey = append(storageKey, address...)
	storageKey = append(storageKey, key...)
	return i.storageIndex.Delete(storageKey)
}

// Clear 清空索引
func (i *StateIndex) Clear() error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if err := i.accountIndex.Clear(); err != nil {
		return err
	}

	if err := i.storageIndex.Clear(); err != nil {
		return err
	}

	return nil
}

// Size 获取索引大小
func (i *StateIndex) Size() int {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return i.accountIndex.Size() + i.storageIndex.Size()
}

// uint64ToBytes 将 uint64 转换为字节数组
func uint64ToBytes(n uint64) []byte {
	return []byte{
		byte(n >> 56),
		byte(n >> 48),
		byte(n >> 40),
		byte(n >> 32),
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
}
