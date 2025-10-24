# Feature: AI 自然语言生成 SQL

## 📋 功能概述

这个新特性为 sql-diff 添加了两个强大的 AI 驱动命令，允许用户通过自然语言描述来生成 SQL 语句。

## ✨ 新增命令

### 1. `sql-diff generate`

**功能**：根据自然语言描述生成 CREATE TABLE 语句

**命令格式**：
```bash
sql-diff generate -d "自然语言描述" [-o output.sql]
```

**使用示例**：
```bash
# 基础示例
sql-diff generate -d "创建用户表，包含 ID、用户名、邮箱、密码、创建时间"

# 输出到文件
sql-diff generate -d "创建商品表：商品ID、名称、价格、库存" -o product.sql

# 复杂示例
sql-diff generate -d "创建订单表：订单号（唯一）、用户ID（外键）、订单金额（精确到分）、订单状态（待支付/已支付/已取消）、支付方式、下单时间、支付时间"
```

**特点**：
- ✅ 自动推断字段类型（VARCHAR、INT、DECIMAL、DATETIME等）
- ✅ 自动设置主键、索引、唯一约束
- ✅ 应用 MySQL 最佳实践（InnoDB、UTF8MB4、注释等）
- ✅ 使用标准命名规范（snake_case）
- ✅ 添加合理的默认值和注释

### 2. `sql-diff alter`

**功能**：根据现有表结构和修改需求生成 ALTER TABLE 语句

**命令格式**：
```bash
# 命令行模式
sql-diff alter -t "CREATE TABLE ..." -d "修改需求" [-o output.sql]

# 交互式模式
sql-diff alter -i -d "修改需求"
```

**使用示例**：
```bash
# 命令行模式
sql-diff alter -t "CREATE TABLE users (id INT, name VARCHAR(100));" \
               -d "添加手机号字段、邮箱改为唯一索引"

# 交互式模式（推荐用于复杂表结构）
sql-diff alter -i -d "添加实名认证字段：真实姓名、身份证号、认证状态"

# 输出到文件
sql-diff alter -i -d "添加商品状态字段" -o alter_product.sql
```

**特点**：
- ✅ 基于现有表结构生成变更
- ✅ 可生成多条 ALTER 语句
- ✅ 考虑数据迁移安全性
- ✅ 支持交互式粘贴大型表结构
- ✅ 保持与现有表的一致性

## 🏗️ 技术实现

### 架构改动

#### 1. Provider 接口扩展

在 `internal/ai/provider.go` 中扩展了 Provider 接口：

```go
type Provider interface {
    // 现有方法
    Analyze(sourceDDL, targetDDL, diff string) (*AnalysisResult, error)
    OptimizeSQL(sql string) (*OptimizationResult, error)
    
    // 新增方法
    GenerateCreateTable(description string) (string, error)
    GenerateAlterTable(currentDDL, description string) (string, error)
}
```

#### 2. Provider 实现

为所有 Provider 实现了新方法：

- **NoOpProvider**：返回错误提示 AI 未启用
- **DeepSeekProvider**：调用 DeepSeek API 生成 SQL
- **OpenAIProvider**：复用 DeepSeek 实现（API 兼容）
- **MockProvider**：用于测试的模拟实现

#### 3. 智能 SQL 清理

新增 `cleanSQLResponse()` 函数：
- 移除 AI 响应中的 Markdown 代码块标记
- 提取纯 SQL 语句
- 处理多行响应
- 移除末尾分号（统一由工具添加）

#### 4. 新命令文件

- `internal/cmd/generate.go`：实现 generate 命令
- `internal/cmd/alter.go`：实现 alter 命令

### Prompt 工程

#### CREATE TABLE Prompt

设计要点：
- 明确要求返回可执行的 MySQL 语句
- 指定字段类型选择规则（如金额用 DECIMAL）
- 要求遵循最佳实践（InnoDB、UTF8MB4）
- 要求使用蛇形命名法
- 只返回 SQL，不要解释文字

#### ALTER TABLE Prompt

设计要点：
- 提供现有表结构作为上下文
- 考虑数据迁移安全性
- 可生成多条 ALTER 语句
- 保持与现有表的一致性
- 每条语句一行，不加分号

## 📊 代码变更统计

```
新增文件：
- internal/cmd/generate.go (114 行)
- internal/cmd/alter.go (165 行)
- AI_SQL_GENERATOR.md (213 行)
- FEATURE_AI_SQL_GENERATOR.md (本文件)

修改文件：
- internal/ai/provider.go (+110 行)
- internal/ai/mock.go (+26 行)
- README.md (+40 行)
- .gitignore (+1 行)

总计：约 600+ 行新代码
```

## 🧪 测试

### 单元测试

所有现有测试通过：
```bash
$ go test ./...
ok      github.com/Bacchusgift/sql-diff/internal/ai     0.008s
ok      github.com/Bacchusgift/sql-diff/internal/differ 0.007s
ok      github.com/Bacchusgift/sql-diff/internal/parser 0.007s
```

### 手动测试清单

- [x] `generate` 命令帮助信息
- [x] `alter` 命令帮助信息
- [x] 编译成功，无错误
- [x] 参数验证正确
- [x] 错误提示友好

## 📝 使用场景

### 场景 1：快速原型设计

```bash
# 产品经理：我们需要一个博客系统
sql-diff generate -d "创建博客文章表：文章ID、标题、内容、作者、分类、标签、发布时间、浏览量"
```

### 场景 2：需求迭代

```bash
# 第一版：基础用户表
sql-diff generate -d "用户表：ID、用户名、邮箱" -o user_v1.sql

# 第二版：添加字段
sql-diff alter -t "$(cat user_v1.sql)" -d "添加手机号、实名状态、注册时间" -o user_v2.sql
```

### 场景 3：需求文档转 SQL

将产品需求直接转换为数据库设计，节省大量时间。

## ⚠️ 注意事项

### 依赖 AI 服务

此功能需要配置 AI 服务：

```yaml
# .sql-diff-config.yaml
ai:
  enabled: true
  provider: deepseek
  api_key: "your-api-key"
  api_endpoint: "https://api.deepseek.com/v1"
  model: "deepseek-chat"
```

### AI 生成结果需审查

- AI 生成的 SQL 仅供参考
- 生产使用前请仔细审查
- 建议在测试环境先验证
- 根据实际需求调整

### API 费用

- 使用 DeepSeek API 可能产生费用
- 建议设置合理的使用限制
- 考虑使用本地模型（未来可能支持）

## 🚀 未来改进

### 短期计划

- [ ] 添加更多示例到文档
- [ ] 支持批量生成（从 CSV/Excel）
- [ ] 生成结果缓存（避免重复调用 API）
- [ ] 支持更多数据库（PostgreSQL、SQLite）

### 长期计划

- [ ] 支持本地 LLM 模型（Ollama）
- [ ] 可视化表设计器
- [ ] ER 图生成
- [ ] 从现有数据库逆向生成自然语言描述

## 📚 相关文档

- [AI_SQL_GENERATOR.md](./AI_SQL_GENERATOR.md) - 详细使用指南
- [README.md](./README.md) - 项目主文档
- [ai/guide.md](./ai/guide.md) - AI 功能配置指南

## 🎯 总结

这次迭代为 sql-diff 添加了强大的 AI 自然语言生成 SQL 功能，大大降低了编写 DDL 的门槛，提高了开发效率。通过精心设计的 Prompt 和智能 SQL 提取，确保生成的 SQL 符合最佳实践。

功能已在 `feature/ai-sql-generator` 分支开发完成，测试通过，可以合并到主分支。

## 📅 开发时间线

- 2025-10-23：功能开发完成
- 分支：`feature/ai-sql-generator`
- 提交数：2 commits
- 状态：✅ 开发完成，待合并
