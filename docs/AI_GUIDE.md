# AI 功能使用指南

SQL-Diff 支持集成 AI 大模型（如 DeepSeek）进行智能分析和优化建议。本指南将帮助您配置和使用 AI 功能。

## 目录

- [功能概述](#功能概述)
- [支持的 AI 提供商](#支持的-ai-提供商)
- [配置 AI](#配置-ai)
- [使用 AI 分析](#使用-ai-分析)
- [AI 输出解读](#ai-输出解读)
- [常见问题](#常见问题)
- [最佳实践](#最佳实践)

## 功能概述

AI 功能可以为您提供:

- 📊 **智能差异分析**: 深入解读表结构变更的影响和意义
- ✨ **优化建议**: 针对性的 SQL 优化和改进建议
- ⚠️ **风险提示**: 识别潜在的性能问题和数据风险
- 📖 **最佳实践**: 行业标准和推荐做法

## 支持的 AI 提供商

### DeepSeek

DeepSeek 是一个高性能的大语言模型，特别适合代码分析和技术建议。

**优势**:
- 专注于技术和代码领域
- 响应速度快
- API 价格实惠
- 支持中文

**获取 API Key**:
1. 访问 [https://platform.deepseek.com](https://platform.deepseek.com)
2. 注册账号
3. 在控制台创建 API Key
4. 查看定价和配额

### OpenAI（兼容）

也支持 OpenAI 的 API 格式，可以使用 GPT-4 等模型。

## 配置 AI

### 1. 创建配置文件

在项目根目录创建 `.sql-diff-config.yaml`:

```yaml
ai:
  # 是否启用 AI 功能
  enabled: true
  
  # AI 提供商: deepseek, openai
  provider: deepseek
  
  # API 密钥
  api_key: sk-your-api-key-here
  
  # API 端点
  api_endpoint: https://api.deepseek.com/v1
  
  # 使用的模型
  model: deepseek-chat
  
  # 请求超时时间（秒）
  timeout: 30
```

### 2. 保护您的 API Key

⚠️ **重要**: 不要将包含真实 API Key 的配置文件提交到版本控制系统！

配置文件 `.sql-diff-config.yaml` 已经在 `.gitignore` 中，确保不会被误提交。

### 3. 验证配置

运行以下命令验证配置是否正确:

```bash
./bin/sql-diff \
  -s "CREATE TABLE test (id INT)" \
  -t "CREATE TABLE test (id INT, name VARCHAR(100))" \
  --ai
```

如果看到 AI 分析结果，说明配置成功！

## 使用 AI 分析

### 基础用法

在命令中添加 `--ai` 参数即可启用 AI 分析:

```bash
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))" \
  --ai
```

### 指定配置文件

如果配置文件不在默认位置，可以指定路径:

```bash
sql-diff \
  -s "..." \
  -t "..." \
  --config ./my-config.yaml \
  --ai
```

### 完整示例

```bash
sql-diff \
  -s "CREATE TABLE products (
    id INT PRIMARY KEY,
    name TEXT,
    price FLOAT,
    stock INT
  )" \
  -t "CREATE TABLE products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    stock INT UNSIGNED DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_created (created_at)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4" \
  --ai \
  -o migration.sql
```

## AI 输出解读

AI 分析结果包含四个部分:

### 📊 差异分析

对表结构变更的整体描述和影响评估。

**示例**:
```
检测到复杂的表结构重构，主要包括数据类型优化、
约束加强和索引优化。这些变更将提升数据一致性和查询性能。
```

### ✨ 优化建议

针对当前变更的改进建议。

**示例**:
```
1. 建议为 email 字段添加唯一索引确保邮箱唯一性
2. created_at 和 updated_at 应成对使用，便于追踪记录生命周期
3. 考虑为 name 字段添加全文索引以支持搜索功能
```

### ⚠️ 潜在风险

可能遇到的问题和需要注意的事项。

**示例**:
```
1. TEXT 改为 VARCHAR(200) 可能截断现有数据，需要先检查
2. 大表添加索引会锁表，建议使用 pt-online-schema-change
3. FLOAT 改为 DECIMAL 可能导致精度变化，需要数据迁移
```

### 📖 最佳实践

行业推荐的做法和注意事项。

**示例**:
```
1. 使用 DECIMAL 而不是 FLOAT 存储金额，避免精度问题
2. 为时间字段添加索引，提升按时间查询的性能
3. 在测试环境完整验证后再应用到生产
4. 准备回滚脚本以应对意外情况
```

## 常见问题

### Q: AI 分析需要多长时间？

通常 5-15 秒，取决于:
- 表结构的复杂度
- AI 模型的响应速度
- 网络延迟

可以通过配置文件调整 `timeout` 值（默认 30 秒）。

### Q: AI 分析会消耗多少 tokens？

典型情况下，一次分析大约消耗:
- DeepSeek: 200-500 tokens
- 成本非常低，通常不到 0.01 元

### Q: 如果 API 请求失败怎么办？

工具会显示警告信息，但不会影响核心的 DDL 生成功能:

```
⚠ AI 分析失败: API 返回错误 429: Rate limit exceeded
```

解决方法:
1. 检查 API Key 是否有效
2. 确认是否超出配额限制
3. 增加 timeout 值重试
4. 暂时禁用 AI 功能

### Q: 可以离线使用 AI 功能吗？

不可以，AI 功能需要调用在线 API。但核心的表结构比对和 DDL 生成功能完全可以离线使用。

### Q: AI 建议是否一定正确？

AI 建议仅供参考，请根据实际业务场景判断:
- ✅ 大部分建议基于最佳实践，值得采纳
- ⚠️ 某些建议可能不适合特定场景
- 🔍 始终在测试环境验证后再应用

## 最佳实践

### 1. 安全使用 API Key

```bash
# ✅ 正确: 使用配置文件
cp .sql-diff-config.example.yaml .sql-diff-config.yaml
# 编辑文件填入 API Key

# ❌ 错误: 不要在脚本或命令行中硬编码 API Key
```

### 2. 结合人工审查

```
AI 分析 → 人工审查 → 测试验证 → 生产应用
```

AI 建议是辅助工具，不应替代人工判断。

### 3. 针对性使用

对于简单的变更（如新增一个字段），可能不需要 AI 分析。
AI 在以下场景最有价值:

- 🔄 复杂的表结构重构
- ⚡ 性能优化相关的变更
- 🔒 涉及数据安全的修改
- 📊 大型系统的架构调整

### 4. 保存分析结果

将 AI 分析结果保存下来，作为文档的一部分:

```bash
sql-diff -s "..." -t "..." --ai | tee ai-analysis.txt
```

### 5. 团队协作

在代码审查时分享 AI 分析结果:

```bash
# 生成 DDL 和 AI 分析
sql-diff -s "..." -t "..." --ai -o migration.sql

# 将分析结果添加到 PR 描述中
```

## 高级配置

### 自定义提示词

未来版本将支持自定义 AI 提示词，以适应特定的业务场景。

### 批量分析

对多个表进行批量分析:

```bash
for table in users products orders; do
  echo "分析表: $table"
  sql-diff \
    -s "$(cat schema/old/$table.sql)" \
    -t "$(cat schema/new/$table.sql)" \
    --ai
done
```

## 成本估算

### DeepSeek 价格参考（2024年）

- 输入: ¥0.001/1K tokens
- 输出: ¥0.002/1K tokens

**示例成本**:
- 一次简单分析: ~300 tokens = ¥0.0006
- 一次复杂分析: ~800 tokens = ¥0.0015
- 100次分析: 约 ¥0.10

成本非常低，适合日常使用。

## 故障排查

### 错误: API 返回 401

**原因**: API Key 无效

**解决**:
1. 检查 API Key 是否正确
2. 确认 API Key 未过期
3. 验证 API Key 的权限

### 错误: API 返回 429

**原因**: 超出速率限制

**解决**:
1. 等待几分钟后重试
2. 升级 API 套餐
3. 降低请求频率

### 错误: 请求超时

**原因**: 网络延迟或响应慢

**解决**:
```yaml
ai:
  timeout: 60  # 增加超时时间到60秒
```

## 示例场景

### 场景 1: 数据库迁移

在进行数据库版本升级时，使用 AI 分析变更影响:

```bash
sql-diff \
  -s "$(mysqldump --no-data old_db users)" \
  -t "$(cat new_schema.sql)" \
  --ai
```

### 场景 2: 性能优化

分析索引优化方案:

```bash
sql-diff \
  -s "CREATE TABLE logs (...)" \
  -t "CREATE TABLE logs (..., INDEX idx_created (created_at), INDEX idx_user (user_id))" \
  --ai
```

### 场景 3: 代码审查

在 CI/CD 中集成 AI 分析:

```yaml
# .github/workflows/review.yml
- name: Analyze Schema Changes
  run: |
    sql-diff \
      -s "$(cat schema/before.sql)" \
      -t "$(cat schema/after.sql)" \
      --ai > ai-review.txt
    cat ai-review.txt >> $GITHUB_STEP_SUMMARY
```

---

💡 **提示**: AI 功能是可选的强大补充，即使不使用 AI，sql-diff 的核心功能依然完整可用！

🔗 **相关链接**:
- [DeepSeek API 文档](https://platform.deepseek.com/api-docs)
- [OpenAI API 文档](https://platform.openai.com/docs)
- [项目主页](../README.md)
