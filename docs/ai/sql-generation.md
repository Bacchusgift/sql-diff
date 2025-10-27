# AI 生成 SQL 语句

SQL-Diff 提供了强大的 AI 生成 SQL 功能，让你可以用自然语言描述需求，自动生成标准的 SQL 建表和变更语句。

## 🌟 核心功能

### 两种生成模式

1. **生成 CREATE TABLE** - 根据自然语言描述创建完整的建表语句
2. **生成 ALTER TABLE** - 基于现有表结构生成变更语句

## 🚀 快速开始

### 前提条件

使用 AI 生成功能前，需要先启用 AI：

```bash
# 方式 1: 命令行参数
sql-diff -i --ai

# 方式 2: 配置文件
sql-diff config --ai-enabled --api-key YOUR_KEY
```

详细配置请查看 [AI 功能指南](/ai/guide)。

## 📋 生成 CREATE TABLE

### 使用交互式模式

```bash
# 启动交互式模式并启用 AI
sql-diff -i --ai

# 使用方向键选择
📋 请选择功能模式
   SQL 表结构比对
👉 AI 生成 CREATE TABLE
   AI 生成 ALTER TABLE

# 按 Enter 确认
```

### 操作流程

1. **选择功能**：使用方向键选择 "AI 生成 CREATE TABLE"
2. **输入描述**：用自然语言描述你要创建的表
3. **生成 SQL**：AI 自动生成完整的建表语句
4. **保存结果**：可选择保存到文件

### 示例 1：创建用户表

**输入描述**：
```
创建用户表，包含 ID、用户名、邮箱、密码、创建时间
```

**AI 生成**：
```sql
CREATE TABLE users (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  username VARCHAR(50) NOT NULL COMMENT '用户名',
  email VARCHAR(255) NOT NULL COMMENT '邮箱',
  password VARCHAR(255) NOT NULL COMMENT '密码',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (id),
  UNIQUE INDEX uk_username (username),
  UNIQUE INDEX uk_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

**AI 自动优化**：
- ✅ 使用 BIGINT 类型的自增主键
- ✅ 为用户名和邮箱添加唯一索引
- ✅ 合适的字段长度
- ✅ 完整的注释
- ✅ InnoDB 引擎和 UTF8MB4 编码

### 示例 2：创建商品表

**输入描述**：
```
创建商品表：商品ID、名称、价格、库存、分类、状态（上架/下架）、创建时间、更新时间
```

**AI 生成**：
```sql
CREATE TABLE products (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  name VARCHAR(255) NOT NULL COMMENT '商品名称',
  price DECIMAL(10, 2) NOT NULL COMMENT '价格',
  stock INT NOT NULL DEFAULT 0 COMMENT '库存',
  category VARCHAR(100) COMMENT '分类',
  status ENUM('online', 'offline') NOT NULL DEFAULT 'offline' COMMENT '状态：online-上架，offline-下架',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id),
  INDEX idx_category (category),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';
```

**AI 优化亮点**：
- ✅ DECIMAL 类型存储价格（精确计算）
- ✅ 使用 ENUM 类型管理状态
- ✅ 自动更新时间戳
- ✅ 为常用查询字段添加索引
- ✅ 合理的默认值

### 示例 3：创建订单表

**输入描述**：
```
创建订单表，包含订单号、用户ID、商品ID、数量、总价、订单状态、支付时间、发货时间、完成时间
```

**AI 生成**：
```sql
CREATE TABLE orders (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  order_no VARCHAR(32) NOT NULL COMMENT '订单号',
  user_id BIGINT NOT NULL COMMENT '用户ID',
  product_id BIGINT NOT NULL COMMENT '商品ID',
  quantity INT NOT NULL DEFAULT 1 COMMENT '数量',
  total_price DECIMAL(10, 2) NOT NULL COMMENT '总价',
  status ENUM('pending', 'paid', 'shipped', 'completed', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '订单状态',
  paid_at TIMESTAMP NULL COMMENT '支付时间',
  shipped_at TIMESTAMP NULL COMMENT '发货时间',
  completed_at TIMESTAMP NULL COMMENT '完成时间',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE INDEX uk_order_no (order_no),
  INDEX idx_user_id (user_id),
  INDEX idx_product_id (product_id),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';
```

**AI 优化亮点**：
- ✅ 订单号使用唯一索引
- ✅ 完整的订单状态流转
- ✅ 多个时间戳字段记录关键节点
- ✅ 外键字段添加索引（user_id, product_id）
- ✅ 合理的状态 ENUM 设计

## 🔧 生成 ALTER TABLE

### 使用交互式模式

```bash
# 启动交互式模式并启用 AI
sql-diff -i --ai

# 使用方向键选择
📋 请选择功能模式
   SQL 表结构比对
   AI 生成 CREATE TABLE
👉 AI 生成 ALTER TABLE

# 按 Enter 确认
```

### 操作流程

1. **选择功能**：使用方向键选择 "AI 生成 ALTER TABLE"
2. **粘贴现有表结构**：粘贴当前的 CREATE TABLE 语句
3. **输入变更需求**：用自然语言描述要做的修改
4. **生成 DDL**：AI 自动生成变更语句
5. **保存结果**：可选择保存到文件

### 示例 1：添加字段

**现有表结构**：
```sql
CREATE TABLE users (
  id BIGINT NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255),
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**变更需求**：
```
添加手机号字段，邮箱改为唯一索引，新增创建时间和更新时间
```

**AI 生成**：
```sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20) COMMENT '手机号';
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
ALTER TABLE users ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';
ALTER TABLE users ADD UNIQUE INDEX uk_email (email);
ALTER TABLE users ADD INDEX idx_phone (phone);
```

### 示例 2：修改字段

**现有表结构**：
```sql
CREATE TABLE products (
  id INT PRIMARY KEY,
  name VARCHAR(50),
  price FLOAT,
  stock INT
) ENGINE=InnoDB;
```

**变更需求**：
```
ID 改为 BIGINT 自增，价格改为 DECIMAL 精确类型，名称长度扩展到 255，添加商品描述字段
```

**AI 生成**：
```sql
ALTER TABLE products MODIFY COLUMN id BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID';
ALTER TABLE products MODIFY COLUMN name VARCHAR(255) NOT NULL COMMENT '商品名称';
ALTER TABLE products MODIFY COLUMN price DECIMAL(10, 2) NOT NULL COMMENT '价格';
ALTER TABLE products ADD COLUMN description TEXT COMMENT '商品描述';
ALTER TABLE products ADD INDEX idx_name (name);
```

### 示例 3：优化索引

**现有表结构**：
```sql
CREATE TABLE orders (
  id BIGINT PRIMARY KEY,
  user_id BIGINT,
  product_id BIGINT,
  status VARCHAR(20),
  created_at TIMESTAMP
) ENGINE=InnoDB;
```

**变更需求**：
```
为 user_id 和 product_id 添加索引，为 status 和 created_at 添加复合索引，添加订单号字段
```

**AI 生成**：
```sql
ALTER TABLE orders ADD COLUMN order_no VARCHAR(32) NOT NULL COMMENT '订单号';
ALTER TABLE orders ADD UNIQUE INDEX uk_order_no (order_no);
ALTER TABLE orders ADD INDEX idx_user_id (user_id);
ALTER TABLE orders ADD INDEX idx_product_id (product_id);
ALTER TABLE orders ADD INDEX idx_status_created (status, created_at);
```

**AI 优化亮点**：
- ✅ 识别需要复合索引的场景
- ✅ 为订单号添加唯一索引
- ✅ 合理的索引命名规范

## 🎯 使用场景

### 1. 快速原型开发

**场景**：需要快速创建数据库表结构

```bash
# 连续创建多个表
sql-diff -i --ai

# 1. 创建用户表
描述: 创建用户表...

# 2. 创建商品表
描述: 创建商品表...

# 3. 创建订单表
描述: 创建订单表...
```

### 2. 需求变更

**场景**：产品经理提出新需求，需要修改表结构

```
产品需求：用户表增加手机号验证功能
```

```bash
sql-diff -i --ai
# 选择 AI 生成 ALTER TABLE
# 粘贴现有用户表结构
# 描述：添加手机号字段和验证码字段
```

### 3. 学习标准 SQL

**场景**：学习如何编写标准的建表语句

```bash
# 输入简单描述，查看 AI 生成的标准 SQL
描述: 创建一个简单的博客文章表
# AI 会生成包含最佳实践的标准 SQL
```

### 4. 代码审查准备

**场景**：准备 SQL 变更的 PR 文档

```bash
sql-diff -i --ai -o migration_001.sql
# 生成的 SQL 可以直接用于 PR
# 包含完整的注释和规范的格式
```

## 💡 最佳实践

### 1. 清晰的描述

**推荐** ✅：
```
创建用户表，包含 ID、用户名、邮箱、密码、手机号、创建时间、更新时间
```

**不推荐** ❌：
```
用户表
```

### 2. 指定关键约束

**推荐** ✅：
```
创建订单表，订单号唯一，用户ID和商品ID需要索引
```

**普通** ⚪：
```
创建订单表，包含订单号、用户ID、商品ID
```

### 3. 说明字段用途

**推荐** ✅：
```
添加状态字段，用于标记订单状态：待支付、已支付、已发货、已完成
```

**普通** ⚪：
```
添加状态字段
```

### 4. 验证生成结果

生成 SQL 后，务必：
- ✅ 检查字段类型是否合适
- ✅ 确认索引设计是否合理
- ✅ 验证字段长度是否足够
- ✅ 在测试环境执行验证

### 5. 保存到文件

```bash
# 生成后选择保存
是否保存到文件? [y/N]: y
请输入文件名: migrations/001_create_users.sql
✓ SQL 已保存到: migrations/001_create_users.sql
```

或使用 `-o` 参数：
```bash
sql-diff -i --ai -o migrations/001_create_users.sql
```

## ⚠️ 注意事项

### 1. 人工审查

AI 生成的 SQL 仅供参考，务必：
- ⚠️ 仔细审查生成的 SQL
- ⚠️ 根据实际业务调整
- ⚠️ 在测试环境验证
- ⚠️ 考虑现有数据兼容性

### 2. 业务理解

AI 基于通用最佳实践，可能需要调整：
- ⚠️ 字段长度根据实际业务调整
- ⚠️ 索引策略根据查询模式优化
- ⚠️ 数据类型根据精度要求选择

### 3. 环境差异

不同数据库版本可能有差异：
- ⚠️ 检查 MySQL 版本兼容性
- ⚠️ 确认字符集设置
- ⚠️ 验证存储引擎支持

## 🔧 高级技巧

### 1. 复杂表结构

**描述技巧**：分步骤描述

```
创建电商订单表：
基础字段：订单ID、订单号、用户ID
商品信息：商品ID、商品名称、数量、单价、总价
状态管理：订单状态（待支付、已支付、已发货、已完成、已取消）
时间记录：创建时间、支付时间、发货时间、完成时间
其他：备注、收货地址
```

### 2. 关联表设计

**分别创建，再建立关联**：

```bash
# 1. 创建用户表
描述: 创建用户表...

# 2. 创建订单表
描述: 创建订单表，关联用户表的 user_id...

# 3. 添加外键（可选）
描述: 为订单表的 user_id 添加外键约束
```

### 3. 批量修改

**一次描述多个变更**：

```
对用户表进行以下修改：
1. 添加手机号字段
2. 邮箱改为唯一索引
3. 添加用户状态字段（正常、冻结、注销）
4. 添加最后登录时间字段
```

## 📊 性能和成本

### API 调用

每次生成会调用一次 AI API：

- **CREATE TABLE**：约 1000-2000 tokens
- **ALTER TABLE**：约 1500-3000 tokens

### 成本估算

以 DeepSeek 为例：
- 单次生成：约 ¥0.002-0.006
- 每月 50 次生成：约 ¥0.1-0.3

详见 [AI 功能指南 - 性能和成本](/ai/guide#性能和成本)

## 🆚 对比传统方式

| 方面 | 传统手写 SQL | AI 生成 SQL | 优势 |
|------|-------------|------------|------|
| **速度** | 5-15 分钟 | 2-5 秒 | ⚡ 效率提升 100 倍 |
| **规范性** | 依赖个人经验 | 统一最佳实践 | ✅ 代码质量稳定 |
| **索引设计** | 容易遗漏 | 自动添加 | ✅ 性能优化 |
| **注释** | 经常缺失 | 自动生成 | ✅ 可维护性强 |
| **学习成本** | 需要学习 SQL | 自然语言即可 | ✅ 降低门槛 |

## 🎓 学习建议

### 初学者

1. 先用简单描述生成 SQL
2. 学习 AI 生成的标准格式
3. 理解各个字段类型的选择
4. 逐步掌握索引设计

### 进阶用户

1. 对比 AI 生成和自己设计的差异
2. 学习 AI 的优化思路
3. 根据业务调整 AI 生成结果
4. 总结适合自己项目的模式

## 📚 相关文档

- [交互式光标选择](/guide/interactive) - 了解交互式界面
- [AI 功能指南](/ai/guide) - AI 功能配置
- [命令行工具](/guide/cli) - CLI 完整用法
- [最佳实践](/ai/best-practices) - AI 使用最佳实践
- [示例](/examples/basic) - 更多实际示例

## 🤝 反馈

如果 AI 生成的 SQL 有问题或有改进建议，欢迎：
- [提交 Issue](https://github.com/Bacchusgift/sql-diff/issues)
- [贡献代码](https://github.com/Bacchusgift/sql-diff/blob/main/CONTRIBUTING.md)
