package memory

import ()

// Memory 以太坊虚拟机内存实现
type Memory struct {
	data     []byte
	capacity int
}

// NewMemory 创建新的内存实例
func NewMemory() *Memory {
	return &Memory{
		data:     make([]byte, 0),
		capacity: 0,
	}
}

// Size 获取当前内存大小
func (m *Memory) Size() int {
	return len(m.data)
}

// Capacity 获取内存容量
func (m *Memory) Capacity() int {
	return m.capacity
}

// Resize 调整内存大小
func (m *Memory) Resize(size int) {
	if size > m.capacity {
		// 扩展内存
		newCapacity := size
		if newCapacity < m.capacity*2 {
			newCapacity = m.capacity * 2
		}
		if newCapacity < 32 {
			newCapacity = 32
		}

		newData := make([]byte, newCapacity)
		copy(newData, m.data)
		m.data = newData
		m.capacity = newCapacity
	}
}

// Set 设置内存数据
func (m *Memory) Set(offset int, data []byte) {
	end := offset + len(data)
	if end > m.capacity {
		m.Resize(end)
	}

	copy(m.data[offset:end], data)
}

// SetByte 设置单个字节
func (m *Memory) SetByte(offset int, value byte) {
	if offset >= m.capacity {
		m.Resize(offset + 1)
	}
	m.data[offset] = value
}

// Get 获取内存数据
func (m *Memory) Get(offset, size int) []byte {
	end := offset + size
	if end > len(m.data) {
		// 返回零填充的数据
		result := make([]byte, size)
		copySize := len(m.data) - offset
		if copySize > 0 {
			copy(result, m.data[offset:])
		}
		return result
	}

	return m.data[offset:end]
}

// GetByte 获取单个字节
func (m *Memory) GetByte(offset int) byte {
	if offset >= len(m.data) {
		return 0
	}
	return m.data[offset]
}

// Load32Bytes 加载32字节数据
func (m *Memory) Load32Bytes(offset int) []byte {
	return m.Get(offset, 32)
}

// Store32Bytes 存储32字节数据
func (m *Memory) Store32Bytes(offset int, data []byte) {
	if len(data) != 32 {
		// 填充或截断到32字节
		padded := make([]byte, 32)
		copySize := len(data)
		if copySize > 32 {
			copySize = 32
		}
		copy(padded, data[:copySize])
		data = padded
	}
	m.Set(offset, data)
}

// Reset 重置内存
func (m *Memory) Reset() {
	m.data = make([]byte, 0)
	m.capacity = 0
}

// Copy 复制内存数据
func (m *Memory) Copy() []byte {
	result := make([]byte, len(m.data))
	copy(result, m.data)
	return result
}

// CalculateGasCost 计算内存扩展的Gas成本
func (m *Memory) CalculateGasCost(size int) uint64 {
	if size <= m.capacity {
		return 0
	}

	// 计算新的内存大小（按32字节对齐）
	newSize := (size + 31) / 32
	oldSize := m.capacity / 32

	// 内存Gas计算：3 * newSize^2 / 512 - 3 * oldSize^2 / 512
	var cost uint64
	if newSize > 0 {
		cost = uint64(3 * newSize * newSize / 512)
	}
	if oldSize > 0 {
		cost -= uint64(3 * oldSize * oldSize / 512)
	}

	return cost
}
