package rpc

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DebugService represents the Debug RPC service
type DebugService struct{}

// NewDebugService creates a new Debug service
func NewDebugService() *DebugService {
	return &DebugService{}
}

// AccountRange returns the account range
func (s *DebugService) AccountRange(blockHash string, startAddr string, maxResults hexutil.Uint) map[string]interface{} {
	return nil
}

// BacktraceAt returns the backtrace at a given position
func (s *DebugService) BacktraceAt(data hexutil.Bytes, offset hexutil.Uint) []interface{} {
	return []interface{}{}
}

// BlockProfile returns the block profile
func (s *DebugService) BlockProfile(file string, n hexutil.Uint) bool {
	return false
}

// CPUProfile returns the CPU profile
func (s *DebugService) CPUProfile(file string, n hexutil.Uint) bool {
	return false
}

// ChaindbCompact compacts the chain database
func (s *DebugService) ChaindbCompact() bool {
	return false
}

// ChaindbProperty returns a chain database property
func (s *DebugService) ChaindbProperty(property string) string {
	return ""
}

// DumpBlock dumps a block
func (s *DebugService) DumpBlock(blockNum hexutil.Uint64) map[string]interface{} {
	return make(map[string]interface{})
}

// GcStats returns the GC stats
func (s *DebugService) GcStats() map[string]interface{} {
	return nil
}

// GoTrace returns the Go trace
func (s *DebugService) GoTrace(file string, n hexutil.Uint) bool {
	return false
}

// MemStats returns the memory stats
func (s *DebugService) MemStats() map[string]interface{} {
	return nil
}

// PrintBlock prints a block
func (s *DebugService) PrintBlock(blockNum hexutil.Uint64) string {
	return ""
}

// SeedHash returns the seed hash
func (s *DebugService) SeedHash(blockNum hexutil.Uint64) string {
	return ""
}

// SetHead sets the head block
func (s *DebugService) SetHead(blockNum hexutil.Uint64) bool {
	return false
}

// Stacks returns the stacks
func (s *DebugService) Stacks() string {
	return ""
}

// StartCPUProfile starts the CPU profile
func (s *DebugService) StartCPUProfile(file string) bool {
	return false
}

// StopCPUProfile stops the CPU profile
func (s *DebugService) StopCPUProfile() bool {
	return false
}

// StopGoTrace stops the Go trace
func (s *DebugService) StopGoTrace() bool {
	return false
}

// TraceBlock traces a block
func (s *DebugService) TraceBlock(blockHash string, options map[string]interface{}) []interface{} {
	return []interface{}{}
}

// TraceBlockByNumber traces a block by number
func (s *DebugService) TraceBlockByNumber(blockNum string, options map[string]interface{}) []interface{} {
	return []interface{}{}
}

// TraceBlockFromFile traces a block from a file
func (s *DebugService) TraceBlockFromFile(file string, options map[string]interface{}) []interface{} {
	return []interface{}{}
}

// TraceCall traces a call
func (s *DebugService) TraceCall(call map[string]interface{}, block string, options map[string]interface{}) map[string]interface{} {
	return nil
}

// TraceTransaction traces a transaction
func (s *DebugService) TraceTransaction(txHash string, options map[string]interface{}) []interface{} {
	return []interface{}{}
}

// Verbosity sets the verbosity
func (s *DebugService) Verbosity(level hexutil.Uint) bool {
	return false
}

// Vmodule sets the vmodule
func (s *DebugService) Vmodule(pattern string) bool {
	return false
}

// WriteBlockProfile writes the block profile
func (s *DebugService) WriteBlockProfile(file string) bool {
	return false
}

// WriteMemProfile writes the memory profile
func (s *DebugService) WriteMemProfile(file string) bool {
	return false
}
