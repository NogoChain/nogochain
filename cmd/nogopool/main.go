package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"nogochain/core/blockchain"
	"nogochain/core/types"
	"nogochain/metrics"
	"nogochain/miner/stratum"
)

// PoolConfig 矿池配置
type PoolConfig struct {
	Pool struct {
		Name           string `json:"name"`
		Fee            float64 `json:"fee"`
		MinPayout      float64 `json:"minPayout"`
		PayoutInterval int     `json:"payoutInterval"`
		Host           string `json:"host"`
		Port           int    `json:"port"`
		SslPort        int    `json:"sslPort"`
		Difficulty     int    `json:"difficulty"`
		MaxDifficulty  int    `json:"maxDifficulty"`
		VarDiff        struct {
			Enabled       bool    `json:"enabled"`
			MinDifficulty int     `json:"minDifficulty"`
			MaxDifficulty int     `json:"maxDifficulty"`
			TargetTime    int     `json:"targetTime"`
			RetargetTime  int     `json:"retargetTime"`
			Variance      float64 `json:"variance"`
		} `json:"varDiff"`
	} `json:"pool"`
	Node struct {
		RpcUrl              string `json:"rpcUrl"`
		RpcTimeout          int    `json:"rpcTimeout"`
		MaxConns            int    `json:"maxConns"`
		BlockRefreshInterval int    `json:"blockRefreshInterval"`
		JobRebroadcastInterval int    `json:"jobRebroadcastInterval"`
	} `json:"node"`
	Workers struct {
		MaxWorkersPerIP     int `json:"maxWorkersPerIP"`
		HashrateExpiration  int `json:"hashrateExpiration"`
		InvalidSharesBan    int `json:"invalidSharesBan"`
		BanDuration         int `json:"banDuration"`
	} `json:"workers"`
	Payment struct {
		Enabled               bool    `json:"enabled"`
		Type                  string  `json:"type"`
		PayoutThreshold       float64 `json:"payoutThreshold"`
		PayoutsPerRound       int     `json:"payoutsPerRound"`
		MinimumConfirmations  int     `json:"minimumConfirmations"`
	} `json:"payment"`
	Stats struct {
		Enabled       bool `json:"enabled"`
		UpdateInterval int  `json:"updateInterval"`
		History       struct {
			Enabled bool `json:"enabled"`
			Length  int  `json:"length"`
		} `json:"history"`
	} `json:"stats"`
	Log struct {
		Level    string `json:"level"`
		File     string `json:"file"`
		MaxSize  int    `json:"maxSize"`
		MaxAge   int    `json:"maxAge"`
		Compress bool   `json:"compress"`
	} `json:"log"`
}

// 全局变量
var (
	configFile string
	config     PoolConfig
	bc         *blockchain.Blockchain
	stratumSrv *stratum.Server
)

// initLogger 初始化日志系统
func initLogger(cfg *PoolConfig) {
	// 创建日志目录
	logDir := filepath.Dir(cfg.Log.File)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	// 配置日志轮转
	logWriter := &lumberjack.Logger{
		Filename: cfg.Log.File,
		MaxSize:  cfg.Log.MaxSize, // MB
		MaxAge:   cfg.Log.MaxAge,  // 天
		Compress: cfg.Log.Compress,
	}

	// 设置zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(logWriter).With().Timestamp().Caller().Logger()

	// 设置日志级别
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}

// initMetrics 初始化Prometheus监控
func initMetrics() {
	// 初始化监控指标
	metrics.InitMetrics()

	// 启动系统资源监控
	metrics.StartSystemMetrics()

	// 启动监控服务器
	addr := "127.0.0.1:9092"
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Info().Str("addr", addr).Msg("Starting Prometheus metrics server")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Error().Err(err).Msg("Failed to start metrics server")
		}
	}()
}

// getWork 获取挖矿任务
func getWork() (*types.Block, error) {
	// 获取当前头部区块
	currentHead := bc.CurrentHead()
	
	// 生成新的区块
	newBlock := types.NewBlock(
		currentHead.Hash(),
		currentHead.Coinbase(),
		common.Hash{},
		common.Hash{},
		common.Hash{},
		currentHead.Header.Difficulty,
		new(big.Int).Add(currentHead.Header.Number, big.NewInt(1)),
		currentHead.GasLimit(),
		0,
		uint64(time.Now().Unix()),
		[]byte("NogoChain Block"),
		common.Hash{},
		0,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)
	
	return newBlock, nil
}

// submitWork 提交挖矿结果
func submitWork(header *types.BlockHeader) error {
	// 创建区块
	block := types.NewBlock(
		header.ParentHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
		header.MixDigest,
		header.Nonce,
		[]*types.Transaction{},
		[]*types.BlockHeader{},
	)
	
	// 添加区块
	err := bc.AddBlock(block)
	if err != nil {
		return err
	}
	log.Info().Uint64("height", header.Number.Uint64()).Str("hash", header.Hash().String()).Msg("Block submitted successfully")
	return nil
}

func main() {
	fmt.Println("NogoChain Mining Pool")
	fmt.Println("Starting NogoChain mining pool...")

	// 解析命令行参数
	flag.StringVar(&configFile, "config", "testnet/config/mining_pool_config.json", "Path to pool config file")
	flag.Parse()

	// 打印配置文件路径
	fmt.Printf("Using config file: %s\n", configFile)

	// 加载配置文件
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Config file loaded successfully")

	// 解析配置
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Config parsed successfully")

	// 初始化日志系统
	initLogger(&config)
	fmt.Println("Logger initialized")

	fmt.Println("NogoChain mining pool starting...")

	// 初始化区块链
	bc = blockchain.NewBlockchain(nil)
	fmt.Printf("Blockchain initialized. Genesis block: %s\n", bc.Genesis().Hash().String())
	fmt.Printf("Current head: %s, height: %d\n", bc.CurrentHead().Hash().String(), bc.CurrentHead().NumberU64())

	// 初始化Stratum服务器
	addr := fmt.Sprintf("%s:%d", config.Pool.Host, config.Pool.Port)
	stratumSrv = stratum.NewServer(addr)
	stratumSrv.SetGetWorkFn(getWork)
	stratumSrv.SetSubmitFn(submitWork)
	fmt.Printf("Stratum server initialized. Address: %s\n", addr)

	// 启动Stratum服务器
	if err := stratumSrv.Start(); err != nil {
		fmt.Printf("Failed to start Stratum server: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Stratum server started successfully. Address: %s\n", addr)

	// 初始化监控
	initMetrics()
	fmt.Println("Metrics initialized")

	// 启动矿工连接检查
	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println("Mining pool is running...")
			fmt.Println("Waiting for miners to connect...")
		}
	}()

	fmt.Println("Mining pool started successfully!")
	fmt.Println("NogoChain mining pool is ready for miners")

	// 保持程序运行
	select {}
}
