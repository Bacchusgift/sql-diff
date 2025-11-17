package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
)

//go:embed static
var staticFiles embed.FS

// Server Web 服务器
type Server struct {
	port       int
	httpServer *http.Server
	detector   *ServiceDetector
}

// NewServer 创建 Web 服务器
func NewServer(port int) (*Server, error) {
	detector, err := NewServiceDetector()
	if err != nil {
		return nil, err
	}

	return &Server{
		port:     port,
		detector: detector,
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	// 创建路由
	mux := http.NewServeMux()

	// 静态资源
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("加载静态资源失败: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API 路由
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/config/export", s.handleConfigExport)
	mux.HandleFunc("/api/compare", s.handleCompare)
	mux.HandleFunc("/api/parse-create", s.handleParseCreate)
	mux.HandleFunc("/api/ai/generate-create", s.handleAIGenerateCreate)
	mux.HandleFunc("/api/ai/generate-alter", s.handleAIGenerateAlter)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.corsMiddleware(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 保存服务信息
	if err := s.detector.SaveServiceInfo(os.Getpid(), s.port); err != nil {
		return fmt.Errorf("保存服务信息失败: %w", err)
	}

	// 启动服务器
	go func() {
		successColor := color.New(color.FgGreen, color.Bold)
		infoColor := color.New(color.FgCyan)

		fmt.Println()
		successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		successColor.Println("       SQL-Diff Web 服务已启动")
		successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		infoColor.Printf("🌐 访问地址: http://%s\n", addr)
		infoColor.Printf("📝 PID: %d\n", os.Getpid())
		fmt.Println()
		infoColor.Println("提示: 按 Ctrl+C 停止服务")
		fmt.Println()
		successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorColor := color.New(color.FgRed, color.Bold)
			errorColor.Printf("✗ 服务器错误: %v\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	return s.Shutdown()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown() error {
	infoColor := color.New(color.FgCyan)
	infoColor.Println("\n正在关闭服务...")

	// 清理 PID 文件
	if err := s.detector.Cleanup(); err != nil {
		return err
	}

	// 关闭 HTTP 服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("关闭服务器失败: %w", err)
	}

	successColor := color.New(color.FgGreen, color.Bold)
	successColor.Println("✓ 服务已关闭")
	return nil
}

// corsMiddleware CORS 中间件
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 仅允许本地访问
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:"+fmt.Sprint(s.port))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleHealth 健康检查
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
