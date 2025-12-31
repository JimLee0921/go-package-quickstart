package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 表单参数 Post 请求的 form 表单
	r := gin.Default()

	r.POST("upload_file", func(c *gin.Context) {
		// 表单文件上传，py request 使用 request.post(flies=file)
		fileHeader, err := c.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fileHeader.Header, fileHeader.Size)

		// 传统手艺保存文件
		//file, err := fileHeader.Open()
		//byteData, _ := io.ReadAll(file)
		//err = os.WriteFile(fileHeader.Filename, byteData, 0666)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}

		// 新手艺保存文件
		err = c.SaveUploadedFile(fileHeader, "uploads/"+fileHeader.Filename)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	r.Run(":8000")
}
