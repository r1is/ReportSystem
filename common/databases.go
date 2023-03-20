package common

import (
	"fmt"
	"github.com/r1is/ReportSystem/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 数据库配置
const (
	DBUsername = "root"
	DBPassword = "123.bmk"
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "ReportSystem"
)

var DB *gorm.DB

// InitDB 初始化数据库
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
		log.Fatalln("failed to connect to database")
	}
	_ = db.AutoMigrate(&model.User{})
	DB = db
	return db
}

// GetDB 获取一个数据库实例
func GetDB() *gorm.DB {
	return DB
}
