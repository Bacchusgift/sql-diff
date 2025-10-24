package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Bacchusgift/sql-diff/internal/config"
)

// Provider AI 提供商接口
type Provider interface {
	// Analyze 分析表结构差异并提供建议
	Analyze(sourceDDL, targetDDL, diff string) (*AnalysisResult, error)

	// OptimizeSQL 优化 SQL 语句
	OptimizeSQL(sql string) (*OptimizationResult, error)

	// GenerateCreateTable 根据自然语言描述生成 CREATE TABLE 语句
	GenerateCreateTable(description string) (string, error)

	// GenerateAlterTable 根据自然语言描述和现有表结构生成 ALTER TABLE 语句
	GenerateAlterTable(currentDDL, description string) (string, error)
}

// AnalysisResult AI 分析结果
type AnalysisResult struct {
	Summary      string   // 差异摘要
	Suggestions  []string // 优化建议
	Risks        []string // 潜在风险
	BestPractice []string // 最佳实践建议
}

// OptimizationResult SQL 优化结果
type OptimizationResult struct {
	OriginalSQL  string   // 原始 SQL
	OptimizedSQL string   // 优化后的 SQL
	Improvements []string // 改进说明
}

// NewProvider 根据配置创建 AI 提供商
func NewProvider(cfg *config.AIConfig) (Provider, error) {
	if !cfg.Enabled {
		return &NoOpProvider{}, nil
	}

	switch cfg.Provider {
	case "deepseek":
		return NewDeepSeekProvider(cfg), nil
	case "openai":
		return NewOpenAIProvider(cfg), nil
	default:
		return nil, fmt.Errorf("不支持的 AI 提供商: %s", cfg.Provider)
	}
}

// NoOpProvider 空操作提供商（AI 未启用时使用）
type NoOpProvider struct{}

func (p *NoOpProvider) Analyze(sourceDDL, targetDDL, diff string) (*AnalysisResult, error) {
	return &AnalysisResult{
		Summary: "AI 功能未启用",
	}, nil
}

func (p *NoOpProvider) OptimizeSQL(sql string) (*OptimizationResult, error) {
	return &OptimizationResult{
		OriginalSQL:  sql,
		OptimizedSQL: sql,
	}, nil
}

func (p *NoOpProvider) GenerateCreateTable(description string) (string, error) {
	return "", fmt.Errorf("AI 功能未启用，无法生成 CREATE TABLE 语句")
}

func (p *NoOpProvider) GenerateAlterTable(currentDDL, description string) (string, error) {
	return "", fmt.Errorf("AI 功能未启用，无法生成 ALTER TABLE 语句")
}

// DeepSeekProvider DeepSeek AI 提供商
type DeepSeekProvider struct {
	config *config.AIConfig
	client *http.Client
}

// NewDeepSeekProvider 创建 DeepSeek 提供商
func NewDeepSeekProvider(cfg *config.AIConfig) *DeepSeekProvider {
	return &DeepSeekProvider{
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// ChatRequest DeepSeek API 请求结构
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse DeepSeek API 响应结构
type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// Analyze 分析表结构差异
func (p *DeepSeekProvider) Analyze(sourceDDL, targetDDL, diff string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(`你是一个资深的数据库架构师和 SQL 专家。请分析以下两个 MySQL 表结构的差异，并提供专业建议。

【源表结构】
%s

【目标表结构】
%s

【检测到的差异】
%s

请按以下格式提供分析（使用 Markdown 格式）：

## 差异分析
[简要总结主要差异]

## 优化建议
- [建议1]
- [建议2]
- [建议3]

## 潜在风险
- [风险1]
- [风险2]

## 最佳实践
- [实践1]
- [实践2]

请用中文回答，建议具体且可操作。`, sourceDDL, targetDDL, diff)

	response, err := p.chat(prompt)
	if err != nil {
		return nil, err
	}

	// 解析响应内容
	result := parseAnalysisResponse(response)
	return result, nil
}

// OptimizeSQL 优化 SQL 语句
func (p *DeepSeekProvider) OptimizeSQL(sql string) (*OptimizationResult, error) {
	prompt := fmt.Sprintf(`请优化以下 SQL DDL 语句，使其更符合最佳实践：

%s

请直接返回优化后的 SQL 语句，并说明改进之处。`, sql)

	response, err := p.chat(prompt)
	if err != nil {
		return nil, err
	}

	return &OptimizationResult{
		OriginalSQL:  sql,
		OptimizedSQL: response,
		Improvements: []string{"请查看 AI 优化建议"},
	}, nil
}

// GenerateCreateTable 根据自然语言描述生成 CREATE TABLE 语句
func (p *DeepSeekProvider) GenerateCreateTable(description string) (string, error) {
	prompt := fmt.Sprintf(`你是一个专业的数据库设计专家。请根据以下自然语言描述，生成一个标准的 MySQL CREATE TABLE 语句。

需求描述：
%s

要求：
1. 生成完整的、可直接执行的 MySQL CREATE TABLE 语句
2. 字段类型要合理（如 VARCHAR 指定长度、金额用 DECIMAL 等）
3. 主键、索引、默认值、注释等都要完善
4. 遵循 MySQL 最佳实践（如使用 InnoDB、UTF8MB4 编码等）
5. 只返回 SQL 语句，不要有其他解释文字
6. 表名和字段名使用蛇形命名法（snake_case）

请直接返回完整的 CREATE TABLE 语句：`, description)

	response, err := p.chat(prompt)
	if err != nil {
		return "", err
	}

	// 清理响应，提取 SQL 语句
	sql := cleanSQLResponse(response)
	return sql, nil
}

// GenerateAlterTable 根据自然语言描述和现有表结构生成 ALTER TABLE 语句
func (p *DeepSeekProvider) GenerateAlterTable(currentDDL, description string) (string, error) {
	prompt := fmt.Sprintf(`你是一个专业的数据库设计专家。请根据以下现有表结构和修改需求，生成相应的 MySQL ALTER TABLE 语句。

【现有表结构】
%s

【修改需求】
%s

要求：
1. 生成完整的、可直接执行的 MySQL ALTER TABLE 语句
2. 如果有多个修改，可以生成多条 ALTER TABLE 语句，每条一行
3. 考虑数据迁移的安全性（如修改字段类型要兼容现有数据）
4. 遵循 MySQL 最佳实践
5. 只返回 SQL 语句，不要有其他解释文字
6. 每条语句后不需要加分号，工具会自动添加

请直接返回 ALTER TABLE 语句：`, currentDDL, description)

	response, err := p.chat(prompt)
	if err != nil {
		return "", err
	}

	// 清理响应，提取 SQL 语句
	sql := cleanSQLResponse(response)
	return sql, nil
}

// chat 调用 DeepSeek Chat API
func (p *DeepSeekProvider) chat(prompt string) (string, error) {
	reqBody := ChatRequest{
		Model: p.config.Model,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "你是一个专业的数据库架构师和 SQL 专家，擅长表结构设计和优化。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	endpoint := p.config.APIEndpoint + "/chat/completions"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("API 返回空响应")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// OpenAIProvider OpenAI 兼容的提供商
type OpenAIProvider struct {
	*DeepSeekProvider // 复用 DeepSeek 实现（API 兼容）
}

// NewOpenAIProvider 创建 OpenAI 提供商
func NewOpenAIProvider(cfg *config.AIConfig) *OpenAIProvider {
	return &OpenAIProvider{
		DeepSeekProvider: NewDeepSeekProvider(cfg),
	}
}

// parseAnalysisResponse 解析 AI 分析响应
func parseAnalysisResponse(response string) *AnalysisResult {
	result := &AnalysisResult{
		Summary:      "",
		Suggestions:  make([]string, 0),
		Risks:        make([]string, 0),
		BestPractice: make([]string, 0),
	}

	// 使用简单的标记解析响应
	lines := strings.Split(response, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 识别章节
		if strings.HasPrefix(line, "##") {
			if strings.Contains(line, "差异分析") || strings.Contains(line, "差异总结") {
				currentSection = "summary"
			} else if strings.Contains(line, "优化建议") {
				currentSection = "suggestions"
			} else if strings.Contains(line, "潜在风险") || strings.Contains(line, "风险") {
				currentSection = "risks"
			} else if strings.Contains(line, "最佳实践") {
				currentSection = "best_practice"
			}
			continue
		}

		// 解析内容
		switch currentSection {
		case "summary":
			if !strings.HasPrefix(line, "#") {
				if result.Summary == "" {
					result.Summary = line
				} else {
					result.Summary += "\n" + line
				}
			}
		case "suggestions":
			if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				content := strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*")
				result.Suggestions = append(result.Suggestions, strings.TrimSpace(content))
			}
		case "risks":
			if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				content := strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*")
				result.Risks = append(result.Risks, strings.TrimSpace(content))
			}
		case "best_practice":
			if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				content := strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*")
				result.BestPractice = append(result.BestPractice, strings.TrimSpace(content))
			}
		}
	}

	// 如果没有解析到结构化内容，将整个响应作为摘要
	if result.Summary == "" && len(result.Suggestions) == 0 {
		result.Summary = response
	}

	return result
}

// cleanSQLResponse 清理 AI 响应，提取 SQL 语句
func cleanSQLResponse(response string) string {
	// 移除 Markdown 代码块标记
	response = strings.ReplaceAll(response, "```sql", "")
	response = strings.ReplaceAll(response, "```mysql", "")
	response = strings.ReplaceAll(response, "```", "")

	// 移除首尾空白
	response = strings.TrimSpace(response)

	// 如果响应中包含多行，尝试提取 CREATE TABLE 或 ALTER TABLE 部分
	lines := strings.Split(response, "\n")
	var sqlLines []string
	inSQL := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		upperLine := strings.ToUpper(trimmedLine)

		// 检测 SQL 语句开始
		if strings.HasPrefix(upperLine, "CREATE TABLE") ||
			strings.HasPrefix(upperLine, "ALTER TABLE") {
			inSQL = true
		}

		if inSQL {
			sqlLines = append(sqlLines, line)

			// 检测 SQL 语句结束（遇到分号）
			if strings.HasSuffix(trimmedLine, ";") {
				break
			}
		}
	}

	if len(sqlLines) > 0 {
		response = strings.Join(sqlLines, "\n")
	}

	// 移除末尾的分号（工具会统一添加）
	response = strings.TrimSuffix(strings.TrimSpace(response), ";")

	return response
}
