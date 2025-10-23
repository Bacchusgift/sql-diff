# SQL-Diff 文档网站

本项目使用 [VitePress](https://vitepress.dev) 构建文档网站。

## 快速开始

### 安装依赖

```bash
npm install
```

### 本地开发

启动开发服务器:

```bash
npm run docs:dev
```

然后访问 http://localhost:5173

### 构建生产版本

```bash
npm run docs:build
```

构建产物在 `.vitepress/dist` 目录。

### 预览生产版本

```bash
npm run docs:preview
```

## 文档结构

```
.
├── index.md                    # 首页
├── .vitepress/
│   ├── config.js              # VitePress 配置
│   └── theme/
│       ├── index.js           # 自定义主题
│       └── style.css          # 自定义样式
├── guide/                     # 指南
│   ├── introduction.md        # 介绍
│   ├── getting-started.md     # 快速开始
│   ├── installation.md        # 安装
│   ├── comparison.md          # 表结构比对
│   ├── ddl-generation.md      # DDL 生成
│   └── cli.md                 # 命令行工具
├── config/                    # 配置
│   ├── environment.md         # 环境变量
│   ├── file.md               # 配置文件
│   └── command.md            # 配置命令
├── ai/                        # AI 功能
│   ├── guide.md              # AI 指南
│   ├── deepseek.md           # DeepSeek 集成
│   └── best-practices.md     # 最佳实践
├── examples/                  # 示例
│   ├── basic.md              # 基础示例
│   ├── advanced.md           # 复杂场景
│   └── ci-cd.md              # CI/CD 集成
└── architecture.md           # 架构设计
```

## 添加新页面

1. 在相应目录下创建 `.md` 文件
2. 在 `.vitepress/config.js` 中添加到导航或侧边栏:

```javascript
// 添加到侧边栏
sidebar: {
  '/guide/': [
    {
      text: '开始',
      items: [
        { text: '新页面', link: '/guide/new-page' }
      ]
    }
  ]
}
```

## 自定义样式

在 `.vitepress/theme/style.css` 中添加自定义样式:

```css
:root {
  --vp-c-brand: #00ADD8;
}
```

## Markdown 增强

VitePress 支持多种 Markdown 扩展:

### 代码块

```bash
npm install
```

### 提示框

::: tip 提示
这是一个提示框
:::

::: warning 警告
这是一个警告框
:::

::: danger 危险
这是一个危险警告框
:::

### Mermaid 图表

```mermaid
graph LR
    A --> B
    B --> C
```

## 部署

### GitHub Pages

在 `.github/workflows/deploy-docs.yml`:

```yaml
name: Deploy Docs

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node
        uses: actions/setup-node@v3
        with:
          node-version: 18
      
      - name: Install dependencies
        run: npm install
      
      - name: Build
        run: npm run docs:build
      
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: .vitepress/dist
```

### Vercel

直接连接 GitHub 仓库,Vercel 会自动检测并部署。

### Netlify

创建 `netlify.toml`:

```toml
[build]
  command = "npm run docs:build"
  publish = ".vitepress/dist"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
```

## 相关链接

- [VitePress 官方文档](https://vitepress.dev)
- [Markdown 语法](https://markdown.com.cn/)
- [Vue 3 文档](https://cn.vuejs.org/)
