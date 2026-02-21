package p2p

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"nogochain/interfaces"
	"nogochain/network/compression"
	"nogochain/network/priority"
	"nogochain/network/types"
)

// Protocol 网络协议
type Protocol struct {
	Name    string
	Version uint
	Length  uint
	Run     func(peer *Peer, rw MsgReadWriter) error
}

// MsgReadWriter 消息读写接口
type MsgReadWriter interface {
	ReadMsg() (types.Msg, error)
	WriteMsg(types.Msg) error
}

// Peer 对等节点连接
type Peer struct {
	ID              interfaces.NodeID
	Node            interface{}
	Conn            net.Conn
	Protocols       map[string]uint
	lastSeen        time.Time
	mu              sync.Mutex
	priorityManager *priority.PriorityManager
}

// NewPeer 创建新的对等节点连接
func NewPeer(id interfaces.NodeID, node interface{}, conn net.Conn) *Peer {
	return &Peer{
		ID:              id,
		Node:            node,
		Conn:            conn,
		Protocols:       make(map[string]uint),
		lastSeen:        time.Now(),
		priorityManager: priority.NewPriorityManager(),
	}
}

// ReadMsg 读取消息
func (p *Peer) ReadMsg() (types.Msg, error) {
	// 读取消息头
	var msg types.Msg
	var header [8]byte

	_, err := io.ReadFull(p.Conn, header[:])
	if err != nil {
		return msg, err
	}

	// 解析消息头
	code := binary.BigEndian.Uint32(header[:4])
	size := binary.BigEndian.Uint32(header[4:])

	// 读取消息体
	compressedPayload := make([]byte, size)
	_, err = io.ReadFull(p.Conn, compressedPayload)
	if err != nil {
		return msg, err
	}

	// 解压消息体
	decompressedPayload, err := compression.Decompress(compressedPayload)
	if err != nil {
		return msg, err
	}

	// 设置消息
	msg.Code = uint64(code)
	msg.Size = uint32(len(decompressedPayload))
	msg.Payload = io.NopCloser(bytes.NewReader(decompressedPayload))

	p.lastSeen = time.Now()
	return msg, nil
}

// WriteMsg 写入消息
func (p *Peer) WriteMsg(msg types.Msg) error {
	// 将消息添加到优先级队列
	peerID := fmt.Sprintf("%x", p.ID)
	p.priorityManager.AddMsg(peerID, &msg, priority.NormalPriority)

	// 延迟发送，以便批量处理
	go func() {
		time.Sleep(50 * time.Millisecond)
		if p.priorityManager.HasMsg(peerID) {
			p.flushBatch(peerID)
		}
	}()

	return nil
}

// flushBatch 刷新批量消息
func (p *Peer) flushBatch(peerID string) error {
	// 从优先级队列获取消息
	var msgs []*types.Msg
	for {
		msg := p.priorityManager.GetNextMsg(peerID)
		if msg == nil {
			break
		}
		msgs = append(msgs, msg)
	}

	if len(msgs) == 0 {
		return nil
	}

	// 批量处理消息
	for _, msg := range msgs {
		// 压缩消息体
		payload, err := io.ReadAll(msg.Payload)
		if err != nil {
			return err
		}

		compressedPayload, err := compression.Compress(payload)
		if err != nil {
			return err
		}

		// 写入消息头
		var header [8]byte
		binary.BigEndian.PutUint32(header[:4], uint32(msg.Code))
		binary.BigEndian.PutUint32(header[4:], uint32(len(compressedPayload)))

		_, err = p.Conn.Write(header[:])
		if err != nil {
			return err
		}

		// 写入压缩后的消息体
		_, err = p.Conn.Write(compressedPayload)
		if err != nil {
			return err
		}
	}

	p.lastSeen = time.Now()
	return nil
}

// Close 关闭连接
func (p *Peer) Close() error {
	return p.Conn.Close()
}

// AddProtocol 添加协议
func (p *Peer) AddProtocol(name string, version uint) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Protocols[name] = version
}

// HasProtocol 检查是否支持协议
func (p *Peer) HasProtocol(name string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, exists := p.Protocols[name]
	return exists
}

// GetProtocolVersion 获取协议版本
func (p *Peer) GetProtocolVersion(name string) uint {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Protocols[name]
}

// LastSeen 获取最后一次通信时间
func (p *Peer) LastSeen() time.Time {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lastSeen
}

// Server P2P服务器
type Server struct {
	addr      string
	peers     map[string]*Peer
	peersMu   sync.RWMutex
	protocols map[string]*Protocol
	ctx       context.Context
	cancel    context.CancelFunc
	isRunning bool
}

// NewServer 创建新的P2P服务器
func NewServer(addr string) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		addr:      addr,
		peers:     make(map[string]*Peer),
		protocols: make(map[string]*Protocol),
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	if s.isRunning {
		return nil
	}

	// 启动监听
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.isRunning = true
	log.Printf("P2P server started on %s", s.addr)

	// 接受连接
	go func() {
		for s.isRunning {
			conn, err := listener.Accept()
			if err != nil {
				if !s.isRunning {
					return
				}
				log.Printf("Accept error: %v", err)
				continue
			}

			// 处理连接
			go s.handleConnection(conn)
		}
	}()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	if !s.isRunning {
		return nil
	}

	s.isRunning = false
	s.cancel()

	// 关闭所有连接
	s.peersMu.Lock()
	for _, peer := range s.peers {
		peer.Close()
	}
	s.peers = make(map[string]*Peer)
	s.peersMu.Unlock()

	log.Println("P2P server stopped")
	return nil
}

// handleConnection 处理新连接
func (s *Server) handleConnection(conn net.Conn) {
	// 这里应该实现节点握手和协议协商
	// 暂时简化处理
	log.Printf("New connection from %s", conn.RemoteAddr())

	// 模拟节点ID
	var id interfaces.NodeID
	copy(id[:], []byte(fmt.Sprintf("%d", time.Now().UnixNano())))

	// 创建对等节点
	peer := NewPeer(id, nil, conn)

	// 添加到对等节点列表
	peerID := fmt.Sprintf("%x", id)
	s.peersMu.Lock()
	s.peers[peerID] = peer
	s.peersMu.Unlock()

	// 启动协议处理
	go s.handlePeer(peer)
}

// handlePeer 处理对等节点
func (s *Server) handlePeer(peer *Peer) {
	// 这里应该实现协议处理逻辑
	// 暂时简化处理
	defer func() {
		peer.Close()
		peerID := fmt.Sprintf("%x", peer.ID)
		s.peersMu.Lock()
		delete(s.peers, peerID)
		s.peersMu.Unlock()
		log.Printf("Peer disconnected: %x", peer.ID)
	}()

	// 循环读取消息
	for s.isRunning {
		msg, err := peer.ReadMsg()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		// 处理消息
		s.handleMsg(peer, msg)
	}
}

// handleMsg 处理消息
func (s *Server) handleMsg(peer *Peer, msg types.Msg) {
	// 这里应该实现消息处理逻辑
	log.Printf("Received msg from peer %x: code=%d, size=%d", peer.ID, msg.Code, msg.Size)
}

// RegisterProtocol 注册协议
func (s *Server) RegisterProtocol(protocol *Protocol) {
	s.protocols[protocol.Name] = protocol
}

// GetPeers 获取所有对等节点
func (s *Server) GetPeers() []*Peer {
	s.peersMu.RLock()
	defer s.peersMu.RUnlock()
	peers := make([]*Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	return peers
}

// GetPeer 通过ID获取对等节点
func (s *Server) GetPeer(id string) *Peer {
	s.peersMu.RLock()
	defer s.peersMu.RUnlock()
	return s.peers[id]
}

// IsRunning 检查服务器是否运行
func (s *Server) IsRunning() bool {
	return s.isRunning
}
