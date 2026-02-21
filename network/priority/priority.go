package priority

import (
	"container/heap"
	"sync"

	"nogochain/network/types"
)

// PriorityLevel 优先级级别
type PriorityLevel int

const (
	LowPriority    PriorityLevel = 0
	NormalPriority PriorityLevel = 1
	HighPriority   PriorityLevel = 2
)

// PriorityMsg 带优先级的消息
type PriorityMsg struct {
	Msg      *types.Msg
	Priority PriorityLevel
	Index    int
}

// PriorityQueue 优先级队列
type PriorityQueue []*PriorityMsg

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	msg := x.(*PriorityMsg)
	msg.Index = n
	*pq = append(*pq, msg)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	msg := old[n-1]
	msg.Index = -1
	*pq = old[0 : n-1]
	return msg
}

// PriorityManager 优先级管理器
type PriorityManager struct {
	queues map[string]*PriorityQueue
	mu     sync.Mutex
}

// NewPriorityManager 创建新的优先级管理器
func NewPriorityManager() *PriorityManager {
	return &PriorityManager{
		queues: make(map[string]*PriorityQueue),
	}
}

// AddMsg 添加消息到优先级队列
func (pm *PriorityManager) AddMsg(peerID string, msg *types.Msg, priority PriorityLevel) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.queues[peerID]; !exists {
		pm.queues[peerID] = &PriorityQueue{}
		heap.Init(pm.queues[peerID])
	}

	heap.Push(pm.queues[peerID], &PriorityMsg{
		Msg:      msg,
		Priority: priority,
	})
}

// GetNextMsg 获取下一个优先级最高的消息
func (pm *PriorityManager) GetNextMsg(peerID string) *types.Msg {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	queue, exists := pm.queues[peerID]
	if !exists || queue.Len() == 0 {
		return nil
	}

	msg := heap.Pop(queue).(*PriorityMsg)
	return msg.Msg
}

// HasMsg 检查是否有消息
func (pm *PriorityManager) HasMsg(peerID string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	queue, exists := pm.queues[peerID]
	if !exists {
		return false
	}

	return queue.Len() > 0
}
