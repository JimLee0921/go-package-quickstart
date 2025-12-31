package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 获取path路径参数
	// /users?id=123 查询参数模式
	// /users/123 路径参数模式
	r := gin.Default()

	r.GET("/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		fmt.Println(userID)
	})
	// /users/123/大傻逼
	r.GET("/users/:id/:name", func(c *gin.Context) {
		userId := c.Param("id")
		userName := c.Param("name")
		fmt.Println(userId, userName)
	})

	r.Run(":8000")
}
