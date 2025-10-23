# GitHub Pages 部署配置指南

## 🎉 代码已成功推送到 GitHub!

仓库地址: https://github.com/Bacchusgift/sql-diff

## 📝 配置 GitHub Pages 步骤

### 1. 进入仓库设置

访问: https://github.com/Bacchusgift/sql-diff/settings

### 2. 配置 Pages

1. 在左侧菜单中找到并点击 **Pages**
2. 在 **Source** 部分:
   - 选择 **GitHub Actions** (而不是 Deploy from a branch)
3. 点击 **Save** 保存设置

### 3. 触发部署

配置完成后,GitHub Actions 会自动部署文档网站。

**方式1**: 等待自动触发
- 下次推送到 `main` 分支时会自动部署

**方式2**: 手动触发
1. 访问: https://github.com/Bacchusgift/sql-diff/actions
2. 点击左侧的 **Deploy Documentation** workflow
3. 点击右侧的 **Run workflow** 按钮
4. 选择 `main` 分支,点击绿色的 **Run workflow** 按钮

### 4. 查看部署状态

1. 访问: https://github.com/Bacchusgift/sql-diff/actions
2. 查看 **Deploy Documentation** workflow 的运行状态
3. 等待所有步骤完成 (显示绿色✓)

### 5. 访问文档网站

部署成功后,文档网站将在以下地址可访问:

**https://bacchusgift.github.io/sql-diff/**

## ✅ 验证部署

访问文档网站后,您应该能看到:

- 🏠 首页 with Hero 区域
- 📚 完整的导航菜单
- 📖 所有文档页面
- 🔍 搜索功能
- 🌓 深色模式切换

## 🔧 故障排查

### ✅ 已解决: Dependencies lock file 错误

**问题**: 
```
Error: Dependencies lock file is not found in /home/runner/work/sql-diff/sql-diff. 
Supported file patterns: package-lock.json,npm-shrinkwrap.json,yarn.lock
```

**解决方案**:
- ✅ 已添加 `package-lock.json` 到版本控制
- ✅ 已从 `.gitignore` 中移除 `package-lock.json`
- ✅ 已将 GitHub Actions 中的 `npm ci` 改为 `npm install`

**状态**: 已修复并推送到 GitHub

### 如果 Actions 失败

1. 检查 **Actions** 权限:
   - 访问: https://github.com/Bacchusgift/sql-diff/settings/actions
   - 确保 **Workflow permissions** 设置为:
     - ✅ Read and write permissions
     - ✅ Allow GitHub Actions to create and approve pull requests

2. 检查 **Pages** 权限:
   - 访问: https://github.com/Bacchusgift/sql-diff/settings/pages
   - 确保 **Source** 设置为 **GitHub Actions**

### 如果页面 404

1. 确认部署成功 (Actions 显示绿色✓)
2. 等待几分钟让 GitHub Pages 生效
3. 清除浏览器缓存并重新访问
4. 检查仓库设置中 Pages 是否启用

### 如果样式丢失

这通常是 base path 问题,已在配置中设置:
```javascript
// .vitepress/config.js
base: '/sql-diff/'
```

## 📋 已完成的配置

✅ 所有 GitHub 链接已更新为 `Bacchusgift/sql-diff`
✅ VitePress 配置了正确的 base path
✅ GitHub Actions workflow 已创建
✅ 代码已推送到 GitHub
✅ Git 仓库已初始化

## 🚀 下一步

1. **配置 GitHub Pages** (按上述步骤)
2. **等待部署完成** (约 2-3 分钟)
3. **访问文档网站**: https://bacchusgift.github.io/sql-diff/
4. **分享给团队** 🎉

## 📚 相关文档

- [GitHub Pages 官方文档](https://docs.github.com/en/pages)
- [VitePress 部署指南](https://vitepress.dev/guide/deploy#github-pages)
- [GitHub Actions 文档](https://docs.github.com/en/actions)

## 💡 提示

- 每次推送到 `main` 分支都会自动重新部署文档
- 可以在 Actions 页面查看部署历史和日志
- 部署通常在 2-3 分钟内完成

---

**文档网站即将上线!** 🎊

配置完成后,您的项目将拥有一个专业的在线文档网站!
