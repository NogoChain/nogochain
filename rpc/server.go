package rpc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang-jwt/jwt/v5"
	"nogochain/network/config"
)

// Server represents the RPC server
type Server struct {
	server     *http.Server
	rpcServer  *rpc.Server
	nonceStore map[string]uint64
	nonceMutex sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	config     *config.RPCConfig
}

// NewServer creates a new RPC server
func NewServer(cfg *config.RPCConfig) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{
		nonceStore: make(map[string]uint64),
		ctx:        ctx,
		cancel:     cancel,
		config:     cfg,
	}

	// Create RPC server
	rpcServer := rpc.NewServer()

	// Register services
	ethService := NewEthService()
	nogService := NewNetService()
	web3Service := NewWeb3Service()
	debugService := NewDebugService()
	nogoService := NewNogoService()

	rpcServer.RegisterName("eth", ethService)
	rpcServer.RegisterName("net", nogService)
	rpcServer.RegisterName("web3", web3Service)
	rpcServer.RegisterName("debug", debugService)
	rpcServer.RegisterName("nogo", nogoService)

	server.rpcServer = rpcServer
	return server
}

// generateJWTToken 生成JWT令牌
func (s *Server) generateJWTToken() (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 365).Unix(), // 1年有效期
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	// 保存令牌到文件
	if err := os.WriteFile(s.config.JWT.TokenFile, []byte(tokenString), 0644); err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateJWTToken 验证JWT令牌
func (s *Server) validateJWTToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	return err == nil && token.Valid
}

// jwtAuthMiddleware JWT认证中间件
func (s *Server) jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 本地连接不需要认证
		if r.RemoteAddr == "127.0.0.1:0" || strings.Contains(r.RemoteAddr, "[::1]") {
			next.ServeHTTP(w, r)
			return
		}

		// 从请求头获取令牌
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 提取令牌
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		if !s.validateJWTToken(tokenString) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Start starts the RPC server
func (s *Server) Start() error {
	// 生成JWT令牌（如果启用且密钥已设置）
	if s.config.JWT.Enabled && s.config.JWT.Secret != "" {
		_, err := s.generateJWTToken()
		if err != nil {
			return fmt.Errorf("failed to generate JWT token: %v", err)
		}
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.rpcServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// 如果启用JWT认证，添加中间件
	if s.config.JWT.Enabled && s.config.JWT.Secret != "" {
		s.server.Handler = s.jwtAuthMiddleware(s.server.Handler)
	}

	fmt.Printf("RPC server started on %s\n", addr)
	return s.server.ListenAndServe()
}

// Stop stops the RPC server
func (s *Server) Stop() error {
	s.cancel()
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// GetNonce gets the nonce for an address
func (s *Server) GetNonce(addr string) uint64 {
	s.nonceMutex.Lock()
	defer s.nonceMutex.Unlock()
	return s.nonceStore[addr]
}

// SetNonce sets the nonce for an address
func (s *Server) SetNonce(addr string, nonce uint64) {
	s.nonceMutex.Lock()
	defer s.nonceMutex.Unlock()
	s.nonceStore[addr] = nonce
}
