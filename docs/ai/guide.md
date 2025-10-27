# AI 功能指南

SQL-Diff 集成了 AI 大模型,可以智能分析表结构变更,提供优化建议和风险提示。

## 概述

AI 功能是可选的增强功能,可以帮助您:

- 📊 **智能分析差异**: 理解结构变更的业务含义
- 💡 **优化建议**: 提供索引、性能、数据类型等优化建议
- ⚠️ **风险识别**: 识别潜在的数据丢失、性能问题等风险
- ✅ **最佳实践**: 根据行业标准提供最佳实践建议

## 快速开始

### 1. 配置 API Key

首先需要获取 DeepSeek API Key (或其他支持的提供商):

```bash
sql-diff config \
  --ai-enabled=true \
  --api-key="sk-your-api-key-here" \
  --provider=deepseek
```

将输出保存到 shell 配置文件:

```bash
sql-diff config --ai-enabled=true --api-key="sk-xxx" >> ~/.zshrc
source ~/.zshrc
```

### 2. 使用 AI 分析

在比对命令中添加 `--ai` 参数:

```bash
sql-diff -s "CREATE TABLE users (id INT);" \
         -t "CREATE TABLE users (id INT, email VARCHAR(255));" \
         --ai
```

## AI 分析输出

AI 会返回四个部分的分析结果:

### 📊 差异分析

对检测到的差异进行高层次总结:

```
📊 差异分析:
检测到表结构新增了 1 个列字段 (email),这是一个向后兼容的变更。
新增字段为 VARCHAR 类型,适合存储电子邮件地址。
```

### 💡 优化建议

提供具体的优化建议:

```
💡 优化建议:
- 建议为 email 字段添加唯一索引,确保邮箱地址唯一性
- 可以添加邮箱格式验证约束
- 考虑添加 created_at 和 updated_at 时间戳字段
- 如果邮箱可能为空,建议使用 NULL 而非空字符串
```

### ⚠️ 潜在风险

识别可能的风险:

```
⚠️  潜在风险:
- 如果表数据量很大,添加 NOT NULL 列需要设置默认值
- 新增列会增加行大小,可能影响性能
- 建议在低峰期执行变更
```

### ✅ 最佳实践

推荐行业最佳实践:

```
✅ 最佳实践:
- 为新字段添加清晰的注释说明用途
- 在测试环境充分验证后再部署到生产环境
- 考虑使用在线 DDL 工具避免锁表
- 保留旧版本 schema 以便回滚
```

## 配置选项

### 环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `SQL_DIFF_AI_ENABLED` | 是否启用 AI | `true` |
| `SQL_DIFF_AI_API_KEY` | API 密钥 | `sk-xxx` |
| `SQL_DIFF_AI_PROVIDER` | AI 提供商 | `deepseek` |
| `SQL_DIFF_AI_MODEL` | 模型名称 | `deepseek-chat` |
| `SQL_DIFF_AI_API_URL` | API 地址 | 自定义 URL |

### 配置文件

创建 `.sql-diff-config.yaml`:

```yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-your-api-key
  model: deepseek-chat
  api_url: https://api.deepseek.com  # 可选
```

## 支持的 AI 提供商

### DeepSeek (推荐)

**特点**:
- 🚀 高性能,响应快
- 💰 价格实惠
- 🇨🇳 支持中文
- 🎯 针对代码优化

**配置示例**:
```bash
sql-diff config \
  --provider=deepseek \
  --api-key="sk-xxx" \
  --model="deepseek-chat"
```

**获取 API Key**: [DeepSeek 平台](https://platform.deepseek.com)

### OpenAI (即将支持)

```bash
sql-diff config \
  --provider=openai \
  --api-key="sk-xxx" \
  --model="gpt-4"
```

### 自定义提供商

支持任何兼容 OpenAI API 格式的提供商:

```bash
sql-diff config \
  --provider=custom \
  --api-url="https://your-api.com/v1/chat/completions" \
  --api-key="your-key" \
  --model="your-model"
```

## 使用场景

### 1. 复杂变更分析

对于包含多个变更的复杂场景:

```bash
sql-diff \
  -s "$(cat old_schema.sql)" \
  -t "$(cat new_schema.sql)" \
  --ai
```

AI 会综合分析所有变更,提供整体建议。

### 2. 生产环境变更评估

在执行生产环境变更前:

```bash
sql-diff \
  -s "$(mysqldump --no-data prod_db users)" \
  -t "$(cat new_users_schema.sql)" \
  --ai \
  --output review.txt
```

查看 AI 的风险评估,决定是否执行。

### 3. Code Review

在 PR 中使用 AI 分析:

```bash
# 在 CI 中运行
sql-diff -s old.sql -t new.sql --ai --format json > ai_review.json

# 提取关键信息
jq '.ai_analysis.risks' ai_review.json
```

### 4. 学习最佳实践

对于初学者,AI 建议可以帮助学习数据库设计:

```bash
sql-diff \
  -s "CREATE TABLE users (id INT);" \
  -t "CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(50));" \
  --ai
```

AI 会解释为什么需要主键、AUTO_INCREMENT 的作用等。

## 高级功能

### 自定义提示词 (即将支持)

未来版本将支持自定义 AI 提示词:

```yaml
# .sql-diff-config.yaml
ai:
  enabled: true
  custom_prompt: |
    你是一个资深的数据库架构师,专注于:
    1. 性能优化
    2. 数据安全
    3. 高可用设计
    
    请分析以下表结构变更...
```

### 批量分析

分析多个表的变更:

```bash
#!/bin/bash
for table in users products orders; do
  sql-diff \
    -s "$(cat old/${table}.sql)" \
    -t "$(cat new/${table}.sql)" \
    --ai > reports/${table}_analysis.txt
done

# 汇总所有风险
grep -h "⚠️" reports/*.txt > all_risks.txt
```

### 成本控制

AI 调用会产生费用,可以通过以下方式控制:

```bash
# 只在需要时启用 AI
sql-diff -s "..." -t "..." --ai  # 启用
sql-diff -s "..." -t "..."       # 不启用

# 使用环境变量控制
export SQL_DIFF_AI_ENABLED=false  # 全局禁用
```

## 性能和成本

### 响应时间

典型的 AI 分析响应时间:

- 简单变更 (1-3 个差异): ~2-5 秒
- 中等复杂度 (4-10 个差异): ~5-10 秒
- 复杂变更 (10+ 个差异): ~10-20 秒

### API 成本

以 DeepSeek 为例:

- 输入: ¥0.001/1K tokens
- 输出: ¥0.002/1K tokens

典型使用成本:
- 单次分析: ~¥0.001-0.01
- 每月 100 次分析: ~¥0.1-1

::: tip 成本优化
只在复杂变更或生产环境变更时使用 AI 功能,日常开发可以不启用。
:::

## 隐私和安全

### 数据处理

- ✅ 只发送 DDL 语句,不包含实际数据
- ✅ 不存储您的 API Key
- ✅ 不记录您的 SQL 语句
- ✅ 使用 HTTPS 加密传输

### API Key 安全

::: warning 安全提示
- 不要在代码中硬编码 API Key
- 使用环境变量或配置文件管理密钥
- 定期轮换 API Key
- 不要将包含 API Key 的配置文件提交到版本控制
:::

推荐做法:

```bash
# .gitignore
.sql-diff-config.yaml
.env

# 使用环境变量
export SQL_DIFF_AI_API_KEY=$(cat ~/.secrets/deepseek_key)
```

## 故障排查

### AI 功能不工作

检查配置:

```bash
# 查看当前配置
sql-diff config --show

# 应该显示:
# ✓ AI 功能: 已启用
# ✓ API Key: sk-xxx***
# ✓ 提供商: deepseek
```

### API 调用失败

常见错误和解决方案:

**错误: Invalid API Key**
```bash
# 检查 API Key 是否正确
echo $SQL_DIFF_AI_API_KEY

# 重新配置
sql-diff config --api-key="correct-key"
```

**错误: Rate Limit Exceeded**
```bash
# 等待一段时间后重试
# 或升级 API 套餐
```

**错误: Network Timeout**
```bash
# 检查网络连接
curl -I https://api.deepseek.com

# 设置代理 (如需要)
export HTTPS_PROXY=http://proxy.example.com:8080
```

### 输出格式问题

如果 AI 输出格式异常:

```bash
# 使用 verbose 模式查看原始响应
sql-diff -s "..." -t "..." --ai --verbose

# 查看完整的 API 响应
```

## 最佳实践

1. **选择性使用**: 只在复杂变更或重要变更时使用 AI
2. **结合人工审查**: AI 建议仅供参考,最终决策需人工判断
3. **保存分析结果**: 将 AI 分析保存为文档,便于后续查阅
4. **定期更新**: 保持工具更新以获得最新的 AI 能力
5. **测试验证**: AI 建议需在测试环境充分验证

## 下一步

- [DeepSeek 集成](/ai/deepseek) - DeepSeek 详细配置
- [最佳实践](/ai/best-practices) - AI 使用最佳实践
- [示例](/examples/advanced) - 实际使用示例
