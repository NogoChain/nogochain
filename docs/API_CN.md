# NogoChain API 文档

## 1. 概述

NogoChain 提供标准的 EVM 兼容 RPC 接口，同时扩展了 NogoChain 特有的功能接口。所有接口均通过 HTTP/JSON-RPC 协议访问，默认端口为 8545。

## 2. 标准接口

### 2.1 eth_* 接口

| 方法 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| eth_accounts | 无 | ArrayString> | 获取本地账户列表 |
| eth_blockNumber | 无 | String (Quantity) | 获取最新区块号 |
| eth_call | Object, String | String | 调用合约方法（只读） |
| eth_chainId | 无 | String (Quantity) | 获取链 ID（318） |
| eth_coinbase | 无 | String | 获取挖矿地址 |
| eth_estimateGas | Object | String (Quantity) | 估算交易 gas 消耗 |
| eth_gasPrice | 无 | String (Quantity) | 获取当前 gas 价格 |
| eth_getBalance | String, String | String (Quantity) | 获取账户余额 |
| eth_getBlockByHash | String, Boolean | Object | 通过哈希获取区块 |
| eth_getBlockByNumber | String, Boolean | Object | 通过编号获取区块 |
| eth_getBlockTransactionCountByHash | String | String (Quantity) | 获取区块交易数 |
| eth_getBlockTransactionCountByNumber | String | String (Quantity) | 获取区块交易数 |
| eth_getCode | String, String | String | 获取合约代码 |
| eth_getLogs | Object | ArrayObject> | 获取日志 |
| eth_getStorageAt | String, String, String | String | 获取存储值 |
| eth_getTransactionByHash | String | Object | 通过哈希获取交易 |
| eth_getTransactionByBlockHashAndIndex | String, String | Object | 通过区块哈希和索引获取交易 |
| eth_getTransactionByBlockNumberAndIndex | String, String | Object | 通过区块编号和索引获取交易 |
| eth_getTransactionCount | String, String | String (Quantity) | 获取交易计数（nonce） |
| eth_getTransactionReceipt | String | Object | 获取交易收据 |
| eth_hashrate | 无 | String (Quantity) | 获取哈希率 |
| eth_mining | 无 | Boolean | 检查是否在挖矿 |
| eth_sendRawTransaction | String | String | 发送原始交易 |
| eth_submitHashrate | String, String | Boolean | 提交哈希率 |
| eth_submitWork | String, String, String | Boolean | 提交挖矿结果 |

### 2.2 net_* 接口

| 方法 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| net_listening | 无 | Boolean | 检查节点是否在监听 |
| net_peerCount | 无 | String (Quantity) | 获取连接的对等节点数 |
| net_version | 无 | String | 获取网络版本 |

### 2.3 web3_* 接口

| 方法 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| web3_clientVersion | 无 | String | 获取客户端版本 |
| web3_sha3 | String | String | 计算数据的 Keccak-256 哈希 |

### 2.4 debug_* 接口

| 方法 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| debug_accountRange | String, String, Number | Object | 获取账户范围 |
| debug_backtraceAt | String | Array | 获取指定位置的回溯 |
| debug_blockProfile | String, Number | Boolean | 生成区块分析文件 |
| debug_cpuProfile | String, Number | Boolean | 生成 CPU 分析文件 |
| debug_dumpBlock | String | Object | 导出区块数据 |
| debug_gcStats | 无 | Object | 获取垃圾回收统计 |
| debug_getBlockRlp | String | String | 获取区块的 RLP 编码 |
| debug_memStats | 无 | Object | 获取内存统计 |
| debug_seedHash | Number | String | 获取指定区块的种子哈希 |
| debug_setHead | String | Boolean | 设置链头 |
| debug_standardTraceBadBlockToFile | String, String | Boolean | 追踪坏区块到文件 |
| debug_standardTraceBlockToFile | String, String | Boolean | 追踪区块到文件 |
| debug_startCPUProfile | String | Boolean | 开始 CPU 分析 |
| debug_stopCPUProfile | 无 | Boolean | 停止 CPU 分析 |
| debug_traceBlock | String, Object | Object | 追踪区块执行 |
| debug_traceBlockByNumber | String, Object | Object | 追踪指定编号区块 |
| debug_traceBlockByHash | String, Object | Object | 追踪指定哈希区块 |
| debug_traceTransaction | String, Object | Object | 追踪交易执行 |
| debug_verbosity | Number | Boolean | 设置日志级别 |

## 3. NogoChain 特有接口 (nogo_*)

### 3.1 共识相关

| 方法 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| nogo_getDifficulty | 无 | String (Quantity) | 获取当前难度 |
| nogo_getReward | 无 | String (Quantity) | 获取当前区块奖励 |
| nogo_getMiningInfo | 无 | Object | 获取挖矿信息（难度、哈希率、奖励） |
| nogo_getWork | 无 | ArrayString> | 获取挖矿工作（NogoPow） |
| nogo_submitWork | String, String, String | Boolean | 提交挖矿工作（NogoPow） |
| nogo_submitHashrate | String, String | Boolean | 提交哈希率（NogoPow） |

### 3.2 链信息

| 方法 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| nogo_getChainInfo | 无 | Object | 获取链信息（链 ID、符号、共识算法等） |

## 4. 接口示例

### 4.1 示例请求

#### 4.1.1 标准接口示例

```json
{
  "jsonrpc": "2.0",
  "method": "eth_blockNumber",
  "params": [],
  "id": 1
}
```

#### 4.1.2 NogoChain 特有接口示例

```json
{
  "jsonrpc": "2.0",
  "method": "nogo_getChainInfo",
  "params": [],
  "id": 1
}
```

```json
{
  "jsonrpc": "2.0",
  "method": "nogo_getMiningInfo",
  "params": [],
  "id": 1
}
```

### 4.2 示例响应

#### 4.2.1 标准接口响应

```json
{
  "jsonrpc": "2.0",
  "result": "0x1234",
  "id": 1
}
```

#### 4.2.2 NogoChain 特有接口响应

```json
{
  "jsonrpc": "2.0",
  "result": {
    "chainId": 318,
    "symbol": "NOGO",
    "decimals": 18,
    "consensus": "NogoPow",
    "difficulty": "1000000"
  },
  "id": 1
}
```

```json
{
  "jsonrpc": "2.0",
  "result": {
    "difficulty": "1000000",
    "hashrate": "0",
    "miner": "0x0000000000000000000000000000000000000000",
    "networkHashrate": "0"
  },
  "id": 1
}
```

### 4.3 错误响应

```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32601,
    "message": "Method not found"
  },
  "id": 1
}
```

## 5. RPC 配置

### 5.1 默认配置

- **HTTP**：`http://127.0.0.1:8545`
- **WebSocket**：`ws://127.0.0.1:8546`
- **认证**：本地访问无需认证，公网访问需 JWT

### 5.2 安全建议

- 生产环境应限制 RPC 访问 IP
- 启用 JWT 认证保护公网 RPC
- 避免在 RPC 中暴露敏感信息

## 6. 工具集成

### 6.1 开发工具

- **Hardhat**：使用标准 EVM 配置
- **Truffle**：配置网络为 NogoChain
- **Foundry**：使用 `--rpc-url` 指定 NogoChain RPC
- **Remix**：通过 "Deploy & Run Transactions" 连接 NogoChain

### 6.2 钱包集成

- **MetaMask**：添加自定义网络（链 ID: 318）
- **Trust Wallet**：支持 NogoChain
- **Ledger/Trezor**：通过标准 EVM 接口支持

## 7. 速率限制

- 标准 RPC 请求：每 IP 1000 请求/分钟
- 调试接口：每 IP 100 请求/分钟
- 可通过配置文件调整限制

## 8. 版本兼容性

| API 版本 | 兼容链版本 | 备注 |
|----------|------------|------|
| v1.0 | 所有版本 | 基础接口 |
| v1.1 | ≥ 1.0.0 | 新增 nogo_ 接口 |
| v1.2 | ≥ 1.1.0 | 优化性能 |

## 9. 常见问题

### 9.1 连接问题
- 确保 RPC 服务已启动
- 检查防火墙设置
- 验证网络配置

### 9.2 性能问题
- 对于高频请求，建议使用 WebSocket
- 批量处理多个请求
- 避免使用昂贵的调试接口

### 9.3 错误处理
- 检查参数格式是否正确
- 验证交易签名
- 确认账户余额充足

## 10. 参考资料

- [以太坊 JSON-RPC 文档](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [NogoChain 源码](https://github.com/nogochain/nogochain)