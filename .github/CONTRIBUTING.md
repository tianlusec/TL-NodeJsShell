# 贡献指南

[English](CONTRIBUTING_EN.md) | 简体中文

---

感谢您对 TL-NodeJsShell 项目的贡献兴趣！本文档提供了项目贡献指南。

## 行为准则

参与本项目即表示您同意为所有贡献者维护一个尊重和包容的环境。

## 如何贡献

### 报告错误

在创建错误报告之前，请检查现有问题以避免重复。

**报告错误时，请包含：**
- 清晰的描述性标题
- 重现问题的步骤
- 预期行为
- 实际行为
- 截图（如适用）
- 环境详情（操作系统、Go 版本、Node.js 版本）
- 错误消息或日志

### 建议改进

欢迎提出改进建议！请提供：
- 清晰的改进描述
- 使用场景和好处
- 可能的实现方法
- 任何相关示例或模型

### 拉取请求

1. **Fork 仓库**
   ```bash
   git clone https://github.com/YOUR_USERNAME/TL-NodeJsShell.git
   cd TL-NodeJsShell
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **进行更改**
   - 遵循编码标准
   - 编写清晰的提交消息
   - 如适用，添加测试
   - 更新文档

4. **测试更改**
   ```bash
   # 后端测试
   cd backend
   go test ./...
   
   # 前端测试
   cd frontend
   npm run test
   ```

5. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加您的功能描述"
   ```

6. **推送到您的 Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **创建拉取请求**
   - 提供清晰的更改描述
   - 引用任何相关问题
   - 确保所有测试通过

## 编码标准

### Go（后端）

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 指南
- 使用 `gofmt` 进行代码格式化
- 编写有意义的变量和函数名
- 为复杂逻辑添加注释
- 保持函数专注和简洁

**示例：**
```go
// GetShellByID 通过 ID 检索 shell 配置
func (h *ShellHandler) GetShellByID(id uint) (*database.Shell, error) {
    var shell database.Shell
    if err := h.db.First(&shell, id).Error; err != nil {
        return nil, fmt.Errorf("shell not found: %w", err)
    }
    return &shell, nil
}
```

### TypeScript/Vue（前端）

- 遵循 [Vue.js 风格指南](https://cn.vuejs.org/style-guide/)
- 使用 TypeScript 确保类型安全
- 使用组合式 API 编写 Vue 组件
- 遵循 ESLint 规则
- 编写自文档化代码

**示例：**
```typescript
// 定义清晰的接口
interface ShellConfig {id: number
  url: string
  password: string
  encodeType: string
}

// 使用组合式 API
const { shells, loading, error } = useShells()
```

### 提交消息约定

遵循 [约定式提交](https://www.conventionalcommits.org/zh-hans/)：

- `feat:` 新功能
- `fix:` 错误修复
- `docs:` 文档更改
- `style:` 代码样式更改（格式化等）
- `refactor:` 代码重构
- `test:` 添加或更新测试
- `chore:` 维护任务

**示例：**
```
feat: 添加代理认证支持
fix: 解决文件上传分块问题
docs: 更新 API 文档
refactor: 改进 shell 处理器的错误处理
```

## 项目结构

```
TL-NodeJsShell/
├── backend/
│   ├── app/              # 应用初始化
│   ├── config/           # 配置管理
│   ├── core/             # 核心业务逻辑
│   ├── database/         # 数据库模型和操作
│   ├── handlers/         # HTTP 请求处理器
│   └── main.go           # 入口点
├── frontend/
│   ├── src/
│   │   ├── api/         # API 客户端函数
│   │   ├── components/  # 可重用的 Vue 组件
│   │   ├── views/       # 页面组件
│   │   ├── stores/      # Pinia 状态管理
│   │   └── router/      # Vue Router 配置
│   └── public/          # 静态资源
└── docs/                # 文档
```

## 开发环境设置

### 后端开发

```bash
cd backend
go mod download
go run main.go
```

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

## 测试

### 后端测试

```bash
cd backend
go test ./... -v
go test -cover ./...
```

### 前端测试

```bash
cd frontend
npm run test
npm run test:coverage
```

## 文档

- 为面向用户的更改更新 README.md
- 在 docs/API.md 中更新 API 文档
- 为复杂代码添加内联注释
- 为重要更改更新 CHANGELOG.md

## 安全

- **永远不要提交敏感数据**（密码、密钥、令牌）
- 私下报告安全漏洞
- 遵循安全编码实践
- 验证所有用户输入

## 有问题？

欢迎：
- 开启 issue 进行讨论
- 加入我们的社区讨论
- 直接联系维护者

---

<div align="center">

**感谢您的贡献！**

</div>