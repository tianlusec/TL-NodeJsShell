package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"NodeJsshell/internal/core/payload"
	"NodeJsshell/internal/core/proxy"
	"NodeJsshell/internal/core/transport"
	"NodeJsshell/internal/database"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShellHandler struct {
	db *database.DB
}

func NewShellHandler(db *database.DB) *ShellHandler {
	return &ShellHandler{db: db}
}

func (h *ShellHandler) List(c *gin.Context) {
	var shells []database.Shell
	if err := h.db.Find(&shells).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shells)
}

func (h *ShellHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var shell database.Shell
	if err := h.db.First(&shell, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shell not found"})
		return
	}
	c.JSON(http.StatusOK, shell)
}

func (h *ShellHandler) Create(c *gin.Context) {
	var req struct {
		URL           string            `json:"url" binding:"required"`
		Password      string            `json:"password" binding:"required"`
		EncodeType    string            `json:"encode_type"`
		Method        string            `json:"method"`
		Group         string            `json:"group"`
		Name          string            `json:"name"`
		ProxyID       uint              `json:"proxy_id"`
		CustomHeaders map[string]string `json:"custom_headers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Method == "" {
		req.Method = "POST"
	}
	if req.Name == "" {
		req.Name = uuid.New().String()[:8]
	}
	var customHeadersJSON string
	if len(req.CustomHeaders) > 0 {
		headersJSON, _ := json.Marshal(req.CustomHeaders)
		customHeadersJSON = string(headersJSON)
	}
	shell := database.Shell{
		URL:           req.URL,
		Password:      req.Password,
		EncodeType:    req.EncodeType,
		Protocol:      "multipart",
		Method:        req.Method,
		Group:         req.Group,
		Name:          req.Name,
		Status:        "offline",
		CustomHeaders: customHeadersJSON,
		ProxyID:       req.ProxyID,
	}
	if err := h.db.Create(&shell).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, shell)
}

func (h *ShellHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var shell database.Shell
	if err := h.db.First(&shell, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shell not found"})
		return
	}
	var req struct {
		URL           string      `json:"url"`
		Password      string      `json:"password"`
		EncodeType    string      `json:"encode_type"`
		Method        string      `json:"method"`
		Group         string      `json:"group"`
		Name          string      `json:"name"`
		ProxyID       *uint       `json:"proxy_id"`
		CustomHeaders interface{} `json:"custom_headers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.URL != "" {
		shell.URL = req.URL
	}
	if req.Password != "" {
		shell.Password = req.Password
	}
	shell.EncodeType = req.EncodeType
	if req.Method != "" {
		shell.Method = req.Method
	}
	if req.Group != "" {
		shell.Group = req.Group
	}
	if req.Name != "" {
		shell.Name = req.Name
	}
	if req.ProxyID != nil {
		shell.ProxyID = *req.ProxyID
	}
	if req.CustomHeaders != nil {
		var headersJSON string
		// 处理 custom_headers 可以是字符串或 map 两种情况
		switch v := req.CustomHeaders.(type) {
		case string:
			// 如果已经是字符串，直接使用
			headersJSON = v
		case map[string]interface{}:
			// 如果是 map，转换为 JSON 字符串
			bytes, err := json.Marshal(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom_headers format: " + err.Error()})
				return
			}
			headersJSON = string(bytes)
		default:
			// 尝试直接序列化
			bytes, err := json.Marshal(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom_headers format: " + err.Error()})
				return
			}
			headersJSON = string(bytes)
		}
		shell.CustomHeaders = headersJSON
	}
	if err := h.db.Save(&shell).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shell)
}

func (h *ShellHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.Shell{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shell deleted"})
}

func (h *ShellHandler) Test(c *gin.Context) {
	id := c.Param("id")
	var shell database.Shell
	if err := h.db.First(&shell, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shell not found"})
		return
	}
	httpClient, err := h.createHTTPClientWithProxy(shell.ProxyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy client: " + err.Error()})
		return
	}
	success, latency, err := h.testShellConnection(&shell, httpClient)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if success {
		shell.LastActive = time.Now()
		shell.Status = "online"
		shell.Latency = latency
	} else {
		shell.Status = "offline"
	}
	h.db.Save(&shell)
	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"latency": latency,
		"status":  shell.Status,
	})
}

func (h *ShellHandler) Execute(c *gin.Context) {
	// 路由参数名是 :id，不是 :shellId
	shellId := c.Param("id")
	var req struct {
		Command string `json:"command" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var shell database.Shell
	if err := h.db.First(&shell, shellId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shell not found"})
		return
	}
	log.Printf("[Execute] 查询Shell - ID参数: %s, 查询到的Shell ID: %d, ProxyID: %d", shellId, shell.ID, shell.ProxyID)
	log.Printf("[Execute] Shell ID=%d, ProxyID=%d, Command=%s", shell.ID, shell.ProxyID, req.Command)
	httpClient, err := h.createHTTPClientWithProxy(shell.ProxyID)
	if err != nil {
		log.Printf("[Execute] ❌ 创建代理客户端失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy client: " + err.Error()})
		return
	}
	log.Printf("[Execute] ✅ HTTP客户端创建完成，准备发送请求到: %s", shell.URL)
	resp, err := SendRequest(httpClient, &shell, req.Command)
	if err != nil {
		log.Printf("请求执行失败 (Shell ID=%d, ProxyID=%d): %v", shell.ID, shell.ProxyID, err)
	}
	if resp != nil && !resp.Success {
		log.Printf("请求失败响应 (Shell ID=%d, ProxyID=%d): Success=%v, Error=%s, Data=%s", shell.ID, shell.ProxyID, resp.Success, resp.Error, resp.Data)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	history := database.History{
		ShellID: shell.ID,
		Type:    "exec",
		Command: req.Command,
		Result:  resp.Data,
		Success: resp.Success,
	}
	h.db.Create(&history)
	shell.LastActive = time.Now()
	if resp.Success {
		shell.Status = "online"
	}
	h.db.Save(&shell)
	c.JSON(http.StatusOK, resp)
}

func (h *ShellHandler) testShellConnection(shell *database.Shell, httpClient *transport.HTTPClient) (bool, int, error) {
	start := time.Now()
	var resp *transport.Response
	var err error
	var customHeaders map[string]string
	if shell.CustomHeaders != "" {
		json.Unmarshal([]byte(shell.CustomHeaders), &customHeaders)
	}
	if customHeaders == nil {
		customHeaders = make(map[string]string)
	}
	resp, err = SendRequest(httpClient, shell, "echo test")
	if err != nil {
		return false, 0, err
	}
	latency := int(time.Since(start).Milliseconds())
	return resp.Success, latency, nil
}

func (h *ShellHandler) GetInfo(c *gin.Context) {
	id := c.Param("id")
	var shell database.Shell
	if err := h.db.First(&shell, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shell not found"})
		return
	}

	forceRefresh := c.Query("refresh") == "true"

	if !forceRefresh && shell.SystemInfoJSON != "" {
		var cachedSystemInfo map[string]interface{}
		if err := json.Unmarshal([]byte(shell.SystemInfoJSON), &cachedSystemInfo); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"system_info": cachedSystemInfo,
				"shell":       shell,
				"cached":      true,
			})
			return
		}
	}

	httpClient, err := h.createHTTPClientWithProxy(shell.ProxyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy client: " + err.Error()})
		return
	}
	templateCode, err := payload.GenerateFunctionTemplate("systemInfo", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate systemInfo template: " + err.Error()})
		return
	}
	var resp *transport.Response
	var customHeaders map[string]string
	if shell.CustomHeaders != "" {
		json.Unmarshal([]byte(shell.CustomHeaders), &customHeaders)
	}
	if customHeaders == nil {
		customHeaders = make(map[string]string)
	}
	if shell.Protocol == "multipart" {
		customHeaders = getCustomHeaders(&shell)
		if shell.Method == "GET" {
			resp, err = httpClient.SendGetRequest(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders)
		} else {
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		req := transport.BuildRequest(shell.Password, shell.EncodeType, "exec", templateCode)
		resp, err = httpClient.SendRequest(shell.URL, req)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !resp.Success {
		errorMsg := resp.Error
		if errorMsg == "" {
			errorMsg = resp.Data
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell request failed: " + errorMsg})
		return
	}

	if resp.Data == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned empty response"})
		return
	}

	dataTrimmed := strings.TrimSpace(resp.Data)
	if !strings.HasPrefix(dataTrimmed, "{") && !strings.HasPrefix(dataTrimmed, "[") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned invalid response: " + resp.Data})
		return
	}

	var systemInfoResp struct {
		Ok                bool                   `json:"ok"`
		Platform          string                 `json:"platform"`
		Arch              string                 `json:"arch"`
		Hostname          string                 `json:"hostname"`
		Type              string                 `json:"type"`
		Release           string                 `json:"release"`
		Uptime            float64                `json:"uptime"`
		Totalmem          int64                  `json:"totalmem"`
		Freemem           int64                  `json:"freemem"`
		Cpus              int                    `json:"cpus"`
		UserInfo          map[string]interface{} `json:"userInfo"`
		EnvVars           map[string]string      `json:"envVars"`
		Hosts             string                 `json:"hosts"`
		NetworkInterfaces []struct {
			Interface string `json:"interface"`
			Address   string `json:"address"`
			Family    string `json:"family"`
			Internal  bool   `json:"internal"`
			Netmask   string `json:"netmask,omitempty"`
			Mac       string `json:"mac,omitempty"`
		} `json:"networkInterfaces"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(resp.Data), &systemInfoResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse shell response: " + err.Error() + ", Response data: " + resp.Data})
		return
	}
	if !systemInfoResp.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned error: " + systemInfoResp.Error})
		return
	}

	systemInfoStr := fmt.Sprintf("OS: %s %s\nArch: %s\nHostname: %s\nUptime: %.0f seconds\nTotal Memory: %d bytes\nFree Memory: %d bytes",
		systemInfoResp.Platform, systemInfoResp.Release, systemInfoResp.Arch, systemInfoResp.Hostname, systemInfoResp.Uptime, systemInfoResp.Totalmem, systemInfoResp.Freemem)

	if len(systemInfoResp.NetworkInterfaces) > 0 {
		systemInfoStr += "\n\nNetwork Interfaces:"
		for _, iface := range systemInfoResp.NetworkInterfaces {
			systemInfoStr += fmt.Sprintf("\n  %s: %s (%s)", iface.Interface, iface.Address, iface.Family)
		}
	}

	systemInfoJSONBytes, _ := json.Marshal(systemInfoResp)
	shell.SystemInfo = systemInfoStr
	shell.SystemInfoJSON = string(systemInfoJSONBytes)
	h.db.Save(&shell)
	c.JSON(http.StatusOK, gin.H{
		"system_info": systemInfoResp,
		"shell":       shell,
		"cached":      false,
	})
}

func (h *ShellHandler) createHTTPClientWithProxy(proxyID uint) (*transport.HTTPClient, error) {
	log.Printf("[createHTTPClientWithProxy] 调用 - ProxyID: %d", proxyID)
	if proxyID > 0 {
		var proxyConfig database.Proxy
		if err := h.db.First(&proxyConfig, proxyID).Error; err != nil {
			log.Printf("[createHTTPClientWithProxy] ❌ Proxy ID %d not found: %v", proxyID, err)
			return transport.NewHTTPClient(), nil
		}
		log.Printf("[createHTTPClientWithProxy] 查询到代理配置: %s %s:%d, Enabled: %v", proxyConfig.Type, proxyConfig.Host, proxyConfig.Port, proxyConfig.Enabled)
		if !proxyConfig.Enabled {
			log.Printf("[createHTTPClientWithProxy] ⚠️  Proxy ID %d is disabled, 不使用代理", proxyID)
			return transport.NewHTTPClient(), nil
		}
		log.Printf("[createHTTPClientWithProxy] ✅ 使用代理: %s %s:%d (用户名: %s)", proxyConfig.Type, proxyConfig.Host, proxyConfig.Port, proxyConfig.Username)
		proxyCfg := &proxy.ProxyConfig{
			Type:     proxyConfig.Type,
			Host:     proxyConfig.Host,
			Port:     proxyConfig.Port,
			Username: proxyConfig.Username,
			Password: proxyConfig.Password,
		}
		client, err := transport.NewHTTPClientWithProxy(proxyCfg)
		if err != nil {
			log.Printf("[createHTTPClientWithProxy] ❌ 创建代理客户端失败: %v", err)
			return nil, err
		}
		log.Printf("[createHTTPClientWithProxy] ✅ 代理客户端创建成功")
		return client, nil
	}
	log.Printf("[createHTTPClientWithProxy] ⚠️  ProxyID为0，不使用代理")
	return transport.NewHTTPClient(), nil
}
