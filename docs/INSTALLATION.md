# Installation Guide

[English](#english) | [中文](#中文)

---

## English

### System Requirements

#### Minimum Requirements
- **Operating System**: Windows 10/11, Linux (Ubuntu 20.04+, CentOS 8+), macOS 10.15+
- **Go**: Version 1.21 or higher
- **Node.js**: Version 16 or higher
- **RAM**: 2GB minimum, 4GB recommended
- **Disk Space**: 500MB for application and dependencies

#### Recommended Requirements
- **RAM**: 8GB or more
- **CPU**: Multi-core processor
- **Network**: Stable internet connection for package downloads

### Prerequisites

#### 1. Install Go

**Linux/macOS:**
```bash
# Download and install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

**Windows:**
1. Download installer from [https://go.dev/dl/](https://go.dev/dl/)
2. Run the installer
3. Verify: Open Command Prompt and run `go version`

#### 2. Install Node.js

**Linux (Ubuntu/Debian):**
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

**macOS:**
```bash
brew install node
```

**Windows:**
1. Download installer from [https://nodejs.org/](https://nodejs.org/)
2. Run the installer
3. Verify: `node --version` and `npm --version`

#### 3. Install Git

**Linux:**
```bash
sudo apt-get install git # Ubuntu/Debian
sudo yum install git  # CentOS/RHEL
```

**macOS:**
```bash
brew install git
```

**Windows:**
Download from [https://git-scm.com/download/win](https://git-scm.com/download/win)

---

### Installation Methods

#### Method 1: From Source (Recommended)

**Step 1: Clone the Repository**
```bash
git clone https://github.com/tianlusec/TL-NodeJsShell.git
cd TL-NodeJsShell
```

**Step 2: Build Backend**
```bash
cd backend
go mod download
go build -o NodeJsshell main.go

# For Windows
go build -o NodeJsshell.exe main.go
```

**Step 3: Build Frontend**
```bash
cd ../frontend
npm install
npm run build
```

**Step 4: Run the Application**
```bash
# From backend directory
cd ../backend
./NodeJsshell

# For Windows
NodeJsshell.exe
```

**Step 5: Access the Application**
Open your browser and navigate to:
```
http://localhost:8080
```

#### Method 2: Development Mode

For development with hot-reload:

**Terminal 1 - Backend:**
```bash
cd backend
go run main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

Access at: `http://localhost:5173` (Vite dev server)

#### Method 3: Docker (Coming Soon)

Docker support will be added in future releases.

---

### Configuration

#### Backend Configuration

Edit [`backend/config/config.go`](../backend/config/config.go):

```go
type Config struct {
 Port string // Default: "8080"
 Host string // Default: "0.0.0.0"
}
```

**Custom Port:**
```go
func Load() *Config {
 return &Config{
  Port: "3000",  // Change port
  Host: "127.0.0.1", // Bind to localhost only
 }
}
```

#### Frontend Configuration

Edit [`frontend/vite.config.ts`](../frontend/vite.config.ts) for proxy settings:

```typescript
export default defineConfig({
 server: {
 proxy: {
  '/api': {
  target: 'http://localhost:8080',
  changeOrigin: true
  }
 }
 }
})
```

---

### Troubleshooting

#### Common Issues

**1. Port Already in Use**
```
Error: listen tcp :8080: bind: address already in use
```
**Solution:**
- Change the port in config
- Or kill the process using the port:
 ```bash
 # Linux/macOS
 lsof -ti:8080 | xargs kill -9
 
 # Windows
 netstat -ano | findstr :8080
 taskkill /PID <PID> /F
 ```

**2. Go Module Download Fails**
```
Error: go: module lookup disabled
```
**Solution:**
```bash
export GOPROXY=https://goproxy.io,direct
go mod download
```

**3. npm Install Fails**
```
Error: EACCES: permission denied
```
**Solution:**
```bash
# Linux/macOS
sudo npm install -g npm
npm install --unsafe-perm

# Or use nvm (recommended)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
```

**4. Frontend Build Fails**
```
Error: JavaScript heap out of memory
```
**Solution:**
```bash
export NODE_OPTIONS="--max-old-space-size=4096"
npm run build
```

**5. Database Permission Error**
```
Error: unable to open database file
```
**Solution:**
```bash
# Ensure write permissions
chmod 755 backend/
chmod 644 backend/*.db
```

---

### Platform-Specific Notes

#### Linux

**SELinux Issues:**
```bash
# Temporarily disable
sudo setenforce 0

# Or add exception
sudo chcon -R -t httpd_sys_content_t /path/to/TL-NodeJsShell
```

**Firewall:**
```bash
# UFW
sudo ufw allow 8080/tcp

# firewalld
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

#### macOS

**Gatekeeper Warning:**
```bash
# Allow unsigned binary
xattr -d com.apple.quarantine NodeJsshell
```

#### Windows

**Windows Defender:**
- Add exception for the application directory
- Or temporarily disable real-time protection during installation

**PowerShell Execution Policy:**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

---

### Verification

After installation, verify everything works:

```bash
# Check backend
curl http://localhost:8080/api/shells

# Expected response: []
```

---

### Updating

To update to the latest version:

```bash
cd TL-NodeJsShell
git pull origin main

# Rebuild backend
cd backend
go build -o NodeJsshell main.go

# Rebuild frontend
cd ../frontend
npm install
npm run build
```

---

### Uninstallation

To completely remove TL-NodeJsShell:

```bash
# Stop the application
# Then remove directory
cd ..
rm -rf TL-NodeJsShell

# Remove database (if needed)
rm -f ~/.TL-NodeJsShell/*.db
```

---

## 中文

### 系统要求

#### 最低要求
- **操作系统**: Windows 10/11, Linux (Ubuntu 20.04+, CentOS 8+), macOS 10.15+
- **Go**: 1.21 或更高版本
- **Node.js**: 16 或更高版本
- **内存**: 最低 2GB，推荐 4GB
- **磁盘空间**: 应用程序和依赖项需要 500MB

#### 推荐配置
- **内存**: 8GB 或更多
- **CPU**: 多核处理器
- **网络**: 稳定的互联网连接以下载包

### 前置条件

#### 1. 安装 Go

**Linux/macOS:**
```bash
# 下载并安装 Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 添加到 PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

**Windows:**
1. 从 [https://go.dev/dl/](https://go.dev/dl/) 下载安装程序
2. 运行安装程序
3. 验证：打开命令提示符并运行 `go version`

#### 2. 安装 Node.js

**Linux (Ubuntu/Debian):**
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

**macOS:**
```bash
brew install node
```

**Windows:**
1. 从 [https://nodejs.org/](https://nodejs.org/) 下载安装程序
2. 运行安装程序
3. 验证：`node --version` 和 `npm --version`

#### 3. 安装 Git

**Linux:**
```bash
sudo apt-get install git # Ubuntu/Debian
sudo yum install git  # CentOS/RHEL
```

**macOS:**
```bash
brew install git
```

**Windows:**
从 [https://git-scm.com/download/win](https://git-scm.com/download/win) 下载

---

### 安装方法

#### 方法 1：从源码安装（推荐）

**步骤 1：克隆仓库**
```bash
git clone https://github.com/tianlusec/TL-NodeJsShell.git
cd TL-NodeJsShell
```

**步骤 2：构建后端**
```bash
cd backend
go mod download
go build -o NodeJsshell main.go

# Windows 系统
go build -o NodeJsshell.exe main.go
```

**步骤 3：构建前端**
```bash
cd ../frontend
npm install
npm run build
```

**步骤 4：运行应用程序**
```bash
# 从 backend 目录
cd ../backend
./NodeJsshell

# Windows 系统
NodeJsshell.exe
```

**步骤 5：访问应用程序**
在浏览器中打开：
```
http://localhost:8080
```

#### 方法 2：开发模式

用于开发的热重载模式：

**终端 1 - 后端：**
```bash
cd backend
go run main.go
```

**终端 2 - 前端：**
```bash
cd frontend
npm run dev
```

访问：`http://localhost:5173`（Vite 开发服务器）

#### 方法 3：Docker（即将推出）

Docker 支持将在未来版本中添加。

---

### 配置

#### 后端配置

编辑 [`backend/config/config.go`](../backend/config/config.go)：

```go
type Config struct {
 Port string // 默认: "8080"
 Host string // 默认: "0.0.0.0"
}
```

**自定义端口：**
```go
func Load() *Config {
 return &Config{
  Port: "3000",  // 更改端口
  Host: "127.0.0.1", // 仅绑定到本地主机
 }
}
```

---

### 故障排除

#### 常见问题

**1. 端口已被占用**
```
Error: listen tcp :8080: bind: address already in use
```
**解决方案：**
- 在配置中更改端口
- 或终止占用端口的进程：
 ```bash
 # Linux/macOS
 lsof -ti:8080 | xargs kill -9
 
 # Windows
 netstat -ano | findstr :8080
 taskkill /PID <PID> /F
 ```

**2. Go 模块下载失败**
```
Error: go: module lookup disabled
```
**解决方案：**
```bash
export GOPROXY=https://goproxy.io,direct
go mod download
```

**3. npm 安装失败**
```
Error: EACCES: permission denied
```
**解决方案：**
```bash
# Linux/macOS
sudo npm install -g npm
npm install --unsafe-perm

# 或使用 nvm（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
```

---

### 验证

安装后，验证一切正常：

```bash
# 检查后端
curl http://localhost:8080/api/shells

# 预期响应: []
```

---

### 更新

更新到最新版本：

```bash
cd TL-NodeJsShell
git pull origin main

# 重新构建后端
cd backend
go build -o NodeJsshell main.go

# 重新构建前端
cd ../frontend
npm install
npm run build
```

---

### 卸载

完全删除 TL-NodeJsShell：

```bash
# 停止应用程序
# 然后删除目录
cd ..
rm -rf TL-NodeJsShell

# 删除数据库（如需要）
rm -f ~/.TL-NodeJsShell/*.db