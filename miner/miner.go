package miner

import (
	"context"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"nogochain/consensus/nogopow"
	"nogochain/core/types"
)

// Config 挖矿配置
type Config struct {
	Enabled          bool
	Coinbase         common.Address
	ExtraData        []byte
	MinGasPrice      *big.Int
	MaxGasLimit      uint64
	RecommitInterval time.Duration
	NumThreads       int
}

// Miner 挖矿实例
type Miner struct {
	config  *Config
	engine  *nogopow.NogoPow
	chain   interface{} // 区块链接口
	stopCh  chan struct{}
	startCh chan struct{}
	quitCh  chan struct{}
	wg      sync.WaitGroup

	// 状态
	isRunning bool
	mu        sync.Mutex
}

// NewMiner 创建新的挖矿实例
func NewMiner(config *Config, engine *nogopow.NogoPow) *Miner {
	if config.NumThreads <= 0 {
		config.NumThreads = runtime.GOMAXPROCS(0)
	}

	return &Miner{
		config:  config,
		engine:  engine,
		stopCh:  make(chan struct{}),
		startCh: make(chan struct{}),
		quitCh:  make(chan struct{}),
	}
}

// SetChain 设置区块链接口
func (m *Miner) SetChain(chain interface{}) {
	m.chain = chain
}

// Start 启动挖矿
func (m *Miner) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return nil
	}

	m.isRunning = true
	close(m.startCh)
	m.startCh = make(chan struct{})

	m.wg.Add(1)
	go m.miningLoop()

	return nil
}

// Stop 停止挖矿
func (m *Miner) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	m.isRunning = false
	close(m.stopCh)
	m.stopCh = make(chan struct{})

	m.wg.Wait()

	return nil
}

// IsRunning 检查挖矿是否运行中
func (m *Miner) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.isRunning
}

// miningLoop 挖矿主循环
func (m *Miner) miningLoop() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.config.RecommitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.minePending()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// minePending 挖矿待处理区块
func (m *Miner) minePending() {
	// TODO: 实现区块构建和挖矿逻辑
	// 1. 构建区块头
	// 2. 初始化NogoPow
	// 3. 执行挖矿
	// 4. 提交区块
}

// Seal  sealing区块
func (m *Miner) Seal(ctx context.Context, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// 初始化NogoPow
	headerBytes := m.serializeHeader(block.Header)
	seed := m.calculateSeed(block.Header.Number)
	m.engine.Initialize(seed)

	// 启动挖矿
	found := make(chan *types.Block, 1)

	go func() {
		nonce, _, mixDigest, ok := m.engine.MineParallel(headerBytes, block.Header.Difficulty, 1000000)
		if ok {
			// 更新区块头
			newHeader := *block.Header
			newHeader.Nonce = nonce
			newHeader.MixDigest = common.BytesToHash(mixDigest)

			newBlock := &types.Block{
				Header:       &newHeader,
				Transactions: block.Transactions,
				Uncles:       block.Uncles,
			}

			select {
			case found <- newBlock:
			case <-stop:
			}
		}
	}()

	select {
	case result := <-found:
		results <- result
	case <-stop:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// serializeHeader 序列化区块头
func (m *Miner) serializeHeader(header *types.BlockHeader) []byte {
	return SerializeHeader(header)
}

// calculateSeed 计算种子
func (m *Miner) calculateSeed(number *big.Int) []byte {
	return CalculateSeed(number)
}

// GetHashRate 获取哈希率
func (m *Miner) GetHashRate() float64 {
	// TODO: 实现哈希率计算
	return 0
}

// SetExtra 设置额外数据
func (m *Miner) SetExtra(extra []byte) {
	m.config.ExtraData = extra
}

// SetCoinbase 设置矿工地址
func (m *Miner) SetCoinbase(addr common.Address) {
	m.config.Coinbase = addr
}
