package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 获取查询参数 ?name=jimlee&age=3&hobby=dance&hobby=gaming
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		// 获取 name 查询参数的值
		name := c.Query("name")
		// 可以传默认值
		age := c.DefaultQuery("age", "3")
		// 获取多个一样的查询参数的值形成切片
		hobbies := c.QueryArray("hobby")

		fmt.Println(name, age, hobbies)
	})

	r.Run(":8000")
}
