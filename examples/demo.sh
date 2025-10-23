#!/bin/bash

# SQL-Diff 演示脚本

echo "======================================"
echo "SQL-Diff 工具演示"
echo "======================================"
echo ""

# 示例 1: 基础用法 - 新增字段
echo "示例 1: 新增字段"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"

echo ""
echo ""

# 示例 2: 修改字段
echo "示例 2: 修改字段属性"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE products (id INT, name VARCHAR(100), price DECIMAL(10,2))" \
  -t "CREATE TABLE products (id INT PRIMARY KEY, name VARCHAR(200) NOT NULL, price DECIMAL(12,2) DEFAULT 0.00)"

echo ""
echo ""

# 示例 3: 新增索引
echo "示例 3: 新增索引"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE posts (id INT PRIMARY KEY, title VARCHAR(200), author_id INT)" \
  -t "CREATE TABLE posts (id INT PRIMARY KEY, title VARCHAR(200), author_id INT, INDEX idx_author (author_id))"

echo ""
echo ""

# 示例 4: 复杂表结构比对
echo "示例 4: 复杂表结构比对"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE orders (id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, user_id INT NOT NULL, total_amount DECIMAL(12,2) DEFAULT 0.00, status VARCHAR(20) DEFAULT 'pending', created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, INDEX idx_user_id (user_id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4" \
  -t "CREATE TABLE orders (id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, user_id INT NOT NULL, total_amount DECIMAL(12,2) DEFAULT 0.00, status VARCHAR(20) DEFAULT 'pending' NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, remark TEXT, INDEX idx_user_id (user_id), INDEX idx_status (status)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

echo ""
echo ""

# 示例 5: 输出到文件
echo "示例 5: 输出到文件"
echo "--------------------------------------"
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))" \
  -o /tmp/migration.sql

if [ -f /tmp/migration.sql ]; then
  echo "✓ DDL 已保存到 /tmp/migration.sql"
  echo "文件内容:"
  cat /tmp/migration.sql
  rm /tmp/migration.sql
fi

echo ""
echo "======================================"
echo "演示完成！"
echo "======================================"
