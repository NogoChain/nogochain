package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// Message Stratum消息
type Message struct {
	ID     interface{}      `json:"id"`
	Method string           `json:"method,omitempty"`
	Params json.RawMessage  `json:"params,omitempty"`
	Result interface{}      `json:"result,omitempty"`
	Error  *json.RawMessage `json:"error,omitempty"`
}

// Job 挖矿任务
type Job struct {
	ID        string `json:"id"`
	Header    string `json:"header"`
	Seed      string `json:"seed"`
	Target    string `json:"target"`
	Height    uint64 `json:"height"`
	Timestamp uint64 `json:"timestamp"`
}

func main() {
	fmt.Println("NogoChain Stratum Miner")
	fmt.Println("Starting NogoChain stratum miner...")

	// 解析命令行参数
	poolAddr := flag.String("pool", "127.0.0.1:3333", "Mining pool address")
	workerName := flag.String("worker", "miner1", "Worker name")
	password := flag.String("password", "", "Worker password")
	threads := flag.Int("threads", 1, "Number of mining threads")
	flag.Parse()

	// 连接到矿池
	conn, err := net.Dial("tcp", *poolAddr)
	if err != nil {
		fmt.Printf("Failed to connect to pool: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to pool: %s\n", *poolAddr)

	// 创建JSON编码器和解码器
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// 1. 订阅挖矿任务
	fmt.Println("Subscribing to mining jobs...")
	subscribeMsg := Message{
		ID:     1,
		Method: "mining.subscribe",
		Params: json.RawMessage(`[]`),
	}

	if err := encoder.Encode(subscribeMsg); err != nil {
		fmt.Printf("Failed to send subscribe message: %v\n", err)
		os.Exit(1)
	}

	// 接收订阅响应
	var subscribeResp Message
	if err := decoder.Decode(&subscribeResp); err != nil {
		fmt.Printf("Failed to receive subscribe response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Subscribed successfully!")

	// 2. 授权
	fmt.Println("Authorizing worker...")
	authorizeMsg := Message{
		ID:     2,
		Method: "mining.authorize",
		Params: json.RawMessage(`["` + *workerName + `", "` + *password + `"]`),
	}

	if err := encoder.Encode(authorizeMsg); err != nil {
		fmt.Printf("Failed to send authorize message: %v\n", err)
		os.Exit(1)
	}

	// 接收授权响应
	var authorizeResp Message
	if err := decoder.Decode(&authorizeResp); err != nil {
		fmt.Printf("Failed to receive authorize response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Authorized successfully!")

	// 3. 开始挖矿
	fmt.Printf("Starting miner with %d threads...\n", *threads)
	fmt.Println("Press Ctrl+C to stop")

	// 模拟挖矿过程
	for i := 1; ; i++ {
		fmt.Printf("Mining... Job #%d\r", i)
		
		// 接收挖矿任务
		var jobMsg Message
		if err := decoder.Decode(&jobMsg); err != nil {
			// 忽略错误，继续挖矿
			continue
		}

		// 处理挖矿任务
		if jobMsg.Method == "mining.notify" {
			fmt.Printf("Received new job: %s\n", jobMsg.Params)
			
			// 模拟挖矿
			time.Sleep(1 * time.Second)
			
			// 提交份额
			submitMsg := Message{
				ID:     3,
				Method: "mining.submit",
				Params: json.RawMessage(`["` + *workerName + `", "job1", "nonce1", "mix1"]`),
			}

			if err := encoder.Encode(submitMsg); err != nil {
				fmt.Printf("Failed to send submit message: %v\n", err)
				continue
			}
			
			// 接收提交响应
			var submitResp Message
			if err := decoder.Decode(&submitResp); err != nil {
				fmt.Printf("Failed to receive submit response: %v\n", err)
				continue
			}
			
			fmt.Println("Share submitted successfully!")
		}
	}
}
