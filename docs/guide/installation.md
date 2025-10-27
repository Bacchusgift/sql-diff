# 安装

SQL-Diff 提供多种安装方式,您可以根据自己的需求选择合适的方式。

## 前置要求

- Go 1.21 或更高版本
- Git

## 方式一: 从源码安装 (推荐)

### 1. 克隆仓库

```bash
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff
```

### 2. 编译安装

```bash
# 编译到 bin 目录
make build

# 或者直接安装到 GOPATH/bin
make install
```

### 3. 验证安装

```bash
sql-diff --version
```

## 方式二: 使用 go install

如果您的系统已经配置好了 Go 环境:

```bash
go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
```

确保 `$GOPATH/bin` 已添加到系统 PATH 中。

## 方式三: 下载预编译二进制 (即将支持)

我们计划在 GitHub Releases 页面提供预编译的二进制文件:

```bash
# macOS (Intel)
curl -L https://github.com/Bacchusgift/sql-diff/releases/download/v1.0.0/sql-diff-darwin-amd64 -o sql-diff

# macOS (Apple Silicon)
curl -L https://github.com/Bacchusgift/sql-diff/releases/download/v1.0.0/sql-diff-darwin-arm64 -o sql-diff

# Linux
curl -L https://github.com/Bacchusgift/sql-diff/releases/download/v1.0.0/sql-diff-linux-amd64 -o sql-diff

# 添加执行权限
chmod +x sql-diff

# 移动到系统路径
sudo mv sql-diff /usr/local/bin/
```

## 验证安装

安装完成后,运行以下命令验证:

```bash
# 查看版本
sql-diff --version

# 查看帮助
sql-diff --help
```

## 配置 AI 功能 (可选)

如果您需要使用 AI 增强功能,需要配置 DeepSeek API Key:

```bash
sql-diff config \
  --ai-enabled=true \
  --api-key="your-api-key" \
  --provider=deepseek
```

然后在您的 shell 配置文件中应用环境变量:

```bash
# 运行配置命令并保存输出
sql-diff config --ai-enabled=true --api-key="sk-xxx" > ~/.sql-diff-env

# 在 ~/.bashrc 或 ~/.zshrc 中添加
echo "source ~/.sql-diff-env" >> ~/.zshrc
source ~/.zshrc
```

详细配置说明请参考 [环境变量配置](/config/environment)。

## 常见问题

### 权限错误

如果遇到权限错误:

```bash
# macOS/Linux
sudo make install
```

### 找不到命令

确保安装路径在系统 PATH 中:

```bash
# 查看 GOPATH
go env GOPATH

# 添加到 PATH (添加到 ~/.bashrc 或 ~/.zshrc)
export PATH=$PATH:$(go env GOPATH)/bin
```

### 网络问题

如果无法访问 GitHub:

```bash
# 使用代理
export GOPROXY=https://goproxy.cn,direct
go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
```

## 下一步

- [快速开始](/guide/getting-started) - 学习基本使用
- [配置](/config/environment) - 配置 AI 功能
- [命令行工具](/guide/cli) - 了解所有命令
