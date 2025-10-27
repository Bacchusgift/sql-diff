# AI 使用最佳实践

本文档总结了在 SQL-Diff 中高效使用 AI 功能的最佳实践和技巧。

## 核心原则

### 1. 选择性使用

**不是所有场景都需要 AI**

✅ **推荐使用 AI**:
- 复杂的多列变更
- 生产环境的重要变更
- 不确定是否有风险的变更
- 学习数据库设计的场景

❌ **不需要 AI**:
- 简单的单列添加/删除
- 明显的日常开发变更
- 重复性的批量操作
- 时间敏感的快速检查

示例:

```bash
# 简单变更 - 不需要 AI
sql-diff -s "CREATE TABLE t (id INT);" -t "CREATE TABLE t (id INT, name VARCHAR(50));"

# 复杂变更 - 使用 AI
sql-diff -s "$(cat old_schema.sql)" -t "$(cat new_schema.sql)" --ai
```

### 2. 结合人工判断

AI 建议仅供参考,不应盲目执行:

::: warning 重要
AI 的建议需要结合实际业务场景判断:
- AI 不了解您的业务逻辑
- AI 不知道您的数据量和访问模式
- AI 可能给出过于保守或激进的建议
:::

最佳实践:

```bash
# 1. 获取 AI 分析
sql-diff -s "..." -t "..." --ai > analysis.txt

# 2. 人工审查
cat analysis.txt

# 3. 结合业务判断
# - 是否需要这个索引?
# - 数据量是否支持这个变更?
# - 是否有更好的替代方案?

# 4. 在测试环境验证
mysql -h test-db < migration.sql

# 5. 确认后执行到生产
```

### 3. 保存分析结果

将 AI 分析作为文档保存:

```bash
# 保存到变更文档
sql-diff -s "..." -t "..." --ai > docs/migrations/20251022_users_schema.md

# 提交到版本控制
git add docs/migrations/20251022_users_schema.md
git commit -m "Add schema migration analysis for users table"
```

## 工作流程最佳实践

### 开发环境

在开发阶段,快速迭代不需要 AI:

```bash
# 快速检查差异
sql-diff -s old.sql -t new.sql

# 只有在最终确定方案时才使用 AI
sql-diff -s old.sql -t final.sql --ai > final_analysis.txt
```

### 测试环境

在测试环境,使用 AI 评估变更:

```bash
# 生成变更脚本和分析
sql-diff \
  -s "$(mysqldump --no-data test_db users)" \
  -t "$(cat new_users.sql)" \
  --ai \
  --output test_migration.sql

# 应用到测试环境
mysql -h test-db -u user -p test_db < test_migration.sql

# 执行测试用例验证
npm run test:integration
```

### 生产环境

生产环境变更必须使用 AI 评估:

```bash
#!/bin/bash
# prod_migration.sh

# 1. 导出当前生产结构
mysqldump --no-data -h prod-db users > current_prod.sql

# 2. AI 分析
sql-diff \
  -s current_prod.sql \
  -t new_users.sql \
  --ai \
  --verbose > migration_analysis.txt

# 3. 人工审查
cat migration_analysis.txt
read -p "Review complete. Proceed? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Migration cancelled."
  exit 1
fi

# 4. 生成最终脚本
sql-diff -s current_prod.sql -t new_users.sql > migration.sql

# 5. 备份
mysqldump -h prod-db users > backup_$(date +%Y%m%d_%H%M%S).sql

# 6. 执行 (建议在低峰期)
mysql -h prod-db -u user -p prod_db < migration.sql

echo "Migration completed!"
```

## 成本优化

### 1. 避免重复分析

使用缓存避免重复的 AI 调用:

```bash
#!/bin/bash
# 智能缓存脚本

SOURCE_HASH=$(echo "$SOURCE_SQL" | md5)
TARGET_HASH=$(echo "$TARGET_SQL" | md5)
CACHE_KEY="${SOURCE_HASH}_${TARGET_HASH}"
CACHE_FILE="~/.sql-diff-cache/${CACHE_KEY}.txt"

if [ -f "$CACHE_FILE" ]; then
  echo "Using cached analysis..."
  cat "$CACHE_FILE"
else
  echo "Generating new analysis..."
  sql-diff -s "$SOURCE_SQL" -t "$TARGET_SQL" --ai | tee "$CACHE_FILE"
fi
```

### 2. 批量分析优化

对于多个表的分析,合并请求:

```bash
# ❌ 低效: 每个表单独调用
for table in users products orders; do
  sql-diff -s old/${table}.sql -t new/${table}.sql --ai
done

# ✅ 高效: 合并分析
cat > batch_analysis.txt << EOF
Table: users
$(sql-diff -s old/users.sql -t new/users.sql)

Table: products
$(sql-diff -s old/products.sql -t new/products.sql)

Table: orders
$(sql-diff -s old/orders.sql -t new/orders.sql)
EOF

# 一次 AI 调用分析所有变更
sql-diff --analyze-file batch_analysis.txt --ai
```

### 3. 按需启用

只在真正需要时启用 AI:

```bash
# 环境变量控制
export SQL_DIFF_AI_ENABLED=false  # 默认关闭

# 只在需要时启用
SQL_DIFF_AI_ENABLED=true sql-diff -s "..." -t "..." --ai

# 或使用 function
sql-diff-ai() {
  SQL_DIFF_AI_ENABLED=true sql-diff "$@" --ai
}
```

## 提示词优化

### 1. 明确分析重点

通过环境变量自定义提示词:

```bash
# 性能优化重点
export SQL_DIFF_AI_SYSTEM_PROMPT="你是数据库性能专家。分析时重点关注:
1. 索引使用是否合理
2. 大表变更的性能影响
3. 查询优化建议"

# 安全重点
export SQL_DIFF_AI_SYSTEM_PROMPT="你是数据库安全专家。分析时重点关注:
1. 数据丢失风险
2. 权限和访问控制
3. 数据加密需求"
```

### 2. 上下文信息

提供额外的上下文获得更好的建议:

```bash
# 在注释中添加上下文
cat > new_schema.sql << EOF
-- 业务背景: 用户表,预计 1000 万行数据
-- 查询模式: 主要通过 email 查询,每秒 1000 QPS
-- 性能要求: 查询响应时间 < 100ms

CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  email VARCHAR(255) NOT NULL,
  name VARCHAR(100),
  INDEX idx_email (email)
) ENGINE=InnoDB;
EOF

sql-diff -s old.sql -t new_schema.sql --ai
```

## 安全和隐私

### 1. 敏感信息脱敏

避免在 SQL 中包含敏感信息:

```bash
# ❌ 不好: 包含敏感默认值
CREATE TABLE users (
  api_secret VARCHAR(255) DEFAULT 'prod-secret-key-123'
);

# ✅ 好: 使用占位符
CREATE TABLE users (
  api_secret VARCHAR(255) DEFAULT 'PLACEHOLDER'
);
```

### 2. API Key 管理

安全管理 API Key:

```bash
# 1. 使用专用的密钥文件
echo "export SQL_DIFF_AI_API_KEY=sk-xxx" > ~/.sql-diff-secrets
chmod 600 ~/.sql-diff-secrets

# 2. 在 .gitignore 中忽略
echo ".sql-diff-secrets" >> ~/.gitignore

# 3. 在 shell 配置中引用
echo "[ -f ~/.sql-diff-secrets ] && source ~/.sql-diff-secrets" >> ~/.zshrc

# 4. 定期轮换
# 每 90 天更新一次 API Key
```

### 3. 审计日志

记录 AI 使用情况:

```bash
# 创建审计日志
sql-diff-audit() {
  local LOG_FILE=~/.sql-diff-audit.log
  echo "$(date): AI analysis requested" >> $LOG_FILE
  echo "User: $(whoami)" >> $LOG_FILE
  echo "Tables: $@" >> $LOG_FILE
  
  sql-diff "$@" --ai | tee -a $LOG_FILE
}
```

## 性能优化

### 1. 超时设置

根据复杂度调整超时:

```bash
# 简单分析: 短超时
SQL_DIFF_AI_TIMEOUT=15 sql-diff -s simple1.sql -t simple2.sql --ai

# 复杂分析: 长超时
SQL_DIFF_AI_TIMEOUT=60 sql-diff -s complex1.sql -t complex2.sql --ai
```

### 2. 并发控制

避免并发调用导致速率限制:

```bash
# ❌ 不好: 并发调用
for table in ${TABLES[@]}; do
  sql-diff -s old/${table}.sql -t new/${table}.sql --ai &
done
wait

# ✅ 好: 顺序调用,添加延迟
for table in ${TABLES[@]}; do
  sql-diff -s old/${table}.sql -t new/${table}.sql --ai
  sleep 2  # 避免速率限制
done
```

### 3. 结果缓存

缓存 AI 分析结果:

```bash
#!/bin/bash
# ai_cache.sh

cache_key() {
  echo "$1$2" | md5sum | cut -d' ' -f1
}

cached_sql_diff() {
  local source="$1"
  local target="$2"
  local cache_dir=~/.sql-diff-cache
  local key=$(cache_key "$source" "$target")
  local cache_file="$cache_dir/$key"
  
  mkdir -p "$cache_dir"
  
  if [ -f "$cache_file" ] && [ $(find "$cache_file" -mmin -60) ]; then
    # 缓存有效期 60 分钟
    cat "$cache_file"
  else
    sql-diff -s "$source" -t "$target" --ai | tee "$cache_file"
  fi
}
```

## CI/CD 集成

### GitHub Actions

```yaml
# .github/workflows/schema-check.yml
name: Schema Change Review

on:
  pull_request:
    paths:
      - 'db/schema/**'

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      
      - name: Install SQL-Diff
        run: go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
      
      - name: Analyze Schema Changes
        id: analyze
        run: |
          sql-diff \
            -s "$(cat db/schema/current.sql)" \
            -t "$(cat db/schema/new.sql)" \
            --ai \
            --format json > analysis.json
        env:
          SQL_DIFF_AI_API_KEY: ${{ secrets.DEEPSEEK_API_KEY }}
      
      - name: Check for Risks
        run: |
          RISKS=$(jq '.ai_analysis.risks | length' analysis.json)
          if [ "$RISKS" -gt 0 ]; then
            echo "⚠️  Detected $RISKS potential risks"
            jq '.ai_analysis.risks' analysis.json
            exit 1
          fi
      
      - name: Comment PR
        uses: actions/github-script@v5
        with:
          script: |
            const fs = require('fs');
            const analysis = JSON.parse(fs.readFileSync('analysis.json'));
            const comment = `## 🤖 AI Schema Analysis\n\n${analysis.ai_analysis.summary}`;
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });
```

### GitLab CI

```yaml
# .gitlab-ci.yml
schema-review:
  stage: test
  image: golang:1.21
  script:
    - go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
    - |
      sql-diff \
        -s "$(cat db/current.sql)" \
        -t "$(cat db/new.sql)" \
        --ai > analysis.txt
    - cat analysis.txt
  only:
    changes:
      - db/schema/**
  variables:
    SQL_DIFF_AI_API_KEY: $DEEPSEEK_API_KEY
```

## 团队协作

### 1. 统一配置

团队共享配置模板:

```yaml
# team-config.yaml
ai:
  enabled: true
  provider: deepseek
  model: deepseek-chat
  # API Key 由个人配置
```

团队成员各自设置 API Key:

```bash
cp team-config.yaml ~/.sql-diff-config.yaml
sql-diff config --api-key="personal-key"
```

### 2. 分析报告模板

统一分析报告格式:

```bash
# generate_report.sh
#!/bin/bash

cat > migration_report.md << EOF
# Schema Migration Report

**Date**: $(date)
**Author**: $(git config user.name)
**Branch**: $(git branch --show-current)

## Changes Summary

$(sql-diff -s old.sql -t new.sql)

## AI Analysis

$(sql-diff -s old.sql -t new.sql --ai)

## Review Checklist

- [ ] Reviewed by DBA
- [ ] Tested in staging environment
- [ ] Performance impact assessed
- [ ] Rollback plan prepared
- [ ] Documentation updated
EOF
```

### 3. Code Review 流程

PR 审查时的最佳实践:

1. **开发者**: 提交 PR 时附带 AI 分析
2. **审查者**: 检查 AI 识别的风险
3. **DBA**: 对生产变更进行终审
4. **测试**: 在 staging 环境验证

## 常见错误及避免

### 1. 过度依赖 AI

❌ **错误做法**:
```bash
# 盲目执行 AI 建议
sql-diff -s old.sql -t new.sql --ai | grep "建议" | execute.sh
```

✅ **正确做法**:
```bash
# 人工审查后决定
sql-diff -s old.sql -t new.sql --ai > review.txt
less review.txt
# 人工判断后再执行
```

### 2. 忽略成本

❌ **错误做法**:
```bash
# 在循环中无限制使用 AI
while true; do
  sql-diff -s old.sql -t new.sql --ai
done
```

✅ **正确做法**:
```bash
# 使用缓存,控制频率
sql-diff -s old.sql -t new.sql --ai
# 缓存结果,重复使用
```

### 3. 暴露敏感信息

❌ **错误做法**:
```bash
# API Key 硬编码
export SQL_DIFF_AI_API_KEY=sk-1234567890
git add .bashrc
```

✅ **正确做法**:
```bash
# 使用独立的密钥文件
echo "export SQL_DIFF_AI_API_KEY=sk-xxx" > ~/.secrets
chmod 600 ~/.secrets
# .gitignore 中忽略 .secrets
```

## 学习和改进

### 1. 收集反馈

记录 AI 建议的准确性:

```bash
# feedback.sh
sql-diff -s old.sql -t new.sql --ai > ai_suggestions.txt

# 执行后记录反馈
cat >> feedback.log << EOF
Date: $(date)
Suggestions: $(cat ai_suggestions.txt)
Outcome: [Helpful/Not Helpful/Partially Helpful]
Notes: [实际执行效果]
EOF
```

### 2. 持续优化提示词

根据反馈优化提示词:

```bash
# 初始提示词
export SQL_DIFF_AI_SYSTEM_PROMPT="分析表结构变更"

# 优化后
export SQL_DIFF_AI_SYSTEM_PROMPT="作为资深 DBA,分析表结构变更,重点关注:
1. 性能影响 (索引、锁表时间)
2. 数据安全 (备份、回滚)
3. 业务影响 (停机时间、兼容性)"
```

## 总结

关键要点:

1. ✅ **选择性使用**: 不是所有场景都需要 AI
2. ✅ **人工审查**: AI 建议需人工判断
3. ✅ **成本控制**: 使用缓存,避免重复调用
4. ✅ **安全第一**: 保护 API Key,脱敏敏感信息
5. ✅ **持续改进**: 收集反馈,优化使用方式

## 下一步

- [示例](/examples/advanced) - 查看实际使用案例
- [DeepSeek 集成](/ai/deepseek) - 深入了解 DeepSeek
- [CLI 工具](/guide/cli) - 掌握命令行选项
