# CI/CD 集成指南

SQL-Diff 可以轻松集成到 CI/CD 流程中，实现自动化的数据库 schema 变更检测和审查。

## 快速开始

### 核心价值

✅ **自动化审查**：每次 PR 自动检测数据库表结构变更  
✅ **风险提前发现**：AI 智能分析潜在风险，避免生产事故  
✅ **规范化流程**：强制 DBA Review，确保变更质量  
✅ **可追溯性**：SQL 文件纳入版本控制，变更历史一目了然

### 前置准备

在开始之前，你需要：

1. **准备 SQL 文件目录**：将表结构 SQL 文件纳入项目 Git 仓库
2. **获取 AI API Key**：注册 DeepSeek 账号并获取 API Key（可选，用于 AI 分析）
3. **配置 CI/CD 平台**：根据你使用的平台（GitHub/GitLab/Jenkins）进行配置

---

## 📂 第一步：组织项目结构

### 推荐的目录结构

```bash
your-project/
├── database/
│   └── schema/              # ⭐ 关键目录：存放表结构定义
│       ├── users.sql        # 用户表
│       ├── products.sql     # 商品表
│       ├── orders.sql       # 订单表
│       └── order_items.sql  # 订单明细表
├── .github/
│   └── workflows/
│       └── schema-check.yml # GitHub Actions 配置
├── .gitlab-ci.yml           # GitLab CI 配置
└── Jenkinsfile              # Jenkins Pipeline 配置
```

### 导出表结构 SQL

#### 方式 1：导出所有表（每个表一个文件，推荐）

```bash
# 创建目录
mkdir -p database/schema

# 导出每个表的结构（不含数据）
for table in $(mysql -u root -p your_database -sN -e "SHOW TABLES"); do
    mysqldump --no-data \
        --skip-add-drop-table \
        --skip-comments \
        -u root -p \
        your_database $table > database/schema/${table}.sql
    echo "✓ Exported $table"
done
```

#### 方式 2：导出整个数据库结构（单文件）

```bash
mysqldump --no-data \
    --skip-add-drop-table \
    --skip-comments \
    -u root -p \
    your_database > database/schema/all_tables.sql
```

#### 方式 3：创建自动导出脚本

```bash
# 创建导出脚本
cat > scripts/export_schema.sh << 'EOF'
#!/bin/bash
set -e

DB_HOST="${DB_HOST:-localhost}"
DB_USER="${DB_USER:-root}"
DB_NAME="${DB_NAME:-your_database}"
OUTPUT_DIR="database/schema"

mkdir -p "$OUTPUT_DIR"

echo "📤 Exporting schema from $DB_NAME..."

# 获取所有表
tables=$(mysql -h "$DB_HOST" -u "$DB_USER" -p "$DB_NAME" -sN -e "SHOW TABLES")

for table in $tables; do
    echo "  → Exporting $table..."
    mysqldump --no-data \
        --skip-add-drop-table \
        --skip-comments \
        -h "$DB_HOST" \
        -u "$DB_USER" \
        -p \
        "$DB_NAME" "$table" > "$OUTPUT_DIR/${table}.sql"
done

echo "✅ Schema exported to $OUTPUT_DIR"
EOF

chmod +x scripts/export_schema.sh
```

### 提交到 Git

```bash
# 添加 SQL 文件到版本控制
git add database/schema/
git commit -m "feat: 初始化数据库表结构定义"
git push origin main
```

**⚠️ 重要提示：**
- SQL 文件必须提交到 Git 仓库，作为项目代码的一部分
- 建议每个表一个文件，便于追踪具体表的变更历史
- 使用 `--no-data` 参数，只导出表结构，不含数据

---

## 💙 GitHub Actions 集成

### 第二步：配置 GitHub Actions

#### 2.1 创建工作流文件

创建 `.github/workflows/schema-check.yml`：

```yaml
name: Database Schema Check

on:
  pull_request:
    paths:
      - 'database/schema/**'  # 只监控 SQL 文件变更

jobs:
  schema-analysis:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write  # 允许评论 PR
    
    steps:
      - name: 检出代码
        uses: actions/checkout@v4
        with:
          fetch-depth: 2  # 需要获取前一个版本用于对比
      
      - name: 安装 SQL-Diff
        run: |
          # 使用 brew 安装（推荐）
          /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
          brew install Bacchusgift/tap/sql-diff
          
          # 或者使用 Go Install
          # go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
      
      - name: 分析 Schema 变更
        id: analyze
        run: |
          mkdir -p /tmp/reports
          
          # 获取变更的 SQL 文件
          CHANGED_FILES=$(git diff --name-only HEAD^ HEAD | grep -E '\.sql$' || true)
          
          if [ -z "$CHANGED_FILES" ]; then
            echo "has_changes=false" >> $GITHUB_OUTPUT
            echo "ℹ️ No SQL files changed"
            exit 0
          fi
          
          echo "has_changes=true" >> $GITHUB_OUTPUT
          echo "🔍 Found changed SQL files:"
          echo "$CHANGED_FILES"
          echo ""
          
          RISK_COUNT=0
          TOTAL_CHANGES=0
          
          # 分析每个变更的文件
          for file in $CHANGED_FILES; do
            table_name=$(basename $file .sql)
            echo "=================================="
            echo "📋 Analyzing: $table_name"
            echo "=================================="
            
            # 获取旧版本（如果存在）
            git show HEAD^:$file > /tmp/old_${table_name}.sql 2>/dev/null || echo "" > /tmp/old_${table_name}.sql
            
            # 运行 SQL-Diff 分析
            sql-diff \
              -s /tmp/old_${table_name}.sql \
              -t $file \
              --ai \
              --output /tmp/reports/${table_name}.txt
            
            # 显示分析结果
            cat /tmp/reports/${table_name}.txt
            echo ""
            
            # 统计变更和风险
            TOTAL_CHANGES=$((TOTAL_CHANGES + 1))
            
            # 检查是否有高风险关键词
            if grep -qi "DROP\|TRUNCATE\|DELETE" $file; then
              RISK_COUNT=$((RISK_COUNT + 1))
            fi
          done
          
          echo "total_changes=$TOTAL_CHANGES" >> $GITHUB_OUTPUT
          echo "risk_count=$RISK_COUNT" >> $GITHUB_OUTPUT
          
          echo ""
          echo "✅ Analysis completed!"
          echo "   Total changes: $TOTAL_CHANGES"
          echo "   High-risk operations: $RISK_COUNT"
        env:
          SQL_DIFF_AI_ENABLED: true
          SQL_DIFF_AI_PROVIDER: deepseek
          SQL_DIFF_AI_API_KEY: ${{ secrets.DEEPSEEK_API_KEY }}
          SQL_DIFF_AI_ENDPOINT: https://api.deepseek.com/v1
          SQL_DIFF_AI_MODEL: deepseek-chat
      
      - name: 生成 PR 评论
        if: steps.analyze.outputs.has_changes == 'true'
        run: |
          cat > /tmp/pr_comment.md << 'HEADER'
          ## 🔍 Database Schema Change Analysis
          
          **⚠️ This PR modifies database schema - Please review carefully!**
          
          HEADER
          
          # 添加每个表的分析结果
          for report in /tmp/reports/*.txt; do
            if [ -f "$report" ]; then
              table=$(basename $report .txt)
              echo "" >> /tmp/pr_comment.md
              echo "### 📋 Table: \`$table\`" >> /tmp/pr_comment.md
              echo "" >> /tmp/pr_comment.md
              echo '```' >> /tmp/pr_comment.md
              cat "$report" >> /tmp/pr_comment.md
              echo '```' >> /tmp/pr_comment.md
              echo "---" >> /tmp/pr_comment.md
            fi
          done
          
          cat >> /tmp/pr_comment.md << EOF
          
          ### 📊 Summary
          
          - **Total Changed Tables:** ${{ steps.analyze.outputs.total_changes }}
          - **High-Risk Operations:** ${{ steps.analyze.outputs.risk_count }}
          
          ### ✅ Review Checklist
          
          - [ ] Schema changes reviewed by DBA
          - [ ] Migration tested in staging environment
          - [ ] Rollback plan prepared
          - [ ] Performance impact assessed
          - [ ] Backward compatibility verified
          
          ---
          
          > 🤖 This analysis was automatically generated by [SQL-Diff](https://github.com/Bacchusgift/sql-diff) with AI assistance.
          EOF
      
      - name: 评论 PR
        if: steps.analyze.outputs.has_changes == 'true'
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const comment = fs.readFileSync('/tmp/pr_comment.md', 'utf8');
            
            // 查找已存在的评论
            const { data: comments } = await github.rest.issues.listComments({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
            });
            
            const botComment = comments.find(c => 
              c.user.type === 'Bot' && 
              c.body.includes('Database Schema Change Analysis')
            );
            
            if (botComment) {
              // 更新已存在的评论
              await github.rest.issues.updateComment({
                owner: context.repo.owner,
                repo: context.repo.repo,
                comment_id: botComment.id,
                body: comment
              });
              console.log('✅ Updated existing comment');
            } else {
              // 创建新评论
              await github.rest.issues.createComment({
                owner: context.repo.owner,
                repo: context.repo.repo,
                issue_number: context.issue.number,
                body: comment
              });
              console.log('✅ Created new comment');
            }
      
      - name: 检查风险阈值
        if: steps.analyze.outputs.has_changes == 'true'
        run: |
          RISK_COUNT=${{ steps.analyze.outputs.risk_count }}
          RISK_THRESHOLD=3  # 设置风险阈值
          
          if [ $RISK_COUNT -gt $RISK_THRESHOLD ]; then
            echo "❌ High-risk operations detected: $RISK_COUNT (threshold: $RISK_THRESHOLD)"
            echo "⚠️ Please get DBA approval before merging!"
            exit 1  # 失败，阻止合并
          else
            echo "✅ Risk count within acceptable threshold"
          fi
```

#### 2.2 配置 Secrets

在 GitHub 仓库设置中添加：

1. 进入 `Settings` → `Secrets and variables` → `Actions`
2. 点击 `New repository secret`
3. 添加：
   - **Name:** `DEEPSEEK_API_KEY`
   - **Value:** 你的 DeepSeek API 密钥

#### 2.3 设置分支保护规则

1. 进入 `Settings` → `Branches`
2. 点击 `Add rule`
3. 配置：
   - Branch name pattern: `main`
   - ☑️ Require status checks to pass before merging
   - 选择: `schema-analysis`
   - ☑️ Require review from Code Owners

### 第三步：使用流程

#### 3.1 开发者修改表结构

```bash
# 1. 创建功能分支
git checkout -b feat/add-user-email

# 2. 修改表结构
vim database/schema/users.sql
# 添加新字段：
# ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;

# 3. 提交变更
git add database/schema/users.sql
git commit -m "feat: add email field to users table"
git push origin feat/add-user-email
```

#### 3.2 创建 Pull Request

1. 在 GitHub 上创建 PR
2. GitHub Actions 自动触发，开始分析
3. 等待分析完成，查看评论中的分析结果

#### 3.3 Review 分析结果

GitHub Actions 会在 PR 中自动评论，包含：

- 📋 **生成的 DDL 语句**：自动生成的迁移 SQL
- 🤖 **AI 分析报告**：
  - 总结说明
  - 优化建议
  - 风险提示
- 📊 **变更统计**：表数量、风险等级

#### 3.4 DBA Review

1. DBA 查看 AI 分析报告
2. 评估风险和影响
3. 如果通过，批准 PR
4. 合并到 main 分支

---

## 🦊 GitLab CI 集成

### 配置 GitLab CI

创建 `.gitlab-ci.yml`：

```yaml
stages:
  - validate
  - review

variables:
  GOPATH: $CI_PROJECT_DIR/.go
  SQL_DIFF_AI_ENABLED: "true"
  SQL_DIFF_AI_PROVIDER: "deepseek"
  SQL_DIFF_AI_ENDPOINT: "https://api.deepseek.com/v1"
  SQL_DIFF_AI_MODEL: "deepseek-chat"

before_script:
  - mkdir -p .go
  - export PATH=$PATH:$GOPATH/bin

schema-validate:
  stage: validate
  image: golang:1.21
  only:
    changes:
      - database/schema/**
  script:
    - echo "🔍 Validating SQL schema files..."
    - |
      for file in database/schema/*.sql; do
        echo "Checking $file..."
        # 可以添加 SQL 语法验证
      done
  cache:
    paths:
      - .go/pkg/mod/

schema-review:
  stage: review
  image: golang:1.21
  only:
    changes:
      - database/schema/**
  script:
    - echo "📦 Installing SQL-Diff..."
    - go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
    
    - echo "🔍 Analyzing schema changes..."
    - mkdir -p reports
    
    - |
      # 获取变更的 SQL 文件
      CHANGED_FILES=$(git diff --name-only $CI_MERGE_REQUEST_DIFF_BASE_SHA $CI_COMMIT_SHA | grep -E '\.sql$' || true)
      
      if [ -z "$CHANGED_FILES" ]; then
        echo "ℹ️ No SQL files changed"
        exit 0
      fi
      
      echo "Found changes:"
      echo "$CHANGED_FILES"
      echo ""
      
      RISK_COUNT=0
      
      for file in $CHANGED_FILES; do
        table=$(basename $file .sql)
        echo "=================================="
        echo "📋 Analyzing: $table"
        echo "=================================="
        
        # 获取旧版本
        git show $CI_MERGE_REQUEST_DIFF_BASE_SHA:$file > old_${table}.sql 2>/dev/null || echo "" > old_${table}.sql
        
        # 运行分析
        sql-diff \
          -s old_${table}.sql \
          -t $file \
          --ai > reports/${table}_analysis.txt
        
        # 显示结果
        cat reports/${table}_analysis.txt
        echo ""
        
        # 统计风险
        if grep -qi "DROP\|TRUNCATE" $file; then
          RISK_COUNT=$((RISK_COUNT + 1))
        fi
      done
      
      echo ""
      echo "✅ Analysis completed!"
      echo "   High-risk operations: $RISK_COUNT"
      
      # 检查风险阈值
      if [ $RISK_COUNT -gt 3 ]; then
        echo "❌ Too many high-risk operations!"
        exit 1
      fi
  
  variables:
    SQL_DIFF_AI_API_KEY: $DEEPSEEK_API_KEY
  
  artifacts:
    paths:
      - reports/
    expire_in: 1 week
  
  allow_failure: false  # 失败则阻止合并
```

### 配置 GitLab CI/CD 变量

1. 进入项目 `Settings` → `CI/CD` → `Variables`
2. 添加变量：
   - **Key:** `DEEPSEEK_API_KEY`
   - **Value:** 你的 DeepSeek API 密钥
   - **Protected:** ✅
   - **Masked:** ✅

### 使用流程

#### 1. 创建 Merge Request

```bash
git checkout -b feat/add-product-category
vim database/schema/products.sql
git add database/schema/products.sql
git commit -m "feat: add category_id to products"
git push origin feat/add-product-category
```

#### 2. 查看 Pipeline

1. 在 GitLab MR 页面查看 Pipeline 状态
2. 点击 `schema-review` job 查看详细分析结果
3. 下载 artifacts 中的分析报告

#### 3. DBA 审查

1. DBA 查看 Pipeline 输出和 artifacts
2. 评估变更影响
3. 批准或拒绝 MR

---

## 🔧 Jenkins Pipeline 集成

### 配置 Jenkinsfile

创建 `Jenkinsfile`：

```groovy
pipeline {
    agent any
    
    environment {
        GOPATH = "${WORKSPACE}/.go"
        PATH = "${PATH}:${GOPATH}/bin"
        SQL_DIFF_AI_ENABLED = "true"
        SQL_DIFF_AI_PROVIDER = "deepseek"
        SQL_DIFF_AI_ENDPOINT = "https://api.deepseek.com/v1"
        SQL_DIFF_AI_MODEL = "deepseek-chat"
        SQL_DIFF_AI_API_KEY = credentials('deepseek-api-key')
    }
    
    triggers {
        // 监听 Pull Request
        issueCommentTrigger('.*test this please.*')
    }
    
    stages {
        stage('Setup') {
            steps {
                echo '📦 Installing SQL-Diff...'
                sh '''
                    mkdir -p ${GOPATH}
                    go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
                    sql-diff --version
                '''
            }
        }
        
        stage('Detect Changes') {
            steps {
                script {
                    echo '🔍 Detecting SQL file changes...'
                    
                    // 获取变更的 SQL 文件
                    def changedFiles = sh(
                        script: """
                            git diff --name-only HEAD~1 HEAD | grep -E '\\.sql\$' || true
                        """,
                        returnStdout: true
                    ).trim()
                    
                    if (!changedFiles) {
                        echo 'ℹ️ No SQL files changed, skipping analysis'
                        currentBuild.result = 'SUCCESS'
                        return
                    }
                    
                    env.CHANGED_SQL_FILES = changedFiles
                    echo "Found changes:\n${changedFiles}"
                }
            }
        }
        
        stage('Schema Analysis') {
            when {
                expression { env.CHANGED_SQL_FILES }
            }
            steps {
                script {
                    echo '🔍 Analyzing schema changes...'
                    
                    def files = env.CHANGED_SQL_FILES.split('\n')
                    def riskCount = 0
                    
                    for (file in files) {
                        def tableName = file.tokenize('/').last().replaceAll('\\.sql$', '')
                        
                        echo "=================================="
                        echo "📋 Analyzing: ${tableName}"
                        echo "=================================="
                        
                        sh """
                            # 获取旧版本
                            git show HEAD~1:${file} > old_${tableName}.sql 2>/dev/null || echo "" > old_${tableName}.sql
                            
                            # 运行分析
                            sql-diff \
                                -s old_${tableName}.sql \
                                -t ${file} \
                                --ai \
                                --output ${tableName}_analysis.txt
                            
                            # 显示结果
                            cat ${tableName}_analysis.txt
                        """
                        
                        // 检查风险
                        def fileContent = readFile(file)
                        if (fileContent =~ /(?i)(DROP|TRUNCATE|DELETE)/) {
                            riskCount++
                        }
                    }
                    
                    env.RISK_COUNT = riskCount.toString()
                    
                    echo ""
                    echo "✅ Analysis completed!"
                    echo "   High-risk operations: ${riskCount}"
                }
            }
        }
        
        stage('Risk Assessment') {
            when {
                expression { env.CHANGED_SQL_FILES }
            }
            steps {
                script {
                    def riskThreshold = 3
                    def riskCount = env.RISK_COUNT as Integer
                    
                    if (riskCount > riskThreshold) {
                        echo "❌ High-risk operations detected: ${riskCount} (threshold: ${riskThreshold})"
                        echo "⚠️ Please get DBA approval before merging!"
                        
                        // 发送邮件通知
                        emailext (
                            subject: "⚠️ High-Risk Database Schema Change Detected",
                            body: """
                                <h2>Database Schema Change Alert</h2>
                                <p>Build: ${env.BUILD_URL}</p>
                                <p>High-risk operations detected: ${riskCount}</p>
                                <p>Please review the schema changes carefully.</p>
                            """,
                            to: 'dba-team@example.com',
                            mimeType: 'text/html'
                        )
                        
                        // 设置为不稳定，需要手动批准
                        input message: 'High-risk changes detected. DBA approval required.', ok: 'Approve'
                    } else {
                        echo "✅ Risk count within acceptable threshold"
                    }
                }
            }
        }
    }
    
    post {
        always {
            echo '📦 Archiving analysis reports...'
            archiveArtifacts artifacts: '*_analysis.txt', allowEmptyArchive: true
        }
        success {
            echo '✅ Schema analysis completed successfully!'
        }
        failure {
            echo '❌ Schema analysis failed!'
            emailext (
                subject: "❌ Schema Analysis Failed: ${env.JOB_NAME} - ${env.BUILD_NUMBER}",
                body: "Check console output at ${env.BUILD_URL}",
                to: 'dev-team@example.com'
            )
        }
    }
}
```

### 配置 Jenkins 凭据

1. 进入 `Manage Jenkins` → `Manage Credentials`
2. 添加凭据：
   - **Kind:** Secret text
   - **Secret:** 你的 DeepSeek API 密钥
   - **ID:** `deepseek-api-key`

### 配置邮件通知

在 `Manage Jenkins` → `Configure System` 中配置 SMTP 服务器。

### 使用流程

#### 1. 提交变更触发构建

```bash
git add database/schema/orders.sql
git commit -m "feat: add status_updated_at to orders"
git push origin develop
```

#### 2. 查看 Jenkins 构建

1. 访问 Jenkins 查看构建进度
2. 查看控制台输出的分析结果
3. 下载归档的分析报告文件

#### 3. 高风险审批

如果检测到高风险操作：
1. Jenkins 会暂停并等待审批
2. DBA 收到邮件通知
3. DBA 点击 "Approve" 继续，或 "Abort" 中止

---

## 📊 使用场景对比

| 平台 | 适用场景 | 优势 | 注意事项 |
|------|---------|------|---------|
| **GitHub Actions** | 开源项目、小型团队 | 配置简单、与 GitHub 深度集成 | 需要配置 Secrets |
| **GitLab CI** | 企业内部项目 | 完整的 DevOps 平台、私有部署 | 需要配置 Runner |
| **Jenkins** | 大型企业、复杂流程 | 高度可定制、插件丰富 | 配置复杂、需要维护 |

---

## 🎯 最佳实践

### 1. SQL 文件管理

```bash
# ✅ 推荐：每个表一个文件
database/schema/
├── users.sql
├── products.sql
└── orders.sql

# ❌ 不推荐：所有表在一个文件
database/schema/
└── all_tables.sql
```

### 2. 分支策略

```bash
# 功能开发
feature/* → develop → staging → main

# Schema 变更
schema/* → develop (自动检查) → staging (测试) → main (生产)
```

### 3. 风险等级定义

| 等级 | 操作类型 | 处理方式 |
|------|---------|---------|
| 🟢 **低风险** | ADD COLUMN | 自动通过 |
| 🟡 **中风险** | MODIFY COLUMN | DBA Review |
| 🔴 **高风险** | DROP/TRUNCATE | DBA 批准 + 备份 |

### 4. 审查清单

在 PR 模板中添加：

```markdown
## Database Schema Changes

- [ ] Schema changes reviewed by DBA
- [ ] AI analysis risks addressed
- [ ] Migration tested in staging
- [ ] Rollback plan prepared
- [ ] Performance impact assessed
- [ ] Data backup verified
- [ ] Backward compatibility confirmed
```

---

## 🔔 告警和通知

### Slack 通知（GitHub Actions）

```yaml
- name: Notify Slack
  if: failure()
  uses: slackapi/slack-github-action@v1
  with:
    payload: |
      {
        "text": "⚠️ Schema migration failed",
        "blocks": [
          {
            "type": "section",
            "text": {
              "type": "mrkdwn",
              "text": "*Schema Migration Alert*\n\nFailed in PR #${{ github.event.pull_request.number }}"
            }
          }
        ]
      }
  env:
    SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

### 企业微信通知（GitLab CI）

```yaml
after_script:
  - |
    if [ $CI_JOB_STATUS == 'failed' ]; then
      curl -X POST "$WECHAT_WEBHOOK" \
        -H 'Content-Type: application/json' \
        -d "{
          \"msgtype\": \"text\",
          \"text\": {
            \"content\": \"⚠️ Schema 检查失败\nProject: $CI_PROJECT_NAME\nBranch: $CI_COMMIT_REF_NAME\"
          }
        }"
    fi
```

---

## 📚 相关文档

- [基础示例](/examples/basic) - 学习基础用法
- [高级示例](/examples/advanced) - 复杂场景示例
- [AI 最佳实践](/ai/best-practices) - 优化 AI 使用
- [配置管理](/config/environment) - 环境变量配置

---

## ❓ 常见问题

### Q1: 如何处理首次导入的表？

```bash
# 创建空的旧版本文件
echo "" > old.sql

# 运行分析
sql-diff -s old.sql -t database/schema/new_table.sql --ai
```

### Q2: 如何跳过 CI 检查？

在 commit message 中添加：
```bash
git commit -m "docs: update README [skip ci]"
```

### Q3: 如何调整风险阈值？

修改工作流中的 `RISK_THRESHOLD` 变量：
```yaml
- name: 检查风险阈值
  run: |
    RISK_THRESHOLD=5  # 根据团队需求调整
```

### Q4: 如何禁用 AI 分析？

移除环境变量或设置为 false：
```yaml
env:
  SQL_DIFF_AI_ENABLED: false
```
