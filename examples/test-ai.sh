#!/bin/bash

# AI 功能集成测试脚本
# 注意：需要有效的 DeepSeek API Key

echo "======================================"
echo "SQL-Diff AI 功能集成测试"
echo "======================================"
echo ""

# 检查配置文件
if [ ! -f .sql-diff-config.yaml ]; then
    echo "⚠️  未找到配置文件 .sql-diff-config.yaml"
    echo "请先创建配置文件并填入有效的 API Key"
    echo ""
    echo "示例配置:"
    echo "---"
    echo "ai:"
    echo "  enabled: true"
    echo "  provider: deepseek"
    echo "  api_key: sk-your-api-key-here"
    echo "  api_endpoint: https://api.deepseek.com/v1"
    echo "  model: deepseek-chat"
    echo "  timeout: 30"
    echo "---"
    exit 1
fi

echo "✓ 找到配置文件"
echo ""

# 测试用例 1: 简单的新增字段
echo "测试 1: 新增字段 + AI 分析"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)" \
  --ai

echo ""
echo ""

# 测试用例 2: 复杂的表结构变更
echo "测试 2: 复杂表结构变更 + AI 分析"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE products (id INT, name TEXT, price FLOAT, stock INT)" \
  -t "CREATE TABLE products (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(200) NOT NULL, description TEXT, price DECIMAL(10,2) NOT NULL DEFAULT 0.00, stock INT UNSIGNED DEFAULT 0, status TINYINT DEFAULT 1, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, INDEX idx_status (status), INDEX idx_created (created_at)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4" \
  --ai

echo ""
echo ""

# 测试用例 3: 索引优化
echo "测试 3: 索引变更 + AI 分析"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE orders (id BIGINT PRIMARY KEY, user_id INT, product_id INT, total DECIMAL(10,2), created_at TIMESTAMP)" \
  -t "CREATE TABLE orders (id BIGINT PRIMARY KEY, user_id INT NOT NULL, product_id INT NOT NULL, total DECIMAL(12,2) NOT NULL, status VARCHAR(20) NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, INDEX idx_user_id (user_id), INDEX idx_product_id (product_id), INDEX idx_status (status), INDEX idx_created (created_at))" \
  --ai

echo ""
echo "======================================"
echo "测试完成！"
echo "======================================"
echo ""
echo "💡 提示:"
echo "- 如果看到 AI 分析结果，说明集成成功"
echo "- 如果显示错误，请检查 API Key 是否有效"
echo "- 如果超时，可以增加配置文件中的 timeout 值"
