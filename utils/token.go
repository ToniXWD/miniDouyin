package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func GenerateToken(username, password string) (string, error) {
	// 选择一个秘钥，用于生成HMAC
	secretKey := []byte("test")

	// 获取当前时间戳（以秒为单位）
	currentTime := time.Now().Unix()

	// 将用户名、密码和时间戳组合在一起，并将其哈希
	data := fmt.Sprintf("%s:%s:%d", username, password, currentTime)
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(data))

	// 计算哈希值
	hashedData := h.Sum(nil)

	// 使用Base64编码哈希值，以生成Token
	token := base64.StdEncoding.EncodeToString(hashedData)

	return token, nil
}
