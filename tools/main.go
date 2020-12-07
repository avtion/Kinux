package tools

import (
	"crypto/rand"
	"fmt"
)

// 生成随机字符串
func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	_, _ = rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
