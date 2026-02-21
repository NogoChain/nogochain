# NogoChain Mining Documentation

## 1. Mining Overview

NogoChain uses the self-developed NogoPow consensus mechanism, specifically designed for CPU/GPU mining, resistant to ASICs and performance-friendly.

### 1.1 Core Parameters

- **Consensus Algorithm**: NogoPow (Proof of Work)
- **Block Time**: Target 20 seconds/block
- **Difficulty Adjustment**: Every 10 blocks, ±50%
- **Block Reward**: Initial 8 NOGO, 20% reduction every 5 million blocks, minimum 0.1 NOGO
- **Transaction Fees**: All go to miners

### 1.2 Mining Advantages

- **ASIC Resistance**: Random memory access pattern, dataset ≤2GB
- **CPU/GPU Friendly**: No special instructions, multi-threaded optimization
- **Low Power Consumption**: Lower power consumption compared to traditional PoW algorithms
- **Fairness**: Ordinary hardware can participate, high degree of decentralization

## 2. Mining Methods

### 2.1 Built-in Mining

NogoChain nodes have built-in mining functionality that can be enabled through configuration files.

**Using Startup Script**:
```bash
# Windows
start-mainnet-node.bat --miner
```

**Manual Configuration**:
Modify `mainnet/config/node_config.json` file to add mining configuration:

**Startup Command**:
```bash
# Windows
./nogochain.exe --datadir mainnet/data --config mainnet/config/node_config.json --mine --miner.threads=4 --miner.coinbase=0xYOUR_ADDRESS
```

### 2.2 External Mining

NogoChain supports the standard Stratum protocol and can use third-party mining software.

**Mining Software Configuration**:
- URL: `stratum+tcp://YOUR_NODE_IP:3333`
- Username: Wallet address
- Password: Any (can be used to identify miners)

### 2.3 Mining Pool

Use the official mining pool components to set up a mining pool:

**Quick Deployment**:
```bash
# Windows
start-mainnet-pool.bat
```

**Manual Configuration**:
Modify `mainnet/config/mining_pool_config.json` file to configure pool parameters:

**Startup Command**:
```bash
# Windows
./nogopool.exe --datadir mainnet/data --config mainnet/config/mining_pool_config.json
```

**Miner Connection**:
- Stratum URL: `stratum+tcp://POOL_IP:3333`
- Username: Wallet address
- Password: Miner identifier

## 3. Hardware Requirements

### 3.1 CPU Mining

- **Recommended Configuration**:
  - 4+ core CPU
  - 8GB+ memory
  - SSD storage
- **Performance Reference**:
  - Intel i5: ~100-200 H/s
  - Intel i7: ~200-400 H/s
  - AMD Ryzen 5: ~300-500 H/s

### 3.2 GPU Mining

- **Recommended Configuration**:
  - NVIDIA GTX 1060 or AMD RX 580 or above
  - 8GB+ VRAM
  - Proper cooling
- **Performance Reference**:
  - NVIDIA GTX 1660: ~1000-1500 H/s
  - NVIDIA RTX 2060: ~2000-3000 H/s
  - AMD RX 5700 XT: ~3000-4000 H/s

### 3.3 Network Requirements

- **Bandwidth**: At least 1Mbps upload/download
- **Latency**: Lower is better, recommended < 50ms
- **Stability**: 7x24 hours stable operation

## 4. Mining Optimization

### 4.1 Software Optimization

- **Thread Adjustment**: Set according to CPU core count, usually physical core count
- **Memory Optimization**: Close unnecessary applications to free up memory
- **System Optimization**:
  - Disable system auto-updates
  - Disable unnecessary services
  - Optimize power plan to high performance

### 4.2 Hardware Optimization

- **Cooling**: Ensure good device cooling to avoid thermal throttling
- **Power Supply**: Use stable power supply to avoid voltage fluctuations
- **Network**: Use wired network connection to reduce network fluctuations

### 4.3 Mining Strategy

- **Continuous Operation**: Mining is a long-term process, recommended 24/7 operation
- **Regular Checks**: Monitor mining status and earnings
- **Risk Control**: Configure hardware reasonably to avoid over-investment

## 5. Mining Monitoring

### 5.1 Built-in Monitoring

NogoChain nodes provide the following monitoring interfaces:

- **RPC Interface**: `nogo_getMiningInfo` to get mining status
- **RPC Interface**: `nogo_getDifficulty` to get current difficulty
- **RPC Interface**: `nogo_getWork` to get mining work
- **Logs**: View node logs to understand mining status
- **Metrics**: Monitor mining performance through Prometheus

### 5.2 External Monitoring

- **Mining Software**: Use mining software that supports NogoChain for monitoring
- **Pool Monitoring**: If joining a pool, use monitoring tools provided by the pool
- **Third-party Tools**: Can use Grafana + Prometheus to build monitoring dashboards
- **Block Explorer**: Monitor mining status and earnings through block explorers

### 5.3 Monitoring Metrics

**Key Metrics**:
- **Hashrate**: Number of hashes calculated per second
- **Difficulty**: Current network difficulty
- **Block Height**: Current block height
- **Confirmation Time**: Transaction confirmation time
- **Earnings**: Mining earnings statistics

## 6. Earnings Calculation

### 6.1 Theoretical Earnings

**Calculation Formula**:
```
Daily Earnings = (Your Hashrate / Network Total Hashrate) * Daily Total Reward
Daily Total Reward = 8 NOGO * 4320 blocks/day * Reduction Factor
```

**Example**:
- Network Hashrate: 100 MH/s
- Your Hashrate: 1 MH/s
- Reduction Factor: 1.0 (initial stage)
- Daily Earnings: (1/100) * 8 * 4320 = 345.6 NOGO

### 6.2 Actual Earnings

Actual earnings may be affected by the following factors:
- Network difficulty adjustment
- Block confirmation time fluctuations
- Transaction fees
- Mining software efficiency
- Hardware stability

## 7. Common Issues

### 7.1 Mining Startup Issues

- **Issue**: Mining not starting
  **Solution**: Check configuration file, ensure miner.enabled = true

- **Issue**: Hashrate is 0
  **Solution**: Check hardware resources, ensure CPU/GPU is working properly

### 7.2 Connection Issues

- **Issue**: Cannot connect to Stratum server
  **Solution**: Check network connection, confirm Stratum service is started

- **Issue**: Node synchronization is slow
  **Solution**: Ensure sufficient network bandwidth, use fast sync mode

### 7.3 Earnings Issues

- **Issue**: Earnings lower than expected
  **Solution**: Check network hashrate changes, optimize mining settings

- **Issue**: No rewards received
  **Solution**: Confirm blocks are confirmed, check if wallet address is correct

## 8. Security Considerations

### 8.1 Hardware Security

- Avoid overclocking to prevent hardware damage
- Ensure good device ventilation to avoid fire risks
- Regularly check hardware status to prevent failures

### 8.2 Software Security

- Only download mining software from official channels
- Regularly update mining software and node versions
- Be vigilant against phishing websites and malware

### 8.3 Wallet Security

- Use hardware wallets to store large amounts of NOGO
- Regularly back up wallet private keys
- Avoid accessing wallets on public devices

## 9. Future Development

### 9.1 Mining Algorithm Upgrades

The NogoChain team will continue to optimize the NogoPow algorithm to maintain its ASIC resistance and efficiency.

### 9.2 Ecosystem Incentives

- Future mining incentive programs may be launched
- Support for decentralized application ecosystem development
- Exploration of sustainable mining economic models

## 10. References

- [NogoChain Official Documentation](https://docs.nogochain.org)
- [Production Environment Deployment Tutorial](DEPLOYMENT_PRODUCTION.md) - Detailed production environment deployment guide
- [Test Network Deployment Tutorial](DEPLOYMENT_TESTNET.md) - Quick test network deployment guide
- [API Documentation](API.md) - RPC interface details
- [Architecture Documentation](ARCHITECTURE.md) - System architecture design
- [Deployment Documentation](DEPLOYMENT.md) - Node deployment and configuration
- [Mining Software Recommendations](https://github.com/nogochain/mining-tools)
- [Hardware Compatibility List](https://github.com/nogochain/hardware-compatibility)