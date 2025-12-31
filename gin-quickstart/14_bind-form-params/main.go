package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/user/login", func(c *gin.Context) {
		type LoginForm struct {
			UserName   string    `form:"user_name,default=admin"`
			Password   string    `form:"password" binding:"required"`
			VerifyCode int       `form:"verify_code" binding:"required"`
			Addresses  [2]string `form:"address,default=chain;shanghai"` // 设置默认值
			LapTimes   []int     `form:"lap_times,default=1;2;3"`
		}
		var loginForm LoginForm
		if err := c.ShouldBind(&loginForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}
		if loginForm.VerifyCode != 66666 {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "codeWrong",
			})
			return
		}
		if loginForm.UserName != "admin" || loginForm.Password != "dsb666..." {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "username or password wrong",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":  "login successful",
			"form": loginForm,
		})
	})
	router.Run(":8000")
}
