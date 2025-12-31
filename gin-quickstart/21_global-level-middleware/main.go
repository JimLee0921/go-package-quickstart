package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 全局中间件就是作用到路由组，给路由组添加中间件
	router := gin.Default() // Default 默认使用的 Logger 和 Recovery 中间件
	ApiRouter := router.Group("api", AuthMiddleware, MiddleWareOne, MiddleWareTwo)

	UserGroup(ApiRouter)

	router.Run(":8000")
}

func UserGroup(r *gin.RouterGroup) {
	r.GET("users", UserView)
	r.POST("users", UserView)
	r.PUT("users", UserView)
	r.DELETE("users", UserView)
}

func UserView(c *gin.Context) {
	fmt.Println(c.Request.Method, c.Request.URL)
	// 同样能拿到中间件中设置的值
	// 请求部份
	v, ok := c.Get("auth")
	if ok && v == "successful" {
		c.JSON(http.StatusOK, gin.H{
			"auth": v,
		})
	}
}

func MiddleWareOne(c *gin.Context) {
	// 请求部份
	v, ok := c.Get("auth")
	if ok && v == "successful" {
		fmt.Println("MiddleWareOne request... verify auth: ", v)
	}
	c.Next()
	// 响应部份
	fmt.Println("MiddleWareOne response...")
}

func MiddleWareTwo(c *gin.Context) {
	// 请求部份
	fmt.Println("MiddleWareOne request...")
	c.Next()
	// 响应部份
	fmt.Println("MiddleWareOne response...")
}

// AuthMiddleware 简单模拟校验请求头中的 token
func AuthMiddleware(c *gin.Context) {
	fmt.Println(c.Request.Header)
	token := c.Request.Header.Get("token")
	fmt.Println(token)
	if token != "my_token" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "用户校验失败",
		})
	}
	// 中间件之间传值
	c.Set("auth", "successful")
	c.Next()
}
