package crypto

import (
	"encoding/base64"
)

func Base64Encode(data string, layers int) string {
	result := data
	for i := 0; i < layers; i++ {
		result = base64.StdEncoding.EncodeToString([]byte(result))
	}
	return result
}

func Base64Decode(data string, layers int) string {
	result := data
	for i := 0; i < layers; i++ {
		decoded, err := base64.StdEncoding.DecodeString(result)
		if err != nil {
			return data
		}
		result = string(decoded)
	}
	return result
}



