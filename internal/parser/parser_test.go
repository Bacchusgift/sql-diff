package parser

import (
	"testing"
)

func TestParseSimpleTable(t *testing.T) {
	sql := `CREATE TABLE users (
		id INT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	parser := NewParser()
	schema, err := parser.Parse(sql)

	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if schema.Name != "users" {
		t.Errorf("表名错误，期望 users，得到 %s", schema.Name)
	}

	if len(schema.Columns) != 4 {
		t.Errorf("列数错误，期望 4，得到 %d", len(schema.Columns))
	}

	// 检查第一列
	if schema.Columns[0].Name != "id" {
		t.Errorf("第一列名错误，期望 id，得到 %s", schema.Columns[0].Name)
	}

	if schema.Columns[0].Type != "INT" {
		t.Errorf("第一列类型错误，期望 INT，得到 %s", schema.Columns[0].Type)
	}

	// 检查主键
	if len(schema.PrimaryKeys) != 1 || schema.PrimaryKeys[0] != "id" {
		t.Errorf("主键错误，期望 [id]，得到 %v", schema.PrimaryKeys)
	}
}

func TestParseTableWithIndex(t *testing.T) {
	sql := `CREATE TABLE products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(200) NOT NULL,
		price DECIMAL(10,2),
		INDEX idx_name (name),
		UNIQUE INDEX idx_price (price)
	)`

	parser := NewParser()
	schema, err := parser.Parse(sql)

	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if len(schema.Indexes) != 2 {
		t.Errorf("索引数量错误，期望 2，得到 %d", len(schema.Indexes))
	}

	// 检查普通索引
	if schema.Indexes[0].Name != "idx_name" {
		t.Errorf("索引名错误，期望 idx_name，得到 %s", schema.Indexes[0].Name)
	}

	// 检查唯一索引
	if schema.Indexes[1].Type != "UNIQUE" {
		t.Errorf("索引类型错误，期望 UNIQUE，得到 %s", schema.Indexes[1].Type)
	}
}

func TestParseComplexTable(t *testing.T) {
	sql := `CREATE TABLE orders (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		total_amount DECIMAL(12,2) DEFAULT 0.00,
		status VARCHAR(20) DEFAULT 'pending' NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_user_id (user_id),
		INDEX idx_status (status),
		INDEX idx_created (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

	parser := NewParser()
	schema, err := parser.Parse(sql)

	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if schema.Name != "orders" {
		t.Errorf("表名错误")
	}

	// 检查列
	if len(schema.Columns) != 6 {
		t.Errorf("列数错误，期望 6，得到 %d", len(schema.Columns))
	}

	// 检查索引
	if len(schema.Indexes) != 3 {
		t.Errorf("索引数量错误，期望 3，得到 %d", len(schema.Indexes))
	}

	// 检查表选项
	if schema.Options["ENGINE"] != "InnoDB" {
		t.Errorf("ENGINE 错误，期望 InnoDB，得到 %s", schema.Options["ENGINE"])
	}

	if schema.Options["CHARSET"] != "utf8mb4" {
		t.Errorf("CHARSET 错误，期望 utf8mb4，得到 %s", schema.Options["CHARSET"])
	}
}
