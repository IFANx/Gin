package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxzcvbnmqwertyuiopZXCVBNMASDFGHJKLQWERTYYUIUOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique""`
	Password  string `gorm:"size:255;not null"`
}

func InitDB() *gorm.DB {
	dsn := "root:13458482147@tcp(127.0.0.1:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}

	db.AutoMigrate(&User{})
	return db
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func main() {
	db := InitDB()
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		//数据验证
		if len(telephone) != 11 {
			//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"}) //gin.H是map[string]interface{}的别名
			ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code": 422, "msg": "手机号必须为11位"})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码需要大于6位"}) //gin.H是map[string]interface{}的别名
			return
		}
		if len(name) == 0 {
			name = RandomString(10)
		}
		log.Printf(name, telephone, password)
		//判断手机号数否存在
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code": 422, "msg": "用户手机号已经存在"})
			return
		}
		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)
		//返回结果
		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}
