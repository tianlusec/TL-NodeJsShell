package handlers

import (
	"NodeJsshell/core/proxy"
	"NodeJsshell/core/transport"
	"NodeJsshell/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CmdHandler struct {
	db *database.DB
}

func NewCmdHandler(db *database.DB) *CmdHandler {
	return &CmdHandler{db: db}
}

func (h *CmdHandler) Execute(c *gin.Context) {
	shellId := c.Param("shellId")
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
	httpClient, err := h.createHTTPClientWithProxy(shell.ProxyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy client: " + err.Error()})
		return
	}
	resp, err := SendRequest(httpClient, &shell, req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *CmdHandler) createHTTPClientWithProxy(proxyID uint) (*transport.HTTPClient, error) {
	if proxyID > 0 {
		var proxyConfig database.Proxy
		if err := h.db.First(&proxyConfig, proxyID).Error; err != nil {
			return transport.NewHTTPClient(), nil
		}
		if !proxyConfig.Enabled {
			return transport.NewHTTPClient(), nil
		}
		proxyCfg := &proxy.ProxyConfig{
			Type:     proxyConfig.Type,
			Host:     proxyConfig.Host,
			Port:     proxyConfig.Port,
			Username: proxyConfig.Username,
			Password: proxyConfig.Password,
		}
		return transport.NewHTTPClientWithProxy(proxyCfg)
	}
	return transport.NewHTTPClient(), nil
}


