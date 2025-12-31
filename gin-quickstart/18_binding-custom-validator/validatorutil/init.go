package validatorutil

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var translator ut.Translator

// Init 初始化全局校验器设置语言为
func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 创建中文翻译器
		zhT := zh.New()
		uni := ut.New(zhT)
		trans, _ := uni.GetTranslator("zh")

		//保存到全局
		translator = trans

		// 注册中文翻译
		_ = zhTranslations.RegisterDefaultTranslations(v, trans)
		// 这里仅用于翻译时显示字段名返回 label，label和json tag 拼接在 errors 中完成
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			if label := fld.Tag.Get("label"); label != "" {
				return label
			}
			// 没有 label 标签再获取 json tag 避免空字段
			name := fld.Tag.Get("json")
			if name == "-" {
				return ""
			}
			return name
		})

		// 自定义校验
		registerCustomValidators(v, trans)
	}
}
