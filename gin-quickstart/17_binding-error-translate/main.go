package main

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var translator ut.Translator

// InitValidator 初始化全局校验器设置语言为
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 创建中文翻译器
		zhT := zh.New()
		uni := ut.New(zhT)
		trans, _ := uni.GetTranslator("zh")

		// 注册中文翻译
		_ = zhTranslations.RegisterDefaultTranslations(v, trans)
		// 这里默认都是 json 传递
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			// 把 label 和 json tag 都拿到，最后返回错误信息进行拼接更规则：
			/*
				{
					"errors" : {
						"name": "性名最长为xxx",
						"age": "age年龄必须大于或等于18"
					}
				}
			*/
			// 字段名优先使用 label
			label := fld.Tag.Get("label")
			// 其次使用 json tage
			name := fld.Tag.Get("json")
			if name == "-" {
				name = ""
			}
			if label == "" {
				label = name
			}
			return name + ":" + label
		})

		//保存到全局
		translator = trans
	}
}

// 判断是否为校验错误并将错误信息等包装拼接成字符串

func ValidateError(err error) any {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	// 非校验错误，直接返回错误信息字符串
	if !ok {
		return err.Error()
	}

	// 校验错误，翻译对错误列表进行遍历拼接返回
	errMessages := make(map[string]string)
	for _, msg := range errs {
		errMessage := msg.Translate(translator)
		errSlice := strings.Split(errMessage, ":")
		errMessages[errSlice[0]] = errSlice[1]
	}
	return errMessages

}

func main() {
	// main 中初始化一次
	InitValidator()

	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		type User struct {
			// label 用于显示信息时把字段也显示为中文
			UserName string `json:"user_name" binding:"required" label:"用户名"`
			Password string `json:"password" binding:"required" label:"密码"`
			Sex      string `json:"sex" binding:"required,oneof=man woman" label:"性别"`
			Age      int    `json:"age" binding:"omitempty,gte=18,lte=65" label:"年龄"` // 传了才校验，可以不传
		}
		var user User
		//if err := c.ShouldBindJSON(&user); err != nil {
		//	// 字段校验错误情况
		//	var errs validator.ValidationErrors
		//	if errors.As(err, &errs) {
		//		c.JSON(http.StatusBadRequest, gin.H{
		//			// 返回中文错误信息
		//			"error": errs.Translate(translator),
		//		})
		//		return
		//	}
		//	// 非字段校验错误情况
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ValidateError(err),
			})
			return

		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "successful",
			"user": user,
		})

	})

	router.Run(":8000")
}
