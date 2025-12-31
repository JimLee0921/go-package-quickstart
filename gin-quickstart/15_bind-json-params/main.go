package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	type User struct {
		UserName string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
		Sex      string `json:"sex"`
		Age      int    `json:"age"`
	}

	router.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}
		if user.Sex == "" {
			user.Sex = "man"
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":  "successful",
			"user": user,
		})

	})

	router.Run(":8000")
}
