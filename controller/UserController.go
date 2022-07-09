package controller

import (
	"Gin/common"
	"Gin/dto"
	"Gin/model"
	"Gin/response"
	"Gin/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"}) //gin.H是map[string]interface{}的别名
		//ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code": 422, "msg": "手机号必须为11位"})
		//统一封装response格式调用
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码需要大于6位"}) //gin.H是map[string]interface{}的别名
		//统一封装response格式调用
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码需要大于6位")
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Printf(name, telephone, password)
	//判断手机号数否存在
	if isTelephoneExist(db, telephone) {
		//ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code": 422, "msg": "用户手机号已经存在"})
		//统一封装response格式调用
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户手机号已经存在")
		return
	}
	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		//ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code": 422, "msg": "加密错误"})
		//统一封装response格式调用
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	db.Create(&newUser)
	//返回结果
	//ctx.JSON(200, gin.H{
	//	"code":    200,
	//	"message": "注册成功",
	//})
	//统一封装response格式调用
	response.Response(ctx, http.StatusUnprocessableEntity, 200, nil, "注册成功")
}
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	//	获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//	数据验证
	if len(telephone) != 11 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"}) //gin.H是map[string]interface{}的别名
		//ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code": 422, "msg": "手机号必须为11位"})
		//统一封装response格式调用
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码需要大于6位"}) //gin.H是map[string]interface{}的别名
		//统一封装response格式调用
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码需要大于6位")
		return
	}

	//	判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"}) //gin.H是map[string]interface{}的别名
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	//	判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 400, "msg": "密码错误"}) //gin.H是map[string]interface{}的别名
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码错误")
		return
	}
	//	发送token
	token, err := common.ReleaseToken(user)
	if err != nil {
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "生成token失败"}) //gin.H是map[string]interface{}的别名
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "生成token失败")
		log.Printf("生成token失败：%v", err)
		return
	}
	//	返回结果
	//ctx.JSON(200, gin.H{
	//	"code":    200,
	//	"message": "登录成功",
	//	"data":    gin.H{"token": token},
	//})
	response.Response(ctx, http.StatusUnprocessableEntity, 200, gin.H{"token": token}, "登录成功")
}
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 200, "msg": "用户信息查询成功", "data": gin.H{"user": dto.ToUserDto(user.(model.User))}}) //gin.H是map[string]interface{}的别名
	response.Response(ctx, http.StatusUnprocessableEntity, 200, gin.H{"user": dto.ToUserDto(user.(model.User))}, "用户信息查询成功")
}
