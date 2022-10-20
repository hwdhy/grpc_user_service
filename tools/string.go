package tools

import "math/rand"

// RandomString 生成随机字符串
func RandomString(n int) string {
	str := "1234567890qwertyuopasdfghjkzxcvbnm"

	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = str[rand.Intn(len(str))]
	}
	return string(res)
}
