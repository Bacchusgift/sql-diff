package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Bacchusgift/sql-diff/internal/web"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	webPort int
)

// webCmd Web 命令
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动 Web 可视化界面",
	Long: `启动 SQL-Diff Web 可视化界面服务。

功能特点:
  ✓ 可视化配置管理
  ✓ Web 界面进行表结构比对
  ✓ Web 界面使用 AI 功能
  ✓ 自动打开浏览器
  ✓ 服务单例检测`,
	Example: `  # 启动 Web 服务
  sql-diff web
  
  # 指定端口启动
  sql-diff web --port 9000`,
	RunE: runWeb,
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntVar(&webPort, "port", 8848, "Web 服务端口")
}

func runWeb(cmd *cobra.Command, args []string) error {
	infoColor := color.New(color.FgCyan)
	successColor := color.New(color.FgGreen, color.Bold)
	errorColor := color.New(color.FgRed, color.Bold)

	fmt.Println()
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       SQL-Diff Web 启动")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 创建服务检测器
	detector, err := web.NewServiceDetector()
	if err != nil {
		errorColor.Printf("✗ 初始化失败: %v\n", err)
		return err
	}

	// 检测服务是否已运行
	running, info, err := detector.IsRunning()
	if err != nil {
		errorColor.Printf("✗ 检测服务状态失败: %v\n", err)
		return err
	}

	if running {
		// 服务已运行,直接打开浏览器
		successColor.Printf("✓ 服务已在运行中 (PID: %d, 端口: %d)\n", info.PID, info.Port)
		fmt.Println()
		infoColor.Println("🌐 正在打开浏览器...")

		url := fmt.Sprintf("http://127.0.0.1:%d", info.Port)
		if err := openBrowser(url); err != nil {
			fmt.Println()
			infoColor.Printf("请手动访问: %s\n", url)
		} else {
			fmt.Println()
			successColor.Println("✓ 浏览器已打开")
		}

		return nil
	}

	// 服务未运行,启动新服务
	port := webPort

	// 查找可用端口
	availablePort, err := web.FindAvailablePort(port, 10)
	if err != nil {
		errorColor.Printf("✗ %v\n", err)
		return err
	}

	if availablePort != port {
		infoColor.Printf("ℹ️  端口 %d 已被占用,使用端口 %d\n", port, availablePort)
		fmt.Println()
	}

	// 创建服务器
	server, err := web.NewServer(availablePort)
	if err != nil {
		errorColor.Printf("✗ 创建服务器失败: %v\n", err)
		return err
	}

	// 打开浏览器
	url := fmt.Sprintf("http://127.0.0.1:%d", availablePort)
	go func() {
		// 延迟一秒等待服务器启动
		// time.Sleep(1 * time.Second)
		if err := openBrowser(url); err != nil {
			infoColor.Printf("ℹ️  请手动访问: %s\n", url)
		} else {
			successColor.Println("✓ 浏览器已打开")
		}
	}()

	// 启动服务器(阻塞直到收到中断信号)
	return server.Start()
}

// openBrowser 打开默认浏览器
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
