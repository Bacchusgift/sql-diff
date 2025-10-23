package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// TableSchema 表结构定义
type TableSchema struct {
	Name        string            // 表名
	Columns     []*Column         // 列定义列表
	PrimaryKeys []string          // 主键列名
	Indexes     []*Index          // 索引定义
	Constraints []*Constraint     // 约束定义
	Options     map[string]string // 表选项（ENGINE, CHARSET 等）
}

// Column 列定义
type Column struct {
	Name         string // 列名
	Type         string // 数据类型
	Length       string // 长度
	NotNull      bool   // 是否非空
	DefaultValue string // 默认值
	AutoInc      bool   // 是否自增
	Comment      string // 注释
	Unsigned     bool   // 是否无符号
}

// Index 索引定义
type Index struct {
	Name    string   // 索引名
	Columns []string // 索引列
	Type    string   // 索引类型：INDEX, UNIQUE, FULLTEXT
}

// Constraint 约束定义
type Constraint struct {
	Name       string // 约束名
	Type       string // 约束类型：PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK
	Definition string // 约束定义
}

// Parser SQL 解析器接口
type Parser interface {
	Parse(sql string) (*TableSchema, error)
}

// SimpleParser 简单的 SQL 解析器实现
// 注意：这是一个简化版本，生产环境建议使用更完善的 SQL 解析库
type SimpleParser struct{}

// NewParser 创建新的解析器
func NewParser() Parser {
	return &SimpleParser{}
}

// Parse 解析 CREATE TABLE 语句
func (p *SimpleParser) Parse(sql string) (*TableSchema, error) {
	sql = strings.TrimSpace(sql)

	// 提取表名
	tableName, err := extractTableName(sql)
	if err != nil {
		return nil, err
	}

	schema := &TableSchema{
		Name:    tableName,
		Columns: make([]*Column, 0),
		Options: make(map[string]string),
	}

	// 提取列定义部分
	columnsDef, err := extractColumnsDefinition(sql)
	if err != nil {
		return nil, err
	}

	// 解析每一列
	lines := splitColumnDefinitions(columnsDef)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析主键
		if strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") {
			keys := extractPrimaryKeys(line)
			schema.PrimaryKeys = append(schema.PrimaryKeys, keys...)
			continue
		}

		// 解析索引
		if isIndex(line) {
			index := parseIndex(line)
			if index != nil {
				schema.Indexes = append(schema.Indexes, index)
			}
			continue
		}

		// 解析普通列
		column := parseColumn(line)
		if column != nil {
			schema.Columns = append(schema.Columns, column)
			// 检查列定义中是否包含 PRIMARY KEY
			if strings.Contains(strings.ToUpper(line), "PRIMARY KEY") {
				schema.PrimaryKeys = append(schema.PrimaryKeys, column.Name)
			}
		}
	}

	// 提取表选项
	schema.Options = extractTableOptions(sql)

	return schema, nil
}

// extractTableName 提取表名
func extractTableName(sql string) (string, error) {
	re := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?` + "`?" + `([a-zA-Z0-9_]+)` + "`?")
	matches := re.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法提取表名")
	}
	return matches[1], nil
}

// extractColumnsDefinition 提取列定义部分
func extractColumnsDefinition(sql string) (string, error) {
	// 找到第一个 ( 和最后一个 )
	start := strings.Index(sql, "(")
	end := strings.LastIndex(sql, ")")

	if start == -1 || end == -1 || start >= end {
		return "", fmt.Errorf("无效的 CREATE TABLE 语法")
	}

	return sql[start+1 : end], nil
}

// splitColumnDefinitions 分割列定义
func splitColumnDefinitions(def string) []string {
	var result []string
	var current strings.Builder
	var parenDepth int

	for _, ch := range def {
		switch ch {
		case '(':
			parenDepth++
			current.WriteRune(ch)
		case ')':
			parenDepth--
			current.WriteRune(ch)
		case ',':
			if parenDepth == 0 {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// parseColumn 解析列定义
func parseColumn(line string) *Column {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)

	if len(parts) < 2 {
		return nil
	}

	column := &Column{
		Name: strings.Trim(parts[0], "`"),
	}

	// 解析数据类型和长度
	typeStr := parts[1]
	if strings.Contains(typeStr, "(") {
		re := regexp.MustCompile(`([A-Z]+)\(([^)]+)\)`)
		matches := re.FindStringSubmatch(strings.ToUpper(typeStr))
		if len(matches) >= 3 {
			column.Type = matches[1]
			column.Length = matches[2]
		}
	} else {
		column.Type = strings.ToUpper(typeStr)
	}

	// 解析其他属性
	upperLine := strings.ToUpper(line)
	column.NotNull = strings.Contains(upperLine, "NOT NULL")
	column.AutoInc = strings.Contains(upperLine, "AUTO_INCREMENT")
	column.Unsigned = strings.Contains(upperLine, "UNSIGNED")

	// 解析默认值
	if strings.Contains(upperLine, "DEFAULT") {
		re := regexp.MustCompile(`(?i)DEFAULT\s+('([^']*)'|"([^"]*)"|([^\s,]+))`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			for i := 2; i < len(matches); i++ {
				if matches[i] != "" {
					column.DefaultValue = matches[i]
					break
				}
			}
		}
	}

	// 解析注释
	if strings.Contains(upperLine, "COMMENT") {
		re := regexp.MustCompile(`(?i)COMMENT\s+'([^']*)'`)
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 2 {
			column.Comment = matches[1]
		}
	}

	return column
}

// extractPrimaryKeys 提取主键列名
func extractPrimaryKeys(line string) []string {
	re := regexp.MustCompile(`PRIMARY\s+KEY\s*\(([^)]+)\)`)
	matches := re.FindStringSubmatch(strings.ToUpper(line))
	if len(matches) < 2 {
		return nil
	}

	keys := strings.Split(matches[1], ",")
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, strings.Trim(strings.TrimSpace(key), "`"))
	}
	return result
}

// isIndex 判断是否是索引定义
func isIndex(line string) bool {
	upper := strings.ToUpper(line)
	return strings.HasPrefix(upper, "INDEX") ||
		strings.HasPrefix(upper, "KEY") ||
		strings.HasPrefix(upper, "UNIQUE") ||
		strings.HasPrefix(upper, "FULLTEXT")
}

// parseIndex 解析索引定义
func parseIndex(line string) *Index {
	upper := strings.ToUpper(line)

	index := &Index{
		Type: "INDEX",
	}

	if strings.HasPrefix(upper, "UNIQUE") {
		index.Type = "UNIQUE"
	} else if strings.HasPrefix(upper, "FULLTEXT") {
		index.Type = "FULLTEXT"
	}

	// 提取索引名和列
	re := regexp.MustCompile(`(?:INDEX|KEY|UNIQUE|FULLTEXT)\s+(?:INDEX|KEY)?\s*` + "`?" + `([a-zA-Z0-9_]+)` + "`?" + `\s*\(([^)]+)\)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 3 {
		index.Name = matches[1]
		columns := strings.Split(matches[2], ",")
		for _, col := range columns {
			index.Columns = append(index.Columns, strings.Trim(strings.TrimSpace(col), "`"))
		}
	}

	return index
}

// extractTableOptions 提取表选项
func extractTableOptions(sql string) map[string]string {
	options := make(map[string]string)

	// 提取 ENGINE
	re := regexp.MustCompile(`(?i)ENGINE\s*=\s*([a-zA-Z0-9]+)`)
	if matches := re.FindStringSubmatch(sql); len(matches) >= 2 {
		options["ENGINE"] = matches[1]
	}

	// 提取 CHARSET
	re = regexp.MustCompile(`(?i)(?:DEFAULT\s+)?CHARSET\s*=\s*([a-zA-Z0-9]+)`)
	if matches := re.FindStringSubmatch(sql); len(matches) >= 2 {
		options["CHARSET"] = matches[1]
	}

	// 提取 COLLATE
	re = regexp.MustCompile(`(?i)COLLATE\s*=\s*([a-zA-Z0-9_]+)`)
	if matches := re.FindStringSubmatch(sql); len(matches) >= 2 {
		options["COLLATE"] = matches[1]
	}

	return options
}
