# Homebrew 快速部署指南

## 🚀 5 分钟完成 Homebrew 部署

### Step 1: 创建 Homebrew Tap 仓库 (1 分钟)

1. 访问 GitHub: https://github.com/new
2. 填写信息：
   - **Repository name**: `homebrew-tap`
   - **Description**: `Homebrew formulae for sql-diff`
   - **Public** ✅
3. 点击 "Create repository"

### Step 2: 创建 GitHub Release (2 分钟)

1. 访问：https://github.com/Bacchusgift/sql-diff/releases/new
2. 填写 Release 信息：
   - **Choose a tag**: `v1.0.1` (已存在)
   - **Release title**: `v1.0.1 - 交互式输入模式`
   - **Description**: 可以使用自动生成的 Release Notes
3. 点击 "Publish release"

### Step 3: 计算 SHA256 (30 秒)

等待 Release 创建完成后，在终端运行：

```bash
curl -sL https://github.com/Bacchusgift/sql-diff/archive/refs/tags/v1.0.1.tar.gz | shasum -a 256
```

复制输出的 SHA256 值（一串 64 位的十六进制字符）

### Step 4: 更新并推送 Formula (1 分钟)

```bash
# 克隆 tap 仓库
git clone https://github.com/Bacchusgift/homebrew-tap.git
cd homebrew-tap

# 创建目录
mkdir -p Formula

# 复制 Formula 文件
cp ../sql-diff/Formula/sql-diff.rb Formula/

# 编辑 Formula，将 Step 3 计算的 SHA256 填入
# 找到这一行：sha256 "" 
# 改为：sha256 "你的SHA256值"
vim Formula/sql-diff.rb
# 或使用其他编辑器

# 提交并推送
git add Formula/sql-diff.rb
git commit -m "Add sql-diff v1.0.1"
git push origin main
```

### Step 5: 测试安装 (30 秒)

```bash
# 添加 tap
brew tap Bacchusgift/tap

# 安装
brew install sql-diff

# 验证
sql-diff --version
sql-diff --help
```

## ✅ 完成！

现在任何人都可以通过以下方式安装：

```bash
brew install Bacchusgift/tap/sql-diff
```

---

## 📋 完整命令速查

### 一次性执行（假设你已经有 SHA256）

```bash
# 克隆 tap 仓库
git clone https://github.com/Bacchusgift/homebrew-tap.git
cd homebrew-tap

# 复制 Formula
mkdir -p Formula
cp ../sql-diff/Formula/sql-diff.rb Formula/

# 更新 SHA256（替换为实际值）
sed -i '' 's/sha256 ""/sha256 "YOUR_SHA256_HERE"/' Formula/sql-diff.rb

# 提交
git add Formula/sql-diff.rb
git commit -m "Add sql-diff v1.0.1"
git push origin main

# 测试
brew tap Bacchusgift/tap
brew install sql-diff
sql-diff --version
```

---

## 🔄 后续版本更新

每次发布新版本（例如 v1.0.2）时：

```bash
# 1. 打新 tag
git tag -a v1.0.2 -m "Release v1.0.2"
git push origin v1.0.2

# 2. 在 GitHub 创建 Release（会自动触发 Actions）

# 3. 计算新版本的 SHA256
curl -sL https://github.com/Bacchusgift/sql-diff/archive/refs/tags/v1.0.2.tar.gz | shasum -a 256

# 4. 更新 homebrew-tap
cd homebrew-tap
vim Formula/sql-diff.rb
# 更新 url 和 sha256

git add Formula/sql-diff.rb
git commit -m "Update sql-diff to v1.0.2"
git push origin main

# 5. 用户更新
brew upgrade sql-diff
```

---

## 🎯 验证清单

- [ ] Tap 仓库已创建：`https://github.com/Bacchusgift/homebrew-tap`
- [ ] GitHub Release 已发布：`https://github.com/Bacchusgift/sql-diff/releases/tag/v1.0.1`
- [ ] SHA256 已计算并填入 Formula
- [ ] Formula 已推送到 tap 仓库
- [ ] 本地测试安装成功：`brew install Bacchusgift/tap/sql-diff`
- [ ] 命令可以正常运行：`sql-diff --version`
- [ ] README 已更新安装说明

---

## 💡 提示

1. **SHA256 为空时的错误**
   - 错误信息：`Error: SHA256 mismatch`
   - 解决：确保填入了正确的 SHA256 值

2. **Formula 语法错误**
   ```bash
   # 审查 formula
   brew audit --strict Formula/sql-diff.rb
   ```

3. **测试 Formula**
   ```bash
   # 从本地文件安装测试
   brew install --build-from-source Formula/sql-diff.rb
   
   # 运行测试
   brew test sql-diff
   ```

4. **查看安装日志**
   ```bash
   brew install --verbose sql-diff
   ```

---

## 🎉 宣传

在项目 README 和文档中添加：

```markdown
## 安装

### macOS (Homebrew)

\`\`\`bash
brew install Bacchusgift/tap/sql-diff
\`\`\`

### Linux / Windows

参见[安装文档](https://bacchusgift.github.io/sql-diff/guide/installation)
```

---

需要帮助？查看完整文档：[HOMEBREW.md](./HOMEBREW.md)
