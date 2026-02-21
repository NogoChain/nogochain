package storage

import (
	"math/big"
)

// Storage 以太坊虚拟机存储实现
type Storage struct {
	// 存储映射
	data map[string][]byte
	// 快照用于回滚
	snapshots []map[string][]byte
}

// NewStorage 创建新的存储实例
func NewStorage() *Storage {
	return &Storage{
		data:      make(map[string][]byte),
		snapshots: make([]map[string][]byte, 0),
	}
}

// Get 获取存储值
func (s *Storage) Get(key []byte) []byte {
	keyStr := string(key)
	if value, exists := s.data[keyStr]; exists {
		return value
	}
	return make([]byte, 32) // 默认返回32字节的零值
}

// Set 设置存储值
func (s *Storage) Set(key, value []byte) {
	keyStr := string(key)
	// 确保值是32字节
	if len(value) != 32 {
		padded := make([]byte, 32)
		copySize := len(value)
		if copySize > 32 {
			copySize = 32
		}
		copy(padded, value[:copySize])
		value = padded
	}
	s.data[keyStr] = value
}

// Snapshot 创建存储快照
func (s *Storage) Snapshot() int {
	// 复制当前状态
	snapshot := make(map[string][]byte)
	for k, v := range s.data {
		copyValue := make([]byte, len(v))
		copy(copyValue, v)
		snapshot[k] = copyValue
	}
	s.snapshots = append(s.snapshots, snapshot)
	return len(s.snapshots) - 1
}

// RevertToSnapshot 回滚到指定快照
func (s *Storage) RevertToSnapshot(id int) {
	if id < 0 || id >= len(s.snapshots) {
		return
	}

	// 恢复快照状态
	snapshot := s.snapshots[id]
	s.data = make(map[string][]byte)
	for k, v := range snapshot {
		copyValue := make([]byte, len(v))
		copy(copyValue, v)
		s.data[k] = copyValue
	}

	// 清理后续快照
	s.snapshots = s.snapshots[:id+1]
}

// Clear 清空存储
func (s *Storage) Clear() {
	s.data = make(map[string][]byte)
	s.snapshots = make([]map[string][]byte, 0)
}

// Copy 复制存储
func (s *Storage) Copy() *Storage {
	result := NewStorage()
	for k, v := range s.data {
		copyValue := make([]byte, len(v))
		copy(copyValue, v)
		result.data[k] = copyValue
	}
	return result
}

// GetBigInt 获取存储值作为big.Int
func (s *Storage) GetBigInt(key []byte) *big.Int {
	value := s.Get(key)
	return new(big.Int).SetBytes(value)
}

// SetBigInt 设置存储值从big.Int
func (s *Storage) SetBigInt(key []byte, value *big.Int) {
	bytes := value.Bytes()
	s.Set(key, bytes)
}

// Exists 检查键是否存在
func (s *Storage) Exists(key []byte) bool {
	keyStr := string(key)
	_, exists := s.data[keyStr]
	return exists
}

// Delete 删除存储键
func (s *Storage) Delete(key []byte) {
	keyStr := string(key)
	delete(s.data, keyStr)
}

// Size 获取存储大小
func (s *Storage) Size() int {
	return len(s.data)
}
