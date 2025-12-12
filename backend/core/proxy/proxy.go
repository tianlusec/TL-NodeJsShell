package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type ProxyConfig struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
}

func CreateProxyClient(cfg *ProxyConfig) (*http.Client, error) {
	var httpTransport *http.Transport

	if cfg.Type == "http" || cfg.Type == "https" {
		// 使用 url.UserPassword 来正确编码用户名和密码中的特殊字符
		proxyURL := &url.URL{
			Scheme: cfg.Type,
			Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		}

		// 如果提供了用户名和密码，设置UserInfo
		if cfg.Username != "" {
			proxyURL.User = url.UserPassword(cfg.Username, cfg.Password)
		}

		httpTransport = &http.Transport{
			Proxy:                 http.ProxyURL(proxyURL),
			DisableCompression:    true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 120 * time.Second, // 代理场景下，响应头超时时间更长
			ExpectContinueTimeout: 1 * time.Second,
		}
	} else if cfg.Type == "socks5" {
		var auth *proxy.Auth
		if cfg.Username != "" {
			auth = &proxy.Auth{
				User:     cfg.Username,
				Password: cfg.Password,
			}
		}

		dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), auth, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
		}

		// 使用 DialContext 而不是 Dial（Dial已过时）
		// 将 dialer.Dial 包装为 DialContext
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 创建一个channel来接收连接和错误
			type result struct {
				conn net.Conn
				err  error
			}
			done := make(chan result, 1)

			go func() {
				conn, err := dialer.Dial(network, addr)
				done <- result{conn: conn, err: err}
			}()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case res := <-done:
				return res.conn, res.err
			}
		}

		httpTransport = &http.Transport{
			DialContext:           dialContext,
			DisableCompression:    true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 120 * time.Second, // 代理场景下，响应头超时时间更长
			ExpectContinueTimeout: 1 * time.Second,
		}
	} else {
		return nil, fmt.Errorf("unsupported proxy type: %s", cfg.Type)
	}

	return &http.Client{
		Transport: httpTransport,
		Timeout:   180 * time.Second, // 进一步加长超时时间，大文件分片可能需要更长时间
	}, nil
}
