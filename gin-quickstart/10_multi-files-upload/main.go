package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload_files", func(c *gin.Context) {
		filesForm, err := c.MultipartForm()
		if err != nil {
			fmt.Println(err)
		}

		for _, headers := range filesForm.File {
			for _, header := range headers {
				c.SaveUploadedFile(header, "uploads/"+header.Filename)
			}
		}
	})
	r.Run(":8000")
}
