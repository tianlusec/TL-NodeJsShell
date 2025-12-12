package crypto

func XOREncode(data string, key string, layers int) string {
	result := data
	for i := 0; i < layers; i++ {
		result = xor(result, key)
	}
	return result
}

func XORDecode(data string, key string, layers int) string {
	return XOREncode(data, key, layers)
}

func xor(data string, key string) string {
	if key == "" {
		return data
	}
	
	dataBytes := []byte(data)
	keyBytes := []byte(key)
	keyLen := len(keyBytes)
	
	result := make([]byte, len(dataBytes))
	for i := 0; i < len(dataBytes); i++ {
		result[i] = dataBytes[i] ^ keyBytes[i%keyLen]
	}
	return string(result)
}



