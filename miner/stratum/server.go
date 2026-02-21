package stratum

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"nogochain/core/types"
	"nogochain/metrics"
)

// Server Stratum服务器
type Server struct {
	addr      string
	listener  net.Listener
	clients   map[*Client]struct{}
	clientsMu sync.Mutex
	jobs      chan *Job
	shares    chan *Share
	stopCh    chan struct{}
	wg        sync.WaitGroup

	// 回调
	getWorkFn func() (*types.Block, error)
	submitFn  func(header *types.BlockHeader) error
}

// Client Stratum客户端
type Client struct {
	conn       net.Conn
	server     *Server
	id         string
	workerName string
	coinbase   common.Address
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// Job 挖矿任务
type Job struct {
	ID        string             `json:"id"`
	Header    *types.BlockHeader `json:"header"`
	Target    string             `json:"target"`
	Seed      string             `json:"seed"`
	Height    uint64             `json:"height"`
	Timestamp uint64             `json:"timestamp"`
}

// Share 提交的份额
type Share struct {
	JobID      string             `json:"job_id"`
	Nonce      uint64             `json:"nonce"`
	MixDigest  common.Hash        `json:"mix_digest"`
	Header     *types.BlockHeader `json:"header"`
	WorkerName string             `json:"worker_name"`
}

// Message Stratum消息
type Message struct {
	ID     interface{}      `json:"id"`
	Method string           `json:"method,omitempty"`
	Params json.RawMessage  `json:"params,omitempty"`
	Result interface{}      `json:"result,omitempty"`
	Error  *json.RawMessage `json:"error,omitempty"`
}

// NewServer 创建新的Stratum服务器
func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		clients: make(map[*Client]struct{}),
		jobs:    make(chan *Job, 100),
		shares:  make(chan *Share, 100),
		stopCh:  make(chan struct{}),
	}
}

// SetGetWorkFn 设置获取工作的回调
func (s *Server) SetGetWorkFn(fn func() (*types.Block, error)) {
	s.getWorkFn = fn
}

// SetSubmitFn 设置提交工作的回调
func (s *Server) SetSubmitFn(fn func(header *types.BlockHeader) error) {
	s.submitFn = fn
}

// Start 启动服务器
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener

	s.wg.Add(1)
	go s.acceptLoop()

	s.wg.Add(1)
	go s.jobLoop()

	s.wg.Add(1)
	go s.shareLoop()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	close(s.stopCh)

	s.clientsMu.Lock()
	for client := range s.clients {
		client.close()
	}
	s.clientsMu.Unlock()

	if s.listener != nil {
		s.listener.Close()
	}

	s.wg.Wait()

	return nil
}

// acceptLoop 接受连接循环
func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stopCh:
				return
			default:
				continue
			}
		}

		client := &Client{
			conn:   conn,
			server: s,
			id:     fmt.Sprintf("%s", conn.RemoteAddr()),
			stopCh: make(chan struct{}),
		}

		s.clientsMu.Lock()
		s.clients[client] = struct{}{}
		// 更新矿工数量指标
		metrics.PoolMiners.Set(float64(len(s.clients)))
		s.clientsMu.Unlock()

		client.wg.Add(1)
		go client.handleConnection()
	}
}

// jobLoop 任务循环
func (s *Server) jobLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.broadcastJob()
		}
	}
}

// shareLoop 份额处理循环
func (s *Server) shareLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopCh:
			return
		case share := <-s.shares:
			s.handleShare(share)
		}
	}
}

// broadcastJob 广播任务
func (s *Server) broadcastJob() {
	if s.getWorkFn == nil {
		return
	}

	block, err := s.getWorkFn()
	if err != nil {
		return
	}

	job := &Job{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Header:    block.Header,
		Target:    block.Header.Difficulty.Text(16),
		Height:    block.Header.Number.Uint64(),
		Timestamp: block.Header.Time,
	}

	s.clientsMu.Lock()
	for client := range s.clients {
		client.sendJob(job)
	}
	s.clientsMu.Unlock()
}

// handleShare 处理份额
func (s *Server) handleShare(share *Share) {
	if s.submitFn == nil {
		return
	}

	err := s.submitFn(share.Header)
	if err != nil {
		return
	}

	// 更新份额计数指标
	metrics.PoolShares.Inc()

	// TODO: 实现份额验证和奖励计算
}

// removeClient 移除客户端
func (s *Server) removeClient(client *Client) {
	s.clientsMu.Lock()
	delete(s.clients, client)
	// 更新矿工数量指标
	metrics.PoolMiners.Set(float64(len(s.clients)))
	s.clientsMu.Unlock()
}

// handleConnection 处理客户端连接
func (c *Client) handleConnection() {
	defer c.wg.Done()
	defer c.server.removeClient(c)
	defer c.conn.Close()

	decoder := json.NewDecoder(c.conn)
	encoder := json.NewEncoder(c.conn)

	for {
		select {
		case <-c.stopCh:
			return
		default:
			var msg Message
			if err := decoder.Decode(&msg); err != nil {
				return
			}

			c.handleMessage(&msg, encoder)
		}
	}
}

// handleMessage 处理消息
func (c *Client) handleMessage(msg *Message, encoder *json.Encoder) {
	switch msg.Method {
	case "mining.subscribe":
		c.handleSubscribe(msg, encoder)
	case "mining.authorize":
		c.handleAuthorize(msg, encoder)
	case "mining.submit":
		c.handleSubmit(msg, encoder)
	default:
		c.sendError(msg.ID, "Unknown method", encoder)
	}
}

// handleSubscribe 处理订阅
func (c *Client) handleSubscribe(msg *Message, encoder *json.Encoder) {
	response := Message{
		ID: msg.ID,
		Result: []interface{}{
			[]string{"mining.notify", "mining.set_difficulty"},
			"1.0",
		},
	}

	encoder.Encode(response)
}

// handleAuthorize 处理授权
func (c *Client) handleAuthorize(msg *Message, encoder *json.Encoder) {
	// TODO: 实现授权逻辑
	response := Message{
		ID:     msg.ID,
		Result: true,
	}

	encoder.Encode(response)
}

// handleSubmit 处理提交
func (c *Client) handleSubmit(msg *Message, encoder *json.Encoder) {
	// TODO: 实现份额提交处理
	response := Message{
		ID:     msg.ID,
		Result: true,
	}

	encoder.Encode(response)
}

// sendJob 发送任务
func (c *Client) sendJob(job *Job) {
	params := []interface{}{
		job.ID,
		"", // 简化处理，实际需要序列化区块头
		"", // 种子
		job.Target,
	}

	msg := Message{
		Method: "mining.notify",
		Params: toJSON(params),
	}

	encoder := json.NewEncoder(c.conn)
	encoder.Encode(msg)
}

// sendError 发送错误
func (c *Client) sendError(id interface{}, message string, encoder *json.Encoder) {
	err := map[string]interface{}{
		"code":    -1,
		"message": message,
	}

	errJSON, _ := json.Marshal(err)
	errMsg := json.RawMessage(errJSON)

	response := Message{
		ID:    id,
		Error: &errMsg,
	}

	encoder.Encode(response)
}

// close 关闭客户端
func (c *Client) close() {
	close(c.stopCh)
	c.wg.Wait()
}

// toJSON 转换为JSON
func toJSON(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
