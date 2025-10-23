package ai
package ai

import (
	"testing"

	"github.com/Bacchusgift/sql-diff/internal/config"
)

// TestParseAnalysisResponse 测试响应解析
func TestParseAnalysisResponse(t *testing.T) {
	response := `## 差异分析
这是一个表结构变更，主要新增了email和created_at两个字段。

## 优化建议
- 建议为email字段添加唯一索引，确保邮箱唯一性
- created_at字段应添加DEFAULT CURRENT_TIMESTAMP
- 考虑添加updated_at字段用于记录更新时间

## 潜在风险
- 新增NOT NULL字段需要考虑现有数据的兼容性
- 大表添加字段可能会锁表，影响业务

## 最佳实践
- 使用TIMESTAMP而不是DATETIME节省存储空间
- 添加字段注释说明用途
- 考虑使用pt-online-schema-change等工具进行在线DDL变更`

	result := parseAnalysisResponse(response)

	// 验证摘要
	if result.Summary == "" {
		t.Error("应该解析到摘要内容")
	}

	// 验证建议
	if len(result.Suggestions) != 3 {
		t.Errorf("期望解析到3条建议，实际%d条", len(result.Suggestions))
	}

	// 验证风险
	if len(result.Risks) != 2 {
		t.Errorf("期望解析到2条风险，实际%d条", len(result.Risks))
	}

	// 验证最佳实践
	if len(result.BestPractice) != 3 {
		t.Errorf("期望解析到3条最佳实践，实际%d条", len(result.BestPractice))
	}
}

// TestNoOpProvider 测试空操作提供商
func TestNoOpProvider(t *testing.T) {
	provider := &NoOpProvider{}

	// 测试分析
	result, err := provider.Analyze("CREATE TABLE users (id INT)", "CREATE TABLE users (id INT, name VARCHAR(100))", "diff")
	if err != nil {
		t.Errorf("NoOpProvider不应返回错误: %v", err)
	}

	if result.Summary != "AI 功能未启用" {
		t.Errorf("NoOpProvider应返回未启用提示")
	}

	// 测试优化
	optResult, err := provider.OptimizeSQL("ALTER TABLE users ADD COLUMN email VARCHAR(255)")
	if err != nil {
		t.Errorf("NoOpProvider不应返回错误: %v", err)
	}

	if optResult.OptimizedSQL != optResult.OriginalSQL {
		t.Error("NoOpProvider不应修改SQL")
	}
}

// TestNewProvider 测试提供商创建
func TestNewProvider(t *testing.T) {
	// 测试未启用AI
	cfg := &config.AIConfig{
		Enabled: false,
	}
	provider, err := NewProvider(cfg)
	if err != nil {
		t.Errorf("创建Provider失败: %v", err)
	}
	if _, ok := provider.(*NoOpProvider); !ok {
		t.Error("未启用AI时应返回NoOpProvider")
	}

	// 测试DeepSeek提供商
	cfg = &config.AIConfig{
		Enabled:     true,
		Provider:    "deepseek",
		APIKey:      "test-key",
		APIEndpoint: "https://api.deepseek.com/v1",
		Model:       "deepseek-chat",
		Timeout:     30,
	}
	provider, err = NewProvider(cfg)
	if err != nil {
		t.Errorf("创建DeepSeek Provider失败: %v", err)
	}
	if _, ok := provider.(*DeepSeekProvider); !ok {
		t.Error("应返回DeepSeekProvider")
	}

	// 测试不支持的提供商
	cfg.Provider = "unsupported"
	_, err = NewProvider(cfg)
	if err == nil {
		t.Error("不支持的提供商应返回错误")
	}
}

// TestParseComplexResponse 测试复杂响应解析
func TestParseComplexResponse(t *testing.T) {
	// 测试只有摘要的响应
	response1 := "这是一个简单的变更，新增了两个字段。"
	result1 := parseAnalysisResponse(response1)
	if result1.Summary == "" {
		t.Error("应该将整个响应作为摘要")
	}

	// 测试包含多个章节的响应
	response2 := `## 差异分析
主要变更说明

## 优化建议
- 建议1
* 建议2

## 潜在风险
* 风险1
- 风险2`

	result2 := parseAnalysisResponse(response2)
	if len(result2.Suggestions) != 2 {
		t.Errorf("应解析到2条建议，实际%d条", len(result2.Suggestions))
	}
	if len(result2.Risks) != 2 {
		t.Errorf("应解析到2条风险，实际%d条", len(result2.Risks))
	}
}
