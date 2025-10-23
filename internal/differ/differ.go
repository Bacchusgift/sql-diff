package differ

import (
	"fmt"
	"strings"

	"github.com/Bacchusgift/sql-diff/internal/parser"
)

// Differ 表结构差异比对器
type Differ struct {
	source *parser.TableSchema // 源表结构
	target *parser.TableSchema // 目标表结构
}

// Diff 表结构差异
type Diff struct {
	AddedColumns    []*parser.Column // 新增的列
	RemovedColumns  []*parser.Column // 删除的列
	ModifiedColumns []*ColumnDiff    // 修改的列
	AddedIndexes    []*parser.Index  // 新增的索引
	RemovedIndexes  []*parser.Index  // 删除的索引
}

// ColumnDiff 列的差异详情
type ColumnDiff struct {
	Name    string
	Source  *parser.Column
	Target  *parser.Column
	Changes []string // 变更描述
}

// NewDiffer 创建新的差异比对器
func NewDiffer(source, target *parser.TableSchema) *Differ {
	return &Differ{
		source: source,
		target: target,
	}
}

// Compare 比对两个表结构并返回差异
func (d *Differ) Compare() *Diff {
	diff := &Diff{
		AddedColumns:    make([]*parser.Column, 0),
		RemovedColumns:  make([]*parser.Column, 0),
		ModifiedColumns: make([]*ColumnDiff, 0),
		AddedIndexes:    make([]*parser.Index, 0),
		RemovedIndexes:  make([]*parser.Index, 0),
	}

	// 创建列映射便于查找
	sourceColumns := make(map[string]*parser.Column)
	for _, col := range d.source.Columns {
		sourceColumns[col.Name] = col
	}

	targetColumns := make(map[string]*parser.Column)
	for _, col := range d.target.Columns {
		targetColumns[col.Name] = col
	}

	// 查找新增和修改的列
	for _, targetCol := range d.target.Columns {
		if sourceCol, exists := sourceColumns[targetCol.Name]; exists {
			// 列存在，检查是否有修改
			if changes := compareColumns(sourceCol, targetCol); len(changes) > 0 {
				diff.ModifiedColumns = append(diff.ModifiedColumns, &ColumnDiff{
					Name:    targetCol.Name,
					Source:  sourceCol,
					Target:  targetCol,
					Changes: changes,
				})
			}
		} else {
			// 列不存在，是新增的
			diff.AddedColumns = append(diff.AddedColumns, targetCol)
		}
	}

	// 查找删除的列
	for _, sourceCol := range d.source.Columns {
		if _, exists := targetColumns[sourceCol.Name]; !exists {
			diff.RemovedColumns = append(diff.RemovedColumns, sourceCol)
		}
	}

	// 比对索引
	sourceIndexes := make(map[string]*parser.Index)
	for _, idx := range d.source.Indexes {
		sourceIndexes[idx.Name] = idx
	}

	targetIndexes := make(map[string]*parser.Index)
	for _, idx := range d.target.Indexes {
		targetIndexes[idx.Name] = idx
	}

	// 查找新增的索引
	for _, targetIdx := range d.target.Indexes {
		if _, exists := sourceIndexes[targetIdx.Name]; !exists {
			diff.AddedIndexes = append(diff.AddedIndexes, targetIdx)
		}
	}

	// 查找删除的索引
	for _, sourceIdx := range d.source.Indexes {
		if _, exists := targetIndexes[sourceIdx.Name]; !exists {
			diff.RemovedIndexes = append(diff.RemovedIndexes, sourceIdx)
		}
	}

	return diff
}

// compareColumns 比较两个列的差异
func compareColumns(source, target *parser.Column) []string {
	changes := make([]string, 0)

	if source.Type != target.Type {
		changes = append(changes, fmt.Sprintf("类型从 %s 改为 %s", source.Type, target.Type))
	}

	if source.Length != target.Length {
		changes = append(changes, fmt.Sprintf("长度从 %s 改为 %s", source.Length, target.Length))
	}

	if source.NotNull != target.NotNull {
		if target.NotNull {
			changes = append(changes, "添加了 NOT NULL 约束")
		} else {
			changes = append(changes, "移除了 NOT NULL 约束")
		}
	}

	if source.DefaultValue != target.DefaultValue {
		changes = append(changes, fmt.Sprintf("默认值从 %s 改为 %s", source.DefaultValue, target.DefaultValue))
	}

	if source.AutoInc != target.AutoInc {
		if target.AutoInc {
			changes = append(changes, "添加了 AUTO_INCREMENT")
		} else {
			changes = append(changes, "移除了 AUTO_INCREMENT")
		}
	}

	if source.Comment != target.Comment {
		changes = append(changes, fmt.Sprintf("注释从 '%s' 改为 '%s'", source.Comment, target.Comment))
	}

	return changes
}

// GenerateDDL 根据差异生成 DDL 语句
func (d *Diff) GenerateDDL(tableName string) []string {
	ddls := make([]string, 0)

	// 生成新增列的 DDL
	for _, col := range d.AddedColumns {
		ddl := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s",
			tableName,
			col.Name,
			formatColumnDefinition(col))
		ddls = append(ddls, ddl)
	}

	// 生成修改列的 DDL
	for _, colDiff := range d.ModifiedColumns {
		ddl := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s",
			tableName,
			colDiff.Target.Name,
			formatColumnDefinition(colDiff.Target))
		ddls = append(ddls, ddl)
	}

	// 生成删除列的 DDL（注释掉，因为删除操作比较危险）
	for _, col := range d.RemovedColumns {
		ddl := fmt.Sprintf("-- ALTER TABLE %s DROP COLUMN %s", tableName, col.Name)
		ddls = append(ddls, ddl)
	}

	// 生成新增索引的 DDL
	for _, idx := range d.AddedIndexes {
		var ddl string
		if idx.Type == "UNIQUE" {
			ddl = fmt.Sprintf("ALTER TABLE %s ADD UNIQUE INDEX %s (%s)",
				tableName, idx.Name, strings.Join(idx.Columns, ", "))
		} else {
			ddl = fmt.Sprintf("ALTER TABLE %s ADD INDEX %s (%s)",
				tableName, idx.Name, strings.Join(idx.Columns, ", "))
		}
		ddls = append(ddls, ddl)
	}

	// 生成删除索引的 DDL（注释掉）
	for _, idx := range d.RemovedIndexes {
		ddl := fmt.Sprintf("-- ALTER TABLE %s DROP INDEX %s", tableName, idx.Name)
		ddls = append(ddls, ddl)
	}

	return ddls
}

// formatColumnDefinition 格式化列定义
func formatColumnDefinition(col *parser.Column) string {
	var parts []string

	// 数据类型
	typeStr := col.Type
	if col.Length != "" {
		typeStr += fmt.Sprintf("(%s)", col.Length)
	}
	parts = append(parts, typeStr)

	// UNSIGNED
	if col.Unsigned {
		parts = append(parts, "UNSIGNED")
	}

	// NOT NULL
	if col.NotNull {
		parts = append(parts, "NOT NULL")
	}

	// DEFAULT
	if col.DefaultValue != "" {
		if needsQuotes(col.DefaultValue) {
			parts = append(parts, fmt.Sprintf("DEFAULT '%s'", col.DefaultValue))
		} else {
			parts = append(parts, fmt.Sprintf("DEFAULT %s", col.DefaultValue))
		}
	}

	// AUTO_INCREMENT
	if col.AutoInc {
		parts = append(parts, "AUTO_INCREMENT")
	}

	// COMMENT
	if col.Comment != "" {
		parts = append(parts, fmt.Sprintf("COMMENT '%s'", col.Comment))
	}

	return strings.Join(parts, " ")
}

// needsQuotes 判断默认值是否需要引号
func needsQuotes(value string) bool {
	// 如果已经有引号，不需要再加
	if strings.HasPrefix(value, "'") || strings.HasPrefix(value, "\"") {
		return false
	}
	// NULL、CURRENT_TIMESTAMP 等关键字不需要引号
	upper := strings.ToUpper(value)
	keywords := []string{"NULL", "CURRENT_TIMESTAMP", "NOW()", "TRUE", "FALSE"}
	for _, kw := range keywords {
		if upper == kw || strings.HasPrefix(upper, kw) {
			return false
		}
	}
	// 数字不需要引号
	if _, err := fmt.Sscanf(value, "%f", new(float64)); err == nil {
		return false
	}
	return true
}

// HasChanges 判断是否有变更
func (d *Diff) HasChanges() bool {
	return len(d.AddedColumns) > 0 ||
		len(d.RemovedColumns) > 0 ||
		len(d.ModifiedColumns) > 0 ||
		len(d.AddedIndexes) > 0 ||
		len(d.RemovedIndexes) > 0
}

// Summary 返回差异摘要
func (d *Diff) Summary() string {
	var summary strings.Builder

	if len(d.AddedColumns) > 0 {
		summary.WriteString(fmt.Sprintf("新增列: %d 个\n", len(d.AddedColumns)))
		for _, col := range d.AddedColumns {
			summary.WriteString(fmt.Sprintf("  + %s %s\n", col.Name, col.Type))
		}
	}

	if len(d.ModifiedColumns) > 0 {
		summary.WriteString(fmt.Sprintf("修改列: %d 个\n", len(d.ModifiedColumns)))
		for _, colDiff := range d.ModifiedColumns {
			summary.WriteString(fmt.Sprintf("  * %s: %s\n", colDiff.Name, strings.Join(colDiff.Changes, ", ")))
		}
	}

	if len(d.RemovedColumns) > 0 {
		summary.WriteString(fmt.Sprintf("删除列: %d 个\n", len(d.RemovedColumns)))
		for _, col := range d.RemovedColumns {
			summary.WriteString(fmt.Sprintf("  - %s\n", col.Name))
		}
	}

	if len(d.AddedIndexes) > 0 {
		summary.WriteString(fmt.Sprintf("新增索引: %d 个\n", len(d.AddedIndexes)))
		for _, idx := range d.AddedIndexes {
			summary.WriteString(fmt.Sprintf("  + %s (%s)\n", idx.Name, strings.Join(idx.Columns, ", ")))
		}
	}

	if len(d.RemovedIndexes) > 0 {
		summary.WriteString(fmt.Sprintf("删除索引: %d 个\n", len(d.RemovedIndexes)))
		for _, idx := range d.RemovedIndexes {
			summary.WriteString(fmt.Sprintf("  - %s\n", idx.Name))
		}
	}

	if summary.Len() == 0 {
		return "没有发现差异"
	}

	return summary.String()
}
