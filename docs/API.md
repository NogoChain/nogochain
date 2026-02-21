# NogoChain API Documentation

## 1. Overview

NogoChain provides standard EVM-compatible RPC interfaces, along with extended interfaces for NogoChain-specific features. All interfaces are accessed through the HTTP/JSON-RPC protocol, with the default port being 8545.

## 2. Standard Interfaces

### 2.1 eth_* Interfaces

| Method | Parameters | Return Value | Description |
|--------|------------|--------------|-------------|
| eth_accounts | None | Array<String> | Get local account list |
| eth_blockNumber | None | String (Quantity) | Get latest block number |
| eth_call | Object, String | String | Call contract method (read-only) |
| eth_chainId | None | String (Quantity) | Get chain ID (318) |
| eth_coinbase | None | String | Get mining address |
| eth_estimateGas | Object | String (Quantity) | Estimate transaction gas consumption |
| eth_gasPrice | None | String (Quantity) | Get current gas price |
| eth_getBalance | String, String | String (Quantity) | Get account balance |
| eth_getBlockByHash | String, Boolean | Object | Get block by hash |
| eth_getBlockByNumber | String, Boolean | Object | Get block by number |
| eth_getBlockTransactionCountByHash | String | String (Quantity) | Get block transaction count |
| eth_getBlockTransactionCountByNumber | String | String (Quantity) | Get block transaction count |
| eth_getCode | String, String | String | Get contract code |
| eth_getLogs | Object | Array<Object> | Get logs |
| eth_getStorageAt | String, String, String | String | Get storage value |
| eth_getTransactionByHash | String | Object | Get transaction by hash |
| eth_getTransactionByBlockHashAndIndex | String, String | Object | Get transaction by block hash and index |
| eth_getTransactionByBlockNumberAndIndex | String, String | Object | Get transaction by block number and index |
| eth_getTransactionCount | String, String | String (Quantity) | Get transaction count (nonce) |
| eth_getTransactionReceipt | String | Object | Get transaction receipt |
| eth_hashrate | None | String (Quantity) | Get hashrate |
| eth_mining | None | Boolean | Check if mining |
| eth_sendRawTransaction | String | String | Send raw transaction |
| eth_submitHashrate | String, String | Boolean | Submit hashrate |
| eth_submitWork | String, String, String | Boolean | Submit mining result |

### 2.2 net_* Interfaces

| Method | Parameters | Return Value | Description |
|--------|------------|--------------|-------------|
| net_listening | None | Boolean | Check if node is listening |
| net_peerCount | None | String (Quantity) | Get connected peer count |
| net_version | None | String | Get network version |

### 2.3 web3_* Interfaces

| Method | Parameters | Return Value | Description |
|--------|------------|--------------|-------------|
| web3_clientVersion | None | String | Get client version |
| web3_sha3 | String | String | Calculate Keccak-256 hash of data |

### 2.4 debug_* Interfaces

| Method | Parameters | Return Value | Description |
|--------|------------|--------------|-------------|
| debug_accountRange | String, String, Number | Object | Get account range |
| debug_backtraceAt | String | Array | Get backtrace at specified position |
| debug_blockProfile | String, Number | Boolean | Generate block profile file |
| debug_cpuProfile | String, Number | Boolean | Generate CPU profile file |
| debug_dumpBlock | String | Object | Export block data |
| debug_gcStats | None | Object | Get garbage collection statistics |
| debug_getBlockRlp | String | String | Get RLP encoding of block |
| debug_memStats | None | Object | Get memory statistics |
| debug_seedHash | Number | String | Get seed hash of specified block |
| debug_setHead | String | Boolean | Set chain head |
| debug_standardTraceBadBlockToFile | String, String | Boolean | Trace bad block to file |
| debug_standardTraceBlockToFile | String, String | Boolean | Trace block to file |
| debug_startCPUProfile | String | Boolean | Start CPU profiling |
| debug_stopCPUProfile | None | Boolean | Stop CPU profiling |
| debug_traceBlock | String, Object | Object | Trace block execution |
| debug_traceBlockByNumber | String, Object | Object | Trace block by number |
| debug_traceBlockByHash | String, Object | Object | Trace block by hash |
| debug_traceTransaction | String, Object | Object | Trace transaction execution |
| debug_verbosity | Number | Boolean | Set log level |

## 3. NogoChain Specific Interfaces (nogo_*)

### 3.1 Consensus Related

| Method | Parameters | Return Value | Description |
|--------|------------|--------------|-------------|
| nogo_getDifficulty | None | String (Quantity) | Get current difficulty |
| nogo_getReward | None | String (Quantity) | Get current block reward |
| nogo_getMiningInfo | None | Object | Get mining information (difficulty, hashrate, reward) |
| nogo_getWork | None | Array<String> | Get mining work (NogoPow) |
| nogo_submitWork | String, String, String | Boolean | Submit mining work (NogoPow) |
| nogo_submitHashrate | String, String | Boolean | Submit hashrate (NogoPow) |

### 3.2 Chain Information

| Method | Parameters | Return Value | Description |
|--------|------------|--------------|-------------|
| nogo_getChainInfo | None | Object | Get chain information (chain ID, symbol, consensus algorithm, etc.) |

## 4. Interface Examples

### 4.1 Example Requests

#### 4.1.1 Standard Interface Example

```json
{
  "jsonrpc": "2.0",
  "method": "eth_blockNumber",
  "params": [],
  "id": 1
}
```

#### 4.1.2 NogoChain Specific Interface Example

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

### 4.2 Example Responses

#### 4.2.1 Standard Interface Response

```json
{
  "jsonrpc": "2.0",
  "result": "0x1234",
  "id": 1
}
```

#### 4.2.2 NogoChain Specific Interface Response

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

### 4.3 Error Response

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

## 5. RPC Configuration

### 5.1 Default Configuration

- **HTTP**: `http://127.0.0.1:8545`
- **WebSocket**: `ws://127.0.0.1:8546`
- **Authentication**: Local access requires no authentication, public access requires JWT

### 5.2 Security Recommendations

- Restrict RPC access IPs in production environment
- Enable JWT authentication to protect public RPC
- Avoid exposing sensitive information in RPC

## 6. Tool Integration

### 6.1 Development Tools

- **Hardhat**: Use standard EVM configuration
- **Truffle**: Configure network as NogoChain
- **Foundry**: Use `--rpc-url` to specify NogoChain RPC
- **Remix**: Connect to NogoChain via "Deploy & Run Transactions"

### 6.2 Wallet Integration

- **MetaMask**: Add custom network (Chain ID: 318)
- **Trust Wallet**: Supports NogoChain
- **Ledger/Trezor**: Supported through standard EVM interfaces

## 7. Rate Limiting

- Standard RPC requests: 1000 requests/minute per IP
- Debug interfaces: 100 requests/minute per IP
- Limits can be adjusted through configuration file

## 8. Version Compatibility

| API Version | Compatible Chain Version | Notes |
|-------------|--------------------------|-------|
| v1.0 | All versions | Basic interfaces |
| v1.1 | ≥ 1.0.0 | Added nogo_ interfaces |
| v1.2 | ≥ 1.1.0 | Performance optimization |

## 9. Common Issues

### 9.1 Connection Issues
- Ensure RPC service is started
- Check firewall settings
- Verify network configuration

### 9.2 Performance Issues
- For high-frequency requests, recommend using WebSocket
- Batch process multiple requests
- Avoid using expensive debug interfaces

### 9.3 Error Handling
- Check if parameter format is correct
- Verify transaction signature
- Confirm account balance is sufficient

## 10. References

- [Ethereum JSON-RPC Documentation](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [NogoChain Source Code](https://github.com/nogochain/nogochain)