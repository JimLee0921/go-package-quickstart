package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 中间件中 Next() 负责继续执行，Abort() 负责中断链路，AbortWithXxx() 是中断 + 返回状态码/JSON 等信息
	// 这三个核心方法都是只在请求阶段有意义；在响应阶段调用没有实际效果
	// 可以在真正的视图函数前调用多个中间件，执行顺序为：
	/*
		1. 请求进入第一个中间件
		2. 遇到 c.Next() 会把控制权交给下一个处理器
		3. 进入第二个中间件开始执行遇到 c.Next() 继续向后调用
		4. 执行最终 Handler 视图函数
		5. 依次先回到第二个中间件响应阶段
		6. 回到第一个中间件想要阶段
	*/
	// 中间件中必须调用 c.Next() 才会一次执行
	// 如果某个中间件匹配到机制进行拦截进行了 c.Abort() 或 AbortWithSXxxx() 就直接结束了

	router.GET("", BeforeHomeMiddleWareOne, BeforeHomeMiddleWareTwo, HomeView)

	router.Run(":8000")
}

func HomeView(c *gin.Context) {
	fmt.Println("Home View")
}

func BeforeHomeMiddleWareOne(c *gin.Context) {
	// 请求部份
	fmt.Println("BeforeHomeMiddleWareOne request...")
	c.Next()
	// 响应部份
	fmt.Println("BeforeHomeMiddleWareOne response...")
}

func BeforeHomeMiddleWareTwo(c *gin.Context) {
	// 请求部份
	fmt.Println("BeforeHomeMiddleWareTwo request...")
	c.Next()
	// 响应部份
	fmt.Println("BeforeHomeMiddleWareTwo response...")
}
