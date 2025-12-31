package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	// 引入静态文件目录，第一个参数是路由别名（通常也命名为 static），第二个参数是目录
	// 访问 http://127.0.0.1:8000/st/文件名 进行访问
	// 注意静态文件的路径不能再被路由使用，也就是不能再写 r.GET("st", ...)
	r.Static("st", "static")
	// 引入单独静态文件
	// 访问 http://127.0.0.1:8000/note 进行访问
	r.StaticFile("note", "static/note.txt")
	r.Run(":8000")
}
