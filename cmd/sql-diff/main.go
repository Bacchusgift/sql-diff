package main

import (
	"github.com/Bacchusgift/sql-diff/internal/cmd"
)

// 版本信息（由编译时注入）
var (
	Version   = "dev"     // 版本号
	BuildTime = "unknown" // 构建时间
	GitCommit = "unknown" // Git 提交哈希
)

func main() {
	// 传递版本信息给 cmd 包
	cmd.SetVersion(Version, BuildTime, GitCommit)
	cmd.Execute()
}
