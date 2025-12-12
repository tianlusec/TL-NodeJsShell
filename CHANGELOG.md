# 更新日志

[English](CHANGELOG_EN.md) | 简体中文

本文件记录项目的所有重要更改。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/spec/v2.0.0.html)。

## [未发布]

### 计划中
- WebSocket 支持实时通信
- 插件系统以提供扩展性
- 多用户支持与身份验证
- 增强的日志和审计跟踪
- Docker 部署支持

## [1.0.0] - 2024-12-12

### 新增
- TL-NodeJsShell 首次发布
- 内存马注入功能
  - Express 中间件注入
  - Koa 中间件注入
  - 原型链污染技术
- 多种编码支持（Base64、XOR、AES）
- 基于 xterm.js 的交互式虚拟终端
- 综合文件管理系统- 文件浏览器与目录导航
  - 支持分块传输的上传/下载
  - 基于 Monaco 编辑器的文件预览和编辑
- 代理支持（HTTP/HTTPS/SOCKS5）
- 自定义 HTTP 请求头配置
- 实时 Shell 状态监控
- 系统信息收集
- 命令历史记录
- 现代化的 Vue 3 + TypeScript 前端
- 基于 Gin 框架的 Go 后端
- SQLite 数据库持久化
- RESTful API 设计
- 基于 Element Plus 的响应式 UI

### 安全性
- Shell 密码保护
- 多种编码方法用于 Payload 混淆
- 代理支持以提供匿名性
- 输入验证和清理

### 文档
- 包含英文和中文版本的综合 README
- API 文档
- 贡献指南
- 行为准则
- MIT 许可证

---

[未发布]: https://github.com/tianlusec/TL-NodeJsShell/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/tianlusec/TL-NodeJsShell/releases/tag/v1.0.0