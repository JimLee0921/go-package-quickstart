package main

import (
	"gin-quickstart/18_binding-custom-validator/validatorutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

//// InitValidator 初始化全局校验器设置语言为
//func InitValidator() {
//	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//		// 创建中文翻译器
//		zhT := zh.New()
//		uni := ut.New(zhT)
//		trans, _ := uni.GetTranslator("zh")
//
//		// 注册中文翻译
//		_ = zhTranslations.RegisterDefaultTranslations(v, trans)
//		// 这里默认都是 json 传递
//		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
//			// 把 label 和 json tag 都拿到，最后返回错误信息进行拼接更规则：
//			/*
//				{
//					"errors" : {
//						"name": "性名最长为xxx",
//						"age": "age年龄必须大于或等于18"
//					}
//				}
//			*/
//			// 字段名优先使用 label
//			label := fld.Tag.Get("label")
//			// 其次使用 json tage
//			name := fld.Tag.Get("json")
//			if name == "-" {
//				name = ""
//			}
//			if label == "" {
//				label = name
//			}
//			return name + ":" + label
//		})
//
//		// 自定义校验
//		registerCustomValidators(v, trans)
//		//保存到全局
//		translator = trans
//	}
//}

// 判断是否为校验错误并将错误信息等包装拼接成字符串

func main() {
	// main 中做一次自定义 validator 初始化
	validatorutil.Init()

	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		type User struct {
			// label 用于显示信息时把字段也显示为中文
			UserName  string `json:"user_name" binding:"required" label:"用户名"`
			Password  string `json:"password" binding:"required" label:"密码"`
			Sex       string `json:"sex" binding:"required,oneof=man woman" label:"性别"`
			Age       int    `json:"age" binding:"required,is-even" label:"年龄"` // 传了才校验，可以不传
			Telephone string `json:"telephone" binding:"cn-mobile" label:"手机号"`
		}
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validatorutil.ValidateError(err, &user),
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
