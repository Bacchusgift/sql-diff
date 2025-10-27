# 交互式光标选择模式

SQL-Diff 提供了现代化的交互式界面，让你通过光标上下移动选择功能，无需记忆命令参数。

## 🎯 核心特性

### ✨ 光标导航选择

不再需要输入数字 1、2、3，直接使用键盘方向键选择：

- **⬆️⬇️** - 上下箭头移动光标
- **⏎ Enter** - 确认选择
- **^C Ctrl+C** - 取消退出

### 🎨 视觉反馈

清晰的视觉指示，让选择更直观：

| 状态 | 显示效果 | 说明 |
|------|---------|------|
| **当前选中** | 👉 SQL 表结构比对 | 青色加粗 |
| **未选中** | &nbsp;&nbsp;SQL 表结构比对 | 白色普通 |
| **确认选择** | ✅ SQL 表结构比对 | 绿色加粗 |
| **禁用选项** | &nbsp;&nbsp;AI 生成 CREATE TABLE [需要 AI] | 灰色 + 红色标记 |

### 📝 实时功能说明

光标移动时，底部自动显示当前选项的详细说明：

```
--------- 功能说明 ---------
描述: 比较两个表结构差异，自动生成 DDL 补全语句
```

## 🚀 快速开始

### 启动交互模式

```bash
# 基础模式（不启用 AI）
sql-diff -i

# 启用 AI 功能
sql-diff -i --ai
```

### 界面效果

当你运行 `sql-diff -i` 时，会看到这样的界面：

```
📋 请选择功能模式
👉 SQL 表结构比对
   AI 生成 CREATE TABLE [需要 AI]
   AI 生成 ALTER TABLE [需要 AI]

--------- 功能说明 ---------
描述: 比较两个表结构差异，自动生成 DDL 补全语句
```

使用方向键移动到第二项：

```
📋 请选择功能模式
   SQL 表结构比对
👉 AI 生成 CREATE TABLE [需要 AI]
   AI 生成 ALTER TABLE [需要 AI]

--------- 功能说明 ---------
描述: 根据自然语言描述，AI 生成完整的建表语句
要求: 需要启用 AI 功能
```

按 Enter 确认选择：

```
✅ AI 生成 CREATE TABLE
```

## 📋 三种功能模式

### 1️⃣ SQL 表结构比对

**功能**：比较两个表结构差异，自动生成 DDL 补全语句

**使用场景**：
- 数据库版本升级时比对表结构变更
- 开发环境与生产环境的表结构同步
- 代码审查时检查 SQL 变更

**操作流程**：

1. 选择 "SQL 表结构比对"
2. 粘贴源表的 CREATE TABLE 语句
3. 输入 `END` 或连续按两次 Enter 结束
4. 粘贴目标表的 CREATE TABLE 语句
5. 再次输入 `END` 或连续按两次 Enter
6. 自动生成 DDL 变更语句

**示例**：

```bash
$ sql-diff -i

📋 请选择功能模式
👉 SQL 表结构比对
   AI 生成 CREATE TABLE [需要 AI]
   AI 生成 ALTER TABLE [需要 AI]

[按 Enter]

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

✓ 已读取 67 个字符

📋 请粘贴目标表的 CREATE TABLE 语句：
（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）

CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(255),
  created_at TIMESTAMP
);
END

✓ 已读取 112 个字符

🔍 开始比对...

✓ 生成成功！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

➕ 新增列 (2):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255);
  2. ALTER TABLE users ADD COLUMN created_at TIMESTAMP;

📋 完整执行脚本:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN created_at TIMESTAMP;

是否保存到文件? [y/N]:
```

### 2️⃣ AI 生成 CREATE TABLE

**功能**：根据自然语言描述，AI 生成完整的建表语句

**要求**：需要启用 AI 功能

**使用场景**：
- 快速创建新表，无需手写 SQL
- 将产品需求直接转换为数据库表结构
- 学习标准的建表语句写法

**操作流程**：

1. 选择 "AI 生成 CREATE TABLE"
2. 输入自然语言描述（如：创建用户表，包含 ID、用户名、邮箱、密码、创建时间）
3. AI 自动生成标准的 CREATE TABLE 语句
4. 可选择保存到文件

**示例**：

```bash
$ sql-diff -i --ai

📋 请选择功能模式
   SQL 表结构比对
👉 AI 生成 CREATE TABLE
   AI 生成 ALTER TABLE

[按 Enter]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       模式 2: AI 生成 CREATE TABLE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💬 请描述您要创建的表结构：
（示例：创建用户表，包含 ID、用户名、邮箱、密码、创建时间）

描述: 创建商品表，包含商品ID、商品名称、价格、库存、分类、创建时间、更新时间

🤖 正在使用 AI 生成 SQL...

✓ 生成成功！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 生成的 CREATE TABLE 语句:

CREATE TABLE products (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  name VARCHAR(255) NOT NULL COMMENT '商品名称',
  price DECIMAL(10, 2) NOT NULL COMMENT '价格',
  stock INT NOT NULL DEFAULT 0 COMMENT '库存',
  category VARCHAR(100) COMMENT '分类',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id),
  INDEX idx_category (category),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

是否保存到文件? [y/N]:
```

### 3️⃣ AI 生成 ALTER TABLE

**功能**：基于现有表结构 + 自然语言描述，AI 生成 DDL 变更语句

**要求**：需要启用 AI 功能

**使用场景**：
- 快速调整现有表结构
- 根据需求变更生成 DDL
- 获取标准的 ALTER TABLE 语句

**操作流程**：

1. 选择 "AI 生成 ALTER TABLE"
2. 粘贴现有表的 CREATE TABLE 语句
3. 输入 `END` 或连续按两次 Enter 结束
4. 描述要做的修改（如：添加手机号字段、邮箱改为唯一索引）
5. AI 自动生成 ALTER TABLE 语句
6. 可选择保存到文件

**示例**：

```bash
$ sql-diff -i --ai

📋 请选择功能模式
   SQL 表结构比对
   AI 生成 CREATE TABLE
👉 AI 生成 ALTER TABLE

[按 Enter]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       模式 3: AI 生成 ALTER TABLE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 请粘贴现有表的 CREATE TABLE 语句：
（直接粘贴完整 SQL，粘贴完成后输入 'END' 或连续按两次 Enter）

CREATE TABLE users (
  id BIGINT NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255),
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
END

✓ 已读取 158 个字符

💬 请描述您要做的修改：
（示例：添加手机号字段、邮箱改为唯一索引）

描述: 添加手机号字段phone，邮箱改为唯一索引，新增创建时间和更新时间字段

🤖 正在使用 AI 生成 SQL...

✓ 生成成功！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 生成的 ALTER TABLE 语句:

ALTER TABLE users ADD COLUMN phone VARCHAR(20) COMMENT '手机号';
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
ALTER TABLE users ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';
ALTER TABLE users ADD UNIQUE INDEX uk_email (email);
ALTER TABLE users ADD INDEX idx_phone (phone);

是否保存到文件? [y/N]:
```

## 💡 使用技巧

### 1. 多行 SQL 输入

支持直接从数据库工具复制粘贴：

```bash
# ✅ 正确：直接粘贴多行 SQL
CREATE TABLE users (
  id INT,
  name VARCHAR(100)
);

# ✅ 正确：带注释的 SQL
CREATE TABLE users (
  id INT COMMENT '用户ID',
  name VARCHAR(100) COMMENT '用户名'
);

# ✅ 正确：复杂的建表语句
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 2. 结束输入的方式

有两种方式结束多行输入：

**方式 1：输入 END**
```bash
CREATE TABLE users (...);
END
```

**方式 2：连续按两次 Enter**
```bash
CREATE TABLE users (...);
[Enter]
[Enter]
```

### 3. AI 功能状态提示

未启用 AI 时，AI 相关选项会显示 `[需要 AI]` 标记：

```
📋 请选择功能模式
👉 SQL 表结构比对
   AI 生成 CREATE TABLE [需要 AI]  ← 红色标记
   AI 生成 ALTER TABLE [需要 AI]   ← 红色标记
```

如果选择了被禁用的选项，会显示详细提示：

```
✗ 该功能需要启用 AI

请通过以下方式之一启用 AI：
  1. 配置文件: 编辑 .sql-diff-config.yaml，设置 ai.enabled: true
  2. 命令行参数: 使用 --ai 参数启动

配置示例：
  sql-diff config  # 运行配置向导
```

### 4. 保存结果到文件

每次生成结果后，都会询问是否保存：

```bash
是否保存到文件? [y/N]: y
请输入文件名: migration_20250123.sql
✓ SQL 已保存到: migration_20250123.sql
```

也可以在启动时直接指定输出文件：

```bash
sql-diff -i -o migration.sql
```

## 🎨 视觉设计

### 界面元素

| 元素 | 符号 | 颜色 | 用途 |
|------|------|------|------|
| 当前选中 | 👉 | 青色加粗 | 指示光标位置 |
| 确认选择 | ✅ | 绿色加粗 | 确认用户选择 |
| 成功提示 | ✓ | 绿色 | 操作成功 |
| 错误提示 | ✗ | 红色 | 操作失败 |
| 警告标记 | [需要 AI] | 红色 | 功能未启用 |
| 信息标题 | 📋 📊 💬 🤖 | 默认色 | 分区标识 |

### 分隔线风格

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       模式标题
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 功能说明面板

```
--------- 功能说明 ---------
描述: 功能的详细说明
要求: 使用条件（如果有）
```

## ⚙️ 启用 AI 功能

### 方式 1：命令行参数

```bash
sql-diff -i --ai
```

### 方式 2：配置文件

创建 `.sql-diff-config.yaml`：

```yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-your-api-key
  model: deepseek-chat
```

然后直接运行：

```bash
sql-diff -i
```

### 方式 3：环境变量

```bash
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_API_KEY=sk-your-api-key
export SQL_DIFF_AI_PROVIDER=deepseek

sql-diff -i
```

### 方式 4：配置命令

```bash
sql-diff config \
  --ai-enabled \
  --provider deepseek \
  --api-key YOUR_KEY \
  >> ~/.bashrc

source ~/.bashrc
sql-diff -i
```

## 🔧 故障排查

### 光标选择不响应

**问题**：按上下箭头没有反应

**解决**：
- 确保终端支持 ANSI 转义序列
- Windows 用户建议使用 Git Bash 或 WSL
- 检查终端设置，确保未禁用光标移动

### AI 功能提示未启用

**问题**：选择 AI 功能时提示需要启用

**解决**：
```bash
# 检查配置
sql-diff config --show

# 或直接使用 --ai 参数
sql-diff -i --ai
```

### 多行输入被截断

**问题**：粘贴多行 SQL 只识别第一行

**解决**：
- 确保粘贴完成后输入 `END` 或连续按两次 Enter
- 不要使用 Ctrl+D，这会直接结束程序

### 无法取消选择

**问题**：想退出但不知道如何操作

**解决**：
- 按 `Ctrl+C` 或 `Esc` 取消选择
- 或直接关闭终端窗口

## 🎯 最佳实践

1. **准备好 SQL**：在文本编辑器中准备好 SQL，然后直接粘贴
2. **善用 AI**：复杂需求优先使用 AI 生成，提高效率
3. **保存结果**：重要的 DDL 语句记得保存到文件
4. **检查提示**：仔细阅读底部的功能说明面板
5. **快捷操作**：熟练使用方向键和 Enter，提升操作速度

## 📚 相关文档

- [命令行工具](/guide/cli) - 完整的 CLI 使用指南
- [AI 功能](/ai/guide) - AI 功能配置和使用
- [配置](/config/environment) - 环境变量和配置文件
- [示例](/examples/basic) - 更多实际使用示例
