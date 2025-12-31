package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// 方法1
	// 需要设置 Content-Type 且文件路径不存在会返回 404
	r.GET("", func(c *gin.Context) {
		// 表示为文件流，
		c.Header("Content-Type", "application/octet-stream")
		// 用于表示下载下来的文件名，不设置就是传入的文件名
		c.Header("Content-Disposition", "attachment; filename=dsb.txt")
		// 如果文件路径不存在会返回 404
		c.File("static/note.txt")
	})

	// 方法2：更常用，前端请求后端接口然后唤起浏览器下载
	// 相当于 a 标签 href 对应文件地址 download 对应下载文件名

	// 方法3 最好的做法是后端不返回实际文件内容，而是生成一个临时下载地址
	// 前端构造 a 标签，再请求这个接口唤起浏览器下载
	r.Run(":8000")
}
