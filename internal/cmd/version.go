package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// versionCmd 版本命令
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "显示版本信息",
	Long:    `显示 sql-diff 的版本号、构建时间和 Git 提交信息`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

// printVersion 打印版本信息
func printVersion() {
	titleColor := color.New(color.FgCyan, color.Bold)
	labelColor := color.New(color.FgWhite, color.Bold)
	valueColor := color.New(color.FgGreen)

	titleColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	titleColor.Println("       SQL-Diff 版本信息")
	titleColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	
	labelColor.Print("版本号:     ")
	valueColor.Println(version)
	
	labelColor.Print("构建时间:   ")
	valueColor.Println(buildTime)
	
	labelColor.Print("Git提交:    ")
	valueColor.Println(gitCommit)
	
	fmt.Println()
	titleColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
