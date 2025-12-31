package main

import (
	"gin-quickstart/02_json-response-encapsulation/response"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 创建实例
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 2. 绑定路由
	r.GET("/login", func(c *gin.Context) {
		response.OkWithMsg(c, "登陆成功")
	})

	r.GET("/users", func(c *gin.Context) {
		response.OkWithData(c, gin.H{
			"admin": "JimLee",
			"user1": "DSB",
			"user2": "LS",
		})
	})

	r.POST("/users", func(c *gin.Context) {
		response.FailWithMessage(c, "param wrong")
	})
	// 3. 启动
	r.Run(":8000")
}
