package main

import (
	"./jwt"
	db "./utils"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ErrorJson struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func main() {
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
			"signature":    newUser.Signature,
			"longitude":    newUser.Longitude,
			"latitude":     newUser.Latitude,
			"with_token":   ok,
		})
	})
	userRoute := r.Group("/user")
	messageRoute := r.Group("/message")
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
		id := newUser.QueryUserId()
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
	userRoute.GET("/position", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
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
		}
		newUser.QueryUserPosition()
		if newUser.Longitude == 0 && newUser.Latitude == 0 {
			context.JSON(200, gin.H{
				"code":    -2,
				"message": "你还没有登记自己的位置信息哦,先去更新一下吧",
			})
			return
		}
		userList, err := db.GetAllUserInfo(newUser.Longitude,newUser.Latitude,newUser.StudentId)
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "查询位置信息失败",
			})
			return
		}
		context.JSON(200, gin.H{
			"code":      1,
			"message":   "success",
			"user_list": userList,
			"user_lng":  newUser.Longitude,
			"user_lat":  newUser.Latitude,
		})
		return
	})
	userRoute.GET("/update/info", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
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
		longitudeStr := context.Query("longitude")
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if longitude == 0 || err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "经度有误",
			})
			return
		}

		latitudeStr := context.Query("latitude")
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if latitude == 0 || err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "经度有误",
			})
			return
		}

		signature := context.Query("signature")
		if len(signature) > 50 {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "你的个性签名太长啦，换个短一点的吧",
			})
		}
		newUser := &db.User{
			StudentId: id,
			Longitude: longitude,
			Latitude:  latitude,
			Signature: signature,
		}
		err = newUser.UpdateUserInfo()
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
		}
		signature = html.EscapeString(signature)
		context.JSON(200, gin.H{
			"code":      1,
			"message":   "用户信息更新成功",
			"signature": signature,
		})
	})
	userRoute.GET("/posinfo", func(context *gin.Context) {
		position := context.Query("position")
		if strings.TrimSpace(position) == "" {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "位置错误,请稍后再试!",
			})
			return
		}
		student_id := context.Query("student_id")
		if strings.TrimSpace(student_id) == "" {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "账户ID错误,请稍后再试!",
			})
			return
		}
		newUser := &db.User{
			StudentId: student_id,
		}
		newUser.QueryAllInfo()
		//userList := db.Getdistance(position, newUser)
		context.JSON(200, gin.H{
			"code":     1,
			"message":  "success",
			"userinfo": newUser,
		})
		return
	})
	messageRoute.POST("/leave", func(context *gin.Context) {
		receiverId, ok := context.GetPostForm("receiver_id")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误,请稍后再试",
			})
			return
		}
		message, ok := context.GetPostForm("message")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "消息不能为空喔",
			})
			return

		}
		senderId, ok := context.GetPostForm("sender_id")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误,请稍后再试",
			})
			return
		}
		if strings.TrimSpace(message) == "" {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "消息不能为空喔",
			})
			return
		}
		if receiverId == "" || senderId == "" {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误,请稍后再试",
			})
			return
		}
		messageTime := time.Now().Format("2006-01-02 15:04:05")
		newMessage := &db.MessageType{
			SenderId:   senderId,
			ReceiverId: receiverId,
			SendTime:   messageTime,
			Message:    message,
		}
		err := newMessage.AddMessage()
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "数据库错误,请稍后再试",
			})
			return
		}
		context.JSON(200, gin.H{
			"code":    1,
			"message": "留言成功，请耐心等待回复哦",
		})
	})
	messageRoute.GET("/mymessage", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
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
		MyMessageType := &db.MessageType{
			ReceiverId: id,
		}
		MessageList, err := MyMessageType.GetMyMessage()
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "查询信息有误，请稍后再试",
			})
			return
		}
		context.JSON(200, gin.H{
			"code":         1,
			"message":      "查询成功",
			"message_list": MessageList,
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
