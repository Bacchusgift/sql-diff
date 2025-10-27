# 复杂场景示例

本文档展示 SQL-Diff 在复杂场景下的实际应用。

## 大规模表结构重构

### 场景描述

对一个生产环境的用户表进行重大重构:
- 添加多个新字段
- 修改多个字段定义
- 重构索引策略
- 添加全文搜索

### 源表结构

```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  username VARCHAR(50),
  password VARCHAR(32),
  email VARCHAR(100),
  status TINYINT(1),
  register_time INT,
  last_login INT,
  KEY idx_username (username),
  KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
```

**存在的问题**:
- 主键没有自增
- 密码字段太短 (MD5 已不安全)
- 使用整数存储时间戳
- 字符集使用 latin1
- 缺少必要的业务字段

### 目标表结构

```sql
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(100) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  email_verified BOOLEAN DEFAULT FALSE,
  phone VARCHAR(20),
  phone_verified BOOLEAN DEFAULT FALSE,
  status ENUM('active', 'inactive', 'suspended', 'deleted') DEFAULT 'active',
  nickname VARCHAR(100),
  avatar_url VARCHAR(500),
  bio TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  last_login_at TIMESTAMP NULL,
  login_count INT DEFAULT 0,
  UNIQUE KEY uk_username (username),
  UNIQUE KEY uk_email (email),
  KEY idx_phone (phone),
  KEY idx_status (status),
  KEY idx_created (created_at),
  FULLTEXT KEY ft_search (username, nickname, bio)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 执行比对

```bash
sql-diff \
  -s "$(cat old_users.sql)" \
  -t "$(cat new_users.sql)" \
  --ai \
  --verbose > users_migration_analysis.txt
```

### 生成的 DDL

```sql
-- 修改现有列
ALTER TABLE users MODIFY COLUMN id BIGINT PRIMARY KEY AUTO_INCREMENT;
ALTER TABLE users MODIFY COLUMN username VARCHAR(100) NOT NULL;
ALTER TABLE users MODIFY COLUMN password_hash VARCHAR(255) NOT NULL;
ALTER TABLE users MODIFY COLUMN email VARCHAR(255) NOT NULL;
ALTER TABLE users MODIFY COLUMN status ENUM('active', 'inactive', 'suspended', 'deleted') DEFAULT 'active';

-- 新增列
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users ADD COLUMN phone_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN nickname VARCHAR(100);
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(500);
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE users ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
ALTER TABLE users ADD COLUMN last_login_at TIMESTAMP NULL;
ALTER TABLE users ADD COLUMN login_count INT DEFAULT 0;

-- 删除旧列
ALTER TABLE users DROP COLUMN register_time;
ALTER TABLE users DROP COLUMN last_login;

-- 索引变更
ALTER TABLE users DROP INDEX idx_username;
ALTER TABLE users DROP INDEX idx_email;
ALTER TABLE users ADD UNIQUE KEY uk_username (username);
ALTER TABLE users ADD UNIQUE KEY uk_email (email);
ALTER TABLE users ADD KEY idx_phone (phone);
ALTER TABLE users ADD KEY idx_status (status);
ALTER TABLE users ADD KEY idx_created (created_at);
ALTER TABLE users ADD FULLTEXT KEY ft_search (username, nickname, bio);

-- 表选项变更
ALTER TABLE users CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### AI 分析摘要

```
🤖 AI 分析结果:

📊 差异分析:
这是一次大规模的表结构重构,包含多个重要变更:
1. 数据类型升级: INT → BIGINT (支持更大的用户量)
2. 安全增强: 密码字段从 32 字符扩展到 255 字符 (支持现代哈希算法)
3. 时间字段标准化: 从 UNIX 时间戳改为 TIMESTAMP
4. 字符集升级: latin1 → utf8mb4 (支持 emoji 等)
5. 业务功能扩展: 添加手机号、头像、个人简介等字段
6. 索引策略优化: 添加唯一索引和全文索引

💡 优化建议:
- 分步执行此迁移,避免长时间锁表
- 考虑使用 pt-online-schema-change 工具
- 字符集转换前备份数据
- 新增字段都有默认值,向后兼容性好
- 建议添加数据迁移脚本处理 register_time → created_at 的转换

⚠️  潜在风险:
- 字符集转换可能导致锁表,大表需谨慎
- password 字段重命名为 password_hash,需要同步更新应用代码
- 删除 register_time 和 last_login,需确保应用不再使用
- INT → BIGINT 转换在大表上可能耗时较长
- 全文索引创建会消耗大量资源

✅ 最佳实践:
- 在维护窗口期执行
- 使用主从切换策略零停机部署
- 准备回滚方案
- 监控执行过程和性能指标
- 分阶段上线新功能
```

### 实际执行计划

```bash
#!/bin/bash
# production_migration.sh

set -e

echo "=== Production Migration: Users Table Refactoring ==="
echo "Start time: $(date)"

# 配置
DB_HOST="prod-db.example.com"
DB_USER="admin"
DB_NAME="myapp"
BACKUP_DIR="backups/$(date +%Y%m%d_%H%M%S)"

# 1. 创建备份
echo "Step 1: Creating backup..."
mkdir -p $BACKUP_DIR
mysqldump -h $DB_HOST -u $DB_USER -p $DB_NAME users > $BACKUP_DIR/users_backup.sql
echo "✓ Backup saved to $BACKUP_DIR/users_backup.sql"

# 2. 创建数据迁移脚本
cat > $BACKUP_DIR/data_migration.sql << 'EOF'
-- 将 register_time 转换为 created_at
UPDATE users 
SET created_at = FROM_UNIXTIME(register_time)
WHERE register_time > 0;

-- 将 last_login 转换为 last_login_at
UPDATE users 
SET last_login_at = FROM_UNIXTIME(last_login)
WHERE last_login > 0;
EOF

# 3. 使用 pt-osc 执行在线变更 (分步执行)
echo "Step 2: Applying schema changes with pt-online-schema-change..."

# 3.1 修改列类型
pt-online-schema-change \
  --alter "MODIFY COLUMN id BIGINT PRIMARY KEY AUTO_INCREMENT" \
  D=$DB_NAME,t=users \
  --execute \
  --chunk-size=1000 \
  --max-load=Threads_running=30

# 3.2 添加新列
pt-online-schema-change \
  --alter "
    ADD COLUMN email_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN phone VARCHAR(20),
    ADD COLUMN phone_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN nickname VARCHAR(100),
    ADD COLUMN avatar_url VARCHAR(500),
    ADD COLUMN bio TEXT,
    ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    ADD COLUMN last_login_at TIMESTAMP NULL,
    ADD COLUMN login_count INT DEFAULT 0
  " \
  D=$DB_NAME,t=users \
  --execute

# 3.3 数据迁移
echo "Step 3: Migrating data..."
mysql -h $DB_HOST -u $DB_USER -p $DB_NAME < $BACKUP_DIR/data_migration.sql

# 3.4 删除旧列
pt-online-schema-change \
  --alter "DROP COLUMN register_time, DROP COLUMN last_login" \
  D=$DB_NAME,t=users \
  --execute

# 3.5 索引优化
pt-online-schema-change \
  --alter "
    DROP INDEX idx_username,
    DROP INDEX idx_email,
    ADD UNIQUE KEY uk_username (username),
    ADD UNIQUE KEY uk_email (email),
    ADD KEY idx_phone (phone),
    ADD KEY idx_status (status),
    ADD KEY idx_created (created_at)
  " \
  D=$DB_NAME,t=users \
  --execute

# 3.6 添加全文索引 (单独执行,资源消耗大)
echo "Step 4: Creating fulltext index (may take a while)..."
mysql -h $DB_HOST -u $DB_USER -p $DB_NAME -e \
  "ALTER TABLE users ADD FULLTEXT KEY ft_search (username, nickname, bio);"

# 3.7 字符集转换
echo "Step 5: Converting charset..."
pt-online-schema-change \
  --alter "CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" \
  D=$DB_NAME,t=users \
  --execute

echo "✅ Migration completed successfully!"
echo "End time: $(date)"
echo ""
echo "Next steps:"
echo "1. Verify data integrity"
echo "2. Update application code"
echo "3. Deploy new application version"
echo "4. Monitor for 24 hours"
```

## 多表联合重构

### 场景描述

重构订单系统,拆分订单表为订单主表和订单明细表。

### 原始结构

```sql
CREATE TABLE orders (
  order_id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  product_id INT NOT NULL,
  product_name VARCHAR(200),
  quantity INT,
  price DECIMAL(10,2),
  total_amount DECIMAL(10,2),
  status VARCHAR(20),
  created_at DATETIME,
  KEY idx_user (user_id),
  KEY idx_product (product_id)
) ENGINE=InnoDB;
```

### 目标结构

**订单主表**:
```sql
CREATE TABLE orders (
  order_id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(32) NOT NULL,
  user_id BIGINT NOT NULL,
  total_amount DECIMAL(12,4) NOT NULL,
  status ENUM('pending', 'paid', 'shipped', 'completed', 'cancelled') DEFAULT 'pending',
  payment_method VARCHAR(50),
  shipping_address TEXT,
  remark TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  paid_at TIMESTAMP NULL,
  shipped_at TIMESTAMP NULL,
  completed_at TIMESTAMP NULL,
  UNIQUE KEY uk_order_no (order_no),
  KEY idx_user (user_id),
  KEY idx_status (status),
  KEY idx_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**订单明细表**:
```sql
CREATE TABLE order_items (
  item_id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_id BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  product_name VARCHAR(255) NOT NULL,
  product_sku VARCHAR(100),
  quantity INT NOT NULL,
  price DECIMAL(12,4) NOT NULL,
  discount DECIMAL(12,4) DEFAULT 0,
  subtotal DECIMAL(12,4) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY idx_order (order_id),
  KEY idx_product (product_id),
  FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 迁移脚本

```sql
-- 1. 重命名旧表
RENAME TABLE orders TO orders_old;

-- 2. 创建新表
CREATE TABLE orders (
  order_id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(32) NOT NULL,
  user_id BIGINT NOT NULL,
  total_amount DECIMAL(12,4) NOT NULL,
  status ENUM('pending', 'paid', 'shipped', 'completed', 'cancelled') DEFAULT 'pending',
  payment_method VARCHAR(50),
  shipping_address TEXT,
  remark TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  paid_at TIMESTAMP NULL,
  shipped_at TIMESTAMP NULL,
  completed_at TIMESTAMP NULL,
  UNIQUE KEY uk_order_no (order_no),
  KEY idx_user (user_id),
  KEY idx_status (status),
  KEY idx_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE order_items (
  item_id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_id BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  product_name VARCHAR(255) NOT NULL,
  product_sku VARCHAR(100),
  quantity INT NOT NULL,
  price DECIMAL(12,4) NOT NULL,
  discount DECIMAL(12,4) DEFAULT 0,
  subtotal DECIMAL(12,4) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY idx_order (order_id),
  KEY idx_product (product_id),
  FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 迁移数据
INSERT INTO orders (
  order_id, order_no, user_id, total_amount, status, created_at
)
SELECT 
  order_id,
  CONCAT('ORD', LPAD(order_id, 10, '0')) as order_no,
  user_id,
  total_amount,
  CASE status
    WHEN 'pending' THEN 'pending'
    WHEN 'paid' THEN 'paid'
    WHEN 'shipped' THEN 'shipped'
    WHEN 'completed' THEN 'completed'
    ELSE 'cancelled'
  END as status,
  created_at
FROM orders_old
GROUP BY order_id;

INSERT INTO order_items (
  order_id, product_id, product_name, quantity, price, subtotal, created_at
)
SELECT
  order_id,
  product_id,
  product_name,
  quantity,
  price,
  total_amount,
  created_at
FROM orders_old;

-- 4. 验证数据
SELECT 
  (SELECT COUNT(*) FROM orders) as new_orders_count,
  (SELECT COUNT(*) FROM order_items) as new_items_count,
  (SELECT COUNT(DISTINCT order_id) FROM orders_old) as old_orders_count,
  (SELECT COUNT(*) FROM orders_old) as old_items_count;

-- 5. 确认无误后删除旧表
-- DROP TABLE orders_old;
```

## 性能优化场景

### 场景描述

优化一个访问频繁的文章表,添加合理的索引和分区。

### 源表

```sql
CREATE TABLE articles (
  id INT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(200),
  content LONGTEXT,
  author_id INT,
  category_id INT,
  tags VARCHAR(500),
  view_count INT DEFAULT 0,
  like_count INT DEFAULT 0,
  created_at DATETIME,
  updated_at DATETIME
) ENGINE=InnoDB;
```

### 目标表

```sql
CREATE TABLE articles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  content LONGTEXT NOT NULL,
  content_hash CHAR(64),
  author_id BIGINT NOT NULL,
  category_id INT NOT NULL,
  tags JSON,
  view_count INT UNSIGNED DEFAULT 0,
  like_count INT UNSIGNED DEFAULT 0,
  comment_count INT UNSIGNED DEFAULT 0,
  status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
  published_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_author (author_id, status, published_at),
  KEY idx_category (category_id, status, published_at),
  KEY idx_status_published (status, published_at),
  KEY idx_hot (view_count, like_count),
  FULLTEXT KEY ft_title_content (title, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
PARTITION BY RANGE (YEAR(created_at)) (
  PARTITION p2023 VALUES LESS THAN (2024),
  PARTITION p2024 VALUES LESS THAN (2025),
  PARTITION p2025 VALUES LESS THAN (2026),
  PARTITION p_future VALUES LESS THAN MAXVALUE
);
```

### 性能对比测试

```sql
-- 测试查询性能

-- 查询 1: 按作者查询已发布文章
-- 优化前
SELECT * FROM articles WHERE author_id = 123 AND status = 'published' ORDER BY created_at DESC LIMIT 10;
-- 全表扫描,耗时: 2.5s

-- 优化后
SELECT * FROM articles WHERE author_id = 123 AND status = 'published' ORDER BY published_at DESC LIMIT 10;
-- 使用索引 idx_author,耗时: 0.05s

-- 查询 2: 全文搜索
-- 优化前
SELECT * FROM articles WHERE title LIKE '%关键词%' OR content LIKE '%关键词%';
-- 全表扫描,耗时: 15s

-- 优化后
SELECT * FROM articles WHERE MATCH(title, content) AGAINST('关键词' IN NATURAL LANGUAGE MODE);
-- 使用全文索引,耗时: 0.2s

-- 查询 3: 热门文章
-- 优化后新增
SELECT * FROM articles WHERE status = 'published' ORDER BY view_count DESC, like_count DESC LIMIT 20;
-- 使用索引 idx_hot,耗时: 0.1s
```

## CI/CD 集成场景

### GitHub Actions 工作流

```yaml
# .github/workflows/schema-migration.yml
name: Database Schema Migration

on:
  pull_request:
    paths:
      - 'database/schema/**'
  push:
    branches:
      - main
    paths:
      - 'database/schema/**'

jobs:
  schema-review:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
        with:
          fetch-depth: 2
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install SQL-Diff
        run: go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
      
      - name: Get changed files
        id: changed-files
        run: |
          git diff --name-only HEAD^ HEAD > changed_files.txt
          cat changed_files.txt
      
      - name: Analyze schema changes
        id: analyze
        run: |
          mkdir -p migration_reports
          
          for file in $(cat changed_files.txt | grep '\.sql$'); do
            table=$(basename $file .sql)
            echo "Analyzing $table..."
            
            # 获取旧版本
            git show HEAD^:$file > old_${table}.sql 2>/dev/null || echo "" > old_${table}.sql
            
            # 比对并生成报告
            sql-diff \
              -s old_${table}.sql \
              -t $file \
              --ai \
              --format json > migration_reports/${table}_analysis.json
          done
        env:
          SQL_DIFF_AI_API_KEY: ${{ secrets.DEEPSEEK_API_KEY }}
      
      - name: Check for risks
        id: check-risks
        run: |
          TOTAL_RISKS=0
          
          for report in migration_reports/*.json; do
            RISKS=$(jq '.ai_analysis.risks | length' $report)
            TOTAL_RISKS=$((TOTAL_RISKS + RISKS))
            
            if [ $RISKS -gt 0 ]; then
              echo "⚠️  Found $RISKS risks in $(basename $report)"
              jq '.ai_analysis.risks' $report
            fi
          done
          
          echo "total_risks=$TOTAL_RISKS" >> $GITHUB_OUTPUT
          
          if [ $TOTAL_RISKS -gt 5 ]; then
            echo "❌ Too many risks detected! Manual review required."
            exit 1
          fi
      
      - name: Generate migration report
        run: |
          cat > migration_summary.md << 'EOF'
          # 📊 Schema Migration Summary
          
          ## Changed Tables
          EOF
          
          for report in migration_reports/*.json; do
            table=$(basename $report _analysis.json)
            cat >> migration_summary.md << EOF
          
          ### $table
          
          **DDL Statements:**
          \`\`\`sql
          $(jq -r '.ddl_statements[]' $report)
          \`\`\`
          
          **AI Analysis:**
          $(jq -r '.ai_analysis.summary' $report)
          
          **Risks:**
          $(jq -r '.ai_analysis.risks[]' $report)
          
          ---
          EOF
          done
      
      - name: Comment PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const summary = fs.readFileSync('migration_summary.md', 'utf8');
            
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: summary
            });
      
      - name: Upload reports
        uses: actions/upload-artifact@v3
        with:
          name: migration-reports
          path: migration_reports/
  
  apply-to-staging:
    needs: schema-review
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Apply to staging database
        run: |
          # 这里添加实际的数据库迁移逻辑
          echo "Applying migrations to staging..."
      
      - name: Run integration tests
        run: |
          # 运行集成测试验证迁移
          echo "Running tests..."
```

## 下一步

- [CI/CD 集成](/examples/ci-cd) - 完整的 CI/CD 集成示例
- [AI 最佳实践](/ai/best-practices) - AI 使用技巧
- [CLI 工具](/guide/cli) - 命令行详细用法
