package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/interfaces"
	"nogochain/metrics"
	"nogochain/network/config"
	"nogochain/rpc"
)

// Network 网络管理器
type Network struct {
	config     *config.Config
	blockchain *blockchain.Blockchain
	peers      map[string]*interfaces.Peer
	peersMutex sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	isStarted  bool
	rpcServer  *rpc.Server
}

// NewNetwork 创建新的网络管理器
func NewNetwork(cfg *config.Config, bc *blockchain.Blockchain) *Network {
	ctx, cancel := context.WithCancel(context.Background())

	network := &Network{
		config:     cfg,
		blockchain: bc,
		peers:      make(map[string]*interfaces.Peer),
		ctx:        ctx,
		cancel:     cancel,
		isStarted:  false,
	}

	// 初始化RPC服务器
	if cfg.RPC.Enabled {
		network.rpcServer = rpc.NewServer(cfg.RPC)
	}

	return network
}

// Start 启动网络
func (n *Network) Start() error {
	if n.isStarted {
		return fmt.Errorf("network already started")
	}

	// 启动监听
	go n.startListener()

	// 启动节点发现
	go n.startDiscovery()

	// 启动RPC服务器
	if n.config.RPC.Enabled && n.rpcServer != nil {
		go func() {
			if err := n.rpcServer.Start(); err != nil {
				log.Printf("RPC server failed to start: %v", err)
			}
		}()
	}

	n.isStarted = true
	log.Println("Network started successfully")
	return nil
}

// Stop 停止网络
func (n *Network) Stop() error {
	if !n.isStarted {
		return fmt.Errorf("network not started")
	}

	// 停止RPC服务器
	if n.rpcServer != nil {
		if err := n.rpcServer.Stop(); err != nil {
			log.Printf("Error stopping RPC server: %v", err)
		}
	}

	n.cancel()
	n.isStarted = false
	log.Println("Network stopped successfully")
	return nil
}

// startListener 启动监听器
func (n *Network) startListener() {
	addr := fmt.Sprintf(":%d", n.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Failed to start listener: %v", err)
		return
	}
	defer listener.Close()

	log.Printf("Listening for connections on %s", addr)

	for n.isStarted {
		conn, err := listener.Accept()
		if err != nil {
			if !n.isStarted {
				return
			}
			log.Printf("Accept error: %v", err)
			continue
		}

		// 处理连接
		go n.handleConnection(conn)
	}
}

// startDiscovery 启动节点发现
func (n *Network) startDiscovery() {
	log.Println("Starting node discovery")
	// 实现节点发现逻辑
}

// handleConnection 处理新连接
func (n *Network) handleConnection(conn net.Conn) {
	log.Printf("New connection from %s", conn.RemoteAddr())
	// 实现连接处理逻辑
}

// AddPeer 添加对等节点
func (n *Network) AddPeer(peer *interfaces.Peer) {
	peerID := fmt.Sprintf("%x", peer.ID)
	n.peersMutex.Lock()
	defer n.peersMutex.Unlock()
	n.peers[peerID] = peer
	log.Printf("Added peer: %s", peerID)
	// 更新对等节点数量指标
	metrics.PeerCount.Set(float64(len(n.peers)))
}

// RemovePeer 移除对等节点
func (n *Network) RemovePeer(id string) {
	n.peersMutex.Lock()
	defer n.peersMutex.Unlock()
	delete(n.peers, id)
	log.Printf("Removed peer: %s", id)
	// 更新对等节点数量指标
	metrics.PeerCount.Set(float64(len(n.peers)))
}

// GetPeers 获取所有对等节点
func (n *Network) GetPeers() []*interfaces.Peer {
	n.peersMutex.RLock()
	defer n.peersMutex.RUnlock()
	peers := make([]*interfaces.Peer, 0, len(n.peers))
	for _, peer := range n.peers {
		peers = append(peers, peer)
	}
	return peers
}

// GetPeer 通过ID获取对等节点
func (n *Network) GetPeer(id string) *interfaces.Peer {
	n.peersMutex.RLock()
	defer n.peersMutex.RUnlock()
	return n.peers[id]
}

// BroadcastBlock 广播区块
func (n *Network) BroadcastBlock(block *types.Block) {
	peers := n.GetPeers()
	for _, peer := range peers {
		// 向对等节点发送区块
		go n.sendBlockToPeer(peer, block)
	}
}

// BroadcastTransaction 广播交易
func (n *Network) BroadcastTransaction(tx *types.Transaction) {
	peers := n.GetPeers()
	for _, peer := range peers {
		// 向对等节点发送交易
		go n.sendTransactionToPeer(peer, tx)
	}
}

// sendBlockToPeer 向对等节点发送区块
func (n *Network) sendBlockToPeer(peer *interfaces.Peer, block *types.Block) {
	// 实现区块发送逻辑
	log.Printf("Sending block %d to peer", block.NumberU64())
}

// sendTransactionToPeer 向对等节点发送交易
func (n *Network) sendTransactionToPeer(peer *interfaces.Peer, tx *types.Transaction) {
	// 实现交易发送逻辑
	log.Printf("Sending transaction %s to peer", tx.Hash())
}

// SyncBlocks 同步区块
func (n *Network) SyncBlocks() error {
	// 实现区块同步逻辑
	log.Println("Starting block sync")
	return nil
}

// IsStarted 检查网络是否已启动
func (n *Network) IsStarted() bool {
	return n.isStarted
}

// GetNode 获取本地节点
func (n *Network) GetNode() interface{} {
	// 实现获取本地节点的逻辑
	return nil
}

// GetConfig 获取网络配置
func (n *Network) GetConfig() *config.Config {
	return n.config
}
