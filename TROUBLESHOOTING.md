# 🔧 GitHub Actions 部署问题解决方案

## ✅ 已解决的问题

### 1. Dependencies lock file 错误

**完整错误信息**:
```
Error: Dependencies lock file is not found in /home/runner/work/sql-diff/sql-diff. 
Supported file patterns: package-lock.json,npm-shrinkwrap.json,yarn.lock
```

**问题原因**:
- GitHub Actions 的 `actions/setup-node@v4` 在使用 `cache: npm` 时,需要 `package-lock.json` 文件来缓存依赖
- 项目的 `.gitignore` 文件中忽略了 `package-lock.json`
- 使用 `npm ci` 命令但没有 lock 文件会导致失败

**解决步骤**:

1. ✅ **从 .gitignore 中移除 package-lock.json**
   ```diff
   # Node.js
   node_modules/
   - package-lock.json
   .vitepress/cache
   .vitepress/dist
   ```

2. ✅ **生成 package-lock.json**
   ```bash
   npm install
   ```

3. ✅ **修改 GitHub Actions workflow**
   ```yaml
   - name: Install dependencies
     run: npm install  # 改为 npm install (原来是 npm ci)
   ```

4. ✅ **提交并推送**
   ```bash
   git add package-lock.json .github/workflows/deploy-docs.yml .gitignore
   git commit -m "fix: add package-lock.json and update workflow"
   git push
   ```

**状态**: ✅ 已修复并推送到 GitHub

---

## 📋 GitHub Pages 部署检查清单

### 部署前检查

- [x] ✅ package-lock.json 已提交
- [x] ✅ .github/workflows/deploy-docs.yml 已配置
- [x] ✅ VitePress base 路径已设置 (`/sql-diff/`)
- [x] ✅ 所有代码已推送到 GitHub

### GitHub 设置检查

访问: https://github.com/Bacchusgift/sql-diff/settings

#### 1. Pages 设置
- [ ] 进入 Settings → Pages
- [ ] Source 选择: **GitHub Actions**
- [ ] 保存设置

#### 2. Actions 权限
- [ ] 进入 Settings → Actions → General
- [ ] Workflow permissions 设置为:
  - ✅ Read and write permissions
  - ✅ Allow GitHub Actions to create and approve pull requests

#### 3. 触发部署
- [ ] 访问 Actions 页面: https://github.com/Bacchusgift/sql-diff/actions
- [ ] 点击 "Deploy Documentation" workflow
- [ ] 点击 "Run workflow" 手动触发
- [ ] 等待部署完成 (约 2-3 分钟)

#### 4. 验证部署
- [ ] 检查 Actions 运行状态 (应该显示绿色 ✓)
- [ ] 访问文档网站: https://bacchusgift.github.io/sql-diff/
- [ ] 确认所有页面可以正常访问

---

## 🐛 其他常见问题

### 问题 2: 404 Not Found

**症状**: 访问 https://bacchusgift.github.io/sql-diff/ 显示 404

**可能原因**:
1. Pages 还未启用
2. 部署尚未完成
3. base 路径配置错误

**解决方案**:
1. 检查 Pages 设置是否正确
2. 等待 2-3 分钟让部署生效
3. 清除浏览器缓存
4. 检查 `.vitepress/config.js` 中的 `base: '/sql-diff/'` 是否正确

### 问题 3: 样式丢失或资源 404

**症状**: 页面可以访问但样式混乱,控制台显示资源 404

**原因**: VitePress base 路径配置错误

**解决方案**:
检查 `.vitepress/config.js`:
```javascript
export default defineConfig({
  base: '/sql-diff/',  // 必须与仓库名一致
  // ...
})
```

### 问题 4: Actions 权限错误

**症状**: 
```
Error: Resource not accessible by integration
```

**解决方案**:
1. 访问: https://github.com/Bacchusgift/sql-diff/settings/actions
2. 设置 Workflow permissions 为 "Read and write permissions"
3. 勾选 "Allow GitHub Actions to create and approve pull requests"
4. 保存并重新运行 workflow

### 问题 5: Node.js 版本问题

**症状**: 构建失败,提示 Node.js 版本不兼容

**解决方案**:
检查 `.github/workflows/deploy-docs.yml`:
```yaml
- name: Setup Node
  uses: actions/setup-node@v4
  with:
    node-version: 20  # 确保版本号正确
```

---

## 🔍 调试技巧

### 查看详细的构建日志

1. 访问: https://github.com/Bacchusgift/sql-diff/actions
2. 点击失败的 workflow 运行
3. 展开每个步骤查看详细日志
4. 查找红色的错误信息

### 本地测试构建

在推送前,先在本地测试构建:

```bash
# 安装依赖
npm install

# 构建文档
npm run docs:build

# 预览构建结果
npm run docs:preview
```

如果本地构建成功,GitHub Actions 也应该能成功。

### 检查文件是否正确推送

```bash
# 检查远程仓库的文件
git ls-tree -r --name-only HEAD

# 确保以下文件存在:
# - package-lock.json
# - .github/workflows/deploy-docs.yml
# - .vitepress/config.js
```

---

## 📞 获取帮助

如果问题仍未解决:

1. **查看 Actions 日志**: https://github.com/Bacchusgift/sql-diff/actions
2. **检查 Issues**: https://github.com/Bacchusgift/sql-diff/issues
3. **VitePress 文档**: https://vitepress.dev/guide/deploy#github-pages
4. **GitHub Pages 文档**: https://docs.github.com/en/pages

---

## ✅ 部署成功标志

当您看到以下内容时,说明部署成功:

1. ✅ GitHub Actions 显示绿色 ✓
2. ✅ Pages 页面显示部署时间和 URL
3. ✅ 访问 https://bacchusgift.github.io/sql-diff/ 可以看到文档首页
4. ✅ 导航、搜索、主题切换等功能正常工作

---

**当前状态**: 🟢 package-lock.json 问题已修复,可以重新运行 GitHub Actions 部署!
