package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"regexp"
	"time"
	"unicode"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20)"`
	Telephone string `gorm:"type:varchar(30);unique_index"`
	Password  string `gorm:"type:varchar(255)"`
}

// 数据库配置
const (
	DBUsername = "root"
	DBPassword = "123.bmk"
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "ReportSystem"
)

// 验证手机号
func isValidPhone(phone string) bool {
	// 2019 工信部公布的手机号
	regexpPhone := regexp.MustCompile(`^(?:(?:\+|00)86)?1(?:(?:3[\d])|(?:4[5-79])|(?:5[0-35-9])|(?:6[5-7])|(?:7[0-8])|(?:8[\d])|(?:9[1589]))\d{8}$`)
	return regexpPhone.MatchString(phone)
}

func main() {
	db := InitDB()

	r := gin.Default()

	r.POST("/api/auth/register.app", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		// 校验手机号
		if !isValidPhone(telephone) {
			ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "msg": "手机号格式错误"})
			return
		}
		
		//校验密码复杂度
		if !isValidPassword(password) {
			ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "msg": "密码复杂度不合要求: 密码长度不低于8位，必须同时有数字、特殊符号和大小写字母"})
			return
		}

		//检查手机号是否已经被注册
		if isPhoneExist(db, telephone) {
			ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "msg": "该手机号已注册"})
			return
		}

		if len(name) == 0 {
			name = RandString(4)
		}

		//创建用户
		newUser := User{
			Name:      name,
			Password:  hashPassword(password),
			Telephone: telephone,
		}

		// 保存用户
		db.Create(&newUser)

		//返回结果
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "user create success",
		})
	})
	r.POST("/api/auth/login", func(ctx *gin.Context) {
		// 获取参数
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		// 校验手机号
		if !isValidPhone(telephone) {
			ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "error": "手机号格式错误"})
			return
		}

		// 判断用户是否存在并验证密码
		var user User
		db.Where(&User{Telephone: telephone}).First(&user)
		if user.ID == 0 || !checkPasswordHash(password, user.Password) {
			ctx.JSON(200, gin.H{"code": http.StatusUnauthorized, "error": "Unauthorized"})
			return
		}

		// 返回结果
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "login success",
		})
	})
	r.Run(":8081")
	performDBOperations(db)
}

func RandString(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0987654321")
	results := make([]byte, n)

	rand.Seed(time.Now().Unix()) //时间戳随机种子
	for i := range results {
		results[i] = letters[rand.Intn(len(letters))]
	}

	return string(results)

}

// 初始化数据库
func InitDB() *gorm.DB {
	//"username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DBUsername,
		DBPassword,
		DBHost,
		DBPort,
		DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	_ = db.AutoMigrate(&User{})
	return db

}

// 校验密码复杂度
func isValidPassword(password string) bool {
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

// 对密码进行加密
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// 验证密码是否正确
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//检查phone是否已经注册
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user User
	db.Where("telephone = ?", phone).First(&user)
	return user.ID != 0
}

//关闭数据库
func performDBOperations(db *gorm.DB) {
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println("failed to close database")
			return
		}
		sqlDB.Close()
	}()

	// Do something with the database
}
