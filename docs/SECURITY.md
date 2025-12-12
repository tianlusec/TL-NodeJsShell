# Security Policy

[English](#english) | [中文](#中文)

---

## English

### ⚠️ Important Security Notice

**TL-NodeJsShell is a security testing tool intended for authorized use only.**

### Legal and Ethical Use

This tool is designed for:
- ✅ Authorized penetration testing
- ✅ Security research in controlled environments
- ✅ Educational purposes with proper authorization
- ✅ Red team exercises with explicit permission

This tool must NOT be used for:
- ❌ Unauthorized access to systems
- ❌ Malicious activities
- ❌ Illegal hacking
- ❌ Any activity without explicit permission

### User Responsibilities

By using this tool, you agree to:

1. **Obtain Authorization**: Always get written permission before testing any system
2. **Follow Laws**: Comply with all applicable local, national, and international laws
3. **Respect Privacy**: Do not access, modify, or exfiltrate data without authorization
4. **Use Responsibly**: Only use in controlled, authorized environments
5. **Accept Liability**: Take full responsibility for your actions

### Reporting Security Vulnerabilities

We take security seriously. If you discover a security vulnerability in TL-NodeJsShell:

#### Please DO:
- Report it privately via email or GitHub Security Advisory
- Provide detailed information about the vulnerability
- Allow reasonable time for us to address the issue
- Follow responsible disclosure practices

#### Please DO NOT:
- Publicly disclose the vulnerability before it's fixed
- Exploit the vulnerability maliciously
- Test vulnerabilities on systems you don't own

#### How to Report

**Email:** Create a GitHub Security Advisory at:
```
https://github.com/tianlusec/TL-NodeJsShell/security/advisories/new
```

**Include:**
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)
- Your contact information

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity (critical issues prioritized)

### Security Best Practices

When using TL-NodeJsShell:

#### 1. Network Security
- Use VPNs or secure networks
- Configure proxies appropriately
- Avoid exposing the management interface to the internet
- Use firewall rules to restrict access

#### 2. Access Control
- Use strong passwords
- Limit access to authorized personnel only
- Keep logs of all activities
- Regularly review access permissions

#### 3. Data Protection
- Encrypt sensitive data
- Use secure communication channels
- Don't store credentials in plain text
- Regularly backup important data

#### 4. System Hardening
- Keep the tool updated
- Run with minimal privileges
- Use isolated environments for testing
- Monitor for suspicious activities

#### 5. Operational Security
- Document all testing activities
- Maintain audit trails
- Follow your organization's security policies
- Conduct regular security reviews

### Known Security Considerations

#### Current Limitations

1. **No Built-in Authentication**: The current version lacks user authentication
   - **Mitigation**: Use network-level access controls
   - **Future**: Authentication will be added in future versions

2. **Local Storage**: Sensitive data stored in SQLite database
   - **Mitigation**: Encrypt the database file
   - **Mitigation**: Restrict file system permissions

3. **Network Traffic**: Communications may be intercepted
   - **Mitigation**: Use HTTPS for target connections
   - **Mitigation**: Use VPN or secure tunnels

4. **Logging**: Activities are logged locally
   - **Mitigation**: Secure log files appropriately
   - **Mitigation**: Regularly review and rotate logs

### Secure Configuration

#### Recommended Settings

```go
// backend/config/config.go
type Config struct {
    Port string  // Bind to localhost only: "127.0.0.1:8080"
    Host string  // Use "127.0.0.1" instead of "0.0.0.0"
}
```

#### Firewall Rules

```bash
# Linux (iptables)
iptables -A INPUT -p tcp --dport 8080 -s 127.0.0.1 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP

# Windows (PowerShell)
New-NetFirewallRule -DisplayName "TL-NodeJsShell" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow -RemoteAddress 127.0.0.1
```

### Disclaimer

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND. The authors and contributors:

- Are NOT responsible for any misuse or damage
- Do NOT endorse illegal activities
- Provide NO guarantees of security or reliability
- Assume NO liability for your actions

**You are solely responsible for ensuring your use complies with all applicable laws and regulations.**

### Updates and Patches

- Security updates will be released as soon as possible
- Critical vulnerabilities will be prioritized
- Subscribe to releases for notifications
- Check CHANGELOG.md for security-related updates

### Contact

For security-related inquiries:
- GitHub Security Advisory: [Create Advisory](https://github.com/tianlusec/TL-NodeJsShell/security/advisories/new)
- General Issues: [GitHub Issues](https://github.com/tianlusec/TL-NodeJsShell/issues)

---

## 中文

### ⚠️ 重要安全声明

**TL-NodeJsShell 是一个安全测试工具，仅供授权使用。**

### 合法和道德使用

本工具设计用于：
- ✅ 授权的渗透测试
- ✅ 受控环境中的安全研究
- ✅ 具有适当授权的教育目的
- ✅ 获得明确许可的红队演练

本工具不得用于：
- ❌ 未经授权访问系统
- ❌ 恶意活动
- ❌ 非法黑客行为
- ❌ 任何未经明确许可的活动

### 用户责任

使用本工具即表示您同意：

1. **获得授权**：在测试任何系统之前始终获得书面许可
2. **遵守法律**：遵守所有适用的地方、国家和国际法律
3. **尊重隐私**：未经授权不得访问、修改或窃取数据
4. **负责任使用**：仅在受控的授权环境中使用
5. **承担责任**：对您的行为承担全部责任

### 报告安全漏洞

我们非常重视安全。如果您在 TL-NodeJsShell 中发现安全漏洞：

#### 请务必：
- 通过电子邮件或 GitHub 安全公告私下报告
- 提供有关漏洞的详细信息
- 给我们合理的时间来解决问题
- 遵循负责任的披露实践

#### 请勿：
- 在修复之前公开披露漏洞
- 恶意利用漏洞
- 在您不拥有的系统上测试漏洞

#### 如何报告

**创建 GitHub 安全公告：**
```
https://github.com/tianlusec/TL-NodeJsShell/security/advisories/new
```

**包含：**
- 漏洞描述
- 重现步骤
- 潜在影响
- 建议的修复方案（如有）
- 您的联系信息

### 响应时间表

- **初始响应**：48 小时内
- **状态更新**：7 天内
- **修复时间**：取决于严重程度（优先处理关键问题）

### 安全最佳实践

使用 TL-NodeJsShell 时：

#### 1. 网络安全
- 使用 VPN 或安全网络
- 适当配置代理
- 避免将管理界面暴露到互联网
- 使用防火墙规则限制访问

#### 2. 访问控制
- 使用强密码
- 仅限授权人员访问
- 保留所有活动的日志
- 定期审查访问权限

#### 3. 数据保护
- 加密敏感数据
- 使用安全通信渠道
- 不要以明文存储凭据
- 定期备份重要数据

#### 4. 系统加固
- 保持工具更新
- 以最小权限运行
- 使用隔离环境进行测试
- 监控可疑活动

#### 5. 操作安全
- 记录所有测试活动
- 维护审计跟踪
- 遵循组织的安全策略
- 进行定期安全审查

### 已知安全注意事项

#### 当前限制

1. **无内置认证**：当前版本缺少用户认证
   - **缓解措施**：使用网络级访问控制
   - **未来**：将在未来版本中添加认证

2. **本地存储**：敏感数据存储在 SQLite 数据库中
   - **缓解措施**：加密数据库文件
   - **缓解措施**：限制文件系统权限

3. **网络流量**：通信可能被拦截
   - **缓解措施**：对目标连接使用 HTTPS
   - **缓解措施**：使用 VPN 或安全隧道

4. **日志记录**：活动在本地记录
   - **缓解措施**：适当保护日志文件
   - **缓解措施**：定期审查和轮换日志

### 安全配置

#### 推荐设置

```go
// backend/config/config.go
type Config struct {
    Port string  // 仅绑定到本地主机: "127.0.0.1:8080"
    Host string  // 使用 "127.0.0.1" 而不是 "0.0.0.0"
}
```

#### 防火墙规则

```bash
# Linux (iptables)
iptables -A INPUT -p tcp --dport 8080 -s 127.0.0.1 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP

# Windows (PowerShell)
New-NetFirewallRule -DisplayName "TL-NodeJsShell" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow -RemoteAddress 127.0.0.1
```

### 免责声明

软件按"原样"提供，不提供任何形式的保证。作者和贡献者：

- 不对任何滥用或损害负责
- 不支持非法活动
- 不提供安全性或可靠性保证
- 不对您的行为承担任何责任

**您有责任确保您的使用符合所有适用的法律法规。**

### 更新和补丁

- 安全更新将尽快发布
- 关键漏洞将优先处理
- 订阅发布以获取通知
- 查看 CHANGELOG.md 了解与安全相关的更新

### 联系方式

安全相关咨询：
- GitHub 安全公告：[创建公告](https://github.com/tianlusec/TL-NodeJsShell/security/advisories/new)
- 一般问题：[GitHub Issues](https://github.com/tianlusec/TL-NodeJsShell/issues)

---<div align="center">

**Security is everyone's responsibility / 安全是每个人的责任**

</div>