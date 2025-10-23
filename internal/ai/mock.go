package ai

import (
	"fmt"

	"github.com/Bacchusgift/sql-diff/internal/config"
)

// MockProvider 模拟的 AI 提供商（用于测试）
type MockProvider struct {
	AnalyzeFunc     func(sourceDDL, targetDDL, diff string) (*AnalysisResult, error)
	OptimizeSQLFunc func(sql string) (*OptimizationResult, error)
}

// NewMockProvider 创建模拟提供商
func NewMockProvider() *MockProvider {
	return &MockProvider{
		// 默认实现
		AnalyzeFunc: func(sourceDDL, targetDDL, diff string) (*AnalysisResult, error) {
			return &AnalysisResult{
				Summary: "模拟分析：检测到表结构变更",
				Suggestions: []string{
					"建议为新增字段添加索引以提高查询性能",
					"建议为NOT NULL字段设置默认值",
					"考虑使用更合适的数据类型以节省存储空间",
				},
				Risks: []string{
					"大表添加字段可能导致锁表",
					"修改字段类型可能导致数据丢失",
				},
				BestPractice: []string{
					"使用pt-online-schema-change进行在线DDL",
					"在测试环境验证后再执行",
					"建议在业务低峰期执行变更",
				},
			}, nil
		},
		OptimizeSQLFunc: func(sql string) (*OptimizationResult, error) {
			return &OptimizationResult{
				OriginalSQL:  sql,
				OptimizedSQL: fmt.Sprintf("-- 优化后:\n%s\n-- 添加索引以提高性能", sql),
				Improvements: []string{
					"添加了性能优化索引",
					"调整了字段顺序",
				},
			}, nil
		},
	}
}

// Analyze 分析表结构差异
func (m *MockProvider) Analyze(sourceDDL, targetDDL, diff string) (*AnalysisResult, error) {
	if m.AnalyzeFunc != nil {
		return m.AnalyzeFunc(sourceDDL, targetDDL, diff)
	}
	return nil, fmt.Errorf("Analyze not implemented")
}

// OptimizeSQL 优化 SQL 语句
func (m *MockProvider) OptimizeSQL(sql string) (*OptimizationResult, error) {
	if m.OptimizeSQLFunc != nil {
		return m.OptimizeSQLFunc(sql)
	}
	return nil, fmt.Errorf("OptimizeSQL not implemented")
}

// NewMockProviderWithConfig 从配置创建 Mock 提供商
func NewMockProviderWithConfig(cfg *config.AIConfig) Provider {
	if !cfg.Enabled {
		return &NoOpProvider{}
	}
	return NewMockProvider()
}
