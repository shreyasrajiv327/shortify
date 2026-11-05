package utils

import (
	"crypto/rand"
	"fmt"
)

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func EncodeBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	result := ""

	for num > 0 {
		remainder := num % 62
		result = string(base62Chars[remainder]) + result
		num = num / 62
	}

	return result
}

func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("short code length must be positive")
	}

	const randomBytesLength = 32
	validByteLimit := byte(256 - 256%len(base62Chars))
	code := make([]byte, 0, length)
	randomBytes := make([]byte, randomBytesLength)

	for len(code) < length {
		if _, err := rand.Read(randomBytes); err != nil {
			return "", err
		}
		for _, value := range randomBytes {
			if value >= validByteLimit {
				continue
			}
			code = append(code, base62Chars[int(value)%len(base62Chars)])
			if len(code) == length {
				break
			}
		}
	}

	return string(code), nil
}
