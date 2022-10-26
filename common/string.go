package common

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

// RandomString Generate n-bit random string
func RandomString(n int) string {
	str := "1234567890qwertyuopasdfghjkzxcvbnm"

	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = str[rand.Intn(len(str))]
	}
	return string(res)
}

// StringHash password hash encode
func StringHash(password string, salt string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum([]byte(salt)))
}
