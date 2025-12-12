# 项目结构

[English](PROJECT_STRUCTURE_EN.md) | 简体中文

本文档提供 TL-NodeJsShell 项目结构的全面概述。

## 目录树

```
TL-NodeJsShell/
├── .github/                         # GitHub 配置与社区文档
│   ├── ISSUE_TEMPLATE/              # Issue 模板
│   ├── PULL_REQUEST_TEMPLATE.md     # PR 模板
│   ├── CODE_OF_CONDUCT.md           # 行为准则
│   └── CONTRIBUTING.md              # 贡献指南
│
├── docs/                            # 项目文档
│   ├── images/                      # 文档图片资源
│   ├── API.md                       # API 文档
│   ├── INSTALLATION.md              # 安装指南
│   └── SECURITY.md                  # 安全策略
│
├── server/                          # Go 后端服务器
│   ├── cmd/                         # 应用程序入口
│   │   └── api/
│   │       └── main.go              # 主程序入口
│   │
│   ├── internal/                    # 内部应用代码
│   │   ├── app/                     # 应用核心逻辑
│   │   │   ├── middleware/          # HTTP 中间件
│   │   │   └── routes/              # 路由定义
│   │   │
│   │   ├── config/                  # 配置管理
│   │   │
│   │   ├── core/                    # 核心业务逻辑
│   │   │   ├── crypto/              # 加密模块
│   │   │   ├── exploit/             # 漏洞利用模块
│   │   │   ├── payload/             # Payload 生成
│   │   │   ├── proxy/               # 代理功能
│   │   │   └── transport/           # 通信传输
│   │   │
│   │   ├── database/                # 数据库操作
│   │   │
│   │   └── handlers/                # HTTP 请求处理器
│   │
│   └── go.mod                       # Go 模块定义
│
├── web/                             # Vue 前端应用
│   ├── src/                         # 源代码
│   │   ├── api/                     # API 接口
│   │   ├── components/              # Vue 组件
│   │   ├── router/                  # 路由配置
│   │   ├── stores/                  # 状态管理 (Pinia)
│   │   ├── types/                   # TypeScript 类型定义
│   │   └── views/                   # 页面视图
│   │
│   ├── index.html                   # HTML 入口
│   ├── package.json                 # 依赖配置
│   ├── tsconfig.json                # TypeScript 配置
│   └── vite.config.ts               # Vite 构建配置
│
├── CHANGELOG.md                     # 更新日志
├── LICENSE                          # 许可证
└── README.md                        # 项目说明
```
