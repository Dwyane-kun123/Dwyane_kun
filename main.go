package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
	_"github.com/go-sql-driver/mysql"
)

type User struct {
	gorm.Model //代码中定义模型（Models）与数据库中的数据表进行映射
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null"`
	Password string `gorm:"size:255";not null`
}

func main() {

	db := InitDB()
	defer db.Close()

	r := gin.Default() //
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		name:= ctx.PostForm("name")
		telephone:= ctx.PostForm("telephone")
		password:= ctx.PostForm("password")

		if len(telephone) != 11{
			fmt.Println(len(telephone),telephone)
			fmt.Println(password)
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code":422,"msg":"电话必须是11位"})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code":422,"msg":"密码不能小于6位"})
			return
		}
		if len(name) == 0  {
			name = RandomString(10)
		}

		//查看手机号是否存在
		if IsTelephineExit(db, telephone){ //这是在查询数据库了
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code":422,"msg":"用户已经存在，不允许注册"})
			return
		}

		newUser:= User{
			Name: name,
			Telephone: telephone,
			Password: password,
		}
		db.Create(&newUser)
		//返回结果
		ctx.JSON(200,gin.H{"msg":"注册成功"})

		log.Println(name,telephone,password)
		return
	}) //func(c *gin.Context) 是为r.GET这个方法提供具体的操作
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func RandomString(n int) string  {
	letters := "qazwsxedcrfvtgbyhnujmiklopQAZWSXEDCRFVTGBYHNUJMIKOLP"
	res := make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i,_ := range res{
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}

func InitDB() *gorm.DB{
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "xukun"
	username := "root"
	password := "111111"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
	)

	db, err := gorm.Open(driverName,args)
	if err != nil {
		panic("连接数据库错误，err" + err.Error())
	}

	//自动创建数据表
	db.AutoMigrate(&User{})
	return db
}

func IsTelephineExit(db *gorm.DB,telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}