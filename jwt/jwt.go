package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	//"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
	//"net/http"
)
type MyClaims struct {
	StudentId string `json:"student_id"`
	jwt.StandardClaims
}
//token 过期时间
const TokenExpireDuration = time.Hour * 48
var MySecret = []byte("夏天夏天悄悄过去")
// GenToken 生成JWT
func GenToken(StudentId string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		StudentId, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}
// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	fmt.Println("解析Token")
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Token开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2003,
			//	"msg":  "请求头中auth为空",
			//})
			c.Next()
			return
		}
		//// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Token") {
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2004,
			//	"msg":  "请求头中auth格式有误",
			//})
			c.Next()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2005,
			//	"msg":  "无效的Token",
			//})
			//fmt.Println(err.Error())
			c.Next()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("StudentId", mc.StudentId)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}



