package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EthService represents the Ethereum RPC service
type EthService struct{}

// NewEthService creates a new Ethereum service
func NewEthService() *EthService {
	return &EthService{}
}

// ProtocolVersion returns the current Ethereum protocol version
func (s *EthService) ProtocolVersion() hexutil.Uint {
	return hexutil.Uint(67)
}

// Syncing returns the sync status
func (s *EthService) Syncing() interface{} {
	// Return false when not syncing
	return false
}

// Coinbase returns the coinbase address
func (s *EthService) Coinbase() common.Address {
	return common.HexToAddress("0x0000000000000000000000000000000000000000")
}

// Mining returns whether the node is mining
func (s *EthService) Mining() bool {
	return false
}

// Hashrate returns the current hashrate
func (s *EthService) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(0)
}

// GasPrice returns the current gas price
func (s *EthService) GasPrice() hexutil.Big {
	return hexutil.Big{}
}

// Accounts returns the list of accounts
func (s *EthService) Accounts() []common.Address {
	return []common.Address{}
}

// BlockNumber returns the current block number
func (s *EthService) BlockNumber() hexutil.Uint64 {
	return hexutil.Uint64(0)
}

// GetBalance returns the balance of an account
func (s *EthService) GetBalance(address common.Address, block string) hexutil.Big {
	return hexutil.Big{}
}

// GetStorageAt returns the storage at a given position
func (s *EthService) GetStorageAt(address common.Address, position string, block string) hexutil.Bytes {
	return hexutil.Bytes{}
}

// GetTransactionCount returns the transaction count for an account
func (s *EthService) GetTransactionCount(address common.Address, block string) hexutil.Uint64 {
	return hexutil.Uint64(0)
}

// GetBlockTransactionCountByHash returns the transaction count for a block by hash
func (s *EthService) GetBlockTransactionCountByHash(blockhash common.Hash) hexutil.Uint {
	return hexutil.Uint(0)
}

// GetBlockTransactionCountByNumber returns the transaction count for a block by number
func (s *EthService) GetBlockTransactionCountByNumber(block string) hexutil.Uint {
	return hexutil.Uint(0)
}

// GetUncleCountByBlockHash returns the uncle count for a block by hash
func (s *EthService) GetUncleCountByBlockHash(blockhash common.Hash) hexutil.Uint {
	return hexutil.Uint(0)
}

// GetUncleCountByBlockNumber returns the uncle count for a block by number
func (s *EthService) GetUncleCountByBlockNumber(block string) hexutil.Uint {
	return hexutil.Uint(0)
}

// GetCode returns the code at an address
func (s *EthService) GetCode(address common.Address, block string) hexutil.Bytes {
	return hexutil.Bytes{}
}

// Sign signs a message
func (s *EthService) Sign(address common.Address, data hexutil.Bytes) hexutil.Bytes {
	return hexutil.Bytes{}
}

// SendTransaction sends a transaction
func (s *EthService) SendTransaction(params map[string]interface{}) (common.Hash, error) {
	return common.Hash{}, nil
}

// SendRawTransaction sends a raw transaction
func (s *EthService) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	return common.Hash{}, nil
}

// Call executes a call
func (s *EthService) Call(params map[string]interface{}, block string) (hexutil.Bytes, error) {
	return hexutil.Bytes{}, nil
}

// EstimateGas estimates the gas needed for a transaction
func (s *EthService) EstimateGas(params map[string]interface{}) (hexutil.Uint64, error) {
	return hexutil.Uint64(21000), nil
}

// GetBlockByHash returns a block by hash
func (s *EthService) GetBlockByHash(blockhash common.Hash, full bool) map[string]interface{} {
	return nil
}

// GetBlockByNumber returns a block by number
func (s *EthService) GetBlockByNumber(block string, full bool) map[string]interface{} {
	return nil
}

// GetTransactionByHash returns a transaction by hash
func (s *EthService) GetTransactionByHash(txhash common.Hash) map[string]interface{} {
	return nil
}

// GetTransactionByBlockHashAndIndex returns a transaction by block hash and index
func (s *EthService) GetTransactionByBlockHashAndIndex(blockhash common.Hash, index hexutil.Uint) map[string]interface{} {
	return nil
}

// GetTransactionByBlockNumberAndIndex returns a transaction by block number and index
func (s *EthService) GetTransactionByBlockNumberAndIndex(block string, index hexutil.Uint) map[string]interface{} {
	return nil
}

// GetTransactionReceipt returns a transaction receipt
func (s *EthService) GetTransactionReceipt(txhash common.Hash) map[string]interface{} {
	return nil
}

// GetUncleByBlockHashAndIndex returns an uncle by block hash and index
func (s *EthService) GetUncleByBlockHashAndIndex(blockhash common.Hash, index hexutil.Uint) map[string]interface{} {
	return nil
}

// GetUncleByBlockNumberAndIndex returns an uncle by block number and index
func (s *EthService) GetUncleByBlockNumberAndIndex(block string, index hexutil.Uint) map[string]interface{} {
	return nil
}

// GetCompilers returns the list of compilers
func (s *EthService) GetCompilers() []string {
	return []string{}
}

// CompileLLL compiles LLL code
func (s *EthService) CompileLLL(source string) (hexutil.Bytes, error) {
	return hexutil.Bytes{}, nil
}

// CompileSolidity compiles Solidity code
func (s *EthService) CompileSolidity(source string) (map[string]interface{}, error) {
	return nil, nil
}

// CompileSerpent compiles Serpent code
func (s *EthService) CompileSerpent(source string) (hexutil.Bytes, error) {
	return hexutil.Bytes{}, nil
}

// NewFilter creates a new filter
func (s *EthService) NewFilter(params map[string]interface{}) (hexutil.Uint, error) {
	return hexutil.Uint(0), nil
}

// NewBlockFilter creates a new block filter
func (s *EthService) NewBlockFilter() (hexutil.Uint, error) {
	return hexutil.Uint(0), nil
}

// NewPendingTransactionFilter creates a new pending transaction filter
func (s *EthService) NewPendingTransactionFilter() (hexutil.Uint, error) {
	return hexutil.Uint(0), nil
}

// UninstallFilter uninstalls a filter
func (s *EthService) UninstallFilter(id hexutil.Uint) bool {
	return true
}

// GetFilterChanges returns filter changes
func (s *EthService) GetFilterChanges(id hexutil.Uint) interface{} {
	return []interface{}{}
}

// GetFilterLogs returns filter logs
func (s *EthService) GetFilterLogs(id hexutil.Uint) []interface{} {
	return []interface{}{}
}

// GetLogs returns logs
func (s *EthService) GetLogs(params map[string]interface{}) []interface{} {
	return []interface{}{}
}

// GetWork returns mining work
func (s *EthService) GetWork() []string {
	return []string{}
}

// SubmitWork submits mining work
func (s *EthService) SubmitWork(nonce string, powHash string, digest string) bool {
	return false
}

// SubmitHashrate submits hashrate
func (s *EthService) SubmitHashrate(hashrate hexutil.Uint64, id string) bool {
	return false
}

// GetProof returns a proof for a given key and block
func (s *EthService) GetProof(address common.Address, keys []string, block string) map[string]interface{} {
	return nil
}
