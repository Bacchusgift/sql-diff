# Homebrew Formula 安装错误修复指南

## 🐛 问题描述

**错误信息**：
```
Error: An exception occurred within a child process:
  RuntimeError: Not a Git repository: /private/tmp/sql-diff-20251024-36099-dslf4h/sql-diff-1.0.2
```

**原因**：
Homebrew Formula 中使用了 `Utils.git_short_head` 来获取 Git commit hash，但当用户通过 `tar.gz` 源码包安装时，源码不是 Git 仓库，导致报错。

## 🔧 修复步骤

### 第一步：修复 homebrew-tap 仓库中的 Formula

```bash
# 1. 进入 homebrew-tap 仓库
cd /Users/youzi/CascadeProjects/homebrew-tap

# 2. 编辑 Formula 文件
vim Formula/sql-diff.rb
# 或使用 VS Code
code Formula/sql-diff.rb
```

**修改内容**：

找到 `def install` 方法（大约第 10 行），将：

```ruby
  def install
    ldflags = %W[
      -s -w
      -X main.Version=#{version}
      -X main.BuildTime=#{time.iso8601}
      -X main.GitCommit=#{Utils.git_short_head}
    ]

    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/sql-diff"
    generate_completions_from_executable(bin/"sql-diff", "completion")
  end
```

**改为**：

```ruby
  def install
    # 构建标志
    ldflags = %W[
      -s -w
      -X main.Version=#{version}
      -X main.BuildTime=#{time.iso8601}
    ]
    
    # 仅在 HEAD 版本（从 Git 安装）时添加 GitCommit
    if build.head?
      ldflags << "-X main.GitCommit=#{Utils.git_short_head}"
    end

    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/sql-diff"
    generate_completions_from_executable(bin/"sql-diff", "completion")
  end
```

**提交更改**：

```bash
git add Formula/sql-diff.rb
git commit -m "fix: 修复从 tar.gz 安装时的 Git 仓库依赖错误

问题：Utils.git_short_head 在非 Git 仓库环境下会报错
解决：仅在 HEAD (Git) 安装时才添加 GitCommit 信息
"
git push origin main
```

### 第二步：验证修复

```bash
# 1. 卸载当前版本
brew uninstall sql-diff

# 2. 清理缓存
brew cleanup

# 3. 更新 tap
brew update

# 4. 重新安装
brew install Bacchusgift/tap/sql-diff

# 5. 验证
sql-diff --version
```

### 第三步：测试功能

```bash
# 测试基本比对
sql-diff -s "CREATE TABLE users (id INT);" -t "CREATE TABLE users (id INT, name VARCHAR(100));"

# 测试交互式模式
sql-diff -i

# 测试新功能（需要配置 AI）
sql-diff generate -d "创建用户表"
```

## 📝 完整的修复后 Formula

修复后的完整 Formula 应该是这样的：

```ruby
class SqlDiff < Formula
  desc "智能 SQL 表结构比对工具，支持交互式多行输入和 AI 分析"
  homepage "https://bacchusgift.github.io/sql-diff/"
  url "https://github.com/Bacchusgift/sql-diff/archive/refs/tags/v1.0.2.tar.gz"
  sha256 "14916df3412cbb81e1e9a0503196aa8da6e793ca764349cdf52fa31915e3cee7"
  license "MIT"
  head "https://github.com/Bacchusgift/sql-diff.git", branch: "main"

  depends_on "go" => :build

  def install
    # 构建标志
    ldflags = %W[
      -s -w
      -X main.Version=#{version}
      -X main.BuildTime=#{time.iso8601}
    ]
    
    # 仅在 HEAD 版本（从 Git 安装）时添加 GitCommit
    if build.head?
      ldflags << "-X main.GitCommit=#{Utils.git_short_head}"
    end

    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/sql-diff"
    generate_completions_from_executable(bin/"sql-diff", "completion")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/sql-diff --version")
    
    source_sql = "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"
    target_sql = "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255));"
    
    output = shell_output("#{bin}/sql-diff -s '#{source_sql}' -t '#{target_sql}'")
    assert_match "ADD COLUMN email", output
  end
end
```

## ✅ 自动化修复说明

**好消息**：sql-diff 项目的 `.github/workflows/release.yml` 已经修复了这个问题！

下次发布新版本（如 v1.1.0）时，GitHub Actions 会自动使用修复后的 Formula 模板，不会再出现这个错误。

## 🚀 后续步骤

修复完成后，你可以考虑：

1. **发布新版本**（如 v1.0.3）来触发自动更新
2. **或者**保持现状，等待下次功能更新时一起发布

## 📚 相关文件

- `/Users/youzi/CascadeProjects/homebrew-tap/Formula/sql-diff.rb` - 需要手动修复
- `/Users/youzi/CascadeProjects/sql-diff/.github/workflows/release.yml` - 已自动修复 ✅

## 💡 技术说明

**为什么会出现这个问题？**

Homebrew 支持两种安装方式：
1. **从 tar.gz 安装**（默认）：下载源码压缩包，不是 Git 仓库
2. **从 HEAD 安装**（可选）：直接克隆 Git 仓库

`Utils.git_short_head` 只能在 Git 仓库中使用，所以我们使用 `build.head?` 来判断是否是从 Git 安装，只在这种情况下才添加 GitCommit 信息。

**影响**：
- 从 tar.gz 安装：不包含 GitCommit 信息（不影响功能）
- 从 HEAD 安装：包含 GitCommit 信息（开发者友好）
