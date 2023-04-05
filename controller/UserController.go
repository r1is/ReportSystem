package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/r1is/ReportSystem/common"
	"github.com/r1is/ReportSystem/dto"
	"github.com/r1is/ReportSystem/model"
	"github.com/r1is/ReportSystem/util"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Cookie(ctx *gin.Context) {
	a1 := "5o5VXExig66ROU71BBcm.bNNpbqVBN9KTv5jQlKmOqMS9umlXppjBRMvBdihmZOhO.fnL7IlgdFRlKuQW00ABkG"
	a2 := "iyc8wUkTH18hfBIoW2U5h8cga0uiaNgkCXp1zPO7aEzcVlq00URo7A_SeMSVvZTdMfkjz1wq3OsPaK1OJR2o31OlkHRnd5bM8_pEFJoO0Xq9oCFWVruv2RUB0k9Qbcv6loCs4DQFVpzR3IbaXzDPsu3HzVB3Tch.mpTAA0Zq_Dl_UjZf25p.iQoztvN9ZYDcF4gni6Y8S1yY3wPdR6sxjpyEzK1mGhdpqRRTsxAbllP.f252fFW2kTirBas.8PtpbDfXoFBaAGu3TDP0DQZqqWjoPNtnuofRKaGmvFU.Q31b08EbaG7mh2.myZY_KjlvOe2TlmtvwYS5vmR75vcJvXQUFyME.kFnL2XYDT7vD_nfMSB9Ka9wk4w9us1Uh4aowZyXxno34VhSleKp8YdAkjQrBNJ9vJSijC892qTMfgg"
	ctx.JSON(200, gin.H{"data": gin.H{"cookies": gin.H{"ruic2m5TbbKpO": a1, "ruic2m5TbbKpP": a2}}, "status": 1})
}

// Register 注册用户
func Register(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 校验手机号
	if !util.IsValidPhone(telephone) {
		common.Failed(ctx, nil, "手机号格式错误")
		return
	}

	//校验密码复杂度
	if !util.IsValidPassword(password) {
		common.Failed(ctx, nil, "密码复杂度不合要求: 密码长度不低于8位，必须同时有数字、特殊符号和大小写字母")
		return
	}

	//检查手机号是否已经被注册
	if isPhoneExist(db, telephone) {
		common.Failed(ctx, nil, "该手机号已注册")
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
	common.Success(ctx, gin.H{"UserName": name, "Phone": telephone}, "user create success")
	//ctx.JSON(200, gin.H{
	//	"code": 200,
	//	"msg":  "user create success",
	//})
}

// Login 用户登录
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

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(200, gin.H{"code": http.StatusInternalServerError, "msg": "系统错误"})
		log.Fatalf("token genrate error: %v\n", err)
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{"token": token,
			"LoginInfo": gin.H{
				"id":       user.ID,
				"username": user.Name,
				"phone":    user.Telephone,
			}},
		"msg": "login success",
	})

}

// UserInfo 获取用户信息
func UserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(200, gin.H{"code": http.StatusOK, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

//检查phone是否已经注册
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("telephone = ?", phone).First(&user)
	return user.ID != 0
}
