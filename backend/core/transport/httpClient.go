package transport

import (
	"net/http"
	"NodeJsshell/core/proxy"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

type ProxyConfig = proxy.ProxyConfig

func NewHTTPClient() *HTTPClient {
	transport := &http.Transport{
		DisableCompression:    true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second, // 响应头超时时间
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &HTTPClient{
		client: &http.Client{
			Timeout:   180 * time.Second, // 进一步加长超时时间，大文件分片可能需要更长时间
			Transport: transport,
		},
	}
}

func NewHTTPClientWithProxy(proxyConfig *ProxyConfig) (*HTTPClient, error) {
	client, err := proxy.CreateProxyClient(proxyConfig)
	if err != nil {
		return nil, err
	}
	// CreateProxyClient 已经设置了 Timeout 和 DisableCompression
	// 这里只需要确保一致性，如果CreateProxyClient没有设置，则设置
	if client.Timeout == 0 {
		client.Timeout = 180 * time.Second // 与默认超时保持一致
	}
	if transport, ok := client.Transport.(*http.Transport); ok {
		if !transport.DisableCompression {
			transport.DisableCompression = true
		}
		// 确保响应头超时设置合理
		if transport.ResponseHeaderTimeout == 0 {
			transport.ResponseHeaderTimeout = 60 * time.Second
		}
	}
	return &HTTPClient{client: client}, nil
}

// GetClient returns the underlying http.Client for custom requests
func (c *HTTPClient) GetClient() *http.Client {
	return c.client
}
