package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/r1is/ReportSystem/common"
	"github.com/r1is/ReportSystem/router"
	"gorm.io/gorm"
)

func main() {
	db := common.InitDB()

	r := gin.Default()
	router.CollectRoute(r)

	r.Run(":8081")
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
