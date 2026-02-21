package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"nogochain/interfaces"
)

// Config 发现配置
type Config struct {
	Enabled   bool
	TTL       int
	Interval  int
	Bootnodes []string
}

// Node 节点信息
type Node struct {
	ID       interfaces.NodeID
	Addr     net.Addr
	LastSeen time.Time
}

// Discovery 节点发现服务
type Discovery struct {
	localNode  *Node
	cfg        *Config
	nodes      map[string]*Node
	nodesMutex sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	isRunning  bool
}

// NewDiscovery 创建新的发现服务
func NewDiscovery(localNode interface{}, cfg *Config) *Discovery {
	if cfg == nil {
		cfg = &Config{
			Enabled:   true,
			TTL:       30,
			Interval:  30,
			Bootnodes: []string{},
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Discovery{
		localNode: &Node{ID: interfaces.NodeID{}, LastSeen: time.Now()},
		cfg:       cfg,
		nodes:     make(map[string]*Node),
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
	}
}

// Start 启动发现服务
func (d *Discovery) Start() error {
	if !d.cfg.Enabled {
		log.Println("Discovery disabled")
		return nil
	}

	if d.isRunning {
		return nil
	}

	// 启动节点表维护
	go d.maintainNodeTable()

	// 连接到引导节点
	go d.connectToBootnodes()

	d.isRunning = true
	log.Println("Discovery service started")
	return nil
}

// Stop 停止发现服务
func (d *Discovery) Stop() error {
	if !d.isRunning {
		return nil
	}

	d.cancel()
	d.isRunning = false
	log.Println("Discovery service stopped")
	return nil
}

// maintainNodeTable 维护节点表
func (d *Discovery) maintainNodeTable() {
	ticker := time.NewTicker(time.Duration(d.cfg.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.refreshNodeTable()
		}
	}
}

// refreshNodeTable 刷新节点表
func (d *Discovery) refreshNodeTable() {
	// 检查节点是否过期
	d.nodesMutex.Lock()
	for id, node := range d.nodes {
		// 检查节点是否过期
		if time.Since(node.LastSeen) > time.Duration(d.cfg.TTL)*time.Second {
			delete(d.nodes, id)
			log.Printf("Removed expired node: %s", id)
		}
	}
	d.nodesMutex.Unlock()

	// 随机查找节点以保持节点表新鲜
	if len(d.nodes) < 10 {
		// 模拟查找节点
		d.simulateFindNodes()
	}
}

// simulateFindNodes 模拟查找节点
func (d *Discovery) simulateFindNodes() {
	// 模拟发现新节点
	for i := 0; i < 5; i++ {
		node := &Node{
			ID:       interfaces.NodeID{},
			Addr:     &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 30303 + i},
			LastSeen: time.Now(),
		}
		d.AddNode(node)
	}
}

// connectToBootnodes 连接到引导节点
func (d *Discovery) connectToBootnodes() {
	for _, bootnode := range d.cfg.Bootnodes {
		go func(nodeAddr string) {
			// 连接到引导节点
			log.Printf("Connecting to bootnode: %s", nodeAddr)
		}(bootnode)
	}
}

// FindNodes 查找节点
func (d *Discovery) FindNodes(target interfaces.NodeID, count int) []*Node {
	// 实现节点查找逻辑
	var nodes []*Node
	d.nodesMutex.RLock()
	for _, node := range d.nodes {
		nodes = append(nodes, node)
		if len(nodes) >= count {
			break
		}
	}
	d.nodesMutex.RUnlock()
	return nodes
}

// AddNode 添加节点
func (d *Discovery) AddNode(node *Node) {
	// 添加节点到节点表
	nodeID := fmt.Sprintf("%x", node.ID)
	d.nodesMutex.Lock()
	d.nodes[nodeID] = node
	d.nodesMutex.Unlock()
	log.Printf("Added node: %s", nodeID)
}

// RemoveNode 移除节点
func (d *Discovery) RemoveNode(node *Node) {
	// 从节点表移除节点
	nodeID := fmt.Sprintf("%x", node.ID)
	d.nodesMutex.Lock()
	delete(d.nodes, nodeID)
	d.nodesMutex.Unlock()
	log.Printf("Removed node: %s", nodeID)
}

// GetNodes 获取所有节点
func (d *Discovery) GetNodes() []*Node {
	// 获取节点表中的所有节点
	d.nodesMutex.RLock()
	defer d.nodesMutex.RUnlock()
	nodes := make([]*Node, 0, len(d.nodes))
	for _, node := range d.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// UpdateNode 更新节点信息
func (d *Discovery) UpdateNode(node *Node) {
	// 更新节点信息
	nodeID := fmt.Sprintf("%x", node.ID)
	d.nodesMutex.Lock()
	d.nodes[nodeID] = node
	d.nodesMutex.Unlock()
	log.Printf("Updated node: %s", nodeID)
}

// LocalNode 获取本地节点
func (d *Discovery) LocalNode() *Node {
	return d.localNode
}

// IsRunning 检查服务是否运行
func (d *Discovery) IsRunning() bool {
	return d.isRunning
}
