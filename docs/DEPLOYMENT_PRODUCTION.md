# NogoChain Production Environment Deployment Tutorial

## 1. Hardware Requirements

### 1.1 Node Server Requirements
- **CPU**: At least 4-core processor, recommended 8-core or above
- **Memory**: At least 16GB RAM, recommended 32GB or above
- **Storage**: At least 500GB SSD, recommended 1TB NVMe SSD
- **Network**: At least 100Mbps bandwidth, recommended 1Gbps or above
- **Operating System**: Windows Server 2016+ or Linux (Ubuntu 20.04+)

### 1.2 Mining Pool Server Requirements
- **CPU**: At least 8-core processor, recommended 16-core or above
- **Memory**: At least 32GB RAM, recommended 64GB or above
- **Storage**: At least 1TB SSD, recommended 2TB NVMe SSD
- **Network**: At least 1Gbps bandwidth, recommended 10Gbps or above
- **Operating System**: Windows Server 2016+ or Linux (Ubuntu 20.04+)

## 2. Network Configuration

### 2.1 Port Settings
| Service | Default Port | Purpose | Need Public Access |
|---------|-------------|---------|-------------------|
| Node P2P | 30303 | Inter-node communication | Yes |
| Node RPC | 8545 | JSON-RPC interface | No (local access only) |
| Node WebSocket | 8546 | WebSocket interface | No (local access only) |
| Pool P2P | 30304 | Mining pool node communication | Yes |
| Pool RPC | 8547 | Mining pool JSON-RPC interface | No (local access only) |
| Pool Stratum | 3333 | Miner connection port | Yes |

### 2.2 Firewall Configuration
- Allow inbound connections for P2P ports and Stratum port
- Restrict RPC ports to local access only
- Enable firewall logging

## 3. Deployment Preparation

### 3.1 Environment Check
1. Check Go version (requires Go 1.22+)
2. Check system dependencies
3. Verify network connectivity

### 3.2 Directory Structure
```
nogochain/
├── mainnet/
│   ├── config/          # Configuration files
│   ├── data/            # Blockchain data
│   ├── logs/            # Log files
│   ├── start_node.bat   # Node startup script
│   └── start_pool.bat   # Mining pool startup script
├── build/               # Compilation artifacts
├── docs/                # Documentation
└── scripts/             # Tool scripts
```

## 4. Node Deployment Steps

### 4.1 Download and Installation
1. Download the latest version of NogoChain from the official repository
2. Extract to the specified directory
3. Verify file integrity

### 4.2 Configuration File Generation
1. Run `start-mainnet-node.bat` script
2. The script will automatically generate default configuration files
3. Modify configuration parameters according to the actual environment

### 4.3 Initialize Blockchain
1. The script will automatically initialize the genesis block
2. Verify initialization result

### 4.4 Start Node
1. Execute startup script
2. Monitor log output
3. Verify node synchronization status

### 4.5 Node Monitoring
1. Configure Prometheus and Grafana
2. Set up monitoring metrics
3. Configure alert rules

## 5. Mining Pool Deployment Steps

### 5.1 Download and Installation
1. Download the latest version of NogoChain mining pool components from the official repository
2. Extract to the specified directory
3. Verify file integrity

### 5.2 Configuration File Generation
1. Run `start-mainnet-pool.bat` script
2. The script will automatically generate default configuration files
3. Modify configuration parameters according to the actual environment

### 5.3 Connect to Node
1. Configure mining pool to connect to local or remote node
2. Verify connection status

### 5.4 Start Mining Pool
1. Execute startup script
2. Monitor log output
3. Verify mining pool running status

### 5.5 Mining Pool Management
1. Configure miner access parameters
2. Set up mining difficulty adjustment
3. Configure reward distribution rules

## 6. Security Settings

### 6.1 System Security
- Update system patches
- Disable unnecessary services
- Configure strong password policies
- Enable multi-factor authentication

### 6.2 Network Security
- Use firewall to restrict access
- Configure VPN for internal service access
- Enable SSL/TLS encryption
- Perform regular security scans

### 6.3 Blockchain Security
- Use JWT authentication to protect RPC interfaces
- Configure reasonable gas limits
- Monitor abnormal transactions
- Regularly back up wallets and configurations

## 7. Common Issues and Solutions

### 7.1 Node Synchronization Issues
- Check network connection
- Increase memory and storage resources
- Adjust synchronization parameters

### 7.2 Mining Pool Connection Issues
- Check Stratum port configuration
- Verify node connection status
- Adjust mining pool parameters

### 7.3 Performance Optimization
- Adjust database cache size
- Optimize network parameters
- Configure appropriate log levels

### 7.4 Security Issues
- Regularly update software versions
- Monitor abnormal access
- Perform security audits

## 8. Maintenance and Upgrade

### 8.1 Daily Maintenance
- Monitor system resource usage
- Check log files
- Back up important data

### 8.2 Version Upgrade
- Download the latest version
- Stop services
- Replace binary files
- Start services and verify

### 8.3 Data Backup
- Regularly back up blockchain data
- Back up configuration files and wallets
- Test recovery process

## 9. Troubleshooting

### 9.1 Log Analysis
- View error logs
- Analyze warning messages
- Identify abnormal patterns

### 9.2 Common Error Codes
| Error Code | Description | Solution |
|------------|-------------|----------|
| 404 | Resource not found | Check path and configuration |
| 500 | Internal server error | View detailed logs |
| 503 | Service unavailable | Check service status |
| 429 | Too many requests | Adjust request frequency |

### 9.3 Emergency Recovery
- Stop services
- Restore backup data
- Start services
- Verify recovery result

## 10. Best Practices

### 10.1 Hardware Selection
- Use SSD storage for better performance
- Configure sufficient memory
- Choose stable network connection

### 10.2 Network Configuration
- Use static IP addresses
- Configure appropriate MTU values
- Enable network QoS

### 10.3 Monitoring and Alerting
- Configure comprehensive monitoring
- Set reasonable alert thresholds
- Establish response procedures

### 10.4 Security Hardening
- Regularly update systems and software
- Implement least privilege principle
- Encrypt sensitive data

## 11. Contact Information

### 11.1 Technical Support
- Official Documentation: https://github.com/nogochain/docs
- Community Forum: https://forum.nogochain.org
- Technical Support: support@nogochain.org

### 11.2 Emergency Support
- Emergency Contact Phone: +86-123-4567-8910
- Emergency Support Email: emergency@nogochain.org

## 12. Appendix

### 12.1 Configuration File Examples

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

### 12.2 Startup Script Parameters

#### Node Startup Script Parameters
- `--datadir`: Data directory path
- `--config`: Configuration file path
- `--genesis`: Genesis block file path
- `--port`: P2P port
- `--rpcport`: RPC port
- `--wsport`: WebSocket port
- `--metrics`: Enable metrics monitoring
- `--verbosity`: Log verbosity level

#### Mining Pool Startup Script Parameters
- `--datadir`: Data directory path
- `--config`: Configuration file path
- `--nodeargs`: Node parameters
- `--poolargs`: Mining pool parameters
- `--help`: Display help information

### 12.3 Monitoring Metrics

#### Node Monitoring Metrics
- `nogo_block_height`: Current block height
- `nogo_peers_count`: Number of connected peer nodes
- `nogo_gas_used`: Gas used
- `nogo_transactions_count`: Number of transactions
- `nogo_sync_status`: Synchronization status

#### Mining Pool Monitoring Metrics
- `nogo_pool_hashrate`: Mining pool hashrate
- `nogo_pool_miners`: Number of connected miners
- `nogo_pool_shares`: Number of shares submitted
- `nogo_pool_rewards`: Reward amount
- `nogo_pool_payments`: Number of payments

### 12.4 Troubleshooting Commands

#### Node Troubleshooting
- `nogocli.exe status`: View node status
- `nogocli.exe block`: View latest block
- `nogocli.exe peers`: View peer nodes
- `nogocli.exe txpool`: View transaction pool

#### Mining Pool Troubleshooting
- `nogocli.exe pool status`: View mining pool status
- `nogocli.exe pool miners`: View connected miners
- `nogocli.exe pool stats`: View mining pool statistics
- `nogocli.exe pool payments`: View payment records

## 13. Version History

| Version | Date | Changes |
|---------|------|----------|
| v1.0.0 | 2024-01-01 | Initial version |
| v1.0.1 | 2024-01-15 | Performance optimization |
| v1.0.2 | 2024-02-01 | Security update |
| v1.1.0 | 2024-03-01 | Feature enhancement |

## 14. Disclaimer

This tutorial is for reference only. Actual deployment should be adjusted according to specific environments. During deployment, please ensure compliance with relevant laws and regulations, and pay attention to security protection measures. The NogoChain team is not responsible for any losses caused by using this tutorial.

---

**Document Version**: v1.1.0
**Last Updated**: 2024-03-01
**Author**: NogoChain Team