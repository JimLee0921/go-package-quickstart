package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReqHeader struct {
	UserAgent   string `header:"User-Agent"`
	Name        string `header:"Name"`
	ContentType string `header:"Content-Type"`
}

func main() {
	router := gin.Default()
	router.POST("/headers", func(c *gin.Context) {
		var reqHeader ReqHeader
		if err := c.ShouldBindHeader(&reqHeader); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":     "successful",
			"headers": reqHeader,
		})
	})
	router.Run(":8000")
}
