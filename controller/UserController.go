package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/r1is/ReportSystem/common"
	"github.com/r1is/ReportSystem/model"
	"github.com/r1is/ReportSystem/util"
	"gorm.io/gorm"
	"net/http"
)

func Register(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 校验手机号
	if !util.IsValidPhone(telephone) {
		ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "msg": "手机号格式错误"})
		return
	}

	//校验密码复杂度
	if !util.IsValidPassword(password) {
		ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "msg": "密码复杂度不合要求: 密码长度不低于8位，必须同时有数字、特殊符号和大小写字母"})
		return
	}

	//检查手机号是否已经被注册
	if isPhoneExist(db, telephone) {
		ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "msg": "该手机号已注册"})
		return
	}

	if len(name) == 0 {
		name = util.RandString(4)
	}

	//创建用户
	newUser := model.User{
		Name:      name,
		Password:  util.HashPassword(password),
		Telephone: telephone,
	}

	// 保存用户
	db.Create(&newUser)

	//返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "user create success",
	})
}

func Login(ctx *gin.Context) {
	db := common.GetDB()
	// 获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 校验手机号
	if !util.IsValidPhone(telephone) {
		ctx.JSON(200, gin.H{"code": http.StatusBadRequest, "error": "手机号格式错误"})
		return
	}

	// 判断用户是否存在并验证密码
	var user model.User
	db.Where(&model.User{Telephone: telephone}).First(&user)
	if user.ID == 0 || !util.CheckPasswordHash(password, user.Password) {
		ctx.JSON(200, gin.H{"code": http.StatusUnauthorized, "error": "Unauthorized"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "login success",
	})
}

//检查phone是否已经注册
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("telephone = ?", phone).First(&user)
	return user.ID != 0
}
