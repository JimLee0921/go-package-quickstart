package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 加载 templates 下的所有 .html 文件
	// 注意路径问题
	r.LoadHTMLGlob("templates/*")
	// 加载单个 HTML 文件
	//r.LoadHTMLFiles("templates/index.html")
	r.GET("/index", func(c *gin.Context) {
		// 第一个参数为状态码
		// 第二个参数为文件名（不是文件路径，必须先使用有 r.LoadHTMLXxx 或 r.LoadHTMLFiles）
		// 第三个参数是传递给 HTML 文件的参数（类似 django 模板语法，前后端分离使用不多）
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/users", func(c *gin.Context) {
		users := []gin.H{
			{"name": "JimLee", "age": 12},
			{"name": "DSB", "age": 23},
			{"name": "LS", "age": 44},
		}
		c.HTML(http.StatusOK, "users.html", gin.H{
			"title": "用户列表",
			"users": users,
		})
	})

	r.Run(":8000")
}
