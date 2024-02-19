package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 小写
func Md5Encode(str string) string {
	var md = md5.New()
	md.Write([]byte(str))
	var code = hex.EncodeToString(md.Sum(nil))
	return code
}

// 大写
func MD5Encode(str string) string {
	return strings.ToUpper(Md5Encode(str))
}

// 加密
func CryptoPassword(pwd, salt string) string {
	return Md5Encode(pwd + salt)
}

// 加密,随机生成 Salt
func CryptoPasswordWithSalt(pwd string) (string, string) {
	salt := generateRandomSalt(6)
	return Md5Encode(pwd + salt), salt
}

// 验证密码
func ValidPassword(pwd, salt string, source_pwd string) bool {
	code := Md5Encode(pwd + salt)
	fmt.Println("code:", code)
	fmt.Println("source_pwd:", source_pwd)
	return code == source_pwd
}

// 生成随机字符串，n 长度
func generateRandomSalt(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
