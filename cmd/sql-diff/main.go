package main

import (
	"github.com/Bacchusgift/sql-diff/internal/cmd"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// 设置版本信息
	cmd.SetVersion(version, buildTime, gitCommit)
	
	// 执行命令
	cmd.Execute()
}
