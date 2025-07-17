import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Erlang Rebar Config Parser',
  description: 'A Go library for parsing Erlang rebar configuration files',
  
  // Base URL for GitHub Pages
  base: '/erlang-rebar-config-parser/',
  
  // Language configuration
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'Erlang Rebar Config Parser',
      description: 'A Go library for parsing Erlang rebar configuration files',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/guide/getting-started' },
          { text: 'API Reference', link: '/api/' },
          { text: 'Examples', link: '/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/erlang-rebar-config-parser' }
        ],
        sidebar: {
          '/guide/': [
            {
              text: 'Guide',
              items: [
                { text: 'Getting Started', link: '/guide/getting-started' },
                { text: 'Installation', link: '/guide/installation' },
                { text: 'Basic Usage', link: '/guide/basic-usage' },
                { text: 'Advanced Usage', link: '/guide/advanced-usage' }
              ]
            }
          ],
          '/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Overview', link: '/api/' },
                { text: 'Core Functions', link: '/api/core-functions' },
                { text: 'Types', link: '/api/types' },
                { text: 'RebarConfig Methods', link: '/api/rebar-config' },
                { text: 'Term Interface', link: '/api/term-interface' }
              ]
            }
          ],
          '/examples/': [
            {
              text: 'Examples',
              items: [
                { text: 'Overview', link: '/examples/' },
                { text: 'Basic Parsing', link: '/examples/basic-parsing' },
                { text: 'Configuration Access', link: '/examples/config-access' },
                { text: 'Pretty Printing', link: '/examples/pretty-printing' },
                { text: 'Term Comparison', link: '/examples/comparison' },
                { text: 'Complex Analysis', link: '/examples/complex-analysis' }
              ]
            }
          ]
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Erlang Rebar 配置解析器',
      description: '用于解析 Erlang rebar 配置文件的 Go 库',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '指南', link: '/zh/guide/getting-started' },
          { text: 'API 参考', link: '/zh/api/' },
          { text: '示例', link: '/zh/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/erlang-rebar-config-parser' }
        ],
        sidebar: {
          '/zh/guide/': [
            {
              text: '指南',
              items: [
                { text: '快速开始', link: '/zh/guide/getting-started' },
                { text: '安装', link: '/zh/guide/installation' },
                { text: '基本用法', link: '/zh/guide/basic-usage' },
                { text: '高级用法', link: '/zh/guide/advanced-usage' }
              ]
            }
          ],
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '概述', link: '/zh/api/' },
                { text: '核心函数', link: '/zh/api/core-functions' },
                { text: '类型定义', link: '/zh/api/types' },
                { text: 'RebarConfig 方法', link: '/zh/api/rebar-config' },
                { text: 'Term 接口', link: '/zh/api/term-interface' }
              ]
            }
          ],
          '/zh/examples/': [
            {
              text: '示例',
              items: [
                { text: '概述', link: '/zh/examples/' },
                { text: '基本解析', link: '/zh/examples/basic-parsing' },
                { text: '配置访问', link: '/zh/examples/config-access' },
                { text: '美化输出', link: '/zh/examples/pretty-printing' },
                { text: '术语比较', link: '/zh/examples/comparison' },
                { text: '复杂分析', link: '/zh/examples/complex-analysis' }
              ]
            }
          ]
        }
      }
    }
  },
  
  themeConfig: {
    logo: '/logo.svg',
    socialLinks: [
      { icon: 'github', link: 'https://github.com/scagogogo/erlang-rebar-config-parser' }
    ],
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024 scagogogo'
    }
  }
})
