package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Spinner loading 动画显示器
type Spinner struct {
	message  string
	frames   []string
	delay    time.Duration
	active   bool
	stopChan chan bool
	mu       sync.Mutex
}

// NewSpinner 创建一个新的 spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames: []string{
			"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
		},
		delay:    100 * time.Millisecond,
		stopChan: make(chan bool),
	}
}

// Start 开始显示 loading 动画
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		cyan := color.New(color.FgCyan)
		i := 0
		for {
			select {
			case <-s.stopChan:
				return
			default:
				frame := s.frames[i%len(s.frames)]
				fmt.Printf("\r%s %s ", cyan.Sprint(frame), s.message)
				i++
				time.Sleep(s.delay)
			}
		}
	}()
}

// Stop 停止 loading 动画
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	s.active = false
	s.stopChan <- true
	// 清除当前行
	fmt.Print("\r\033[K")
}

// Update 更新 loading 消息
func (s *Spinner) Update(message string) {
	s.mu.Lock()
	s.message = message
	s.mu.Unlock()
}

// Success 显示成功消息并停止
func (s *Spinner) Success(message string) {
	s.Stop()
	green := color.New(color.FgGreen, color.Bold)
	fmt.Printf("%s %s\n", green.Sprint("✓"), message)
}

// Error 显示错误消息并停止
func (s *Spinner) Error(message string) {
	s.Stop()
	red := color.New(color.FgRed, color.Bold)
	fmt.Printf("%s %s\n", red.Sprint("✗"), message)
}

// Warning 显示警告消息并停止
func (s *Spinner) Warning(message string) {
	s.Stop()
	yellow := color.New(color.FgYellow, color.Bold)
	fmt.Printf("%s %s\n", yellow.Sprint("⚠"), message)
}
