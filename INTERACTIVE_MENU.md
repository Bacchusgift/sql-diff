# 交互式菜单功能说明

## 🎯 功能概述

新的交互式模式提供了友好的菜单选择界面，用户无需记忆命令参数，通过光标选择即可使用所有功能。

## 🚀 使用方法

### 启动交互式模式

```bash
sql-diff -i
```

### 界面展示

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       SQL 表结构比对工具 - 交互式模式
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🤖 AI 智能分析: 已启用
   提供商: deepseek
   模型: deepseek-chat

📋 请选择功能模式：

  [1] SQL 表结构比对
      比较两个表结构差异，自动生成 DDL 补全语句

  [2] AI 生成 CREATE TABLE (需要 AI)
      根据自然语言描述，AI 生成完整的建表语句

  [3] AI 生成 ALTER TABLE (需要 AI)
      基于现有表结构 + 自然语言描述，AI 生成 DDL 变更语句

请输入选项编号 [1-3]: 
```

## 📋 三种模式详解

### 模式 1：SQL 表结构比对

**使用场景**：
- 比较开发环境和生产环境的表结构差异
- 数据库迁移前的差异检查
- Code Review 中的表结构变更审查

**操作流程**：
1. 选择 `1`
2. 粘贴源表的 CREATE TABLE 语句
3. 输入 `END` 或连续按两次 Enter
4. 粘贴目标表的 CREATE TABLE 语句
5. 输入 `END` 或连续按两次 Enter
6. 查看生成的 DDL 差异语句

**示例**：
```
请输入选项编号 [1-3]: 1

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       模式 1: SQL 表结构比对
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 请粘贴源表的 CREATE TABLE 语句：
（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）
（提示：建议在文本编辑器中准备好 SQL，然后直接粘贴）

CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
);
END

✓ 已读取 XXX 个字符

📋 请粘贴目标表的 CREATE TABLE 语句：
...
```

---

### 模式 2：AI 生成 CREATE TABLE

**使用场景**：
- 快速原型设计
- 从产品需求直接生成表结构
- 学习数据库设计最佳实践

**操作流程**：
1. 选择 `2`
2. 输入表结构的自然语言描述
3. AI 自动生成完整的 CREATE TABLE 语句
4. 选择是否保存到文件

**示例**：
```
请输入选项编号 [1-3]: 2

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       模式 2: AI 生成 CREATE TABLE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💬 请描述您要创建的表结构：
（示例：创建用户表，包含 ID、用户名、邮箱、密码、创建时间）

描述: 创建订单表，包含订单号（唯一）、用户ID、商品ID、数量、金额（精确到分）、订单状态、下单时间

🤖 正在使用 AI 生成 SQL...

✓ 生成成功！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 生成的 CREATE TABLE 语句:

CREATE TABLE orders (
  order_no VARCHAR(50) NOT NULL UNIQUE COMMENT '订单号',
  user_id INT NOT NULL COMMENT '用户ID',
  product_id INT NOT NULL COMMENT '商品ID',
  quantity INT NOT NULL DEFAULT 1 COMMENT '数量',
  amount DECIMAL(10,2) NOT NULL COMMENT '金额（精确到分）',
  status ENUM('pending', 'paid', 'shipped', 'completed', 'cancelled') DEFAULT 'pending' COMMENT '订单状态',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '下单时间',
  PRIMARY KEY (order_no),
  INDEX idx_user_id (user_id),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

是否保存到文件? [y/N]: y
请输入文件名: orders.sql
✓ SQL 已保存到: orders.sql
```

**AI 生成特点**：
- ✅ 自动推断字段类型（VARCHAR、INT、DECIMAL、ENUM等）
- ✅ 自动添加主键、索引、唯一约束
- ✅ 应用 MySQL 最佳实践（InnoDB、UTF8MB4、注释）
- ✅ 使用蛇形命名法（snake_case）

---

### 模式 3：AI 生成 ALTER TABLE

**使用场景**：
- 表结构迭代升级
- 需求变更时的快速 DDL 生成
- 数据库重构

**操作流程**：
1. 选择 `3`
2. 粘贴现有表的 CREATE TABLE 语句
3. 输入 `END` 或连续按两次 Enter
4. 输入修改需求的自然语言描述
5. AI 自动生成 ALTER TABLE 语句
6. 选择是否保存到文件

**示例**：
```
请输入选项编号 [1-3]: 3

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       模式 3: AI 生成 ALTER TABLE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 请粘贴现有表的 CREATE TABLE 语句：
（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）
（提示：建议在文本编辑器中准备好 SQL，然后直接粘贴）

CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
END

✓ 已读取 XXX 个字符

💬 请描述您要做的修改：
（示例：添加手机号字段、邮箱改为唯一索引）

描述: 添加手机号字段（必填、唯一），添加实名认证字段（真实姓名、身份证号、认证状态），邮箱改为唯一索引

🤖 正在使用 AI 生成 SQL...

✓ 生成成功！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 生成的 ALTER TABLE 语句:

ALTER TABLE users ADD COLUMN mobile VARCHAR(20) NOT NULL UNIQUE COMMENT '手机号';
ALTER TABLE users ADD COLUMN real_name VARCHAR(50) DEFAULT NULL COMMENT '真实姓名';
ALTER TABLE users ADD COLUMN id_card VARCHAR(18) DEFAULT NULL COMMENT '身份证号';
ALTER TABLE users ADD COLUMN verified_status ENUM('unverified', 'pending', 'verified') DEFAULT 'unverified' COMMENT '认证状态';
ALTER TABLE users ADD UNIQUE INDEX uk_email (email);
ALTER TABLE users ADD INDEX idx_id_card (id_card);

是否保存到文件? [y/N]: n
未保存到文件
```

**AI 生成特点**：
- ✅ 基于现有表结构生成变更
- ✅ 可生成多条 ALTER 语句
- ✅ 考虑数据迁移安全性
- ✅ 保持与现有表的一致性

---

## 🎨 界面特点

### 1. AI 状态显示

**AI 已启用**：
```
🤖 AI 智能分析: 已启用
   提供商: deepseek
   模型: deepseek-chat
```

**AI 未启用**：
```
ℹ️  AI 智能分析: 未启用
   （可通过 --ai 参数或配置文件启用）
```

**禁用选项显示**：
当 AI 未启用时，模式 2 和 3 会以灰色显示：
```
  [2] AI 生成 CREATE TABLE (需要 AI) [未启用]
      根据自然语言描述，AI 生成完整的建表语句

  [3] AI 生成 ALTER TABLE (需要 AI) [未启用]
      基于现有表结构 + 自然语言描述，AI 生成 DDL 变更语句
```

如果选择了禁用的选项，会显示提示：
```
✗ 该功能需要启用 AI

请通过以下方式之一启用 AI：
  1. 配置文件: 编辑 .sql-diff-config.yaml，设置 ai.enabled: true
  2. 命令行参数: 使用 --ai 参数启动

配置示例：
  sql-diff config  # 运行配置向导
```

### 2. 彩色输出

- 🔵 **蓝色**：标题、分隔线
- 🟢 **绿色**：成功消息、AI 启用状态、可用选项
- 🟡 **黄色**：输入提示
- ⚪ **白色**：描述文字
- ⚫ **灰色**：禁用选项

### 3. 友好提示

每个步骤都有清晰的提示：
- 📋 输入SQL时的操作说明
- 💬 自然语言描述的示例
- 🤖 AI 处理中的状态提示
- ✓ 成功完成的确认
- ✗ 错误时的详细说明

---

## 💡 使用技巧

### 1. 快速启动

```bash
# 启动交互式模式
sql-diff -i

# 启动交互式模式 + 启用 AI（如果配置文件中未启用）
sql-diff -i --ai
```

### 2. 准备工作

**对于模式 1 和模式 3**（需要粘贴 SQL）：
1. 在文本编辑器中准备好 CREATE TABLE 语句
2. 复制完整的 SQL
3. 在终端中直接粘贴
4. 输入 `END` 或连续按两次 Enter

**对于模式 2 和模式 3**（需要描述）：
1. 提前想好需求描述
2. 尽量详细、具体
3. 可以参考提示的示例格式

### 3. 常见场景

**场景 1：开发新功能，需要创建新表**
```
选择: [2] AI 生成 CREATE TABLE
描述: "创建优惠券表：优惠券ID、优惠券码（唯一）、优惠类型、折扣金额、使用条件、有效期、使用状态"
```

**场景 2：功能迭代，需要修改表结构**
```
选择: [3] AI 生成 ALTER TABLE
粘贴: 现有表结构
描述: "添加会员等级字段、积分字段、上次登录时间"
```

**场景 3：环境同步，检查表结构差异**
```
选择: [1] SQL 表结构比对
粘贴: 开发环境的表结构
粘贴: 生产环境的表结构
```

---

## 🎯 优势

### 与旧版本对比

#### ❌ 旧版本（命令行模式）

```bash
# 表结构比对
sql-diff -s "CREATE TABLE ..." -t "CREATE TABLE ..."

# AI 生成表
sql-diff generate -d "创建用户表"

# AI 生成 DDL
sql-diff alter -i -d "添加字段"
```

**问题**：
- 需要记忆不同的命令和参数
- 切换功能需要退出重新执行
- 不直观，新手不友好

#### ✅ 新版本（交互式菜单）

```bash
# 一条命令启动
sql-diff -i

# 菜单选择功能
[1/2/3] → 按需选择
```

**优势**：
- ✅ 无需记忆命令参数
- ✅ 所有功能一目了然
- ✅ 友好的引导提示
- ✅ 错误提示更清晰
- ✅ 支持保存到文件

---

## 📚 相关文档

- [README.md](./README.md) - 项目主文档
- [AI_SQL_GENERATOR.md](./AI_SQL_GENERATOR.md) - AI 功能详细说明
- [HELP_IMPROVEMENTS.md](./HELP_IMPROVEMENTS.md) - 帮助信息优化说明
- [VERSION_INFO.md](./VERSION_INFO.md) - 版本信息功能

---

## 🎉 总结

新的交互式菜单系统极大地提升了用户体验：

1. **简单直观**：通过数字选择功能，无需记忆命令
2. **引导清晰**：每个步骤都有明确的提示
3. **功能集成**：三大核心功能统一入口
4. **智能提示**：AI 状态实时显示，禁用选项明确标记
5. **友好反馈**：彩色输出，成功/错误清晰区分

**一句话**：让 SQL 表结构管理变得简单、高效、愉悦！🚀
