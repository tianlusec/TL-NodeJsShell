package transport

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

type MultipartResponse struct {
	Success bool   `json:"success"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Error   string `json:"error"`
}

func (c *HTTPClient) SendMultipartRequest(url string, command string, password string, encodeType string, headers map[string]string) (*Response, error) {
	return c.SendMultipartRequestWithProtocol(url, command, password, encodeType, headers, "multipart")
}

func (c *HTTPClient) SendMultipartRequestWithProtocol(url string, command string, password string, encodeType string, headers map[string]string, protocol string) (*Response, error) {
	commandWithNewline := command + "\n"
	var commandValue string
	encodeTypeTrimmed := strings.TrimSpace(encodeType)

	isNextJS := strings.Contains(protocol, "nextjs") || strings.Contains(url, "nextjs") || strings.Contains(protocol, "next")

	if isNextJS {
		commandBytes := []byte(commandWithNewline)
		commandValue = base64.StdEncoding.EncodeToString(commandBytes)
	} else if encodeTypeTrimmed == "" || encodeTypeTrimmed == "none" {
		commandValue = commandWithNewline
	} else if encodeTypeTrimmed == "base64" {
		commandBytes := []byte(commandWithNewline)
		commandValue = base64.StdEncoding.EncodeToString(commandBytes)
	} else {
		commandValue = commandWithNewline
	}

	var requestBody string
	var contentType string
	if isNextJS {
		paramName := "raw"
		if headers != nil && headers["Content-Type"] != "" && strings.Contains(headers["Content-Type"], "application/json") {
			jsonBody := map[string]string{paramName: commandValue}
			bodyBytes, _ := json.Marshal(jsonBody)
			requestBody = string(bodyBytes)
			contentType = "application/json"
		} else {
			requestBody = fmt.Sprintf("%s=%s", neturl.QueryEscape(paramName), neturl.QueryEscape(commandValue))
			contentType = "application/x-www-form-urlencoded"
		}
	} else {
		requestBody = fmt.Sprintf("%s=%s", neturl.QueryEscape(password), neturl.QueryEscape(commandValue))
		contentType = "application/x-www-form-urlencoded"
	}

	bodyBytes := []byte(requestBody)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36 Assetnote/1.0.0")
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}
	httpReq.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))

	for key, value := range headers {
		keyLower := strings.ToLower(key)
		if keyLower != "content-type" && keyLower != "content-length" {
			httpReq.Header.Set(key, value)
		} else if keyLower == "content-type" && !isNextJS {
			httpReq.Header.Set(key, value)
		}
	}

	httpReq.Header.Del("Accept-Encoding")
	if isNextJS {
		httpReq.Header.Del("Connection")
		httpReq.Header.Del("Keep-Alive")
	}

	// 添加请求发送前的日志（立即刷新）
	fmt.Printf("[SendMultipartRequestWithProtocol] 准备发送请求: URL=%s, Method=%s, BodySize=%d bytes\n", url, httpReq.Method, len(bodyBytes))

	// 检查超时设置
	fmt.Printf("[SendMultipartRequestWithProtocol] 客户端超时设置: %v\n", c.client.Timeout)

	if transport, ok := c.client.Transport.(*http.Transport); ok {
		fmt.Printf("[SendMultipartRequestWithProtocol] Transport 类型正确: *http.Transport\n")
		if transport.Proxy != nil {
			proxyURL, proxyErr := transport.Proxy(httpReq)
			if proxyErr != nil {
				fmt.Printf("[SendMultipartRequestWithProtocol] ⚠️  代理函数执行出错: %v\n", proxyErr)
			} else if proxyURL != nil {
				fmt.Printf("[SendMultipartRequestWithProtocol] ✅ 通过代理发送: %s\n", proxyURL.String())
			} else {
				fmt.Printf("[SendMultipartRequestWithProtocol] ⚠️  代理函数返回 nil，可能直接连接\n")
			}
		} else {
			fmt.Printf("[SendMultipartRequestWithProtocol] ⚠️  未配置代理，直接连接\n")
		}
	} else {
		fmt.Printf("[SendMultipartRequestWithProtocol] ⚠️  无法检查代理配置，Transport类型: %T\n", c.client.Transport)
	}

	fmt.Printf("[SendMultipartRequestWithProtocol] 开始执行 Do() 调用...\n")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		fmt.Printf("[SendMultipartRequestWithProtocol] ❌ 请求发送失败: %v\n", err)
		return nil, err
	}
	fmt.Printf("[SendMultipartRequestWithProtocol] ✅ 请求发送成功，响应状态: %d\n", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return &Response{
			Success:   false,
			Data:      string(body),
			Error:     fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			Timestamp: time.Now().Unix(),
		}, nil
	}

	return ParseMultipartResponse(body, resp.StatusCode)
}

func ParseMultipartResponse(data []byte, statusCode int) (*Response, error) {
	if statusCode != http.StatusOK {
		response := &Response{
			Success:   false,
			Timestamp: time.Now().Unix(),
		}
		if len(data) > 0 {
			response.Data = string(data)
			response.Error = fmt.Sprintf("HTTP %d", statusCode)
		} else {
			response.Error = fmt.Sprintf("HTTP %d: Empty response", statusCode)
		}
		return response, nil
	}

	responseText := string(data)
	if strings.Contains(responseText, "Server action not found") || strings.Contains(responseText, "server action not found") {
		response := &Response{
			Success:   false,
			Data:      responseText,
			Error:     "Server action not found",
			Timestamp: time.Now().Unix(),
		}
		return response, nil
	}

	var multipartResp MultipartResponse
	err := json.Unmarshal(data, &multipartResp)
	if err != nil {
		response := &Response{
			Success:   false,
			Data:      string(data),
			Error:     fmt.Sprintf("Failed to parse JSON response: %v. Response body: %s", err, string(data)),
			Timestamp: time.Now().Unix(),
		}
		return response, nil
	}

	if !multipartResp.Success && (strings.Contains(responseText, "Server action not found") || strings.Contains(responseText, "server action not found") || multipartResp.Stdout == "Server action not found.") {
		response := &Response{
			Success:   false,
			Data:      multipartResp.Stdout,
			Error:     "Server action not found",
			Timestamp: time.Now().Unix(),
		}
		return response, nil
	}

	if !multipartResp.Success || strings.Contains(multipartResp.Stdout, "Server action not found") || strings.Contains(multipartResp.Stderr, "Server action not found") {
		if strings.Contains(multipartResp.Stdout, "Server action not found") || strings.Contains(multipartResp.Stderr, "Server action not found") {
			response := &Response{
				Success:   false,
				Data:      multipartResp.Stdout + multipartResp.Stderr,
				Error:     "Server action not found",
				Timestamp: time.Now().Unix(),
			}
			return response, nil
		}
	}

	response := &Response{
		Success:   multipartResp.Success,
		Timestamp: time.Now().Unix(),
	}
	if multipartResp.Stdout != "" {
		response.Data = multipartResp.Stdout
	}
	if multipartResp.Stderr != "" {
		if response.Data != "" {
			if !strings.HasSuffix(response.Data, "\n") {
				response.Data += "\n"
			}
			response.Data += multipartResp.Stderr
		} else {
			response.Data = multipartResp.Stderr
		}
	}
	if response.Data != "" {
		normalized := strings.ReplaceAll(response.Data, "\r\n", "\n")
		normalized = strings.ReplaceAll(normalized, "\r", "\n")
		lines := strings.Split(normalized, "\n")
		normalizedLines := make([]string, 0, len(lines))
		for i, line := range lines {
			if i == len(lines)-1 && line == "" {
				continue
			}
			trimmedLine := strings.TrimRight(line, " \t\r")
			trimmedLine = strings.TrimLeft(trimmedLine, " \t")
			normalizedLines = append(normalizedLines, trimmedLine)
		}
		response.Data = strings.Join(normalizedLines, "\n")
		if len(normalizedLines) > 0 {
			response.Data += "\n"
		}
	}
	if multipartResp.Error != "" {
		response.Error = multipartResp.Error
		if !response.Success {
			response.Data = multipartResp.Error
		}
	}
	return response, nil
}
