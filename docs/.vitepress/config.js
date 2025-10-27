import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/sql-diff/',
  title: 'SQL-Diff',
  description: '基于 AST 的智能 SQL 表结构比对工具',
  lang: 'zh-CN',
  
  // 忽略死链接检查
  ignoreDeadLinks: [
    // 忽略本地开发服务器链接
    /^http:\/\/localhost/,
    /^https:\/\/localhost/,
  ],
  
  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#00ADD8' }]
  ],

  themeConfig: {
    logo: '/logo.svg',
    
    nav: [
      { text: '指南', link: '/guide/getting-started' },
      { text: '配置', link: '/config/environment' },
      { text: 'AI 功能', link: '/ai/guide' },
      { text: '示例', link: '/examples/basic' },
      { 
        text: '更多',
        items: [
          { text: '架构设计', link: '/architecture' },
          { text: '贡献指南', link: '/contributing' },
          { text: 'GitHub', link: 'https://github.com/Bacchusgift/sql-diff' }
        ]
      }
    ],

    sidebar: {
      '/guide/': [
        {
          text: '开始',
          items: [
            { text: '介绍', link: '/guide/introduction' },
            { text: '快速开始', link: '/guide/getting-started' },
            { text: '安装', link: '/guide/installation' }
          ]
        },
        {
          text: '核心功能',
          items: [
            { text: '交互式光标选择', link: '/guide/interactive' },
            { text: '表结构比对', link: '/guide/comparison' },
            { text: 'DDL 生成', link: '/guide/ddl-generation' },
            { text: '命令行工具', link: '/guide/cli' }
          ]
        }
      ],
      '/config/': [
        {
          text: '配置',
          items: [
            { text: '环境变量', link: '/config/environment' },
            { text: '配置文件', link: '/config/file' },
            { text: '配置命令', link: '/config/command' }
          ]
        }
      ],
      '/ai/': [
        {
          text: 'AI 功能',
          items: [
            { text: 'AI 指南', link: '/ai/guide' },
            { text: 'AI 生成 SQL', link: '/ai/sql-generation' },
            { text: 'DeepSeek 集成', link: '/ai/deepseek' },
            { text: '最佳实践', link: '/ai/best-practices' }
          ]
        }
      ],
      '/examples/': [
        {
          text: '示例',
          items: [
            { text: '基础用法', link: '/examples/basic' },
            { text: '复杂场景', link: '/examples/advanced' },
            { text: 'CI/CD 集成', link: '/examples/ci-cd' }
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/Bacchusgift/sql-diff' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025 SQL-Diff Contributors'
    },

    search: {
      provider: 'local'
    },

    editLink: {
      pattern: 'https://github.com/Bacchusgift/sql-diff/edit/main/:path',
      text: '在 GitHub 上编辑此页'
    },

    lastUpdated: {
      text: '最后更新',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'medium'
      }
    }
  }
})
