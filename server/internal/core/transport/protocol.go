package transport

import (
	"encoding/base64"
	"encoding/json"
	"NodeJsshell/internal/core/crypto"
)

type Request struct {
	Password  string `json:"password"`
	EncodeType string `json:"encode_type"`
	Type      string `json:"type"`
	Data      string `json:"data"`
}

type Response struct {
	Success   bool   `json:"success"`
	Data      string `json:"data"`
	Error     string `json:"error"`
	Timestamp int64  `json:"timestamp"`
}

func BuildRequest(password string, encodeType string, reqType string, data string) *Request {
	return &Request{
		Password:   password,
		EncodeType: encodeType,
		Type:       reqType,
		Data:       data,
	}
}

func EncodeData(data string, encodeType string, password string) string {
	switch encodeType {
	case "base64":
		return crypto.Base64Encode(data, 1)
	case "xor":
		return crypto.XOREncode(data, password, 1)
	case "aes":
		key := []byte(password)
		if len(key) < 32 {
			for len(key) < 32 {
				key = append(key, 0)
			}
		} else {
			key = key[:32]
		}
		encoded, err := crypto.AESEncrypt(data, key)
		if err != nil {
			return data
		}
		return encoded
	default:
		return data
	}
}

func (r *Request) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Request) ToBase64() string {
	data, _ := r.ToJSON()
	return base64.StdEncoding.EncodeToString(data)
}



