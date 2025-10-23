# 📝 交互式模式使用指南

## 🎯 优化说明

交互式模式已优化，现在支持：

✅ **粘贴多行 SQL**  
✅ **粘贴后继续编辑、删除**  
✅ **三种灵活的结束方式**

## 🚀 使用方法

### 启动交互模式

```bash
sql-diff -i
```

### 输入 SQL 的三种结束方式

#### 方式 1: 输入单独一行的 `END`（推荐）

```
请粘贴源表的 CREATE TABLE 语句：
（粘贴后可以继续编辑，输入单独一行的 'END' 或 ';' 结束）

CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
);
END           ← 输入 END 并按 Enter
```

#### 方式 2: 输入单独一行的 `;`

```
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
);
;             ← 输入单个分号并按 Enter
```

#### 方式 3: 连续两次按 Enter（输入空行）

```
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
);
              ← 第一次按 Enter（空行）
              ← 第二次按 Enter（结束）
```

## 💡 使用技巧

### 从 Navicat 复制表结构

1. 在 Navicat 中右键表 → 转储 SQL 文件 → 仅结构
2. 复制 CREATE TABLE 语句
3. 在 `sql-diff -i` 中直接粘贴
4. 粘贴后可以修改、删除不需要的部分
5. 输入 `END` 结束

### 处理包含注释的 SQL

```sql
-- 用户表
CREATE TABLE users (
  id INT PRIMARY KEY COMMENT '用户ID',
  name VARCHAR(100) COMMENT '用户名',
  created_at TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
END
```

### 多表处理建议

如果需要比对多个表，建议：

1. 准备好两个 SQL 文件
2. 使用 `cat source.sql | sql-diff -i` 方式
3. 或者分别打开文件复制粘贴

## ⚡ 完整示例

```bash
$ sql-diff -i

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       SQL 表结构比对工具 - 交互式模式
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🤖 AI 智能分析: 未启用
   （可通过 --ai 参数或配置文件启用）

📋 请粘贴源表的 CREATE TABLE 语句：
（粘贴后可以继续编辑，输入单独一行的 'END' 或 ';' 结束）
（也可以连续两次按 Enter（输入空行）结束）

CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
);
END

✓ 已读取 68 个字符

📋 请粘贴目标表的 CREATE TABLE 语句：
（粘贴后可以继续编辑，输入单独一行的 'END' 或 ';' 结束）
（也可以连续两次按 Enter（输入空行）结束）

CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
END

✓ 已读取 142 个字符

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       开始比对
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📖 正在解析源表结构...
✓ 源表: users (2 列)

📖 正在解析目标表结构...
✓ 目标表: users (4 列)

🔍 正在比对表结构...

📊 差异摘要:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
新增列: 2
  - email (VARCHAR(255))
  - created_at (TIMESTAMP)

🔧 生成 DDL 语句...

✓ 生成的 DDL 语句:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

➕ 新增列 (2):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255);
  2. ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

📋 完整执行脚本:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           完成！
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 🎨 启用 AI 分析

```bash
sql-diff -i --ai
```

AI 会提供：
- 💡 智能分析表结构差异
- 📊 SQL 优化建议
- ⚠️ 潜在风险提示
- ✨ 最佳实践建议

## ❓ 常见问题

### Q: 粘贴的 SQL 包含 `END` 关键字怎么办？

A: 使用方式 2（单独一行输入 `;`）或方式 3（连续两个空行）

### Q: 如何取消当前输入？

A: 按 `Ctrl+C` 可以随时中断

### Q: 支持从文件读取吗？

A: 支持！使用管道：
```bash
cat source.sql | sql-diff -i
```

### Q: 输入错了可以重新开始吗？

A: 当前输入完成后会提示错误，可以重新运行 `sql-diff -i`

## 📚 更多信息

- [完整文档](https://bacchusgift.github.io/sql-diff/)
- [命令行参数](https://bacchusgift.github.io/sql-diff/guide/cli)
- [AI 配置](https://bacchusgift.github.io/sql-diff/ai/guide)
