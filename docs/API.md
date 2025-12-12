# API Documentation

[English](#english) | [中文](#中文)

---

## English

### Base URL

```
http://localhost:8080/api
```

### Authentication

Currently, the API does not require authentication. Future versions will include token-based authentication.

---

## Shell Management

### List All Shells

**Endpoint:** `GET /shells`

**Description:** Retrieve a list of all configured shells.

**Response:**
```json
[
 {
 "id": 1,
 "url": "http://example.com/shell.php",
 "password": "password123",
 "encode_type": "base64",
 "protocol": "multipart",
 "method": "POST",
 "group": "production",
 "name": "Web Server 1",
 "status": "online",
 "latency": 150,
 "last_active": "2024-12-12T10:30:00Z",
 "system_info": "OS: Linux...",
 "custom_headers": "{\"User-Agent\":\"Mozilla/5.0\"}",
 "proxy_id": 0
 }
]
```

### Get Shell Details

**Endpoint:** `GET /shells/:id`

**Description:** Retrieve details of a specific shell.

**Parameters:**
- `id` (path) - Shell ID

**Response:**
```json
{
 "id": 1,
 "url": "http://example.com/shell.php",
 "password": "password123",
 "encode_type": "base64",
 "status": "online",
 "latency": 150
}
```

### Create Shell

**Endpoint:** `POST /shells`

**Description:** Add a new shell configuration.

**Request Body:**
```json
{
 "url": "http://example.com/shell.php",
 "password": "password123",
 "encode_type": "base64",
 "method": "POST",
 "group": "production",
 "name": "Web Server 1",
 "proxy_id": 0,
 "custom_headers": {
 "User-Agent": "Mozilla/5.0"
 }
}
```

**Response:**
```json
{
 "id": 1,
 "url": "http://example.com/shell.php",
 "status": "offline",
 "created_at": "2024-12-12T10:30:00Z"
}
```

### Update Shell

**Endpoint:** `PUT /shells/:id`

**Description:** Update an existing shell configuration.

**Parameters:**
- `id` (path) - Shell ID

**Request Body:**
```json
{
 "url": "http://example.com/shell.php",
 "password": "newpassword",
 "encode_type": "xor",
 "name": "Updated Name"
}
```

### Delete Shell

**Endpoint:** `DELETE /shells/:id`

**Description:** Delete a shell configuration.

**Parameters:**
- `id` (path) - Shell ID

**Response:**
```json
{
 "message": "Shell deleted"
}
```

### Test Shell Connection

**Endpoint:** `POST /shells/:id/test`

**Description:** Test connectivity to a shell.

**Parameters:**
- `id` (path) - Shell ID

**Response:**
```json
{
 "success": true,
 "latency": 150,
 "status": "online"
}
```

### Get System Information

**Endpoint:** `GET /shells/:id/info`

**Description:** Retrieve system information from the target.

**Parameters:**
- `id` (path) - Shell ID
- `refresh` (query, optional) - Force refresh cached data

**Response:**
```json
{
 "system_info": {
 "ok": true,
 "platform": "linux",
 "arch": "x64",
 "hostname": "webserver01",
 "type": "Linux",
 "release": "5.4.0-42-generic",
 "uptime": 1234567,
 "totalmem": 8589934592,
 "freemem": 2147483648,
 "cpus": 4,
 "networkInterfaces": [
  {
  "interface": "eth0",
  "address": "192.168.1.100",
  "family": "IPv4",
  "internal": false
  }
 ]
 },
 "cached": false
}
```

---

## Command Execution

### Execute Command

**Endpoint:** `POST /shells/:id/execute`

**Description:** Execute a command on the target system.

**Parameters:**
- `id` (path) - Shell ID

**Request Body:**
```json
{
 "command": "ls -la"
}
```

**Response:**
```json
{
 "success": true,
 "data": "total 48\ndrwxr-xr-x 12 user user 4096 Dec 12 10:30 .\n...",
 "error": ""
}
```

---

## File Management

### List Directory

**Endpoint:** `POST /shells/:id/files/list`

**Description:** List files and directories in a path.

**Request Body:**
```json
{
 "path": "/var/www/html"
}
```

**Response:**
```json
{
 "success": true,
 "data": [
 {
  "name": "index.php",
  "type": "file",
  "size": 1024,
  "modified": "2024-12-12T10:30:00Z",
  "permissions": "rw-r--r--"
 }
 ]
}
```

### Read File

**Endpoint:** `POST /shells/:id/files/read`

**Description:** Read file contents.

**Request Body:**
```json
{
 "path": "/var/www/html/config.php"
}
```

**Response:**
```json
{
 "success": true,
 "data": "<?php\n$config = array();\n..."
}
```

### Upload File

**Endpoint:** `POST /shells/:id/files/upload`

**Description:** Upload a file to the target system.

**Request Body:**
```json
{
 "path": "/var/www/html/upload.txt",
 "content": "base64_encoded_content",
 "chunk_index": 0,
 "total_chunks": 1
}
```

**Response:**
```json
{
 "success": true,
 "message": "File uploaded successfully"
}
```

### Download File

**Endpoint:** `POST /shells/:id/files/download`

**Description:** Download a file from the target system.

**Request Body:**
```json
{
 "path": "/var/www/html/data.txt",
 "chunk_index": 0,
 "chunk_size": 1048576
}
```

**Response:**
```json
{
 "success": true,
 "data": "base64_encoded_content",
 "has_more": false
}
```

### Delete File

**Endpoint:** `POST /shells/:id/files/delete`

**Description:** Delete a file or directory.

**Request Body:**
```json
{
 "path": "/var/www/html/temp.txt"
}
```

---

## Payload Generation

### List Templates

**Endpoint:** `GET /payload/templates`

**Description:** Get available payload templates.

**Response:**
```json
[
 {
 "name": "express-middleware",
 "description": "Express中间件注入"
 },
 {
 "name": "koa-middleware",
 "description": "Koa中间件注入"
 }
]
```

### Generate Payload

**Endpoint:** `POST /payload/generate`

**Description:** Generate a payload from a template.

**Request Body:**
```json
{
 "template": "express-middleware",
 "password": "password123",
 "encode_type": "base64",
 "layers": 1
}
```

**Response:**
```json
{
 "payload": "base64_encoded_payload_string"
}
```

---

## Proxy Management

### List Proxies

**Endpoint:** `GET /proxies`

**Description:** Get all configured proxies.

**Response:**
```json
[
 {
 "id": 1,
 "name": "SOCKS5 Proxy",
 "type": "socks5",
 "host": "127.0.0.1",
 "port": 1080,
 "username": "",
 "password": "",
 "enabled": true
 }
]
```

### Create Proxy

**Endpoint:** `POST /proxies`

**Description:** Add a new proxy configuration.

**Request Body:**
```json
{
 "name": "HTTP Proxy",
 "type": "http",
 "host": "proxy.example.com",
 "port": 8080,
 "username": "user",
 "password": "pass",
 "enabled": true
}
```

### Update Proxy

**Endpoint:** `PUT /proxies/:id`

**Description:** Update proxy configuration.

### Delete Proxy

**Endpoint:** `DELETE /proxies/:id`

**Description:** Delete a proxy configuration.

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
 "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `404` - Not Found
- `500` - Internal Server Error

---

## 中文

### 基础 URL

```
http://localhost:8080/api
```

### 认证

当前 API 不需要认证。未来版本将包含基于令牌的认证。

---

## Shell 管理

### 列出所有 Shell

**端点:** `GET /shells`

**描述:** 获取所有已配置 Shell 的列表。

**响应:**
```json
[
 {
 "id": 1,
 "url": "http://example.com/shell.php",
 "password": "password123",
 "encode_type": "base64",
 "protocol": "multipart",
 "method": "POST",
 "group": "production",
 "name": "Web Server 1",
 "status": "online",
 "latency": 150,
 "last_active": "2024-12-12T10:30:00Z",
 "system_info": "OS: Linux...",
 "custom_headers": "{\"User-Agent\":\"Mozilla/5.0\"}",
 "proxy_id": 0
 }
]
```

### 获取 Shell 详情

**端点:** `GET /shells/:id`

**描述:** 获取特定 Shell 的详细信息。

**参数:**
- `id` (路径) - Shell ID

### 创建 Shell

**端点:** `POST /shells`

**描述:** 添加新的 Shell 配置。

**请求体:**
```json
{
 "url": "http://example.com/shell.php",
 "password": "password123",
 "encode_type": "base64",
 "method": "POST",
 "group": "production",
 "name": "Web Server 1",
 "proxy_id": 0,
 "custom_headers": {
 "User-Agent": "Mozilla/5.0"
 }
}
```

### 更新 Shell

**端点:** `PUT /shells/:id`

**描述:** 更新现有的 Shell 配置。

### 删除 Shell

**端点:** `DELETE /shells/:id`

**描述:** 删除 Shell 配置。

### 测试 Shell 连接

**端点:** `POST /shells/:id/test`

**描述:** 测试与 Shell 的连接。

### 获取系统信息

**端点:** `GET /shells/:id/info`

**描述:** 从目标获取系统信息。

**参数:**
- `id` (路径) - Shell ID
- `refresh` (查询参数，可选) - 强制刷新缓存数据

---

## 命令执行

### 执行命令

**端点:** `POST /shells/:id/execute`

**描述:** 在目标系统上执行命令。

**请求体:**
```json
{
 "command": "ls -la"
}
```

---

## 文件管理

### 列出目录

**端点:** `POST /shells/:id/files/list`

**描述:** 列出路径中的文件和目录。

### 读取文件

**端点:** `POST /shells/:id/files/read`

**描述:** 读取文件内容。

### 上传文件

**端点:** `POST /shells/:id/files/upload`

**描述:** 上传文件到目标系统。

### 下载文件

**端点:** `POST /shells/:id/files/download`

**描述:** 从目标系统下载文件。

### 删除文件

**端点:** `POST /shells/:id/files/delete`

**描述:** 删除文件或目录。

---

## Payload 生成

### 列出模板

**端点:** `GET /payload/templates`

**描述:** 获取可用的 Payload 模板。

### 生成 Payload

**端点:** `POST /payload/generate`

**描述:** 从模板生成 Payload。

---

## 代理管理

### 列出代理

**端点:** `GET /proxies`

**描述:** 获取所有已配置的代理。

### 创建代理

**端点:** `POST /proxies`

**描述:** 添加新的代理配置。

### 更新代理

**端点:** `PUT /proxies/:id`

**描述:** 更新代理配置。

### 删除代理

**端点:** `DELETE /proxies/:id`

**描述:** 删除代理配置。

---

## 错误响应

所有端点可能返回以下格式的错误响应：

```json
{
 "error": "错误消息描述"
}
```

**常见 HTTP 状态码:**
- `200` - 成功
- `201` - 已创建
- `400` - 错误请求
- `404` - 未找到
- `500` - 内部服务器错误