package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("users/:id/:name", func(c *gin.Context) {
		var user struct {
			ID   int    `uri:"id" binding:"required"` // binding 为 required 说明该字段不能为空
			Name string `uri:"name" binding:"required"`
		}
		// 路径参数绑定
		if err := c.ShouldBindUri(&user); err != nil {
			c.JSON(400, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"name": user.Name,
			"id":   user.ID,
		})

		fmt.Println(user)
	})
	router.Run(":8000")
}
