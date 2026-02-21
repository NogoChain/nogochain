package tools

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"nogochain/core/types"
	"nogochain/network"
)

// TxGenerator 交易生成器
type TxGenerator struct {
	network     *network.Network
	rate        int           // 每秒生成的交易数
	duration    time.Duration // 持续时间
	transactions chan *types.Transaction
	mu          sync.Mutex
	stats       TxStats
}

// TxStats 交易统计信息
type TxStats struct {
	TotalTransactions int64
	SuccessfulTransactions int64
	FailedTransactions int64
	StartTime          time.Time
	EndTime            time.Time
}

// NewTxGenerator 创建新的交易生成器
func NewTxGenerator(network *network.Network, rate int, duration time.Duration) *TxGenerator {
	return &TxGenerator{
		network:     network,
		rate:        rate,
		duration:    duration,
		transactions: make(chan *types.Transaction, rate*10),
		stats: TxStats{
			StartTime: time.Now(),
		},
	}
}

// Start 开始生成交易
func (g *TxGenerator) Start() {
	go g.generateTransactions()
	go g.broadcastTransactions()
}

// generateTransactions 生成交易
func (g *TxGenerator) generateTransactions() {
	ticker := time.NewTicker(time.Second / time.Duration(g.rate))
	defer ticker.Stop()

	timer := time.NewTimer(g.duration)
	defer timer.Stop()

	nonce := uint64(0)

	for {
		select {
		case <-ticker.C:
			// 创建交易
			tx := &types.Transaction{
				Nonce:    nonce,
				GasPrice: big.NewInt(1),
				Gas:      21000,
				To:       &common.Address{},
				Value:    big.NewInt(1000),
				Data:     nil,
			}

			g.transactions <- tx
			nonce++

			g.mu.Lock()
			g.stats.TotalTransactions++
			g.mu.Unlock()

		case <-timer.C:
			close(g.transactions)
			g.mu.Lock()
			g.stats.EndTime = time.Now()
			g.mu.Unlock()
			return
		}
	}
}

// broadcastTransactions 广播交易
func (g *TxGenerator) broadcastTransactions() {
	for tx := range g.transactions {
		g.network.BroadcastTransaction(tx)
		g.mu.Lock()
		g.stats.SuccessfulTransactions++
		g.mu.Unlock()
	}
}

// GetStats 获取交易统计信息
func (g *TxGenerator) GetStats() TxStats {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.stats
}

// PrintStats 打印交易统计信息
func (g *TxGenerator) PrintStats() {
	stats := g.GetStats()
	elapsed := stats.EndTime.Sub(stats.StartTime)
	fmt.Printf("交易生成统计信息:\n")
	fmt.Printf("总交易数: %d\n", stats.TotalTransactions)
	fmt.Printf("成功交易数: %d\n", stats.SuccessfulTransactions)
	fmt.Printf("持续时间: %v\n", elapsed)
	fmt.Printf("实际交易速率: %.2f TPS\n", float64(stats.SuccessfulTransactions)/elapsed.Seconds())
}
