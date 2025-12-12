package handlers

import (
	"io"
	"net/http"
	"NodeJsshell/core/proxy"
	"NodeJsshell/core/transport"
	"NodeJsshell/database"
	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	db *database.DB
}

func NewProxyHandler(db *database.DB) *ProxyHandler {
	return &ProxyHandler{db: db}
}

func (h *ProxyHandler) List(c *gin.Context) {
	var proxies []database.Proxy
	if err := h.db.Find(&proxies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proxies)
}

func (h *ProxyHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Type     string `json:"type" binding:"required"`
		Host     string `json:"host" binding:"required"`
		Port     int    `json:"port" binding:"required"`
		Username string `json:"username"`
		Password string `json:"password"`
		Enabled  bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proxy := database.Proxy{
		Name:     req.Name,
		Type:     req.Type,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Enabled:  req.Enabled,
	}
	if err := h.db.Create(&proxy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, proxy)
}

func (h *ProxyHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var proxy database.Proxy
	if err := h.db.First(&proxy, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proxy not found"})
		return
	}
	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Enabled  *bool  `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name != "" {
		proxy.Name = req.Name
	}
	if req.Type != "" {
		proxy.Type = req.Type
	}
	if req.Host != "" {
		proxy.Host = req.Host
	}
	if req.Port > 0 {
		proxy.Port = req.Port
	}
	if req.Username != "" {
		proxy.Username = req.Username
	}
	if req.Password != "" {
		proxy.Password = req.Password
	}
	if req.Enabled != nil {
		proxy.Enabled = *req.Enabled
	}
	if err := h.db.Save(&proxy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proxy)
}

func (h *ProxyHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.Proxy{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Proxy deleted"})
}

func (h *ProxyHandler) Test(c *gin.Context) {
	id := c.Param("id")
	var proxyConfig database.Proxy
	if err := h.db.First(&proxyConfig, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proxy not found"})
		return
	}
	proxyCfg := &proxy.ProxyConfig{
		Type:     proxyConfig.Type,
		Host:     proxyConfig.Host,
		Port:     proxyConfig.Port,
		Username: proxyConfig.Username,
		Password: proxyConfig.Password,
	}
	httpClient, err := transport.NewHTTPClientWithProxy(proxyCfg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	// 创建测试请求：GET /ip HTTP/1.1
	testURL := "http://httpbin.org/ip"
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   "Failed to create request: " + err.Error(),
		})
		return
	}
	
	// 设置请求头
	req.Header.Set("Host", "httpbin.org")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	// 发送请求
	resp, err := httpClient.GetClient().Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   "Failed to send request: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   "Failed to read response: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success":    resp.StatusCode == http.StatusOK,
		"statusCode": resp.StatusCode,
		"data":       string(body),
		"headers":    resp.Header,
	})
}

