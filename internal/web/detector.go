package web

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ServiceDetector 服务检测器
type ServiceDetector struct {
	pidFilePath string
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	PID  int
	Port int
}

// NewServiceDetector 创建服务检测器
func NewServiceDetector() (*ServiceDetector, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户主目录失败: %w", err)
	}

	sqlDiffDir := filepath.Join(homeDir, ".sql-diff")
	if err := os.MkdirAll(sqlDiffDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	return &ServiceDetector{
		pidFilePath: filepath.Join(sqlDiffDir, "web.pid"),
	}, nil
}

// IsRunning 检测服务是否正在运行
// 返回: (是否运行, 服务信息, 错误)
func (d *ServiceDetector) IsRunning() (bool, *ServiceInfo, error) {
	// 检查 PID 文件是否存在
	data, err := os.ReadFile(d.pidFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, fmt.Errorf("读取 PID 文件失败: %w", err)
	}

	// 解析 PID 文件内容: <PID>:<PORT>
	content := strings.TrimSpace(string(data))
	parts := strings.Split(content, ":")
	if len(parts) != 2 {
		// PID 文件格式错误,清理并返回未运行
		_ = d.Cleanup()
		return false, nil, nil
	}

	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		_ = d.Cleanup()
		return false, nil, nil
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		_ = d.Cleanup()
		return false, nil, nil
	}

	info := &ServiceInfo{PID: pid, Port: port}

	// 检查进程是否存在
	if !d.isProcessRunning(pid) {
		// 进程不存在,清理 PID 文件
		_ = d.Cleanup()
		return false, nil, nil
	}

	// 检查端口是否可访问
	if !d.isPortListening(port) {
		// 端口不可访问,清理 PID 文件
		_ = d.Cleanup()
		return false, nil, nil
	}

	return true, info, nil
}

// SaveServiceInfo 保存服务信息
func (d *ServiceDetector) SaveServiceInfo(pid, port int) error {
	content := fmt.Sprintf("%d:%d", pid, port)
	if err := os.WriteFile(d.pidFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 PID 文件失败: %w", err)
	}
	return nil
}

// Cleanup 清理 PID 文件
func (d *ServiceDetector) Cleanup() error {
	if err := os.Remove(d.pidFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("清理 PID 文件失败: %w", err)
	}
	return nil
}

// isProcessRunning 检查进程是否正在运行
func (d *ServiceDetector) isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// 发送信号 0 检查进程是否存在(不会真正发送信号)
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// isPortListening 检查端口是否正在监听
func (d *ServiceDetector) isPortListening(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// FindAvailablePort 查找可用端口
// 从 startPort 开始,尝试最多 maxAttempts 次
func FindAvailablePort(startPort, maxAttempts int) (int, error) {
	for i := 0; i < maxAttempts; i++ {
		port := startPort + i
		if isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("无法找到可用端口 (尝试范围: %d-%d)", startPort, startPort+maxAttempts-1)
}

// isPortAvailable 检查端口是否可用
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
