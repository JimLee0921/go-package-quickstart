package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 真正发送 JSON response 包内使用
func response(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(http.StatusOK, &Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// 定义状态码对应的值
var codeMap = map[int]string{
	1001: "权限错误",
	1002: "角色错误",
}

// ============ 外部使用函数 ============

// Ok
// 默认成功响应，code 为 0
func Ok(ctx *gin.Context, message string, data any) {
	response(ctx, 0, message, data)
}

// OkWithData
// message 为默认值 successful 需要传入 data 数据
func OkWithData(ctx *gin.Context, data any) {
	Ok(ctx, "successful", data)
}

// OkWithMsg 不返回 data 只有 message
// 这里没有数据 data 尽量不要传入 nil
// 使用 map[string]any{} 或 gin.H{}
func OkWithMsg(ctx *gin.Context, message string) {
	Ok(ctx, message, map[string]any{})

}

func Fail(ctx *gin.Context, code int, message string, data any) {
	response(ctx, code, message, data)
}

func FailWithMessage(ctx *gin.Context, message string) {
	Fail(ctx, 1001, message, gin.H{})
}

func FailWithCode(ctx *gin.Context, code int) {
	msg, ok := codeMap[code]
	if !ok {
		msg = "server error"
	}
	Fail(ctx, code, msg, gin.H{})
}
