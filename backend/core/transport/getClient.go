package transport

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *HTTPClient) SendGetRequest(urlStr string, command string, password string, encodeType string, headers map[string]string) (*Response, error) {
	commandWithNewline := command + "\n"
	var commandValue string
	encodeTypeTrimmed := strings.TrimSpace(encodeType)
	if encodeTypeTrimmed == "" || encodeTypeTrimmed == "none" {
		commandValue = commandWithNewline
	} else if encodeTypeTrimmed == "base64" {
		commandBytes := []byte(commandWithNewline)
		commandValue = base64.StdEncoding.EncodeToString(commandBytes)
	} else {
		commandValue = commandWithNewline
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}
	query := parsedURL.Query()
	query.Set(password, commandValue)
	parsedURL.RawQuery = query.Encode()
	httpReq, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36 Assetnote/1.0.0")
	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}

	// 添加请求发送前的日志
	fmt.Printf("[SendGetRequest] 准备发送请求: URL=%s, Method=%s\n", parsedURL.String(), httpReq.Method)
	if transport, ok := c.client.Transport.(*http.Transport); ok {
		if transport.Proxy != nil {
			proxyURL, _ := transport.Proxy(httpReq)
			if proxyURL != nil {
				fmt.Printf("[SendGetRequest] ✅ 通过代理发送: %s\n", proxyURL.String())
			} else {
				fmt.Printf("[SendGetRequest] ⚠️  代理函数返回 nil，可能直接连接\n")
			}
		} else {
			fmt.Printf("[SendGetRequest] ⚠️  未配置代理，直接连接\n")
		}
	} else {
		fmt.Printf("[SendGetRequest] ⚠️  无法检查代理配置，Transport类型: %T\n", c.client.Transport)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		fmt.Printf("[SendGetRequest] ❌ 请求发送失败: %v\n", err)
		return nil, err
	}
	fmt.Printf("[SendGetRequest] ✅ 请求发送成功，响应状态: %d\n", resp.StatusCode)
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
