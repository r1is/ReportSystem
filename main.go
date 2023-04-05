package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/r1is/ReportSystem/common"
	"github.com/r1is/ReportSystem/config"
	"github.com/r1is/ReportSystem/router"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	config.InitConfig()
	db := common.InitDB()

	r := gin.Default()
	router.CollectRoute(r)
	port := viper.GetString("server.port")

	r.Run(":" + port)
	performDBOperations(db)
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
