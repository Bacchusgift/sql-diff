package differ

import (
	"strings"
	"testing"

	"github.com/youzi/sql-diff/internal/parser"
)

func TestDiffAddColumns(t *testing.T) {
	sourceSQL := `CREATE TABLE users (
		id INT PRIMARY KEY,
		name VARCHAR(100)
	)`

	targetSQL := `CREATE TABLE users (
		id INT PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(255),
		created_at TIMESTAMP
	)`

	p := parser.NewParser()
	source, _ := p.Parse(sourceSQL)
	target, _ := p.Parse(targetSQL)

	differ := NewDiffer(source, target)
	diff := differ.Compare()

	if len(diff.AddedColumns) != 2 {
		t.Errorf("期望新增 2 列，实际 %d 列", len(diff.AddedColumns))
	}

	if !diff.HasChanges() {
		t.Error("应该检测到变更")
	}
}

func TestDiffModifyColumns(t *testing.T) {
	sourceSQL := `CREATE TABLE users (
		id INT,
		name VARCHAR(100)
	)`

	targetSQL := `CREATE TABLE users (
		id INT PRIMARY KEY,
		name VARCHAR(200) NOT NULL
	)`

	p := parser.NewParser()
	source, _ := p.Parse(sourceSQL)
	target, _ := p.Parse(targetSQL)

	differ := NewDiffer(source, target)
	diff := differ.Compare()

	if len(diff.ModifiedColumns) != 1 {
		t.Errorf("期望修改 1 列，实际 %d 列", len(diff.ModifiedColumns))
	}

	// 检查修改详情
	if len(diff.ModifiedColumns) > 0 {
		colDiff := diff.ModifiedColumns[0]
		if colDiff.Name != "name" {
			t.Errorf("期望修改列为 name，实际为 %s", colDiff.Name)
		}
		if len(colDiff.Changes) == 0 {
			t.Error("应该记录变更详情")
		}
	}
}

func TestGenerateDDL(t *testing.T) {
	sourceSQL := `CREATE TABLE users (
		id INT PRIMARY KEY
	)`

	targetSQL := `CREATE TABLE users (
		id INT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255)
	)`

	p := parser.NewParser()
	source, _ := p.Parse(sourceSQL)
	target, _ := p.Parse(targetSQL)

	differ := NewDiffer(source, target)
	diff := differ.Compare()
	ddls := diff.GenerateDDL("users")

	if len(ddls) != 2 {
		t.Errorf("期望生成 2 条 DDL，实际 %d 条", len(ddls))
	}

	// 检查 DDL 内容
	ddlStr := strings.Join(ddls, "\n")
	if !strings.Contains(ddlStr, "ALTER TABLE") {
		t.Error("DDL 应包含 ALTER TABLE")
	}
	if !strings.Contains(ddlStr, "ADD COLUMN") {
		t.Error("DDL 应包含 ADD COLUMN")
	}
}

func TestNoChanges(t *testing.T) {
	sql := `CREATE TABLE users (
		id INT PRIMARY KEY,
		name VARCHAR(100)
	)`

	p := parser.NewParser()
	source, _ := p.Parse(sql)
	target, _ := p.Parse(sql)

	differ := NewDiffer(source, target)
	diff := differ.Compare()

	if diff.HasChanges() {
		t.Error("相同的表不应有变更")
	}

	ddls := diff.GenerateDDL("users")
	if len(ddls) != 0 {
		t.Error("相同的表不应生成 DDL")
	}
}

func TestDiffWithIndexes(t *testing.T) {
	sourceSQL := `CREATE TABLE products (
		id INT PRIMARY KEY,
		name VARCHAR(200)
	)`

	targetSQL := `CREATE TABLE products (
		id INT PRIMARY KEY,
		name VARCHAR(200),
		INDEX idx_name (name)
	)`

	p := parser.NewParser()
	source, _ := p.Parse(sourceSQL)
	target, _ := p.Parse(targetSQL)

	differ := NewDiffer(source, target)
	diff := differ.Compare()

	if len(diff.AddedIndexes) != 1 {
		t.Errorf("期望新增 1 个索引，实际 %d 个", len(diff.AddedIndexes))
	}

	ddls := diff.GenerateDDL("products")
	ddlStr := strings.Join(ddls, "\n")
	if !strings.Contains(ddlStr, "ADD INDEX") {
		t.Error("DDL 应包含 ADD INDEX")
	}
}
