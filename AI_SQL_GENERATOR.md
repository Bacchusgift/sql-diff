# AI 自然语言生成 SQL 功能

这是 sql-diff 的新功能，使用 AI 根据自然语言描述生成 SQL 语句。

## 功能概览

### 1. 生成 CREATE TABLE 语句

根据自然语言描述，自动生成标准的 MySQL CREATE TABLE 语句。

**命令：** `sql-diff generate`

**示例：**

```bash
# 生成用户表
sql-diff generate -d "创建一个用户表，包含 ID、用户名、邮箱、密码、创建时间"

# 生成商品表并保存到文件
sql-diff generate -d "创建商品表：商品ID、名称、价格、库存、分类、状态" -o product.sql
```

**特点：**
- ✅ 自动推断字段类型（如：金额用 DECIMAL，日期用 DATETIME）
- ✅ 自动添加主键和索引
- ✅ 遵循 MySQL 最佳实践（InnoDB、UTF8MB4 等）
- ✅ 使用蛇形命名法（snake_case）
- ✅ 添加合理的注释

### 2. 生成 ALTER TABLE 语句

根据现有表结构和修改需求，生成相应的 ALTER TABLE 语句。

**命令：** `sql-diff alter`

**示例：**

```bash
# 命令行模式
sql-diff alter -t "CREATE TABLE users (id INT, name VARCHAR(100));" -d "添加手机号字段、邮箱改为唯一索引"

# 交互式模式（推荐用于复杂表结构）
sql-diff alter -i -d "添加商品状态字段，默认值为上架"
```

**特点：**
- ✅ 基于现有表结构生成变更
- ✅ 考虑数据迁移安全性
- ✅ 支持交互式粘贴表结构
- ✅ 可生成多条 ALTER 语句

## 前置要求

此功能依赖 AI，需要先配置 AI 服务：

### 方法一：使用配置文件（推荐）

创建 `.sql-diff-config.yaml`：

```yaml
ai:
  enabled: true
  provider: deepseek  # 或 openai
  api_key: "your-api-key-here"
  api_endpoint: "https://api.deepseek.com/v1"  # DeepSeek
  # api_endpoint: "https://api.openai.com/v1"  # OpenAI
  model: "deepseek-chat"
  timeout: 30
```

### 方法二：使用命令行参数

```bash
# 使用 --ai 参数临时启用（需要配置文件中有 API Key）
sql-diff generate --ai -d "创建用户表"
```

### 配置 AI 服务

使用交互式配置工具：

```bash
sql-diff config
```

## 使用场景

### 场景 1：快速原型设计

```bash
# 快速生成表结构原型
sql-diff generate -d "创建博客文章表：文章ID、标题、内容、作者、分类、标签、发布时间、浏览量"
```

### 场景 2：表结构迭代

```bash
# 第一步：生成初始表结构
sql-diff generate -d "用户表：ID、用户名、邮箱" -o user_v1.sql

# 第二步：根据需求添加字段
sql-diff alter -t "$(cat user_v1.sql)" -d "添加手机号、实名状态、注册时间" -o user_v2.sql
```

### 场景 3：需求文档转 SQL

将产品需求直接转换为数据库设计：

```bash
sql-diff generate -d "订单表：订单号（唯一）、用户ID（外键）、订单金额（精确到分）、订单状态（待支付/已支付/已取消）、支付方式、下单时间、支付时间"
```

## 完整示例

### 示例 1：生成用户表

```bash
$ sql-diff generate -d "创建用户表，包含：用户ID（自增主键）、用户名（唯一）、手机号（唯一）、邮箱、密码（加密存储）、头像URL、用户状态（正常/冻结）、创建时间、更新时间"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       AI 生成 CREATE TABLE 语句
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📝 需求描述: 创建用户表，包含：用户ID（自增主键）、用户名（唯一）...

🤖 正在使用 AI 生成 SQL...

✓ 生成成功！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 生成的 CREATE TABLE 语句:

CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
  username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
  mobile VARCHAR(20) NOT NULL UNIQUE COMMENT '手机号',
  email VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
  password VARCHAR(255) NOT NULL COMMENT '加密密码',
  avatar_url VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
  status ENUM('normal', 'frozen') DEFAULT 'normal' COMMENT '用户状态',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX idx_mobile (mobile),
  INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 示例 2：修改表结构

```bash
$ sql-diff alter -i -d "添加实名认证字段：真实姓名、身份证号、认证状态、认证时间"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       AI 生成 ALTER TABLE 语句
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 请粘贴现有表的 CREATE TABLE 语句：
（粘贴完成后输入 'END' 或连续按两次 Enter）

[粘贴上面的 users 表结构]
END

✓ 已读取 XXX 个字符

📝 修改需求: 添加实名认证字段：真实姓名、身份证号、认证状态、认证时间

🤖 正在使用 AI 生成 SQL...

✓ 生成成功！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 生成的 ALTER TABLE 语句:

ALTER TABLE users ADD COLUMN real_name VARCHAR(50) DEFAULT NULL COMMENT '真实姓名';
ALTER TABLE users ADD COLUMN id_card VARCHAR(18) DEFAULT NULL COMMENT '身份证号';
ALTER TABLE users ADD COLUMN verified_status ENUM('unverified', 'pending', 'verified') DEFAULT 'unverified' COMMENT '认证状态';
ALTER TABLE users ADD COLUMN verified_at DATETIME DEFAULT NULL COMMENT '认证时间';
ALTER TABLE users ADD INDEX idx_id_card (id_card);
```

## 优势

1. **节省时间**：无需手动编写冗长的 DDL 语句
2. **降低出错**：AI 自动处理字段类型、长度、索引等细节
3. **最佳实践**：自动应用数据库设计最佳实践
4. **快速迭代**：需求变更时快速生成新的表结构
5. **学习工具**：通过 AI 生成的 SQL 学习最佳实践

## 注意事项

1. **仔细审查**：AI 生成的 SQL 仅供参考，请根据实际需求调整
2. **安全第一**：生产环境使用前，请在测试环境验证
3. **API 费用**：使用 AI 功能会调用第三方 API，可能产生费用
4. **网络要求**：需要能够访问 AI 服务提供商的 API

## 支持的 AI 提供商

- **DeepSeek**（推荐，性价比高）
- **OpenAI**（ChatGPT）

配置方法见上文"前置要求"章节。

## 技术实现

- 使用 DeepSeek/OpenAI Chat API
- 精心设计的 Prompt 工程
- 智能 SQL 提取和清理
- 与现有比对功能无缝集成

## 反馈与建议

如果你在使用过程中遇到问题或有改进建议，欢迎提 Issue！
