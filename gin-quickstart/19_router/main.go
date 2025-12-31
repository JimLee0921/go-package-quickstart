package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 创建路由分组
	apiRouterWithMiddleWare := router.Group("api")
	// 路由分组可以注册统一中间件(认证相关)
	apiRouterWithMiddleWare.Use(AuthMiddleWare)
	// 把需要认证的路由分组交给 UserGroup 进行分组注册
	UserGroup(apiRouterWithMiddleWare)

	// 还可以创建同样命名的路由
	apiRouterNoMiddleWare := router.Group("api")
	// 登录路由分组不走中间件
	LoginGroup(apiRouterNoMiddleWare)

	router.Run(":8000")
}

func UserView(context *gin.Context) {
	url := context.Request.URL
	method := context.Request.Method
	fmt.Printf("url：%s; method: %s\n", url, method)
}

// UserGroup 用户路由分组
// RESTFUL 风格路由规范，只是规范，不是标准
// GET    /api/users     用户列表
// POST   /api/users     创建用户
// PUT    /api/users/:id 更新用户
// DELETE /api/users/:id 删除单个用户
// DELETE /api/users	 批量删除
func UserGroup(r *gin.RouterGroup) {
	r.GET("users", UserView)
	r.POST("users", UserView)
	r.PUT("users", UserView)
	r.DELETE("users", UserView)
	// 支持所有请求类型
	// r.Any()
}

func LoginGroup(r *gin.RouterGroup) {
	r.GET("login", UserView)
}

func AuthMiddleWare(context *gin.Context) {
	fmt.Println("this is auth middleware")
}
