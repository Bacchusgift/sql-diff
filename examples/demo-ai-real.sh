#!/bin/bash

# SQL-Diff 真实 AI 功能演示脚本
# 使用真实的 DeepSeek API 进行智能分析

echo "======================================"
echo "SQL-Diff 真实 AI 功能演示"
echo "======================================"
echo ""
echo "✓ 使用真实的 DeepSeek API"
echo "✓ 展示完整的 AI 分析能力"
echo ""

# 确保已构建
if [ ! -f bin/sql-diff ]; then
    echo "正在构建程序..."
    make build
    echo ""
fi

# 检查配置
if [ ! -f .sql-diff-config.yaml ]; then
    echo "❌ 错误：未找到配置文件 .sql-diff-config.yaml"
    echo "请先运行: cp .sql-diff-config.example.yaml .sql-diff-config.yaml"
    exit 1
fi

echo "======================================"
echo "演示 1: 用户表新增字段"
echo "======================================"
echo ""
echo "场景: 为用户表添加 email 和 created_at 字段"
echo ""
read -p "按回车继续..."
echo ""

./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)" \
  --ai

echo ""
echo ""
read -p "按回车继续下一个演示..."
echo ""

echo "======================================"
echo "演示 2: 产品表结构优化"
echo "======================================"
echo ""
echo "场景: 优化数据类型、添加约束和索引"
echo ""
read -p "按回车继续..."
echo ""

./bin/sql-diff \
  -s "CREATE TABLE products (id INT, name TEXT, price FLOAT, stock INT)" \
  -t "CREATE TABLE products (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(200) NOT NULL, description TEXT, price DECIMAL(10,2) NOT NULL DEFAULT 0.00, stock INT UNSIGNED DEFAULT 0, status TINYINT DEFAULT 1, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, INDEX idx_status (status), INDEX idx_created (created_at)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4" \
  --ai

echo ""
echo ""
read -p "按回车继续下一个演示..."
echo ""

echo "======================================"
echo "演示 3: 订单表索引优化 (带文件输出)"
echo "======================================"
echo ""
echo "场景: 添加索引和完整性约束，并保存 DDL 到文件"
echo ""
read -p "按回车继续..."
echo ""

./bin/sql-diff \
  -s "CREATE TABLE orders (id BIGINT PRIMARY KEY, user_id INT, product_id INT, total DECIMAL(10,2), created_at TIMESTAMP)" \
  -t "CREATE TABLE orders (id BIGINT PRIMARY KEY, user_id INT NOT NULL, product_id INT NOT NULL, total DECIMAL(12,2) NOT NULL, status VARCHAR(20) NOT NULL DEFAULT 'pending', created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, INDEX idx_user_id (user_id), INDEX idx_product_id (product_id), INDEX idx_status (status), INDEX idx_created (created_at))" \
  --ai \
  -o /tmp/orders_migration.sql

echo ""
echo "✓ DDL 已保存到 /tmp/orders_migration.sql"
echo ""
echo "文件内容:"
echo "-----------------------------------"
cat /tmp/orders_migration.sql
echo "-----------------------------------"
echo ""

echo "======================================"
echo "演示完成！"
echo "======================================"
echo ""
echo "🎉 主要功能展示:"
echo "  ✓ 智能差异分析"
echo "  ✓ 专业优化建议"
echo "  ✓ 风险提示"
echo "  ✓ 最佳实践推荐"
echo "  ✓ DDL 语句生成和保存"
echo ""
echo "💡 AI 分析亮点:"
echo "  • 复合索引优化建议"
echo "  • 数据类型精度优化"
echo "  • ENUM 类型推荐"
echo "  • 分区表设计建议"
echo "  • 数据迁移风险提示"
echo ""
echo "📊 性能指标:"
echo "  • 平均响应时间: 6-7秒"
echo "  • 平均成本: <¥0.002/次"
echo "  • 分析质量: 5⭐"
echo ""
echo "查看详细测试报告: docs/REAL_AI_TEST_REPORT.md"
echo ""
