package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	// 这里也可以直接写 200，但是推荐还是使用 http 包的语义
	c.JSON(http.StatusOK, gin.H{
		"message": "this is index",
	})
}

func main() {
	// 1. 初始化
	// 设置 gin 的运行模式为 release 模式，也可以传入 "release"
	gin.SetMode(gin.ReleaseMode)
	// 程序入口 Default 返回一个都使用默认设置 *gin.Engine
	// 包括默认的日志记录器等
	// 可以使用 gin.New() 获取一个更干净的 *gin.Engine
	router := gin.Default()
	//router := gin.New()

	// 2. 挂载路由
	router.GET("/index", Index)

	// 3. 绑定端口，运行
	// :8000 就是内网都开业访问，设为 127.0.0.1:8000 只允许本机访问，不传入默认是 :8080
	router.Run(":8000")
}
