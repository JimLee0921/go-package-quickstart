package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("", func(c *gin.Context) {
		byteData, _ := io.ReadAll(c.Request.Body)
		// 阅后即焚
		fmt.Println(string(byteData))

		// 重新写入
		c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	})

	r.Run(":8000")
}
