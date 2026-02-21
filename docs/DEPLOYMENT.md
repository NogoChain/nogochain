# NogoChain Deployment Documentation

## 1. System Requirements

### 1.1 Hardware Requirements

| Node Type | CPU | Memory | Storage | Network | Use Case |
|-----------|-----|--------|---------|---------|----------|
| Full Node | 4+ cores | 8GB+ | 200GB+ SSD | 10Mbps+ | Production Environment |
| Light Node | 2+ cores | 4GB+ | 50GB+ SSD | 5Mbps+ | Development and Testing |
| Mining Node | 4+ cores | 16GB+ | 200GB+ SSD | 10Mbps+ | Consensus Participation |
| RPC Node | 8+ cores | 16GB+ | 500GB+ SSD | 100Mbps+ | Service Provider |

### 1.2 Software Requirements

- **Operating System**: Linux (Ubuntu 20.04+/CentOS 8+), Windows 10+, macOS 10.15+
- **Go Version**: 1.22+ (required for compilation)
- **Dependencies**: See `go.mod` file
- **Network**: Open ports 30303 (P2P), 8545 (RPC), 3333 (Stratum, optional)

## 2. Deployment Methods

### 2.1 Quick Deployment (Recommended)

Use the official startup scripts for quick deployment:

**Production Environment Deployment**:

```bash
# Windows
start-mainnet-node.bat

# Pool Deployment
start-mainnet-pool.bat
```

**Test Network Deployment**:

```bash
# Windows
testnet/start_node.bat

# Pool Deployment
testnet/start_pool.bat
```

The startup scripts automatically perform the following operations:
- Create necessary directory structure
- Generate configuration files
- Initialize blockchain
- Start services and monitor status

### 2.2 Binary Deployment

**Step 1: Download Binary Files**

Download binary files for the corresponding platform from the official release page:
- nogochain (node)
- nogopool (mining pool)
- nogocli (command-line tool)
- nogominer (standalone mining tool)

**Step 2: Configuration Files**

Create JSON configuration files (examples):

#### Node Configuration File (node_config.json)
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

#### Mining Pool Configuration File (mining_pool_config.json)
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

**Step 3: Start Node**

```bash
# Linux/macOS
chmod +x ./nogochain
./nogochain --datadir mainnet/data --config mainnet/config/node_config.json

# Windows
./nogochain.exe --datadir mainnet/data --config mainnet/config/node_config.json
```

**Step 4: Start Mining Pool**

```bash
# Linux/macOS
chmod +x ./nogopool
./nogopool --datadir mainnet/data --config mainnet/config/mining_pool_config.json

# Windows
./nogopool.exe --datadir mainnet/data --config mainnet/config/mining_pool_config.json
```

### 2.3 Source Code Compilation

**Step 1: Clone Code**

```bash
git clone https://github.com/nogochain/nogochain.git
cd nogochain
```

**Step 2: Verify Dependencies**

```bash
./scripts/verify-deps.sh
```

**Step 3: Compile**

```bash
# Linux/macOS
./scripts/build.sh

# Windows
./scripts/build.ps1

# Ubuntu System Detailed Compilation Steps
# 1. Install dependencies
sudo apt update
sudo apt install -y build-essential git curl

# 2. Install Go 1.22+
curl -LO https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 3. Compile to mainnet files
go build -o mainnet/nogochain ./cmd/nogochain
go build -o mainnet/nogocli ./cmd/nogocli
go build -o mainnet/nogominer ./cmd/nogominer
go build -o mainnet/nogopool ./cmd/nogopool
```

**Step 4: Start Node**

```bash
./bin/nogod --config config.toml

# Ubuntu System Startup Commands
# Start mainnet node
./mainnet/nogochain --datadir mainnet/data --config mainnet/config/node_config.json

# Background start mainnet node
nohup ./mainnet/nogochain --datadir mainnet/data --config mainnet/config/node_config.json > mainnet/logs/nogochain.log 2>&1 &
```

### 2.4 Container Deployment

**Step 1: Docker Build**

```bash
docker build -t nogochain .
```

**Step 2: Docker Run**

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

**Step 3: Docker Compose**

Create `docker-compose.yml`:

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

Start:
```bash
docker-compose up -d
```

## 3. Configuration Details

### 3.1 Node Configuration File (node_config.json)

| Configuration | Type | Default | Description |
|---------------|------|---------|-------------|
| chain_id | number | 318 | Chain ID |
| data_dir | string | "mainnet/data" | Data storage directory |
| log_dir | string | "mainnet/logs" | Log file directory |
| p2p_port | number | 30303 | P2P communication port |
| rpc_port | number | 8545 | HTTP RPC port |
| ws_port | number | 8546 | WebSocket port |
| max_peers | number | 50 | Maximum number of connected nodes |
| min_gas_price | string | "1000000000" | Minimum gas price |
| gas_limit | string | "8000000" | Block gas limit |
| metrics | boolean | true | Enable metrics monitoring |
| metrics_port | number | 6060 | Metrics monitoring port |
| jwt_secret | string | "mainnet/jwt-token.txt" | JWT secret file path |

### 3.2 Mining Pool Configuration File (mining_pool_config.json)

| Configuration | Type | Default | Description |
|---------------|------|---------|-------------|
| pool_name | string | "NogoChain Pool" | Mining pool name |
| node_url | string | "http://localhost:8545" | Node RPC address |
| stratum_port | number | 3333 | Stratum service port |
| stratum_tcp_port | number | 3334 | TCP Stratum port |
| stratum_tls_port | number | 3335 | TLS Stratum port |
| payment_interval | number | 600 | Payment interval (seconds) |
| min_payout | string | "1000000000000000000" | Minimum payout amount |
| fee | number | 0.01 | Mining pool fee |
| reward_address | string | "0xYourRewardAddress" | Reward address |
| stats_interval | number | 30 | Statistics update interval (seconds) |
| hashrate_window | number | 600 | Hashrate calculation window (seconds) |
| max_conns_per_ip | number | 10 | Maximum connections per IP |

### 3.3 Startup Script Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| --datadir | string | Script directory | Data directory path |
| --config | string | Auto-generated | Configuration file path |
| --genesis | string | Auto-generated | Genesis block file path |
| --port | number | 30303 | P2P port |
| --rpcport | number | 8545 | RPC port |
| --wsport | number | 8546 | WebSocket port |
| --metrics | boolean | true | Enable metrics monitoring |
| --verbosity | number | 3 | Log verbosity level |

## 4. Network Settings

### 4.1 Port Mapping

| Port | Protocol | Purpose | Must Open |
|------|----------|---------|----------|
| 30303 | TCP/UDP | P2P communication | Yes |
| 8545 | TCP | HTTP RPC | No (required for public access) |
| 8546 | TCP | WebSocket RPC | No (required for public access) |
| 3333 | TCP | Stratum service | No (required for mining) |

### 4.2 Firewall Settings

**Linux (iptables)**:

```bash
iptables -A INPUT -p tcp --dport 30303 -j ACCEPT
iptables -A INPUT -p udp --dport 30303 -j ACCEPT
# To open RPC
iptables -A INPUT -p tcp --dport 8545 -j ACCEPT
```

**Windows (Firewall)**:

Add inbound rules through Windows Firewall Advanced Settings to open the corresponding ports.

### 4.3 Node Discovery

- **Automatic Discovery**: Automatically discover nodes in the network through Kademlia algorithm
- **Manual Addition**: Add nodes using `admin_addPeer` RPC method
- **Bootstrap Nodes**: Set bootstrap nodes in configuration file to accelerate discovery

## 5. Data Management

### 5.1 Data Directory Structure

```
data/
├── chaindata/         # Block data
├── trie/              # State tree data
├── nodes/             # P2P node information
├── keystore/          # Account keys
└── logs/              # Log files
```

### 5.2 Data Backup

**Regular Backup**:
- Back up `keystore` directory (contains account private keys)
- Back up `chaindata` directory (optional, for quick recovery)

**Backup Strategy**:
- Cold Backup: Store backup data offline
- Encrypted Backup: Encrypt backup data
- Multiple Backups: Store in different locations

### 5.3 Data Cleaning

**Clean Old Data**:
```bash
# Stop node
# Delete data directory
rm -rf data/chaindata/*
# Restart node (resync)
```

**Note**: After cleaning data, you need to resync the entire blockchain, which may take a long time.

## 6. Node Management

### 6.1 Start and Stop

**Start**:
```bash
# Frontend start
./nogod --config config.toml

# Background start
nohup ./nogod --config config.toml > nogod.log 2>&1 &
```

**Stop**:
```bash
# Stop via RPC
curl -X POST --data '{"jsonrpc":"2.0","method":"admin_stop","params":[],"id":1}' http://localhost:8545

# Or use process management
kill -SIGINT $(pgrep nogod)
```

### 6.2 Monitoring and Logs

**Log Configuration**:
```toml
[logger]
level = "info"  # debug, info, warn, error
format = "json"  # json, console
file = "logs/nogochain.log"
maxsize = 100  # MB
maxage = 7  # days
```

**View Logs**:
```bash
# Real-time view
 tail -f logs/nogochain.log

# View errors
 grep "ERROR" logs/nogochain.log
```

### 6.3 Health Check

**RPC Check**:
```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545
```

**P2P Check**:
```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' http://localhost:8545
```

## 7. Security Configuration

### 7.1 RPC Security

**Production Environment Recommendations**:
- Restrict RPC access address: `--rpc.http.addr=127.0.0.1`
- Enable JWT authentication: `--rpc.jwtsecret=jwt-token.txt`
- Disable debug interfaces: `--rpc.disable-debug`

**JWT Generation**:
```bash
openssl rand -hex 32 > jwt-token.txt
```

### 7.2 Node Security

**Recommended Measures**:
- Regularly update node version
- Use firewall to restrict port access
- Avoid exposing ports other than P2P to the public network
- Run node with non-root account

### 7.3 Key Security

**Private Key Management**:
- Do not store private keys in plaintext in configuration files
- Use hardware wallet to manage mining addresses
- Regularly change mining addresses

## 8. Common Issues

### 8.1 Startup Issues

- **Issue**: Node cannot start
  **Solution**: Check configuration file, view log error messages

- **Issue**: Port occupied
  **Solution**: Modify port in configuration file, or stop the process occupying the port

### 8.2 Synchronization Issues

- **Issue**: Slow synchronization
  **Solution**: Ensure sufficient network bandwidth, add more peer nodes

- **Issue**: Synchronization stuck
  **Solution**: Restart node, check network connection

### 8.3 Performance Issues

- **Issue**: High CPU usage
  **Solution**: Adjust mining thread count, optimize system configuration

- **Issue**: High memory usage
  **Solution**: Increase memory, adjust cache configuration

### 8.4 Network Issues

- **Issue**: Cannot connect to other nodes
  **Solution**: Check firewall settings, ensure ports are open

- **Issue**: P2P connection count is 0
  **Solution**: Add bootstrap nodes, check network connection

## 9. Upgrade and Maintenance

### 9.1 Version Upgrade

**Step 1: Backup Data**

```bash
cp -r data data_backup
```

**Step 2: Stop Node**

```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"admin_stop","params":[],"id":1}' http://localhost:8545
```

**Step 3: Update Binary Files**

```bash
# Download new version or recompile
```

**Step 4: Start Node**

```bash
./nogod --config config.toml
```

### 9.2 Daily Maintenance

- **Monitoring**: Regularly check node status and performance
- **Updates**: Follow official releases of new versions and security patches
- **Backup**: Regularly back up critical data
- **Optimization**: Adjust configuration based on node operation

## 10. References

- [NogoChain Official Documentation](https://docs.nogochain.org)
- [Production Environment Deployment Tutorial](DEPLOYMENT_PRODUCTION.md) - Detailed production environment deployment guide
- [Test Network Deployment Tutorial](DEPLOYMENT_TESTNET.md) - Quick test network deployment guide
- [API Documentation](API.md) - RPC interface details
- [Architecture Documentation](ARCHITECTURE.md) - System architecture design
- [Mining Documentation](MINING.md) - Mining configuration and optimization
- [Go Language Installation Guide](https://golang.org/doc/install)
- [Docker Official Documentation](https://docs.docker.com)
- [Linux System Optimization](https://wiki.archlinux.org/title/Optimization)