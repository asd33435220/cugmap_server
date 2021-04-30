package main

import (
	db "./db"
	"./jwt"
	"fmt"
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
	placeRoute := r.Group("/place")
	commentRoute := r.Group("/comment")

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
		userList, err := db.GetAllUserInfo(newUser.Longitude, newUser.Latitude, newUser.StudentId)
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
	userRoute.GET("/position2", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
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
		userList, err := db.GetAllUserInfo2(newUser.Longitude, newUser.Latitude, newUser.StudentId)
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
		if strings.TrimSpace(message) == "" || len(message) > 50 {
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
		messageTime := strconv.Itoa(int(time.Now().UnixNano() / 1e6))
		messageTimeStr := time.Now().Format("2006-01-02 15:04:05")
		placeCodeStr, ok := context.GetPostForm("place_code")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误,请稍后再试",
			})
			return
		}
		placeCode, err := strconv.Atoi(placeCodeStr)

		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误,请稍后再试",
			})
			return
		}

		newMessage := &db.MessageType{
			SenderId:    senderId,
			ReceiverId:  receiverId,
			SendTimeStr: messageTimeStr,
			SendTime:    messageTime,
			Message:     message,
			PlaceCode:   int64(placeCode),
		}
		err = newMessage.AddMessage()
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
		ReceiverMessageList, SenderMessageList, err := MyMessageType.GetMyMessage()
		AllMyMessageList := append(ReceiverMessageList, SenderMessageList...)
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
			"message_list": AllMyMessageList,
		})
	})
	messageRoute.GET("/allmymessage", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
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
		messageList, err := MyMessageType.GetAllMyMessage()
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
			"message_list": messageList,
		})
	})

	messageRoute.GET("/read", func(context *gin.Context) {
		ReceiverId := context.Query("ReceiverId")
		SenderId := context.Query("SenderId")
		if ReceiverId == "" || SenderId == "" {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "信息错误",
			})
		}
		MyMessageType := &db.MessageType{
			ReceiverId: ReceiverId,
			SenderId:   SenderId,
		}
		err := MyMessageType.UpdateMessage()
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "查询信息有误，请稍后再试",
			})
			return
		}
		context.JSON(200, gin.H{
			"code":    1,
			"message": "更新已读状态成功",
		})
	})
	placeRoute.POST("/add", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
			})
			return
		}
		Founder, ok := StudentId.(string)
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地点名字有误，请修改后重试",
			})
			return
		}
		Name, ok := context.GetPostForm("Name")
		if len(Name) > 50 || strings.TrimSpace(Name) == "" || !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地物名字有误，请修改后重试",
			})
			return
		}
		Address, ok := context.GetPostForm("Address")
		fmt.Println(len(Address))
		if len(Address) > 200 || strings.TrimSpace(Address) == "" || !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "详细地址有误，请修改后重试",
			})
			return
		}
		Lng, ok := context.GetPostForm("Lng")
		Longitude, err := strconv.ParseFloat(Lng, 64)
		if !ok || err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地物经度有误，请修改后重试",
			})
			return
		}
		Lat, ok := context.GetPostForm("Lat")
		Latitude, err := strconv.ParseFloat(Lat, 64)
		if !ok || err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地物纬度有误，请修改后重试",
			})
			return
		}
		Number, ok := context.GetPostForm("Number")
		if len(Number) > 15 || strings.TrimSpace(Number) == "" || !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "联系方式有误，请修改后重试",
			})
			return
		}
		Image1Url, ok := context.GetPostForm("Image1Url")
		if len(Image1Url) > 3000 || !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "图片1过长，请修改后重试",
			})
			return
		}
		Image2Url, ok := context.GetPostForm("Image2Url")
		if len(Image2Url) > 3000 || !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "图片2过长，请修改后重试",
			})
			return
		}
		Comment, ok := context.GetPostForm("Comment")
		if len(Comment) > 200 || strings.TrimSpace(Comment) == "" || !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "评论长度有误，请修改后重试",
			})
			return
		}
		TypeStr, ok := context.GetPostForm("Type")
		Type, err := strconv.Atoi(TypeStr)
		if !ok || err != nil || Type > 3 || Type < 0 {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地物类型有误，请修改后重试",
			})
			return
		}
		ScoreStr, ok := context.GetPostForm("Score")
		Score, err := strconv.ParseFloat(ScoreStr, 64)
		if !ok || err != nil || Score > 5 || Score < 0 {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地物评分有误，请修改后重试",
			})
			return
		}

		newPlace := &db.Place{
			Name:           Name,
			Address:        Address,
			Longitude:      Longitude,
			Latitude:       Latitude,
			Image1Url:      Image1Url,
			Image2Url:      Image2Url,
			Type:           Type,
			Score:          Score,
			Number:         Number,
			Founder:        Founder,
			FounderComment: Comment,
		}
		err = newPlace.AddPlace()
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误，请重试",
			})
			return
		}

		context.JSON(200, gin.H{
			"code":    1,
			"message": "地物更新成功！",
		})
		return

	})
	placeRoute.GET("/info", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
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
		newUser.QueryAllInfo()
		placeInfoList, err := db.GetPlace(newUser)
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误，请重试",
			})
			return
		}

		context.JSON(200, gin.H{
			"code":       1,
			"message":    "地物查询成功！",
			"place_info": placeInfoList,
		})
		return

	})
	placeRoute.GET("/chat", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
			})
			return
		}
		_, ok = StudentId.(string)
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
			})
			return
		}
		placeCode := context.Query("placeCode")
		placeInfo, err := db.GetOnePlace(placeCode)
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误，请重试",
			})
			return
		}
		context.JSON(200, gin.H{
			"code":       1,
			"message":    "地物查询成功！",
			"place_info": placeInfo,
		})
		return
	})

	commentRoute.GET("/all", func(context *gin.Context) {
		PlaceCodeStr := context.Query("placeCode")
		_, err := strconv.Atoi(PlaceCodeStr)
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地物代号有误!",
			})
		}

		CommentList, err := db.GetComment(PlaceCodeStr)
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误，请稍后重试",
			})
		}
		context.JSON(200, gin.H{
			"code":         1,
			"message":      "地物评论查询成功",
			"comment_list": CommentList,
		})

	})
	commentRoute.POST("/add", jwt.JWTAuthMiddleware(), func(context *gin.Context) {
		StudentId, ok := context.Get("StudentId")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
			})
			return
		}
		Commentator, ok := StudentId.(string)
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "用户状态有误，请重新登陆",
			})
			return
		}
		PlaceCode, ok := context.GetPostForm("placeCode")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "地物代号有误",
			})
			return
		}
		CommentMessage, ok := context.GetPostForm("commentMessage")
		if !ok {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "评论信息有误",
			})
			return
		}
		ScoreStr, ok := context.GetPostForm("score")
		Score, err := strconv.Atoi(ScoreStr)
		if !ok || err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "评分有误",
			})
			return
		}
		CommentTime := strconv.Itoa(int(time.Now().UnixNano() / 1e6))
		CommentTimeStr := time.Now().Format("2006-01-02 15:04:05")
		newComment := &db.Comment{
			PlaceCode:      PlaceCode,
			CommentTime:    CommentTime,
			CommentTimeStr: CommentTimeStr,
			Score:          Score,
			Commentator:    Commentator,
			CommentMessage: CommentMessage,
		}
		err = newComment.AddComment()
		if err != nil {
			context.JSON(200, gin.H{
				"code":    -1,
				"message": "服务器错误，请稍后重试",
			})
			return
		}
		context.JSON(200, gin.H{
			"code":    -1,
			"message": "评论添加成功！",
		})
	})
	//commentRoute.GET("/like", func(context *gin.Context) {
	//
	//})
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
