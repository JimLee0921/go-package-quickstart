package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 表单参数 Post 请求的 form 表单
	r := gin.Default()

	// 创建用户表单
	r.POST("login", func(c *gin.Context) {
		username := c.PostForm("username")           // 直接获取。可能为空
		password, ok := c.GetPostForm("password")    // 安全获取
		hobbies, ok := c.GetPostFormArray("hobbies") // 接收列表
		//friends, ok := c.GetPostFormMap("friends")   // 接收字典，使用python 定义为：
		/*
			{
			    "user[name]": "JimLee",
			    "user[role]": "admin",
			    "user[age]": "30",
			}
		*/

		fmt.Println(username)
		fmt.Println(password, ok)
		fmt.Println(hobbies, ok)
	})

	r.Run(":8000")
}
