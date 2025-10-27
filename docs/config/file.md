# 配置文件

除了环境变量,SQL-Diff 也支持使用配置文件进行配置管理。

## 配置文件格式

SQL-Diff 使用 YAML 格式的配置文件:

```yaml
# .sql-diff-config.yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-your-api-key-here
  model: deepseek-chat
  api_url: https://api.deepseek.com/v1
  timeout: 30
  max_tokens: 2000
  temperature: 0.7
```

## 配置文件位置

SQL-Diff 会按以下顺序查找配置文件:

1. 当前目录: `./.sql-diff-config.yaml`
2. 用户主目录: `~/.sql-diff-config.yaml`
3. 系统配置: `/etc/sql-diff/config.yaml`

## 配置项说明

### AI 配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `ai.enabled` | boolean | `false` | 是否启用 AI 功能 |
| `ai.provider` | string | `deepseek` | AI 提供商 |
| `ai.api_key` | string | - | API 密钥 |
| `ai.model` | string | `deepseek-chat` | 模型名称 |
| `ai.api_url` | string | - | API 地址 (可选) |
| `ai.timeout` | integer | `30` | 请求超时时间(秒) |
| `ai.max_tokens` | integer | `2000` | 最大 token 数 |
| `ai.temperature` | float | `0.7` | 温度参数 (0-1) |

## 配置示例

### 基础配置

```yaml
# .sql-diff-config.yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-your-api-key
```

### 完整配置

```yaml
# .sql-diff-config.yaml
ai:
  # AI 功能开关
  enabled: true
  
  # AI 提供商
  provider: deepseek
  
  # API 密钥
  api_key: sk-your-api-key-here
  
  # 模型配置
  model: deepseek-chat
  api_url: https://api.deepseek.com/v1
  
  # 性能配置
  timeout: 30
  max_tokens: 2000
  temperature: 0.7
  
# 输出配置 (计划中)
output:
  format: text
  color: true
  verbose: false
```

### 团队共享配置

```yaml
# team-config.yaml
# 团队共享的配置模板

ai:
  enabled: true
  provider: deepseek
  model: deepseek-chat
  api_url: https://api.deepseek.com/v1
  timeout: 30
  max_tokens: 2000
  
  # API Key 由个人设置,不提交到版本控制
  # api_key: YOUR_API_KEY_HERE
```

团队成员使用:

```bash
# 1. 复制团队配置
cp team-config.yaml ~/.sql-diff-config.yaml

# 2. 添加个人 API Key
# 编辑 ~/.sql-diff-config.yaml,添加:
# ai:
#   api_key: sk-your-personal-key

# 3. 或使用环境变量
export SQL_DIFF_AI_API_KEY=sk-your-personal-key
```

## 配置优先级

配置加载优先级 (从高到低):

1. **命令行参数**
   ```bash
   sql-diff -s "..." -t "..." --ai --provider=openai
   ```

2. **环境变量**
   ```bash
   export SQL_DIFF_AI_ENABLED=true
   export SQL_DIFF_AI_API_KEY=sk-xxx
   ```

3. **配置文件**
   ```yaml
   # .sql-diff-config.yaml
   ai:
     enabled: true
     api_key: sk-xxx
   ```

4. **默认值**
   ```go
   // 代码中的默认值
   enabled: false
   provider: "deepseek"
   ```

## 使用示例

### 示例 1: 个人开发环境

创建 `~/.sql-diff-config.yaml`:

```yaml
ai:
  enabled: true
  api_key: sk-your-personal-key
  provider: deepseek
  model: deepseek-chat
```

然后直接使用:

```bash
# 自动使用配置文件中的设置
sql-diff -s old.sql -t new.sql --ai
```

### 示例 2: 项目级配置

在项目根目录创建 `.sql-diff-config.yaml`:

```yaml
ai:
  enabled: true
  provider: deepseek
  model: deepseek-chat
  # API Key 通过环境变量提供
```

在 `.env` 文件中:

```bash
SQL_DIFF_AI_API_KEY=sk-project-key
```

在 `.gitignore` 中:

```
.env
.sql-diff-config.yaml
```

### 示例 3: CI/CD 环境

在 CI 中使用配置文件:

```yaml
# .github/workflows/schema-check.yml
- name: Create config file
  run: |
    cat > .sql-diff-config.yaml << EOF
    ai:
      enabled: true
      provider: deepseek
      api_key: ${{ secrets.DEEPSEEK_API_KEY }}
      model: deepseek-chat
      timeout: 60
    EOF

- name: Run SQL-Diff
  run: sql-diff -s old.sql -t new.sql --ai
```

## 安全最佳实践

### 1. 不要提交敏感信息

在 `.gitignore` 中添加:

```
# SQL-Diff 配置
.sql-diff-config.yaml
.sql-diff-config.yml
.sql-diff-env

# 环境变量
.env
.env.local
```

### 2. 使用示例配置文件

提供不包含敏感信息的示例:

```yaml
# .sql-diff-config.example.yaml
ai:
  enabled: true
  provider: deepseek
  model: deepseek-chat
  
  # 请在这里填入您的 API Key
  api_key: YOUR_API_KEY_HERE
  
  # 可选配置
  timeout: 30
  max_tokens: 2000
```

团队成员可以复制并修改:

```bash
cp .sql-diff-config.example.yaml .sql-diff-config.yaml
# 然后编辑 .sql-diff-config.yaml 填入真实的 API Key
```

### 3. 文件权限

确保配置文件权限正确:

```bash
# 只有所有者可读写
chmod 600 ~/.sql-diff-config.yaml

# 验证权限
ls -l ~/.sql-diff-config.yaml
# 应显示: -rw------- 1 user group ...
```

### 4. 使用密钥管理工具

结合密钥管理工具使用:

```bash
# 使用 Pass (密码管理器)
pass insert sql-diff/api-key

# 在配置文件中引用
# (需要手动获取并设置环境变量)
export SQL_DIFF_AI_API_KEY=$(pass show sql-diff/api-key)
```

## 配置验证

### 查看当前配置

```bash
sql-diff config --show
```

输出:

```
📋 当前配置:

AI 配置:
✓ AI 功能: 已启用
✓ API Key: sk-xxx******
✓ 提供商: deepseek
✓ 模型: deepseek-chat
✓ API 地址: https://api.deepseek.com/v1
✓ 超时时间: 30秒
✓ 最大 Tokens: 2000

配置来源:
- 环境变量: SQL_DIFF_AI_ENABLED
- 配置文件: ~/.sql-diff-config.yaml
```

### 测试配置

```bash
# 测试 AI 功能是否正常
sql-diff \
  -s "CREATE TABLE test (id INT);" \
  -t "CREATE TABLE test (id INT, name VARCHAR(100));" \
  --ai
```

## 常见问题

### Q: 配置文件不生效?

检查配置文件位置和格式:

```bash
# 检查文件是否存在
ls -la .sql-diff-config.yaml

# 检查 YAML 格式
cat .sql-diff-config.yaml
```

### Q: 如何覆盖配置文件的设置?

使用环境变量或命令行参数:

```bash
# 临时覆盖
SQL_DIFF_AI_ENABLED=false sql-diff -s old.sql -t new.sql

# 或使用命令行参数 (计划中)
sql-diff -s old.sql -t new.sql --no-ai
```

### Q: 如何在多个项目间共享配置?

使用用户级配置文件:

```bash
# 在主目录创建配置
vim ~/.sql-diff-config.yaml

# 所有项目都会使用这个配置
# 项目级配置文件优先级更高
```

## 迁移指南

### 从配置文件迁移到环境变量

```bash
# 1. 读取现有配置文件
cat .sql-diff-config.yaml

# 2. 转换为环境变量
sql-diff config \
  --api-key="$(yq .ai.api_key .sql-diff-config.yaml)" \
  --ai-enabled=true

# 3. 保存到 shell 配置
sql-diff config --ai-enabled=true --api-key="sk-xxx" >> ~/.zshrc
source ~/.zshrc

# 4. (可选) 删除配置文件
rm .sql-diff-config.yaml
```

### 从环境变量迁移到配置文件

```bash
# 1. 生成配置文件
cat > ~/.sql-diff-config.yaml << EOF
ai:
  enabled: ${SQL_DIFF_AI_ENABLED:-true}
  api_key: ${SQL_DIFF_AI_API_KEY}
  provider: ${SQL_DIFF_AI_PROVIDER:-deepseek}
  model: ${SQL_DIFF_AI_MODEL:-deepseek-chat}
EOF

# 2. (可选) 从 shell 配置中移除环境变量
```

## 下一步

- [环境变量配置](/config/environment) - 使用环境变量配置
- [配置命令](/config/command) - 使用命令行配置
- [快速开始](/guide/getting-started) - 开始使用 SQL-Diff
