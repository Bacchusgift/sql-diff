# 版本信息功能说明

## ✨ 新增功能

sql-diff 现在支持完整的版本信息显示功能，提供多种方式查看版本：

### 1. `--version` 标志（简洁输出）

```bash
$ sql-diff --version
v1.0.2
```

**特点**：
- ✅ 只显示版本号，适合脚本使用
- ✅ 符合 CLI 工具标准
- ✅ 输出干净，易于解析

### 2. `-v` 标志（Cobra 默认）

```bash
$ sql-diff -v
v1.0.2
```

**说明**：
- Cobra 框架默认支持 `-v` 作为 `--version` 的短标志
- 与 `--version` 输出相同

### 3. `version` 命令（详细信息）

```bash
$ sql-diff version
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       SQL-Diff 版本信息
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

版本号:     v1.0.2
构建时间:   2025-10-24_06:00:47
Git提交:    2e857fd

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**特点**：
- ✅ 显示完整的版本信息
- ✅ 包含构建时间和 Git 提交哈希
- ✅ 美观的彩色输出
- ✅ 适合调试和问题报告

### 4. `v` 命令（别名）

```bash
$ sql-diff v
# 输出同 version 命令
```

**特点**：
- ✅ `version` 命令的快捷别名
- ✅ 更快速的输入
- ✅ 显示详细信息

## 📊 对比

| 命令 | 输出内容 | 适用场景 |
|------|---------|---------|
| `--version` | 仅版本号 | 脚本使用、版本检查 |
| `-v` | 仅版本号 | 同 `--version` |
| `version` | 详细信息 | 调试、问题报告、了解构建信息 |
| `v` | 详细信息 | 快速查看详细信息 |

## 🔧 技术实现

### 版本信息注入

版本信息通过编译时的 `-ldflags` 注入：

```makefile
# Makefile
LDFLAGS=-ldflags "-s -w \
  -X main.Version=$(VERSION) \
  -X main.BuildTime=$(BUILD_TIME) \
  -X main.GitCommit=$(GIT_COMMIT)"
```

### 变量定义

在 `cmd/sql-diff/main.go` 中定义：

```go
var (
    Version   = "dev"      // 版本号
    BuildTime = "unknown"  // 构建时间
    GitCommit = "unknown"  // Git 提交哈希
)
```

### 传递给 cmd 包

```go
func main() {
    cmd.SetVersion(Version, BuildTime, GitCommit)
    cmd.Execute()
}
```

### 在 cmd 包中使用

```go
// internal/cmd/root.go
var (
    version   = "dev"
    buildTime = "unknown"
    gitCommit = "unknown"
)

func SetVersion(v, bt, gc string) {
    if v != "" {
        version = v
        rootCmd.Version = v  // 更新 rootCmd
    }
    if bt != "" {
        buildTime = bt
    }
    if gc != "" {
        gitCommit = gc
    }
}
```

### version 命令实现

```go
// internal/cmd/version.go
var versionCmd = &cobra.Command{
    Use:     "version",
    Aliases: []string{"v"},
    Short:   "显示版本信息",
    Run: func(cmd *cobra.Command, args []string) {
        printVersion()
    },
}

func printVersion() {
    // 彩色输出详细版本信息
    fmt.Printf("版本号:     %s\n", version)
    fmt.Printf("构建时间:   %s\n", buildTime)
    fmt.Printf("Git提交:    %s\n", gitCommit)
}
```

## 🎯 使用示例

### 开发阶段

```bash
# 使用 make build 编译（自动注入版本信息）
$ make build
✓ 编译完成: bin/sql-diff
✓ 版本: v1.0.2-7-g2e857fd-dirty

# 查看版本
$ ./bin/sql-diff version
版本号:     v1.0.2-7-g2e857fd-dirty
构建时间:   2025-10-24_06:00:47
Git提交:    2e857fd
```

**说明**：
- `v1.0.2`：最近的 tag
- `-7-g2e857fd`：距离 tag 有 7 个提交
- `-dirty`：有未提交的修改

### 生产版本

```bash
# 发布版本（通过 tag 触发）
$ git tag v1.1.0
$ git push origin v1.1.0

# GitHub Actions 自动构建
# 用户安装后查看版本
$ sql-diff --version
v1.1.0
```

### 脚本中使用

```bash
#!/bin/bash

# 检查 sql-diff 版本
VERSION=$(sql-diff --version)

if [[ "$VERSION" < "v1.0.0" ]]; then
    echo "需要 sql-diff >= v1.0.0"
    exit 1
fi

echo "当前版本: $VERSION"
```

## 📝 更新日志

### 2025-10-24

**新增**：
- ✅ 添加 `--version` 和 `-v` 标志
- ✅ 添加 `version` 命令
- ✅ 添加 `v` 命令别名
- ✅ 支持构建时注入版本信息
- ✅ 美化版本信息输出

**技术细节**：
- 使用 Cobra 的 `Version` 字段支持 `--version`
- 自定义 `SetVersionTemplate` 简化输出
- 创建独立的 `version` 命令提供详细信息

## 🔍 调试技巧

### 查看编译注入的版本信息

```bash
# 查看二进制文件中的版本字符串
$ strings bin/sql-diff | grep -E "v[0-9]+\.[0-9]+\.[0-9]+"
```

### 验证版本信息是否正确注入

```bash
# 使用 make build 编译
$ make build

# 检查所有版本相关命令
$ ./bin/sql-diff --version   # 应显示版本号
$ ./bin/sql-diff version     # 应显示详细信息
```

## 🚀 最佳实践

1. **开发时**：使用 `make build` 编译，自动注入 Git 信息
2. **发布时**：通过 tag 触发 CI/CD，确保版本号一致
3. **脚本中**：使用 `--version` 获取简洁输出
4. **调试时**：使用 `version` 命令查看完整信息
5. **问题报告**：提供 `sql-diff version` 的完整输出

## 📚 相关文件

- `cmd/sql-diff/main.go` - 版本变量定义
- `internal/cmd/root.go` - 版本信息设置和 rootCmd 配置
- `internal/cmd/version.go` - version 命令实现
- `Makefile` - 版本信息注入配置

## 🎉 总结

现在 sql-diff 拥有完整的版本信息功能：

- ✅ **标准化**：符合 CLI 工具的标准实践
- ✅ **灵活性**：提供简洁和详细两种输出方式
- ✅ **易用性**：支持多种命令和别名
- ✅ **可追溯**：包含构建时间和 Git 提交信息
- ✅ **自动化**：通过 Makefile 和 CI/CD 自动注入

这使得版本管理、问题调试和用户支持都变得更加容易！
