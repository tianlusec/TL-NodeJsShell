package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *HTTPClient) SendRequest(url string, req *Request) (*Response, error) {
	jsonData, err := req.ToJSON()
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 添加请求发送前的日志
	fmt.Printf("[SendRequest] 准备发送请求: URL=%s, Method=%s, BodySize=%d bytes\n", url, httpReq.Method, len(jsonData))
	if transport, ok := c.client.Transport.(*http.Transport); ok {
		if transport.Proxy != nil {
			proxyURL, _ := transport.Proxy(httpReq)
			if proxyURL != nil {
				fmt.Printf("[SendRequest] ✅ 通过代理发送: %s\n", proxyURL.String())
			} else {
				fmt.Printf("[SendRequest] ⚠️  代理函数返回 nil，可能直接连接\n")
			}
		} else {
			fmt.Printf("[SendRequest] ⚠️  未配置代理，直接连接\n")
		}
	} else {
		fmt.Printf("[SendRequest] ⚠️  无法检查代理配置，Transport类型: %T\n", c.client.Transport)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		fmt.Printf("[SendRequest] ❌ 请求发送失败: %v\n", err)
		return nil, err
	}
	fmt.Printf("[SendRequest] ✅ 请求发送成功，响应状态: %d\n", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &response, nil
}
