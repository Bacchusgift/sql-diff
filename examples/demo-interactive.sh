#!/bin/bash

# SQL-Diff 交互式模式演示脚本
# 本脚本展示如何使用交互式模式处理多行 SQL

set -e

echo "======================================"
echo "  SQL-Diff 交互式模式演示"
echo "======================================"
echo ""
echo "方式 1: 使用交互式模式（推荐）"
echo "命令: ./bin/sql-diff -i"
echo ""
echo "方式 2: 使用 Here Document 自动输入"
echo "======================================"
echo ""

# 定义源表 SQL
SOURCE_SQL='CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
)'

# 定义目标表 SQL  
TARGET_SQL='CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(255),
  created_at TIMESTAMP
)'

# 使用 here document 自动输入
./bin/sql-diff -i <<EOF
$SOURCE_SQL
EOF
$TARGET_SQL
EOF

echo ""
echo "======================================"
echo "  演示完成！"
echo "======================================"
