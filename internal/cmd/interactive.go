package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Bacchusgift/sql-diff/internal/ai"
	"github.com/Bacchusgift/sql-diff/internal/config"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// ModeOption 功能模式选项
type ModeOption struct {
	Name        string
	Description string
	Value       int
	NeedAI      bool
	Enabled     bool
}

// showModeMenu 显示功能选择菜单（使用光标选择）
func showModeMenu(aiEnabled bool) (int, error) {
	// 定义功能选项
	options := []ModeOption{
		{
			Name:        "SQL 表结构比对",
			Description: "比较两个表结构差异，自动生成 DDL 补全语句",
			Value:       1,
			NeedAI:      false,
			Enabled:     true,
		},
		{
			Name:        "AI 生成 CREATE TABLE",
			Description: "根据自然语言描述，AI 生成完整的建表语句",
			Value:       2,
			NeedAI:      true,
			Enabled:     aiEnabled,
		},
		{
			Name:        "AI 生成 ALTER TABLE",
			Description: "基于现有表结构 + 自然语言描述，AI 生成 DDL 变更语句",
			Value:       3,
			NeedAI:      true,
			Enabled:     aiEnabled,
		},
	}

	// 创建选择模板
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F449 {{ .Name | cyan | bold }}{{ if not .Enabled }} {{ `[需要 AI]` | red }}{{ end }}",
		Inactive: "  {{ .Name | white }}{{ if not .Enabled }} {{ `[需要 AI]` | faint }}{{ end }}",
		Selected: "\U00002705 {{ .Name | green | bold }}",
		Details: `
--------- 功能说明 ---------
{{ "描述:" | faint }} {{ .Description }}{{ if .NeedAI }}
{{ "要求:" | faint }} 需要启用 AI 功能{{ end }}`,
	}

	// 创建选择器
	prompt := promptui.Select{
		Label:     "📋 请选择功能模式",
		Items:     options,
		Templates: templates,
		Size:      3,
		CursorPos: 0,
	}

	// 执行选择
	idx, _, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("选择被取消: %v", err)
	}

	selected := options[idx]

	// 检查是否启用
	if !selected.Enabled {
		fmt.Println()
		errorColor.Println("✗ 该功能需要启用 AI")
		fmt.Println()
		fmt.Println("请通过以下方式之一启用 AI：")
		fmt.Println("  1. 配置文件: 编辑 .sql-diff-config.yaml，设置 ai.enabled: true")
		fmt.Println("  2. 命令行参数: 使用 --ai 参数启动")
		fmt.Println()
		fmt.Println("配置示例：")
		fmt.Println("  sql-diff config  # 运行配置向导")
		return 0, fmt.Errorf("AI 功能未启用")
	}

	fmt.Println()
	return selected.Value, nil
}

// runCompareMode SQL 表结构比对模式
func runCompareMode(cfg *config.Config) error {
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       模式 1: SQL 表结构比对")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 读取源表 SQL
	color.New(color.FgYellow, color.Bold).Println("📋 请粘贴源表的 CREATE TABLE 语句：")
	color.New(color.FgWhite).Println("（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）")
	color.New(color.FgCyan).Println("（提示：建议在文本编辑器中准备好 SQL，然后直接粘贴）")
	fmt.Println()

	sourceSQL, err := readMultilineInput()
	if err != nil {
		return fmt.Errorf("读取源表 SQL 失败: %v", err)
	}

	if strings.TrimSpace(sourceSQL) == "" {
		return fmt.Errorf("源表 SQL 不能为空")
	}

	successColor.Printf("✓ 已读取 %d 个字符\n", len(sourceSQL))
	fmt.Println()

	// 读取目标表 SQL
	color.New(color.FgYellow, color.Bold).Println("📋 请粘贴目标表的 CREATE TABLE 语句：")
	color.New(color.FgWhite).Println("（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）")
	color.New(color.FgCyan).Println("（提示：建议在文本编辑器中准备好 SQL，然后直接粘贴）")
	fmt.Println()

	targetSQL, err := readMultilineInput()
	if err != nil {
		return fmt.Errorf("读取目标表 SQL 失败: %v", err)
	}

	if strings.TrimSpace(targetSQL) == "" {
		return fmt.Errorf("目标表 SQL 不能为空")
	}

	successColor.Printf("✓ 已读取 %d 个字符\n", len(targetSQL))
	fmt.Println()

	// 调用核心比对逻辑
	return processComparison(sourceSQL, targetSQL, cfg)
}

// readSingleLineInput 读取单行输入（支持完整的一行文本，包含空格）
func readSingleLineInput() (string, error) {
	// 使用 bufio.Scanner 更稳定
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("读取输入失败")
	}
	return strings.TrimSpace(scanner.Text()), nil
}

// runGenerateTableMode AI 生成 CREATE TABLE 模式
func runGenerateTableMode(cfg *config.Config) error {
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       模式 2: AI 生成 CREATE TABLE")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 读取自然语言描述
	color.New(color.FgYellow, color.Bold).Println("💬 请描述您要创建的表结构：")
	color.New(color.FgCyan).Println("（示例：创建用户表，包含 ID、用户名、邮箱、密码、创建时间）")
	fmt.Println()

	color.New(color.FgWhite).Print("描述: ")
	description, err := readSingleLineInput()
	if err != nil {
		return fmt.Errorf("读取描述失败: %v", err)
	}

	description = strings.TrimSpace(description)
	if description == "" {
		return fmt.Errorf("描述不能为空")
	}

	fmt.Println()
	infoColor.Println("🤖 正在使用 AI 生成 SQL...")

	// 创建 AI Provider
	provider, err := ai.NewProvider(&cfg.AI)
	if err != nil {
		errorColor.Printf("✗ AI 初始化失败: %v\n", err)
		return err
	}

	// 调用 AI 生成 SQL
	sql, err := provider.GenerateCreateTable(description)
	if err != nil {
		errorColor.Printf("✗ 生成失败: %v\n", err)
		return err
	}

	// 显示结果
	fmt.Println()
	successColor.Println("✓ 生成成功！")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	color.New(color.FgWhite, color.Bold).Println("📋 生成的 CREATE TABLE 语句:")
	fmt.Println()
	fmt.Println(sql + ";")
	fmt.Println()

	// 询问是否保存到文件
	return askSaveToFile(sql + ";\n")
}

// runGenerateAlterMode AI 生成 ALTER TABLE 模式
func runGenerateAlterMode(cfg *config.Config) error {
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	infoColor.Println("       模式 3: AI 生成 ALTER TABLE")
	infoColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 读取现有表结构
	color.New(color.FgYellow, color.Bold).Println("📋 请粘贴现有表的 CREATE TABLE 语句：")
	color.New(color.FgWhite).Println("（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）")
	color.New(color.FgCyan).Println("（提示：建议在文本编辑器中准备好 SQL，然后直接粘贴）")
	fmt.Println()

	currentDDL, err := readMultilineInput()
	if err != nil {
		return fmt.Errorf("读取表结构失败: %v", err)
	}

	if strings.TrimSpace(currentDDL) == "" {
		return fmt.Errorf("表结构不能为空")
	}

	successColor.Printf("✓ 已读取 %d 个字符\n", len(currentDDL))
	fmt.Println()

	// 读取修改需求描述
	color.New(color.FgYellow, color.Bold).Println("💬 请描述您要做的修改：")
	color.New(color.FgCyan).Println("（示例：添加手机号字段、邮箱改为唯一索引）")
	fmt.Println()

	color.New(color.FgWhite).Print("描述: ")
	description, err := readSingleLineInput()
	if err != nil {
		return fmt.Errorf("读取描述失败: %v", err)
	}

	description = strings.TrimSpace(description)
	if description == "" {
		return fmt.Errorf("描述不能为空")
	}

	fmt.Println()
	infoColor.Println("🤖 正在使用 AI 生成 SQL...")

	// 创建 AI Provider
	provider, err := ai.NewProvider(&cfg.AI)
	if err != nil {
		errorColor.Printf("✗ AI 初始化失败: %v\n", err)
		return err
	}

	// 调用 AI 生成 SQL
	sql, err := provider.GenerateAlterTable(currentDDL, description)
	if err != nil {
		errorColor.Printf("✗ 生成失败: %v\n", err)
		return err
	}

	// 显示结果
	fmt.Println()
	successColor.Println("✓ 生成成功！")
	successColor.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	color.New(color.FgWhite, color.Bold).Println("📋 生成的 ALTER TABLE 语句:")
	fmt.Println()

	// 处理多条 SQL 语句
	sqlStatements := strings.Split(sql, "\n")
	var output strings.Builder
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			fmt.Println(stmt + ";")
			output.WriteString(stmt + ";\n")
		}
	}
	fmt.Println()

	// 询问是否保存到文件
	return askSaveToFile(output.String())
}

// askSaveToFile 询问用户是否保存到文件
func askSaveToFile(content string) error {
	color.New(color.FgYellow).Print("是否保存到文件? [y/N]: ")

	input, err := readSingleLineInput()
	if err != nil {
		return nil // 忽略错误，不影响主流程
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		infoColor.Println("未保存到文件")
		return nil
	}

	// 读取文件名
	color.New(color.FgYellow).Print("请输入文件名: ")
	filename, err := readSingleLineInput()
	if err != nil {
		errorColor.Printf("✗ 读取文件名失败: %v\n", err)
		return nil
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		errorColor.Println("✗ 文件名不能为空")
		return nil
	}

	// 写入文件
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		errorColor.Printf("✗ 写入文件失败: %v\n", err)
		return nil
	}

	successColor.Printf("✓ SQL 已保存到: %s\n", filename)
	return nil
}
