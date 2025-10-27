# DeepSeek 集成

SQL-Diff 默认集成了 DeepSeek AI,这是一个高性能、低成本的大语言模型,特别适合代码分析场景。

## 为什么选择 DeepSeek

### 优势

- 🚀 **高性能**: 响应速度快,延迟低
- 💰 **成本低**: 相比 GPT-4 等模型,成本降低 90%+
- 🇨🇳 **中文友好**: 对中文支持优秀
- 🎯 **代码优化**: 针对代码理解和生成优化
- 🔒 **数据安全**: 国内服务,数据合规

### 性能对比

| 指标 | DeepSeek | GPT-4 | GPT-3.5 |
|------|----------|-------|---------|
| 响应速度 | ⚡⚡⚡ | ⚡⚡ | ⚡⚡⚡ |
| 成本 | 💰 | 💰💰💰💰 | 💰💰 |
| 中文能力 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| 代码理解 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

## 快速开始

### 1. 获取 API Key

访问 [DeepSeek 开放平台](https://platform.deepseek.com):

1. 注册账号
2. 进入 API Keys 页面
3. 创建新的 API Key
4. 复制保存 (只显示一次)

### 2. 配置 SQL-Diff

使用命令行配置:

```bash
sql-diff config \
  --ai-enabled=true \
  --provider=deepseek \
  --api-key="sk-your-api-key-here" \
  --model="deepseek-chat"
```

将输出保存到环境变量:

```bash
# 生成环境变量配置
sql-diff config \
  --ai-enabled=true \
  --api-key="sk-xxx" > ~/.sql-diff-env

# 添加到 shell 配置
echo "source ~/.sql-diff-env" >> ~/.zshrc
source ~/.zshrc
```

### 3. 验证配置

```bash
# 查看配置
sql-diff config --show

# 测试 AI 功能
sql-diff \
  -s "CREATE TABLE test (id INT);" \
  -t "CREATE TABLE test (id INT, name VARCHAR(100));" \
  --ai
```

## 配置选项

### 环境变量

```bash
# 必需配置
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_API_KEY=sk-your-api-key
export SQL_DIFF_AI_PROVIDER=deepseek

# 可选配置
export SQL_DIFF_AI_MODEL=deepseek-chat           # 模型名称
export SQL_DIFF_AI_API_URL=https://api.deepseek.com  # API 地址
export SQL_DIFF_AI_TIMEOUT=30                    # 超时时间(秒)
export SQL_DIFF_AI_MAX_TOKENS=2000               # 最大 tokens
```

### 配置文件

创建 `.sql-diff-config.yaml`:

```yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-your-api-key
  model: deepseek-chat
  api_url: https://api.deepseek.com/v1
  timeout: 30
  max_tokens: 2000
```

### 配置优先级

配置加载优先级 (从高到低):

1. 命令行参数
2. 环境变量
3. 配置文件
4. 默认值

## API 详解

### 模型选择

DeepSeek 提供多个模型:

| 模型 | 说明 | 适用场景 |
|------|------|----------|
| `deepseek-chat` | 通用对话模型 | **推荐用于 SQL-Diff** |
| `deepseek-coder` | 代码专用模型 | 代码生成和理解 |

推荐使用 `deepseek-chat`:

```bash
export SQL_DIFF_AI_MODEL=deepseek-chat
```

### API 端点

DeepSeek API 兼容 OpenAI API 格式:

```
POST https://api.deepseek.com/v1/chat/completions
```

请求格式:

```json
{
  "model": "deepseek-chat",
  "messages": [
    {
      "role": "system",
      "content": "你是一个数据库专家..."
    },
    {
      "role": "user",
      "content": "分析以下表结构变更..."
    }
  ],
  "max_tokens": 2000,
  "temperature": 0.7
}
```

### 速率限制

DeepSeek 的速率限制:

| 套餐 | RPM (每分钟请求) | TPM (每分钟 tokens) |
|------|------------------|---------------------|
| 免费 | 60 | 60,000 |
| 标准 | 600 | 600,000 |
| 企业 | 定制 | 定制 |

::: tip
对于个人使用,免费套餐通常足够。如果遇到速率限制,可以添加延迟或升级套餐。
:::

## 定价

DeepSeek 定价 (截至 2025 年):

| 类型 | 价格 |
|------|------|
| 输入 tokens | ¥0.001 / 1K tokens |
| 输出 tokens | ¥0.002 / 1K tokens |

### 成本估算

典型的 SQL-Diff 使用:

**简单比对**:
- 输入: ~500 tokens
- 输出: ~300 tokens
- 成本: ¥0.001

**复杂比对**:
- 输入: ~1500 tokens
- 输出: ~800 tokens
- 成本: ¥0.003

**每月使用** (100 次分析):
- 总成本: ~¥0.1-0.3

::: tip 成本优化
- 只在需要时使用 `--ai` 参数
- 使用缓存避免重复分析
- 批量分析时合并请求
:::

## 高级用法

### 自定义提示词

通过环境变量自定义系统提示词:

```bash
export SQL_DIFF_AI_SYSTEM_PROMPT="你是一个资深的 MySQL DBA,专注于性能优化和数据安全。请分析表结构变更,重点关注索引优化和潜在的性能问题。"
```

### 调整温度参数

温度控制输出的随机性:

```bash
# 更保守的输出 (推荐)
export SQL_DIFF_AI_TEMPERATURE=0.3

# 更有创意的输出
export SQL_DIFF_AI_TEMPERATURE=0.9
```

### 超时设置

设置 API 调用超时:

```bash
# 30 秒超时 (默认)
export SQL_DIFF_AI_TIMEOUT=30

# 对于复杂分析,可以增加超时
export SQL_DIFF_AI_TIMEOUT=60
```

## 实际示例

### 示例 1: 基础分析

```bash
sql-diff \
  -s "CREATE TABLE users (id INT, name VARCHAR(100));" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(255) NOT NULL, email VARCHAR(255));" \
  --ai
```

输出:

```
🤖 AI 分析结果:

📊 差异分析:
检测到以下变更:
1. id 字段添加了 PRIMARY KEY 约束
2. name 字段类型从 VARCHAR(100) 扩展到 VARCHAR(255),并添加 NOT NULL 约束
3. 新增 email 字段

这些都是向后兼容的改进型变更。

💡 优化建议:
- 为 email 字段添加唯一索引: `ADD UNIQUE KEY uk_email (email)`
- 考虑添加邮箱格式验证
- 建议添加 created_at 和 updated_at 时间戳字段

⚠️  潜在风险:
- name 字段添加 NOT NULL 约束,确保现有数据没有空值
- email 字段如果允许为空,建议显式声明 NULL

✅ 最佳实践:
- 主键使用得当
- 建议为字段添加注释
- 考虑添加索引提升查询性能
```

### 示例 2: 复杂场景分析

```bash
sql-diff \
  -s "$(cat production_schema.sql)" \
  -t "$(cat new_schema.sql)" \
  --ai \
  --verbose
```

这会显示详细的 API 交互日志,包括:
- 发送的完整 prompt
- API 响应时间
- Token 使用情况
- 解析过程

### 示例 3: 批量分析

```bash
#!/bin/bash

# 分析多个表并生成报告
for table in users products orders; do
  echo "=== Analyzing $table ===" >> analysis_report.txt
  
  sql-diff \
    -s "$(cat old/${table}.sql)" \
    -t "$(cat new/${table}.sql)" \
    --ai >> analysis_report.txt
  
  echo "" >> analysis_report.txt
done

# 查看所有风险
grep -A 5 "⚠️" analysis_report.txt > all_risks.txt
```

## 故障排查

### API Key 无效

```bash
# 错误信息
Error: Invalid API Key

# 解决方案
# 1. 检查 API Key 是否正确
echo $SQL_DIFF_AI_API_KEY

# 2. 重新获取并配置
sql-diff config --api-key="new-key"
```

### 网络连接失败

```bash
# 测试连接
curl -I https://api.deepseek.com

# 如果无法访问,设置代理
export HTTPS_PROXY=http://proxy.example.com:8080
```

### 速率限制

```bash
# 错误信息
Error: Rate limit exceeded

# 解决方案
# 1. 等待 60 秒后重试
sleep 60

# 2. 或升级 API 套餐
# 访问 https://platform.deepseek.com
```

### 响应超时

```bash
# 增加超时时间
export SQL_DIFF_AI_TIMEOUT=60

# 重试请求
sql-diff -s "..." -t "..." --ai
```

## 安全性

### API Key 保护

::: danger 安全警告
- 永远不要将 API Key 硬编码在代码中
- 不要将包含 API Key 的文件提交到版本控制
- 定期轮换 API Key
:::

最佳实践:

```bash
# .gitignore
.sql-diff-config.yaml
.env
.sql-diff-env

# 使用密钥管理工具
# 1. macOS Keychain
security add-generic-password -s sql-diff -a api-key -w "sk-xxx"
export SQL_DIFF_AI_API_KEY=$(security find-generic-password -s sql-diff -w)

# 2. 环境变量文件
echo "export SQL_DIFF_AI_API_KEY=sk-xxx" > ~/.sql-diff-secrets
chmod 600 ~/.sql-diff-secrets
source ~/.sql-diff-secrets
```

### 数据隐私

SQL-Diff 发送给 DeepSeek 的数据:

✅ **会发送**:
- DDL 语句 (表结构定义)
- 检测到的差异摘要

❌ **不会发送**:
- 实际的表数据
- 数据库连接信息
- 您的业务逻辑代码

### 合规性

DeepSeek 是国内服务,适合需要数据合规的场景:
- 数据不出境
- 符合国内数据安全法规
- 提供企业级 SLA

## 监控和调试

### 启用详细日志

```bash
# 查看完整的 API 交互
sql-diff -s "..." -t "..." --ai --verbose

# 输出包含:
# - 发送的 prompt
# - API 响应
# - Token 统计
# - 处理时间
```

### Token 使用统计

```bash
# 查看 Token 使用情况
sql-diff -s "..." -t "..." --ai --verbose 2>&1 | grep -i token

# 输出示例:
# Tokens used: 823 input, 456 output
# Estimated cost: ¥0.0018
```

### 性能监控

```bash
# 测量响应时间
time sql-diff -s "..." -t "..." --ai

# 输出:
# real    0m3.245s
# user    0m0.123s
# sys     0m0.034s
```

## 下一步

- [AI 最佳实践](/ai/best-practices) - 高效使用 AI 的技巧
- [示例](/examples/advanced) - 更多实际案例
- [CLI 工具](/guide/cli) - 命令行详细用法
