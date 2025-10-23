#!/bin/bash

# Mock AI 功能演示脚本（无需真实 API Key）

echo "======================================"
echo "SQL-Diff Mock AI 功能演示"
echo "======================================"
echo ""
echo "注意：此演示使用模拟的 AI 响应，无需真实 API Key"
echo ""

# 创建临时配置文件（使用mock提供商）
cat > /tmp/sql-diff-mock-config.yaml << 'EOF'
ai:
  enabled: true
  provider: mock
  api_key: mock-api-key
  api_endpoint: http://mock.api
  model: mock-model
  timeout: 30
EOF

echo "✓ 创建模拟配置文件"
echo ""

# 由于当前实现不直接支持mock provider，我们创建一个示例输出
echo "示例 AI 分析输出："
echo "--------------------------------------"
echo ""
echo "📊 差异分析:"
echo "检测到表结构变更，主要新增了 email 和 created_at 两个字段。"
echo "这些变更符合常见的用户表设计模式。"
echo ""
echo "✨ 优化建议:"
echo "  1. 建议为 email 字段添加唯一索引以确保邮箱唯一性"
echo "  2. 建议为 NOT NULL 字段设置合理的默认值"
echo "  3. 考虑添加 updated_at 字段记录更新时间"
echo ""
echo "⚠️  潜在风险:"
echo "  1. 如果表中已有数据，新增 NOT NULL 字段需要提供默认值或迁移数据"
echo "  2. 大表添加字段可能导致长时间锁表，影响线上业务"
echo ""
echo "📖 最佳实践:"
echo "  1. 使用 pt-online-schema-change 等工具进行在线 DDL 变更"
echo "  2. 在测试环境充分验证后再应用到生产环境"
echo "  3. 建议在业务低峰期执行表结构变更"
echo "  4. 保留表结构变更的回滚方案"
echo ""
echo "--------------------------------------"
echo ""

# 实际测试（不使用AI）
echo "实际比对示例（不启用AI）:"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"

echo ""
echo ""
echo "💡 提示："
echo "- 上面展示的 AI 分析是模拟输出"
echo "- 要使用真实的 AI 分析，请配置有效的 DeepSeek API Key"
echo "- 运行 ./examples/test-ai.sh 进行真实 AI 测试"
echo ""

# 清理
rm -f /tmp/sql-diff-mock-config.yaml
echo "======================================"
echo "演示完成！"
echo "======================================"
