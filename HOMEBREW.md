# Homebrew 部署指南

本文档介绍如何将 SQL-Diff 部署到 Homebrew，让用户可以通过 `brew install` 安装。

## 📦 部署方式

### 方式 1：通过 Homebrew Tap（推荐）

这是最快捷的方式，无需等待官方审核。

#### 1. 创建 Homebrew Tap 仓库

```bash
# 创建一个名为 homebrew-tap 的仓库
# GitHub 仓库命名规则：homebrew-<tap名称>
# 例如：homebrew-tap, homebrew-tools 等
```

在 GitHub 上创建仓库：`https://github.com/Bacchusgift/homebrew-tap`

#### 2. 准备 Formula 文件

Formula 文件已经在 `Formula/sql-diff.rb`，需要：

1. **创建 GitHub Release**
   ```bash
   # 确保已经打好 tag
   git tag v1.0.1
   git push origin v1.0.1
   ```

2. **在 GitHub 上创建 Release**
   - 访问：https://github.com/Bacchusgift/sql-diff/releases/new
   - 选择 tag: v1.0.1
   - 填写 Release notes
   - 发布

3. **计算 SHA256**
   ```bash
   # 下载发布的 tar.gz
   curl -sL https://github.com/Bacchusgift/sql-diff/archive/refs/tags/v1.0.1.tar.gz | shasum -a 256
   ```

4. **更新 Formula 中的 SHA256**
   将计算出的值填入 `Formula/sql-diff.rb` 的 `sha256` 字段

#### 3. 将 Formula 推送到 Tap 仓库

```bash
# 克隆你的 tap 仓库
git clone https://github.com/Bacchusgift/homebrew-tap.git
cd homebrew-tap

# 复制 Formula 文件
mkdir -p Formula
cp /path/to/sql-diff/Formula/sql-diff.rb Formula/

# 提交并推送
git add Formula/sql-diff.rb
git commit -m "Add sql-diff formula"
git push origin main
```

#### 4. 用户安装方式

```bash
# 添加 tap
brew tap Bacchusgift/tap

# 安装 sql-diff
brew install sql-diff

# 或者一行命令
brew install Bacchusgift/tap/sql-diff
```

---

### 方式 2：提交到 Homebrew Core（长期目标）

适合成熟稳定的项目，需要满足以下条件：

#### 前置要求

- [ ] 项目有一定的用户量和 Stars
- [ ] 持续维护，有稳定的发布周期
- [ ] 通过所有测试
- [ ] 有完善的文档
- [ ] 开源协议清晰（MIT/Apache 2.0 等）

#### 提交步骤

1. **Fork Homebrew Core**
   ```bash
   # Fork https://github.com/Homebrew/homebrew-core
   ```

2. **创建 Formula**
   ```bash
   brew create https://github.com/Bacchusgift/sql-diff/archive/v1.0.1.tar.gz
   ```

3. **测试 Formula**
   ```bash
   brew install --build-from-source sql-diff
   brew test sql-diff
   brew audit --strict sql-diff
   ```

4. **提交 Pull Request**
   - 确保通过所有 CI 检查
   - 遵循 Homebrew 的贡献指南
   - 等待维护者审核

---

## 🔄 自动化发布流程

### GitHub Actions 自动发布

创建 `.github/workflows/release.yml`：

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build all platforms
        run: make build-all
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            dist/*
          generate_release_notes: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Update Homebrew Tap
        run: |
          # 自动更新 tap 仓库中的 formula
          # TODO: 实现自动更新逻辑
```

---

## 📝 Formula 文件说明

### 关键字段

```ruby
class SqlDiff < Formula
  desc "项目简短描述"
  homepage "项目主页"
  url "源码 tar.gz 地址"
  sha256 "文件的 SHA256 校验和"
  license "开源协议"
  
  depends_on "go" => :build  # 构建依赖
  
  def install
    # 编译和安装逻辑
  end
  
  test do
    # 测试逻辑
  end
end
```

### 版本更新

每次发布新版本时：

1. 更新 `url` 中的版本号
2. 重新计算并更新 `sha256`
3. 提交到 tap 仓库

---

## ✅ 验证安装

### 本地测试

```bash
# 测试 formula 语法
brew audit --strict --online Formula/sql-diff.rb

# 本地安装测试
brew install --build-from-source Formula/sql-diff.rb

# 测试功能
sql-diff --version
sql-diff -i
```

### 卸载

```bash
brew uninstall sql-diff
```

---

## 🎯 快速开始（推荐流程）

### Step 1: 创建 Tap 仓库

```bash
# 在 GitHub 上创建仓库
Repository name: homebrew-tap
Description: Homebrew formulae for Bacchusgift's tools
```

### Step 2: 创建 GitHub Release

```bash
# 本地打 tag（如果还没有）
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1

# 在 GitHub 网页上创建 Release
# https://github.com/Bacchusgift/sql-diff/releases/new
```

### Step 3: 计算 SHA256

```bash
# 等待 Release 创建完成后
curl -sL https://github.com/Bacchusgift/sql-diff/archive/refs/tags/v1.0.1.tar.gz | shasum -a 256
```

### Step 4: 更新并推送 Formula

```bash
# 克隆 tap 仓库
git clone https://github.com/Bacchusgift/homebrew-tap.git
cd homebrew-tap

# 创建目录
mkdir -p Formula

# 复制并更新 Formula（记得填入 SHA256）
cp ../sql-diff/Formula/sql-diff.rb Formula/

# 提交
git add Formula/sql-diff.rb
git commit -m "Add sql-diff v1.0.1"
git push origin main
```

### Step 5: 测试安装

```bash
# 添加 tap
brew tap Bacchusgift/tap

# 安装
brew install sql-diff

# 验证
sql-diff --version
sql-diff --help
```

---

## 📚 更多资源

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Creating Taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)

---

## 🐛 常见问题

### Q: SHA256 怎么计算？

```bash
# 方法 1：从 GitHub Release 下载
curl -sL https://github.com/Bacchusgift/sql-diff/archive/refs/tags/v1.0.1.tar.gz | shasum -a 256

# 方法 2：从本地文件
shasum -a 256 sql-diff-1.0.1.tar.gz
```

### Q: Formula 测试失败？

```bash
# 查看详细错误
brew install --verbose --debug Formula/sql-diff.rb

# 审查 formula
brew audit --strict Formula/sql-diff.rb
```

### Q: 如何支持多个版本？

创建不同的 formula 文件：
- `Formula/sql-diff.rb` - 最新稳定版
- `Formula/sql-diff@1.0.rb` - 1.0 系列
- `Formula/sql-diff@2.0.rb` - 2.0 系列

---

## 🎉 完成！

现在用户可以通过以下方式安装：

```bash
brew install Bacchusgift/tap/sql-diff
```

在 README 中添加安装说明：

```markdown
## 安装

### macOS (Homebrew)

\`\`\`bash
brew install Bacchusgift/tap/sql-diff
\`\`\`

### 其他平台

参见 [安装文档](https://bacchusgift.github.io/sql-diff/guide/installation)
\`\`\`
```
