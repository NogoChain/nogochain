package rpc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"nogochain/network/config"
)

// 测试NewServer函数
func TestNewServer(t *testing.T) {
	// 创建RPC配置
	cfg := &config.RPCConfig{
		Enabled: true,
		Port:    8545,
		Host:    "127.0.0.1",
		JWT: &config.JWTConfig{
			Enabled:   false,
			Secret:    "test-secret",
			TokenFile: "jwt-token.txt",
		},
	}

	// 创建RPC服务器
	server := NewServer(cfg)
	if server == nil {
		t.Errorf("NewServer returned nil")
	}

	if server.config != cfg {
		t.Errorf("Server config mismatch")
	}

	if server.rpcServer == nil {
		t.Errorf("Server rpcServer should not be nil")
	}

	if len(server.nonceStore) != 0 {
		t.Errorf("Initial nonceStore should be empty")
	}
}

// 测试GetNonce和SetNonce方法
func TestNonceMethods(t *testing.T) {
	// 创建RPC服务器
	cfg := &config.RPCConfig{
		Enabled: true,
		Port:    8545,
		Host:    "127.0.0.1",
	}

	server := NewServer(cfg)

	// 测试初始nonce
	addr := "0x01"
	nonce := server.GetNonce(addr)
	if nonce != 0 {
		t.Errorf("Initial nonce should be 0, got %d", nonce)
	}

	// 测试设置nonce
	newNonce := uint64(5)
	server.SetNonce(addr, newNonce)
	nonce = server.GetNonce(addr)
	if nonce != newNonce {
		t.Errorf("Nonce should be %d, got %d", newNonce, nonce)
	}

	// 测试设置另一个地址的nonce
	addr2 := "0x02"
	server.SetNonce(addr2, 10)
	nonce2 := server.GetNonce(addr2)
	if nonce2 != 10 {
		t.Errorf("Nonce for addr2 should be 10, got %d", nonce2)
	}

	// 确保第一个地址的nonce不变
	nonce = server.GetNonce(addr)
	if nonce != newNonce {
		t.Errorf("Nonce for addr should still be %d, got %d", newNonce, nonce)
	}
}

// 测试EthService
func TestEthService(t *testing.T) {
	ethService := NewEthService()
	if ethService == nil {
		t.Errorf("NewEthService returned nil")
	}

	// 测试ProtocolVersion
	protocolVersion := ethService.ProtocolVersion()
	if protocolVersion != hexutil.Uint(67) {
		t.Errorf("ProtocolVersion should be 67, got %d", protocolVersion)
	}

	// 测试Syncing
	syncing := ethService.Syncing()
	if syncing != false {
		t.Errorf("Syncing should be false, got %v", syncing)
	}

	// 测试Coinbase
	coinbase := ethService.Coinbase()
	if coinbase != (common.Address{}) {
		t.Errorf("Coinbase should be zero address, got %v", coinbase)
	}

	// 测试Mining
	mining := ethService.Mining()
	if mining != false {
		t.Errorf("Mining should be false, got %v", mining)
	}

	// 测试Hashrate
	hashrate := ethService.Hashrate()
	if hashrate != hexutil.Uint64(0) {
		t.Errorf("Hashrate should be 0, got %d", hashrate)
	}

	// 测试GasPrice
	gasPrice := ethService.GasPrice()
	if gasPrice.ToInt().Sign() != 0 {
		t.Errorf("GasPrice should be zero, got %v", gasPrice)
	}

	// 测试Accounts
	accounts := ethService.Accounts()
	if len(accounts) != 0 {
		t.Errorf("Accounts should be empty, got %d accounts", len(accounts))
	}

	// 测试BlockNumber
	blockNumber := ethService.BlockNumber()
	if blockNumber != hexutil.Uint64(0) {
		t.Errorf("BlockNumber should be 0, got %d", blockNumber)
	}

	// 测试GetBalance
	address := common.HexToAddress("0x01")
	balance := ethService.GetBalance(address, "latest")
	if balance.ToInt().Sign() != 0 {
		t.Errorf("GetBalance should return zero, got %v", balance)
	}

	// 测试GetStorageAt
	storage := ethService.GetStorageAt(address, "0x01", "latest")
	if len(storage) != 0 {
		t.Errorf("GetStorageAt should return empty bytes, got %v", storage)
	}

	// 测试GetTransactionCount
	txCount := ethService.GetTransactionCount(address, "latest")
	if txCount != hexutil.Uint64(0) {
		t.Errorf("GetTransactionCount should be 0, got %d", txCount)
	}

	// 测试GetBlockTransactionCountByHash
	hash := common.Hash{}
	blockTxCount := ethService.GetBlockTransactionCountByHash(hash)
	if blockTxCount != hexutil.Uint(0) {
		t.Errorf("GetBlockTransactionCountByHash should be 0, got %d", blockTxCount)
	}

	// 测试GetBlockTransactionCountByNumber
	blockTxCountByNumber := ethService.GetBlockTransactionCountByNumber("latest")
	if blockTxCountByNumber != hexutil.Uint(0) {
		t.Errorf("GetBlockTransactionCountByNumber should be 0, got %d", blockTxCountByNumber)
	}

	// 测试GetUncleCountByBlockHash
	uncleCount := ethService.GetUncleCountByBlockHash(hash)
	if uncleCount != hexutil.Uint(0) {
		t.Errorf("GetUncleCountByBlockHash should be 0, got %d", uncleCount)
	}

	// 测试GetUncleCountByBlockNumber
	uncleCountByNumber := ethService.GetUncleCountByBlockNumber("latest")
	if uncleCountByNumber != hexutil.Uint(0) {
		t.Errorf("GetUncleCountByBlockNumber should be 0, got %d", uncleCountByNumber)
	}

	// 测试GetCode
	code := ethService.GetCode(address, "latest")
	if len(code) != 0 {
		t.Errorf("GetCode should return empty bytes, got %v", code)
	}

	// 测试Sign
	data := hexutil.Bytes{0x01, 0x02}
	signature := ethService.Sign(address, data)
	if len(signature) != 0 {
		t.Errorf("Sign should return empty bytes, got %v", signature)
	}

	// 测试SendTransaction
	txParams := map[string]interface{}{}
	txHash, err := ethService.SendTransaction(txParams)
	if err != nil {
		t.Errorf("SendTransaction returned error: %v", err)
	}

	if txHash != (common.Hash{}) {
		t.Errorf("SendTransaction should return zero hash, got %v", txHash)
	}

	// 测试SendRawTransaction
	txData := hexutil.Bytes{0x01, 0x02}
	txHash, err = ethService.SendRawTransaction(txData)
	if err != nil {
		t.Errorf("SendRawTransaction returned error: %v", err)
	}

	if txHash != (common.Hash{}) {
		t.Errorf("SendRawTransaction should return zero hash, got %v", txHash)
	}

	// 测试Call
	callParams := map[string]interface{}{}
	callResult, err := ethService.Call(callParams, "latest")
	if err != nil {
		t.Errorf("Call returned error: %v", err)
	}

	if len(callResult) != 0 {
		t.Errorf("Call should return empty bytes, got %v", callResult)
	}

	// 测试EstimateGas
	gasParams := map[string]interface{}{}
	gasEstimate, err := ethService.EstimateGas(gasParams)
	if err != nil {
		t.Errorf("EstimateGas returned error: %v", err)
	}

	if gasEstimate != hexutil.Uint64(21000) {
		t.Errorf("EstimateGas should be 21000, got %d", gasEstimate)
	}

	// 测试GetBlockByHash
	blockByHash := ethService.GetBlockByHash(hash, false)
	if blockByHash != nil {
		t.Errorf("GetBlockByHash should return nil, got %v", blockByHash)
	}

	// 测试GetBlockByNumber
	blockByNumber := ethService.GetBlockByNumber("latest", false)
	if blockByNumber != nil {
		t.Errorf("GetBlockByNumber should return nil, got %v", blockByNumber)
	}

	// 测试GetTransactionByHash
	txByHash := ethService.GetTransactionByHash(hash)
	if txByHash != nil {
		t.Errorf("GetTransactionByHash should return nil, got %v", txByHash)
	}

	// 测试GetTransactionByBlockHashAndIndex
	txByBlockHash := ethService.GetTransactionByBlockHashAndIndex(hash, hexutil.Uint(0))
	if txByBlockHash != nil {
		t.Errorf("GetTransactionByBlockHashAndIndex should return nil, got %v", txByBlockHash)
	}

	// 测试GetTransactionByBlockNumberAndIndex
	txByBlockNumber := ethService.GetTransactionByBlockNumberAndIndex("latest", hexutil.Uint(0))
	if txByBlockNumber != nil {
		t.Errorf("GetTransactionByBlockNumberAndIndex should return nil, got %v", txByBlockNumber)
	}

	// 测试GetTransactionReceipt
	txReceipt := ethService.GetTransactionReceipt(hash)
	if txReceipt != nil {
		t.Errorf("GetTransactionReceipt should return nil, got %v", txReceipt)
	}

	// 测试GetUncleByBlockHashAndIndex
	uncleByHash := ethService.GetUncleByBlockHashAndIndex(hash, hexutil.Uint(0))
	if uncleByHash != nil {
		t.Errorf("GetUncleByBlockHashAndIndex should return nil, got %v", uncleByHash)
	}

	// 测试GetUncleByBlockNumberAndIndex
	uncleByNumber := ethService.GetUncleByBlockNumberAndIndex("latest", hexutil.Uint(0))
	if uncleByNumber != nil {
		t.Errorf("GetUncleByBlockNumberAndIndex should return nil, got %v", uncleByNumber)
	}

	// 测试GetCompilers
	compilers := ethService.GetCompilers()
	if len(compilers) != 0 {
		t.Errorf("GetCompilers should return empty slice, got %v", compilers)
	}

	// 测试CompileLLL
	lllSource := "(return 1)"
	lllResult, err := ethService.CompileLLL(lllSource)
	if err != nil {
		t.Errorf("CompileLLL returned error: %v", err)
	}

	if len(lllResult) != 0 {
		t.Errorf("CompileLLL should return empty bytes, got %v", lllResult)
	}

	// 测试CompileSolidity
	soliditySource := "pragma solidity ^0.8.0; contract Test {}"
	solidityResult, err := ethService.CompileSolidity(soliditySource)
	if err != nil {
		t.Errorf("CompileSolidity returned error: %v", err)
	}

	if solidityResult != nil {
		t.Errorf("CompileSolidity should return nil, got %v", solidityResult)
	}

	// 测试CompileSerpent
	serpentSource := ""
	serpentResult, err := ethService.CompileSerpent(serpentSource)
	if err != nil {
		t.Errorf("CompileSerpent returned error: %v", err)
	}

	if len(serpentResult) != 0 {
		t.Errorf("CompileSerpent should return empty bytes, got %v", serpentResult)
	}

	// 测试NewFilter
	filterParams := map[string]interface{}{}
	filterID, err := ethService.NewFilter(filterParams)
	if err != nil {
		t.Errorf("NewFilter returned error: %v", err)
	}

	if filterID != hexutil.Uint(0) {
		t.Errorf("NewFilter should return 0, got %d", filterID)
	}

	// 测试NewBlockFilter
	blockFilterID, err := ethService.NewBlockFilter()
	if err != nil {
		t.Errorf("NewBlockFilter returned error: %v", err)
	}

	if blockFilterID != hexutil.Uint(0) {
		t.Errorf("NewBlockFilter should return 0, got %d", blockFilterID)
	}

	// 测试NewPendingTransactionFilter
	pendingTxFilterID, err := ethService.NewPendingTransactionFilter()
	if err != nil {
		t.Errorf("NewPendingTransactionFilter returned error: %v", err)
	}

	if pendingTxFilterID != hexutil.Uint(0) {
		t.Errorf("NewPendingTransactionFilter should return 0, got %d", pendingTxFilterID)
	}

	// 测试UninstallFilter
	uninstallResult := ethService.UninstallFilter(hexutil.Uint(0))
	if !uninstallResult {
		t.Errorf("UninstallFilter should return true, got %v", uninstallResult)
	}

	// 测试GetFilterChanges
	filterChanges := ethService.GetFilterChanges(hexutil.Uint(0))
	if filterChanges == nil {
		t.Errorf("GetFilterChanges should return empty slice, got nil")
	}

	// 测试GetFilterLogs
	filterLogs := ethService.GetFilterLogs(hexutil.Uint(0))
	if len(filterLogs) != 0 {
		t.Errorf("GetFilterLogs should return empty slice, got %v", filterLogs)
	}

	// 测试GetLogs
	logsParams := map[string]interface{}{}
	logs := ethService.GetLogs(logsParams)
	if len(logs) != 0 {
		t.Errorf("GetLogs should return empty slice, got %v", logs)
	}

	// 测试GetWork
	work := ethService.GetWork()
	if len(work) != 0 {
		t.Errorf("GetWork should return empty slice, got %v", work)
	}

	// 测试SubmitWork
	submitWorkResult := ethService.SubmitWork("0x00", "0x00", "0x00")
	if submitWorkResult != false {
		t.Errorf("SubmitWork should return false, got %v", submitWorkResult)
	}

	// 测试SubmitHashrate
	submitHashrateResult := ethService.SubmitHashrate(hexutil.Uint64(1000), "0x00")
	if submitHashrateResult != false {
		t.Errorf("SubmitHashrate should return false, got %v", submitHashrateResult)
	}

	// 测试GetProof
	proofKeys := []string{"0x01"}
	proof := ethService.GetProof(address, proofKeys, "latest")
	if proof != nil {
		t.Errorf("GetProof should return nil, got %v", proof)
	}
}

// 测试NetService
func TestNetService(t *testing.T) {
	netService := NewNetService()
	if netService == nil {
		t.Errorf("NewNetService returned nil")
	}

	// 测试Version
	version := netService.Version()
	if version != "318" {
		t.Errorf("Version should be '318', got '%s'", version)
	}

	// 测试Listening
	listening := netService.Listening()
	if listening != false {
		t.Errorf("Listening should be false, got %v", listening)
	}

	// 测试PeerCount
	peerCount := netService.PeerCount()
	if peerCount != 0 {
		t.Errorf("PeerCount should be 0, got %d", peerCount)
	}
}

// 测试Web3Service
func TestWeb3Service(t *testing.T) {
	web3Service := NewWeb3Service()
	if web3Service == nil {
		t.Errorf("NewWeb3Service returned nil")
	}

	// 测试ClientVersion
	clientVersion := web3Service.ClientVersion()
	if clientVersion != "NogoChain/v1.0.0/go1.22" {
		t.Errorf("ClientVersion should be 'NogoChain/v1.0.0/go1.22', got '%s'", clientVersion)
	}

	// 测试Sha3
	data := hexutil.Bytes{0x01, 0x02}
	hash := web3Service.Sha3(data)
	if len(hash) == 0 {
		t.Errorf("Sha3 should return non-empty bytes")
	}
}

// 测试DebugService
func TestDebugService(t *testing.T) {
	debugService := NewDebugService()
	if debugService == nil {
		t.Errorf("NewDebugService returned nil")
	}

	// 测试DumpBlock
	dumpBlock := debugService.DumpBlock(hexutil.Uint64(0))
	if dumpBlock == nil {
		t.Errorf("DumpBlock should return non-nil")
	}
}

// 测试NogoService
func TestNogoService(t *testing.T) {
	nogoService := NewNogoService()
	if nogoService == nil {
		t.Errorf("NewNogoService returned nil")
	}

	// 测试GetWork
	work := nogoService.GetWork()
	if work == nil {
		t.Errorf("GetWork should return non-nil")
	}

	// 测试SubmitWork
	submitWorkResult := nogoService.SubmitWork("0x00", "0x00", "0x00")
	if submitWorkResult != false {
		t.Errorf("SubmitWork should return false, got %v", submitWorkResult)
	}

	// 测试GetDifficulty
	difficulty := nogoService.GetDifficulty()
	if difficulty == 0 {
		t.Errorf("GetDifficulty should return non-zero value")
	}

	// 测试GetChainInfo
	chainInfo := nogoService.GetChainInfo()
	if chainInfo == nil {
		t.Errorf("GetChainInfo should return non-nil")
	}
}

// 测试JWT令牌生成和验证
func TestJWTToken(t *testing.T) {
	// 创建启用JWT的RPC配置
	cfg := &config.RPCConfig{
		Enabled: true,
		Port:    8545,
		Host:    "127.0.0.1",
		JWT: &config.JWTConfig{
			Enabled:   true,
			Secret:    "test-secret",
			TokenFile: "jwt-token.txt",
		},
	}

	server := NewServer(cfg)

	// 测试生成JWT令牌
	token, err := server.generateJWTToken()
	if err != nil {
		t.Errorf("generateJWTToken returned error: %v", err)
	}

	if token == "" {
		t.Errorf("generateJWTToken should return non-empty string")
	}

	// 测试验证有效的JWT令牌
	valid := server.validateJWTToken(token)
	if !valid {
		t.Errorf("validateJWTToken should return true for valid token")
	}

	// 测试验证无效的JWT令牌
	invalidToken := "invalid-token"
	invalid := server.validateJWTToken(invalidToken)
	if invalid {
		t.Errorf("validateJWTToken should return false for invalid token")
	}
}

// 集成测试：测试完整的RPC服务器流程
func TestRPCIntegration(t *testing.T) {
	// 创建RPC配置
	cfg := &config.RPCConfig{
		Enabled: true,
		Port:    8545,
		Host:    "127.0.0.1",
		JWT: &config.JWTConfig{
			Enabled:   false,
			Secret:    "test-secret",
			TokenFile: "jwt-token.txt",
		},
	}

	// 创建RPC服务器
	server := NewServer(cfg)
	if server == nil {
		t.Errorf("NewServer returned nil")
	}

	// 测试服务器配置
	if server.config != cfg {
		t.Errorf("Server config mismatch")
	}

	// 测试nonce管理
	addr := "0x01"
	server.SetNonce(addr, 5)
	nonce := server.GetNonce(addr)
	if nonce != 5 {
		t.Errorf("Nonce should be 5, got %d", nonce)
	}

	// 测试服务注册
	if server.rpcServer == nil {
		t.Errorf("RPC server should not be nil")
	}

	// 测试JWT功能（如果启用）
	if cfg.JWT.Enabled && cfg.JWT.Secret != "" {
		token, err := server.generateJWTToken()
		if err != nil {
			t.Errorf("generateJWTToken returned error: %v", err)
		}

		if token == "" {
			t.Errorf("generateJWTToken should return non-empty string")
		}

		valid := server.validateJWTToken(token)
		if !valid {
			t.Errorf("validateJWTToken should return true for valid token")
		}
	}

	// 测试所有服务
	ethService := NewEthService()
	if ethService == nil {
		t.Errorf("NewEthService returned nil")
	}

	netService := NewNetService()
	if netService == nil {
		t.Errorf("NewNetService returned nil")
	}

	web3Service := NewWeb3Service()
	if web3Service == nil {
		t.Errorf("NewWeb3Service returned nil")
	}

	debugService := NewDebugService()
	if debugService == nil {
		t.Errorf("NewDebugService returned nil")
	}

	nogoService := NewNogoService()
	if nogoService == nil {
		t.Errorf("NewNogoService returned nil")
	}
}
