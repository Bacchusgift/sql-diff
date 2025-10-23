# 贡献指南

感谢您对 SQL-Diff 项目的关注!我们欢迎所有形式的贡献。

## 如何贡献

### 报告 Bug

如果您发现了 Bug,请创建一个 Issue,并包含:

- 清晰的标题和描述
- 重现步骤
- 预期行为和实际行为
- 您的环境信息 (OS, Go 版本等)
- 相关的错误日志

### 提出新功能

如果您有新功能的想法:

1. 先检查是否已有相关 Issue
2. 创建 Feature Request Issue
3. 详细描述功能需求和使用场景
4. 等待维护者反馈

### 提交代码

#### 开发流程

1. **Fork 仓库**
   ```bash
   # 在 GitHub 上 Fork 项目
   git clone https://github.com/YOUR_USERNAME/sql-diff.git
   cd sql-diff
   ```

2. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

3. **编写代码**
   - 遵循项目的代码风格
   - 添加必要的测试
   - 更新相关文档

4. **运行测试**
   ```bash
   # 运行所有测试
   make test
   
   # 检查代码格式
   make fmt
   
   # 运行 lint
   make lint
   ```

5. **提交变更**
   ```bash
   git add .
   git commit -m "feat: add awesome feature"
   
   # 提交信息规范见下文
   ```

6. **推送到 GitHub**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 填写 PR 模板
   - 等待 CI 通过和代码审查

## 代码规范

### Go 代码风格

遵循 [Effective Go](https://go.dev/doc/effective_go) 和项目现有风格:

```go
// ✅ 好的示例
func ParseCreateTable(sql string) (*TableSchema, error) {
    if sql == "" {
        return nil, errors.New("empty SQL string")
    }
    
    schema := &TableSchema{
        Columns: make([]Column, 0),
    }
    
    // ... 解析逻辑
    
    return schema, nil
}

// ❌ 不好的示例
func parse(s string) *TableSchema {
    // 缺少错误处理
    // 变量名不清晰
    result := &TableSchema{}
    // ...
    return result
}
```

### 命名规范

- **包名**: 小写,单数,简洁 (`parser`, `differ`, `config`)
- **函数名**: 驼峰命名,动词开头 (`ParseCreateTable`, `GenerateDDL`)
- **变量名**: 驼峰命名,清晰表意 (`tableName`, `sourceSchema`)
- **常量**: 大写驼峰 (`DefaultTimeout`, `MaxRetries`)

### 注释规范

```go
// Package parser provides SQL parsing functionality.
package parser

// TableSchema represents the structure of a database table.
// It includes columns, indexes, and table-level options.
type TableSchema struct {
    Name    string            // Table name
    Columns []Column          // List of columns
    Indexes []Index           // List of indexes
}

// ParseCreateTable parses a CREATE TABLE SQL statement and returns
// the corresponding TableSchema. It returns an error if the SQL
// syntax is invalid.
func ParseCreateTable(sql string) (*TableSchema, error) {
    // ...
}
```

## 测试规范

### 编写测试

每个新功能都应该有相应的测试:

```go
// parser_test.go
func TestParseCreateTable(t *testing.T) {
    tests := []struct {
        name    string
        sql     string
        want    *TableSchema
        wantErr bool
    }{
        {
            name: "simple table",
            sql:  "CREATE TABLE users (id INT PRIMARY KEY)",
            want: &TableSchema{
                Name: "users",
                Columns: []Column{
                    {Name: "id", Type: "INT"},
                },
                PrimaryKeys: []string{"id"},
            },
            wantErr: false,
        },
        {
            name:    "invalid SQL",
            sql:     "INVALID SQL",
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseCreateTable(tt.sql)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseCreateTable() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseCreateTable() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/parser

# 显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 提交信息规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档变更
- `style`: 代码格式 (不影响功能)
- `refactor`: 重构 (不是新功能也不是修复)
- `perf`: 性能优化
- `test`: 添加或修改测试
- `chore`: 构建过程或辅助工具的变更

### 示例

```bash
# 新功能
git commit -m "feat(parser): add support for PostgreSQL syntax"

# Bug 修复
git commit -m "fix(differ): correct index comparison logic"

# 文档
git commit -m "docs(readme): update installation instructions"

# 重构
git commit -m "refactor(ai): extract common provider logic"
```

## Pull Request 规范

### PR 标题

使用与提交信息相同的格式:

```
feat(parser): add PostgreSQL support
fix(cli): handle empty input gracefully
docs(guide): improve getting started section
```

### PR 描述模板

```markdown
## 变更类型
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## 变更说明
<!-- 描述您的变更 -->

## 相关 Issue
Closes #123

## 测试
<!-- 描述您如何测试这些变更 -->
- [ ] 添加了新的测试
- [ ] 所有测试通过
- [ ] 手动测试通过

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 已更新相关文档
- [ ] 已添加测试用例
- [ ] 所有测试通过
- [ ] Commit 信息符合规范
```

## 文档贡献

文档同样重要!您可以:

- 修复文档中的错误
- 改进文档的清晰度
- 添加示例
- 翻译文档

文档位于:
- 项目 `README.md`
- `docs/` 目录 (VitePress 文档)
- 代码注释

## 开发环境设置

### 前置要求

- Go 1.21 或更高版本
- Git
- Make (可选)

### 克隆项目

```bash
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff
```

### 安装依赖

```bash
go mod download
```

### 构建项目

```bash
make build
# 或
go build -o bin/sql-diff cmd/sql-diff/main.go
```

### 运行测试

```bash
make test
# 或
go test ./...
```

### 本地安装

```bash
make install
# 或
go install ./cmd/sql-diff
```

## 项目结构

详细的项目结构说明请参考 [架构设计](/architecture)。

```
sql-diff/
├── cmd/
│   └── sql-diff/          # 主程序入口
├── internal/
│   ├── parser/            # SQL 解析器
│   ├── differ/            # 差异比对器
│   ├── ai/               # AI 集成
│   ├── config/           # 配置管理
│   └── cmd/              # CLI 命令
├── docs/                  # 文档源文件
├── examples/              # 示例文件
└── tests/                 # 测试文件
```

## 发布流程

维护者发布新版本的流程:

1. 更新 `VERSION` 文件
2. 更新 `CHANGELOG.md`
3. 创建 Git tag
4. 推送到 GitHub
5. GitHub Actions 自动构建和发布

## 社区

- GitHub Issues: 报告 Bug 和功能请求
- GitHub Discussions: 讨论和提问
- Pull Requests: 贡献代码

## 行为准则

请遵循友好、包容的社区环境准则。

## 许可证

通过贡献代码,您同意您的贡献将在 MIT License 下发布。

## 感谢

感谢所有贡献者对项目的支持! 🙏

查看 [贡献者列表](https://github.com/Bacchusgift/sql-diff/graphs/contributors)。

## 联系方式

如有任何问题,欢迎:
- 创建 Issue
- 发起 Discussion
- 发送邮件至 maintainer@example.com

---

再次感谢您的贡献! ❤️
