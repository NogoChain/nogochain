# NogoChain 部署文档

## 1. 系统要求

### 1.1 硬件要求

| 节点类型 | CPU | 内存 | 存储 | 网络 | 适用场景 |
|----------|-----|------|------|------|----------|
| 全节点 | 4核+ | 8GB+ | 200GB+ SSD | 10Mbps+ | 生产环境 |
| 轻节点 | 2核+ | 4GB+ | 50GB+ SSD | 5Mbps+ | 开发测试 |
| 挖矿节点 | 4核+ | 16GB+ | 200GB+ SSD | 10Mbps+ | 参与共识 |
| RPC 节点 | 8核+ | 16GB+ | 500GB+ SSD | 100Mbps+ | 提供服务 |

### 1.2 软件要求

- **操作系统**：Linux (Ubuntu 20.04+/CentOS 8+), Windows 10+, macOS 10.15+
- **Go 版本**：1.22+（编译时需要）
- **依赖**：见 `go.mod` 文件
- **网络**：开放端口 30303 (P2P), 8545 (RPC), 3333 (Stratum，可选)

## 2. 部署方式

### 2.1 快速部署（推荐）

使用官方提供的启动脚本快速部署：

**生产环境部署**：

```bash
# Windows
start-mainnet-node.bat

# 矿池部署
start-mainnet-pool.bat
```

**测试网络部署**：

```bash
# Windows
testnet/start_node.bat

# 矿池部署
testnet/start_pool.bat
```

启动脚本会自动完成以下操作：
- 创建必要的目录结构
- 生成配置文件
- 初始化区块链
- 启动服务并监控状态

### 2.2 二进制部署

**步骤 1：下载二进制文件**

从官方发布页下载对应平台的二进制文件：
- nogochain（节点）
- nogopool（矿池）
- nogocli（命令行工具）
- nogominer（独立挖矿工具）

**步骤 2：配置文件**

创建 JSON 配置文件（示例）：

#### 节点配置文件 (node_config.json)
```json
{
  "chain_id": 318,
  "data_dir": "mainnet/data",
  "log_dir": "mainnet/logs",
  "p2p_port": 30303,
  "rpc_port": 8545,
  "ws_port": 8546,
  "max_peers": 50,
  "min_gas_price": "1000000000",
  "gas_limit": "8000000",
  "metrics": true,
  "metrics_port": 6060,
  "jwt_secret": "mainnet/jwt-token.txt"
}
```

#### 矿池配置文件 (mining_pool_config.json)
```json
{
  "pool_name": "NogoChain Pool",
  "node_url": "http://localhost:8545",
  "stratum_port": 3333,
  "stratum_tcp_port": 3334,
  "stratum_tls_port": 3335,
  "payment_interval": 600,
  "min_payout": "1000000000000000000",
  "fee": 0.01,
  "reward_address": "0xYourRewardAddress",
  "stats_interval": 30,
  "hashrate_window": 600,
  "max_conns_per_ip": 10
}
```

**步骤 3：启动节点**

```bash
# Linux/macOS
chmod +x ./nogochain
./nogochain --datadir mainnet/data --config mainnet/config/node_config.json

# Windows
./nogochain.exe --datadir mainnet/data --config mainnet/config/node_config.json
```

**步骤 4：启动矿池**

```bash
# Linux/macOS
chmod +x ./nogopool
./nogopool --datadir mainnet/data --config mainnet/config/mining_pool_config.json

# Windows
./nogopool.exe --datadir mainnet/data --config mainnet/config/mining_pool_config.json
```

### 2.3 源码编译

**步骤 1：克隆代码**

```bash
git clone https://github.com/nogochain/nogochain.git
cd nogochain
```

**步骤 2：验证依赖**

```bash
./scripts/verify-deps.sh
```

**步骤 3：编译**

```bash
# Linux/macOS
./scripts/build.sh

# Windows
./scripts/build.ps1

# Ubuntu 系统详细编译步骤
# 1. 安装依赖
sudo apt update
sudo apt install -y build-essential git curl

# 2. 安装 Go 1.22+
curl -LO https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 3. 编译到主网文件
go build -o mainnet/nogochain ./cmd/nogochain
go build -o mainnet/nogocli ./cmd/nogocli
go build -o mainnet/nogominer ./cmd/nogominer
go build -o mainnet/nogopool ./cmd/nogopool
```

**步骤 4：启动节点**

```bash
./bin/nogod --config config.toml

# Ubuntu 系统启动命令
# 启动主网节点
./mainnet/nogochain --datadir mainnet/data --config mainnet/config/node_config.json

# 后台启动主网节点
nohup ./mainnet/nogochain --datadir mainnet/data --config mainnet/config/node_config.json > mainnet/logs/nogochain.log 2>&1 &
```

### 2.4 容器部署

**步骤 1：Docker 构建**

```bash
docker build -t nogochain .
```

**步骤 2：Docker 运行**

```bash
docker run -d \
  --name nogochain \
  -p 30303:30303 \
  -p 8545:8545 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/config.toml:/app/config.toml \
  nogochain \
  --config config.toml
```

**步骤 3：Docker Compose**

创建 `docker-compose.yml`：

```yaml
version: '3'
services:
  nogochain:
    build: .
    ports:
      - "30303:30303"
      - "8545:8545"
    volumes:
      - ./data:/app/data
      - ./config.toml:/app/config.toml
    command: --config config.toml
```

启动：
```bash
docker-compose up -d
```

## 3. 配置详解

### 3.1 节点配置文件 (node_config.json)

| 配置项 | 类型 | 默认值 | 描述 |
|--------|------|--------|------|
| chain_id | number | 318 | 链 ID |
| data_dir | string | "mainnet/data" | 数据存储目录 |
| log_dir | string | "mainnet/logs" | 日志文件目录 |
| p2p_port | number | 30303 | P2P 通信端口 |
| rpc_port | number | 8545 | HTTP RPC 端口 |
| ws_port | number | 8546 | WebSocket 端口 |
| max_peers | number | 50 | 最大连接节点数 |
| min_gas_price | string | "1000000000" | 最低 gas 价格 |
| gas_limit | string | "8000000" | 区块 gas 限制 |
| metrics | boolean | true | 启用指标监控 |
| metrics_port | number | 6060 | 指标监控端口 |
| jwt_secret | string | "mainnet/jwt-token.txt" | JWT 密钥文件路径 |

### 3.2 矿池配置文件 (mining_pool_config.json)

| 配置项 | 类型 | 默认值 | 描述 |
|--------|------|--------|------|
| pool_name | string | "NogoChain Pool" | 矿池名称 |
| node_url | string | "http://localhost:8545" | 节点 RPC 地址 |
| stratum_port | number | 3333 | Stratum 服务端口 |
| stratum_tcp_port | number | 3334 | TCP Stratum 端口 |
| stratum_tls_port | number | 3335 | TLS Stratum 端口 |
| payment_interval | number | 600 | 支付间隔（秒） |
| min_payout | string | "1000000000000000000" | 最小支付金额 |
| fee | number | 0.01 | 矿池手续费 |
| reward_address | string | "0xYourRewardAddress" | 奖励地址 |
| stats_interval | number | 30 | 统计信息更新间隔（秒） |
| hashrate_window | number | 600 | 算力计算窗口（秒） |
| max_conns_per_ip | number | 10 | 每 IP 最大连接数 |

### 3.3 启动脚本参数

| 参数 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| --datadir | string | 脚本目录 | 数据目录路径 |
| --config | string | 自动生成 | 配置文件路径 |
| --genesis | string | 自动生成 | 创世区块文件路径 |
| --port | number | 30303 | P2P 端口 |
| --rpcport | number | 8545 | RPC 端口 |
| --wsport | number | 8546 | WebSocket 端口 |
| --metrics | boolean | true | 启用指标监控 |
| --verbosity | number | 3 | 日志详细程度 |

## 4. 网络设置

### 4.1 端口映射

| 端口 | 协议 | 用途 | 必须开放 |
|------|------|------|----------|
| 30303 | TCP/UDP | P2P 通信 | 是 |
| 8545 | TCP | HTTP RPC | 否（公网访问需要） |
| 8546 | TCP | WebSocket RPC | 否（公网访问需要） |
| 3333 | TCP | Stratum 服务 | 否（挖矿需要） |

### 4.2 防火墙设置

**Linux (iptables)**：

```bash
iptables -A INPUT -p tcp --dport 30303 -j ACCEPT
iptables -A INPUT -p udp --dport 30303 -j ACCEPT
# 如需开放 RPC
iptables -A INPUT -p tcp --dport 8545 -j ACCEPT
```

**Windows (防火墙)**：

通过 Windows 防火墙高级设置添加入站规则，开放对应端口。

### 4.3 节点发现

- 自动发现：通过 Kademlia 算法自动发现网络中的节点
- 手动添加：使用 `admin_addPeer` RPC 方法添加节点
- 引导节点：通过配置文件设置引导节点加速发现

## 5. 数据管理

### 5.1 数据目录结构

```
data/
├── chaindata/         # 区块数据
├── trie/              # 状态树数据
├── nodes/             # P2P 节点信息
├── keystore/          # 账户密钥
└── logs/              # 日志文件
```

### 5.2 数据备份

**定期备份**：
- 备份 `keystore` 目录（包含账户私钥）
- 备份 `chaindata` 目录（可选，用于快速恢复）

**备份策略**：
- 冷备份：离线存储备份数据
- 加密备份：对备份数据进行加密
- 多份备份：存储在不同位置

### 5.3 数据清理

**清理旧数据**：
```bash
# 停止节点
# 删除数据目录
rm -rf data/chaindata/*
# 重启节点（重新同步）
```

**注意**：清理数据后需要重新同步整个区块链，可能需要较长时间。

## 6. 节点管理

### 6.1 启动与停止

**启动**：
```bash
# 前台启动
./nogod --config config.toml

# 后台启动
nohup ./nogod --config config.toml > nogod.log 2>&1 &
```

**停止**：
```bash
# 通过 RPC 停止
curl -X POST --data '{"jsonrpc":"2.0","method":"admin_stop","params":[],"id":1}' http://localhost:8545

# 或使用进程管理
kill -SIGINT $(pgrep nogod)
```

### 6.2 监控与日志

**日志配置**：
```toml
[logger]
level = "info"  # debug, info, warn, error
format = "json"  # json, console
file = "logs/nogochain.log"
maxsize = 100  # MB
maxage = 7  # days
```

**查看日志**：
```bash
# 实时查看
 tail -f logs/nogochain.log

# 查看错误
 grep "ERROR" logs/nogochain.log
```

### 6.3 健康检查

**RPC 检查**：
```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545
```

**P2P 检查**：
```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' http://localhost:8545
```

## 7. 安全配置

### 7.1 RPC 安全

**生产环境建议**：
- 限制 RPC 访问地址：`--rpc.http.addr=127.0.0.1`
- 启用 JWT 认证：`--rpc.jwtsecret=jwt-token.txt`
- 禁用调试接口：`--rpc.disable-debug`

**JWT 生成**：
```bash
openssl rand -hex 32 > jwt-token.txt
```

### 7.2 节点安全

**建议措施**：
- 定期更新节点版本
- 使用防火墙限制端口访问
- 避免在公网暴露 P2P 以外的端口
- 使用非 root 账户运行节点

### 7.3 密钥安全

**私钥管理**：
- 不要在配置文件中明文存储私钥
- 使用硬件钱包管理挖矿地址
- 定期更换挖矿地址

## 8. 常见问题

### 8.1 启动问题

- **问题**：节点无法启动
  **解决**：检查配置文件，查看日志错误信息

- **问题**：端口被占用
  **解决**：修改配置文件中的端口，或停止占用端口的进程

### 8.2 同步问题

- **问题**：同步速度慢
  **解决**：确保网络带宽充足，添加更多对等节点

- **问题**：同步卡住
  **解决**：重启节点，检查网络连接

### 8.3 性能问题

- **问题**：CPU 使用率高
  **解决**：调整挖矿线程数，优化系统配置

- **问题**：内存使用高
  **解决**：增加内存，调整缓存配置

### 8.4 网络问题

- **问题**：无法连接到其他节点
  **解决**：检查防火墙设置，确保端口开放

- **问题**：P2P 连接数为 0
  **解决**：添加引导节点，检查网络连接

## 9. 升级与维护

### 9.1 版本升级

**步骤 1：备份数据**

```bash
cp -r data data_backup
```

**步骤 2：停止节点**

```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"admin_stop","params":[],"id":1}' http://localhost:8545
```

**步骤 3：更新二进制文件**

```bash
# 下载新版本或重新编译
```

**步骤 4：启动节点**

```bash
./nogod --config config.toml
```

### 9.2 日常维护

- **监控**：定期检查节点状态和性能
- **更新**：关注官方发布的新版本和安全补丁
- **备份**：定期备份关键数据
- **优化**：根据节点运行情况调整配置

## 10. 参考资料

- [NogoChain 官方文档](https://docs.nogochain.org)
- [生产环境部署教程](DEPLOYMENT_PRODUCTION.md) - 详细的生产环境部署指南
- [测试网络部署教程](DEPLOYMENT_TESTNET.md) - 快速测试网络部署指南
- [API 文档](API.md) - RPC 接口详细说明
- [架构文档](ARCHITECTURE.md) - 系统架构设计
- [挖矿文档](MINING.md) - 挖矿配置和优化
- [Go 语言安装指南](https://golang.org/doc/install)
- [Docker 官方文档](https://docs.docker.com)
- [Linux 系统优化](https://wiki.archlinux.org/title/Optimization)