package main

import (
	"fmt"
	"ginWeb/jwt"
	db "ginWeb/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ErrorJson struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func addUser() {

}

func main() {
	//fmt.Println(jwt.GenToken("朱宇宸"))
	r := gin.Default()
	r.LoadHTMLFiles()
	r.Use(CROSHandler()) //跨域中间件
	r.GET("/token", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		student_id, ok := StudentId.(string)
		newUser := &db.User{
			StudentId: student_id,
		}
		if ok {
			newUser = newUser.QueryUserInfo()
		}
		context.JSON(200, map[string]interface{}{
			"code":         1,
			"student_id":   StudentId,
			"student_name": newUser.Username,
			"with_token":   ok,
		})
	})
	userRoute := r.Group("/user")

	userRoute.GET("/query", func(context *gin.Context) {

	})
	userRoute.GET("/check", func(context *gin.Context) {
		newUser := &db.User{}
		studentId := context.Query("student_id")
		if len(studentId) != 11 {
			context.JSON(200, map[string]interface{}{
				"code":    -1,
				"message": "请输入正确的学生证",
			})
			return
		}
		newUser.StudentId = studentId
		id := newUser.QueryUser()
		if strings.TrimSpace(id) != "" {
			context.JSON(200, map[string]interface{}{
				"code":    -1,
				"message": "该学生证已经被注册",
			})
			return
		} else {
			context.JSON(200, map[string]interface{}{
				"code":    1,
				"message": "该学生证可以使用",
			})
			return
		}

	})
	userRoute.POST("/login", func(context *gin.Context) {
		newUser := &db.User{}
		errJson := &ErrorJson{
			"nothing",
			0,
		}
		studentId, _ := context.GetPostForm("student_id")
		if len(studentId) != 11 {
			errJson.Message = "学生证位数有误"
			errJson.Code = -1
			context.JSON(200, errJson)
			return
		}
		newUser.StudentId = studentId
		password, _ := context.GetPostForm("password")
		if len(password) > 20 || len(password) < 8 {
			errJson.Message = "密码长度有误"
			errJson.Code = -1
			context.JSON(200, errJson)
			return
		}
		newUser.Password = password
		username, err := newUser.GetUser()
		if err != nil {
			errJson.Message = err.Error()
			errJson.Code = -1
			context.JSON(400, map[string]interface{}{
				"code":    -1,
				"message": "数据库操作失败",
				"error":   err.Error(),
			})
		}
		if username == "" {
			context.JSON(200, map[string]interface{}{
				"code":    -1,
				"message": "账号或密码错误",
			})
			return
		}
		token, err := jwt.GenToken(studentId)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "token生成错误",
				"code":    -1,
			})
			return
		}
		context.JSON(200, gin.H{
			"message": "用户" + username + "登陆成功",
			"code":    1,
			"token":   token,
		})
		return
	})
	userRoute.POST("/register", func(context *gin.Context) {
		newUser := &db.User{}
		errJson := &ErrorJson{
			"nothing",
			0,
		}
		username, ok := context.GetPostForm("username")
		if !ok || strings.TrimSpace(username) == "" {
			errJson.Message = "请输入姓名"
			errJson.Code = -1
			context.JSON(200, errJson)
			return
		} else if len(username) > 20 {
			errJson.Message = "名字过长"
			errJson.Code = -1
			context.JSON(200, errJson)
			return
		}
		newUser.Username = username
		password, ok := context.GetPostForm("password")
		if !ok {
			errJson.Message = "请输入密码"
			errJson.Code = -1
			context.JSON(200, errJson)
			return
		} else if len(password) > 20 || len(password) < 8 {
			errJson.Message = "密码位数在8-20位"
			errJson.Code = -1
			context.JSON(200, errJson)
			return

		}
		newUser.Password = password
		studentId, ok := context.GetPostForm("student_id")
		if !ok {
			errJson.Message = "请输入学生证"
			errJson.Code = -1
			context.JSON(200, errJson)
			return

		} else if len(studentId) != 11 {
			errJson.Message = "请输入正确的学生证(位数不对)"
			errJson.Code = -1
			context.JSON(200, errJson)
			return
		}
		newUser.StudentId = studentId
		err := newUser.AddUser()
		if err != nil {
			errJson.Message = err.Error()
			errJson.Code = -1
			context.JSON(400, map[string]interface{}{
				"code":    -1,
				"message": "数据库操作失败",
				"error":   err.Error(),
			})
		}
		errJson.Message = "用户" + username + "添加成功"
		errJson.Code = 1
		context.JSON(200, errJson)
		return

		//context.JSON(200, map[string]interface{}{
		//	"params": "success",
		//	"code":   200,
		//})
	})
	userRoute.GET("/update/position", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
			})
			return
		}
		longitude := context.Query("longitude")
		if strings.TrimSpace(longitude) == "" {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "经度不存在",
			})
			return
		}
		latitude := context.Query("longitude")
		if strings.TrimSpace(latitude) == "" {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "纬度不存在",
			})
			return
		}
		id, ok := StudentId.(string)
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
			})
			return
		}
		newUser := &db.User{
			StudentId: id,
			Position:  longitude + ";" + latitude,
		}
		err := newUser.UpdateUserPosition()
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
		}
		context.JSON(200, gin.H{
			"code":    1,
			"message": "用户位置更新成功",
		})
	})

	r.Run(":9999")
}

func CROSHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "*")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Expose-Headers", "*")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Header("Access-Control-Request-Headers", "*")

		context.Set("content-type", "application/json")
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		//处理请求
		context.Next()
	}
}
