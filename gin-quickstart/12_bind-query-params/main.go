package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("", func(c *gin.Context) {
		type User struct {
			Name string `form:"name"`
			Age  int    `form:"age"`
		}
		var user User
		// 自动做结构体字段校验
		err := c.ShouldBindQuery(&user)
		
		fmt.Println(user, err)
	})

	r.Run(":8000")
}
