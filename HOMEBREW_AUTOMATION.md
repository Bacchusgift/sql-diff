# 🍺 Homebrew 自动化发布配置指南

本文档说明如何配置 GitHub Actions 自动更新 Homebrew Tap。

## 📋 前置要求

1. ✅ 已创建 `Bacchusgift/homebrew-tap` 仓库
2. ✅ 已在 homebrew-tap 仓库中创建 `Formula/sql-diff.rb` 文件
3. ✅ 已配置 GitHub Personal Access Token

## 🔑 第一步：创建 GitHub Personal Access Token

### 1. 访问 GitHub 设置

访问：https://github.com/settings/tokens/new

### 2. 配置 Token

**Token 名称**: `HOMEBREW_TAP_TOKEN`

**权限选择** (Fine-grained token 推荐):
- ✅ Repository access: **Only select repositories**
  - 选择: `Bacchusgift/homebrew-tap`
- ✅ Repository permissions:
  - **Contents**: Read and write
  - **Metadata**: Read-only (自动包含)

**权限选择** (Classic token):
- ✅ `repo` (完整仓库访问权限)

**过期时间**: 建议选择 90 天或 1 年

### 3. 生成并保存 Token

点击 **Generate token**，复制生成的 token（只显示一次！）

## 🔧 第二步：配置 Repository Secret

### 1. 打开项目设置

访问：https://github.com/Bacchusgift/sql-diff/settings/secrets/actions

### 2. 添加 Secret

点击 **New repository secret**

- **Name**: `HOMEBREW_TAP_TOKEN`
- **Value**: 粘贴刚才复制的 token

点击 **Add secret**

## ✅ 第三步：验证配置

### 方式1: 创建测试 tag

```bash
# 创建并推送一个测试 tag
git tag v1.0.2-test
git push origin v1.0.2-test

# 删除测试 tag
git tag -d v1.0.2-test
git push origin :refs/tags/v1.0.2-test
```

### 方式2: 手动触发 workflow

1. 访问：https://github.com/Bacchusgift/sql-diff/actions
2. 选择 **Release** workflow
3. 查看最新的运行记录

## 🚀 完整发布流程

配置完成后，发布新版本变得非常简单：

```bash
# 1. 确保所有更改已提交
git add .
git commit -m "feat: 新功能"
git push origin main

# 2. 创建并推送 tag（自动触发一切！）
git tag v1.0.3
git push origin v1.0.3
```

**自动化流程会完成：**

1. ⚙️ **运行 CI 测试**
   - 代码格式检查
   - 代码质量检查
   - 单元测试
   - 多平台编译验证

2. 📦 **构建 Release**
   - 跨平台编译（6个平台）
   - 生成 SHA256 校验和
   - 创建 GitHub Release
   - 上传所有二进制文件

3. 🍺 **更新 Homebrew**
   - 自动更新 Formula 文件
   - 更新版本号和 SHA256
   - 自动提交并推送到 homebrew-tap
   - 用户可以立即 `brew upgrade sql-diff`

## 📊 工作流程图

```
推送 Tag (v1.0.3)
    ↓
┌───────────────────┐
│   运行 CI 测试    │ ← 代码质量检查
└─────────┬─────────┘
          ↓ (测试通过)
┌───────────────────┐
│  构建 & 发布      │ ← 跨平台编译
│  GitHub Release   │   创建 Release
└─────────┬─────────┘
          ↓
┌───────────────────┐
│ 更新 Homebrew Tap │ ← 自动更新 Formula
│   自动提交推送    │   用户可立即使用
└───────────────────┘
```

## 🔍 故障排查

### Token 权限不足

**错误信息**: `refusing to allow a Personal Access Token to create or update workflow`

**解决方案**:
1. 确保使用 Fine-grained token
2. 或者在 Classic token 中启用 `workflow` 权限

### Homebrew Tap 仓库不存在

**错误信息**: `repository not found`

**解决方案**:
```bash
# 1. 创建 homebrew-tap 仓库
# 访问 https://github.com/new
# 仓库名: homebrew-tap

# 2. 初始化仓库
mkdir homebrew-tap
cd homebrew-tap
git init
mkdir Formula
touch Formula/.gitkeep
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/Bacchusgift/homebrew-tap.git
git push -u origin main
```

### Formula 更新失败

**错误信息**: `nothing to commit`

**原因**: Formula 内容没有变化

**解决方案**: 这是正常的，说明 Formula 已经是最新的

## 🎯 高级配置

### 仅在主仓库启用自动更新

工作流已配置：只在 `Bacchusgift/sql-diff` 仓库启用自动更新

```yaml
if: github.repository == 'Bacchusgift/sql-diff'
```

### 使用默认 Token (可选)

如果不想创建 Personal Access Token，可以使用默认的 `GITHUB_TOKEN`：

```yaml
token: ${{ secrets.GITHUB_TOKEN }}
```

**限制**:
- 只能访问当前仓库
- 无法推送到其他仓库（如 homebrew-tap）
- **不推荐**用于跨仓库操作

## 📚 相关文档

- [GitHub Actions - Encrypted secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)

## 💡 最佳实践

1. **定期更新 Token**: 设置提醒在 Token 过期前更新
2. **测试环境隔离**: 使用测试 tag 验证工作流
3. **版本命名规范**: 始终使用 `v` 前缀（如 `v1.0.3`）
4. **发布前检查**: 确保本地测试通过再推送 tag
5. **监控工作流**: 订阅 Actions 邮件通知

## 🎉 完成！

配置完成后，您只需要：

```bash
git tag v1.0.3
git push origin v1.0.3
```

剩下的一切都会自动完成！✨
