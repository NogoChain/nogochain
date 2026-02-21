# NogoChain (EVM+NogoPow)

## Overview
NogoChain is an EVM-compatible blockchain with a custom consensus algorithm (NogoPow) designed to be ASIC-resistant and CPU/GPU friendly.

## Core Features
- **ChainID**: 318, **Symbol**: NOGO, **Decimals**: 18
- **Base**: EVM-compatible client, Go 1.22+, No CGO
- **Consensus**: NogoPow (Anti-ASIC, CPU/GPU Friendly)
- **Performance**: ≤5ms/calc (desktop), ≤15ms (laptop)
- **Reward**: 8 NOGO, 5M20 (-20%/5M), min 0.1 NOGO

## Directory Structure
- `cmd/`: Command-line tools and main entry points
- `internal/`: Internal packages
- `pkg/`: Public packages
- `scripts/`: Build and utility scripts
- `docs/`: Documentation
- `tests/`: Test files

## Quick Start

### Prerequisites
- Go 1.22+
- CGO disabled (`CGO_ENABLED=0`)

### Build
1. Verify dependencies:
   ```bash
   ./scripts/verify-deps.sh
   ```

2. Build the binary:
   ```bash
   ./scripts/build.sh
   ```

3. Run the node:
   ```bash
   ./build/nogochain
   ```

## Documentation
- `docs/ARCHITECTURE.md`: System architecture
- `docs/API.md`: API documentation
- `docs/MINING.md`: Mining guide
- `docs/DEPLOYMENT.md`: Deployment instructions
- `docs/SECURITY.md`: Security considerations

## Compatibility
- **Wallets**: MetaMask, Trust Wallet, Ledger, Trezor
- **RPC**: Geth standard + nogo_ extensions
- **Explorers**: Blockscout/Etherscan API compatible
- **Tools**: Hardhat, Truffle, Foundry, Remix
- **Tokens**: ERC-20/721/1155 fully compatible

## Security & Operations
- No key custody, raw tx only
- RPC: 127.0.0.1, JWT for public access
- Log: zerolog JSON, 100MB/7d
- Metrics: Prometheus
- Test coverage ≥80%

## Code Style
- Go standard style
- Chinese comments for exports
- Package names lowercase
- Constant names uppercase
