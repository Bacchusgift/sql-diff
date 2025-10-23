# SQL-Diff 架构设计文档

## 1. 项目概述

SQL-Diff 是一个基于 AST 语法树的 SQL 表结构比对工具，可以智能比对两个表结构的差异，并自动生成 DDL 补全语句。支持可选的 AI 增强功能，提供智能分析和优化建议。

## 2. 项目结构

```
sql-diff/
├── cmd/
│   └── sql-diff/
│       └── main.go              # 程序入口
├── internal/
│   ├── parser/                  # SQL 解析器
│   │   ├── parser.go            # 解析器实现
│   │   └── parser_test.go       # 解析器测试
│   ├── differ/                  # 差异比对器
│   │   ├── differ.go            # 比对器实现
│   │   └── differ_test.go       # 比对器测试
│   ├── ai/                      # AI 模型集成
│   │   └── provider.go          # AI 提供商抽象
│   ├── config/                  # 配置管理
│   │   └── config.go            # 配置加载和验证
│   └── cmd/                     # CLI 命令
│       └── root.go              # 根命令实现
├── docs/                        # 文档
│   └── EXAMPLES.md              # 使用示例
├── examples/                    # 示例代码
│   └── demo.sh                  # 演示脚本
├── bin/                         # 编译输出目录
├── go.mod                       # Go 模块定义
├── go.sum                       # 依赖校验
├── Makefile                     # 构建脚本
├── README.md                    # 项目说明
├── .gitignore                   # Git 忽略配置
└── .sql-diff-config.example.yaml # 配置文件示例
```

## 3. 核心模块设计

### 3.1 Parser 模块

**职责**：解析 SQL CREATE TABLE 语句，构建表结构的 AST 表示

**核心数据结构**：

```go
// TableSchema 表结构定义
type TableSchema struct {
    Name        string            // 表名
    Columns     []*Column         // 列定义列表
    PrimaryKeys []string          // 主键列名
    Indexes     []*Index          // 索引定义
    Constraints []*Constraint     // 约束定义
    Options     map[string]string // 表选项
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
```

**解析流程**：
1. 提取表名
2. 提取列定义部分
3. 逐行解析列定义、主键、索引
4. 提取表选项（ENGINE、CHARSET 等）

### 3.2 Differ 模块

**职责**：比对两个表结构，识别差异并生成 DDL 语句

**核心数据结构**：

```go
// Diff 表结构差异
type Diff struct {
    AddedColumns    []*Column     // 新增的列
    RemovedColumns  []*Column     // 删除的列
    ModifiedColumns []*ColumnDiff // 修改的列
    AddedIndexes    []*Index      // 新增的索引
    RemovedIndexes  []*Index      // 删除的索引
}

// ColumnDiff 列的差异详情
type ColumnDiff struct {
    Name    string
    Source  *Column
    Target  *Column
    Changes []string // 变更描述
}
```

**比对流程**：
1. 将两个表的列转换为映射（Map）便于查找
2. 遍历目标表列，识别新增和修改的列
3. 遍历源表列，识别删除的列
4. 比对索引差异
5. 生成 DDL 语句

**DDL 生成规则**：
- 新增列：`ALTER TABLE xxx ADD COLUMN ...`
- 修改列：`ALTER TABLE xxx MODIFY COLUMN ...`
- 删除列：注释掉（避免误删）
- 新增索引：`ALTER TABLE xxx ADD INDEX ...`
- 删除索引：注释掉（避免误删）

### 3.3 AI 模块

**职责**：提供 AI 智能分析功能，包括差异分析和 SQL 优化建议

**接口设计**：

```go
// Provider AI 提供商接口
type Provider interface {
    // Analyze 分析表结构差异并提供建议
    Analyze(sourceDDL, targetDDL, diff string) (*AnalysisResult, error)
    
    // OptimizeSQL 优化 SQL 语句
    OptimizeSQL(sql string) (*OptimizationResult, error)
}
```

**支持的提供商**：
- DeepSeek
- OpenAI（兼容 OpenAI API 的其他服务）
- 自定义提供商（通过配置）

**工作流程**：
1. 构建分析提示词（包含源表、目标表、差异摘要）
2. 调用 AI API 获取分析结果
3. 解析响应并展示给用户

### 3.4 Config 模块

**职责**：管理应用配置，包括 AI API 密钥等敏感信息

**配置结构**：

```go
type Config struct {
    AI AIConfig `yaml:"ai"`
}

type AIConfig struct {
    Enabled     bool   `yaml:"enabled"`
    Provider    string `yaml:"provider"`
    APIKey      string `yaml:"api_key"`
    APIEndpoint string `yaml:"api_endpoint"`
    Model       string `yaml:"model"`
    Timeout     int    `yaml:"timeout"`
}
```

**配置加载顺序**：
1. 默认配置
2. 配置文件（`.sql-diff-config.yaml`）
3. 命令行参数覆盖

### 3.5 CLI 模块

**职责**：提供命令行接口，处理用户输入和输出

**使用的库**：
- `github.com/spf13/cobra`：命令行框架
- `github.com/fatih/color`：彩色输出

**命令参数**：
- `-s, --source`：源表 SQL
- `-t, --target`：目标表 SQL
- `--ai`：启用 AI 分析
- `--config`：配置文件路径
- `-o, --output`：输出文件

## 4. 工作流程

```
用户输入
   ↓
解析源表 (Parser)
   ↓
解析目标表 (Parser)
   ↓
比对差异 (Differ)
   ↓
生成 DDL (Differ)
   ↓
[可选] AI 分析 (AI Provider)
   ↓
输出结果
```

## 5. 技术选型

- **语言**：Go 1.21+
- **CLI 框架**：Cobra
- **配置解析**：gopkg.in/yaml.v3
- **彩色输出**：fatih/color
- **HTTP 客户端**：标准库 net/http

## 6. 安全考虑

1. **敏感信息保护**：
   - API Key 存储在配置文件中
   - 配置文件加入 .gitignore
   - 不在日志中输出敏感信息

2. **SQL 注入防护**：
   - 只解析 CREATE TABLE 语句
   - 不执行任何 SQL
   - 生成的 DDL 需要人工审核

3. **删除操作保护**：
   - 删除列和索引的 DDL 默认注释
   - 避免误删重要数据

## 7. 扩展性设计

### 7.1 新增 AI 提供商

实现 `Provider` 接口：

```go
type CustomProvider struct {
    config *config.AIConfig
}

func (p *CustomProvider) Analyze(...) (*AnalysisResult, error) {
    // 自定义实现
}
```

在 `NewProvider` 函数中注册：

```go
case "custom":
    return NewCustomProvider(cfg), nil
```

### 7.2 支持更多数据库

当前支持 MySQL 语法，可扩展：
- PostgreSQL
- SQLite
- SQL Server

方法：实现不同的 Parser 和 Differ

## 8. 性能优化

1. **解析性能**：
   - 使用正则表达式预编译
   - 避免重复解析

2. **内存优化**：
   - 使用指针减少拷贝
   - 及时释放大对象

3. **并发处理**：
   - 可并行解析源表和目标表
   - AI 调用可异步处理

## 9. 测试策略

1. **单元测试**：
   - Parser 模块：测试各种 SQL 语法
   - Differ 模块：测试差异识别和 DDL 生成
   - 覆盖率目标：> 80%

2. **集成测试**：
   - 端到端测试完整流程
   - 测试 AI 集成（使用 Mock）

3. **示例测试**：
   - examples/demo.sh 作为冒烟测试

## 10. 未来计划

1. **功能增强**：
   - 支持更多数据库类型
   - 支持表间关系比对
   - 支持数据迁移建议

2. **性能优化**：
   - 大表结构优化
   - 批量比对支持

3. **用户体验**：
   - Web UI 界面
   - 交互式模式
   - 结果可视化

4. **AI 增强**：
   - 更智能的差异分析
   - 自动生成迁移脚本
   - 性能影响评估
