package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"regexp"
	"time"
	"unicode"
)

// RandString 随机名字
func RandString(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0987654321")
	results := make([]byte, n)

	rand.Seed(time.Now().Unix()) //时间戳随机种子
	for i := range results {
		results[i] = letters[rand.Intn(len(letters))]
	}

	return string(results)

}

// 验证密码是否正确
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 对密码进行加密
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
	}
	return string(hash)
}

// 验证手机号
func IsValidPhone(phone string) bool {
	// 2019 工信部公布的手机号
	regexpPhone := regexp.MustCompile(`^(?:(?:\+|00)86)?1(?:(?:3[\d])|(?:4[5-79])|(?:5[0-35-9])|(?:6[5-7])|(?:7[0-8])|(?:8[\d])|(?:9[1589]))\d{8}$`)
	return regexpPhone.MatchString(phone)
}

// 校验密码复杂度
func IsValidPassword(password string) bool {
	var (
		hasNumber    bool
		hasSymbol    bool
		hasLowerCase bool
		hasUpperCase bool
	)

	if len(password) < 8 {
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsUpper(char):
			hasUpperCase = true
		default:
			// do nothing
		}
	}

	return hasNumber && hasSymbol && hasLowerCase && hasUpperCase
}
