package validatorutil

import (
	"reflect"
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// 注册自定义校验器统一入口
func registerCustomValidators(v *validator.Validate, trans ut.Translator) {
	// 注册 is-even 校验规则：整数且必须为偶数
	// 第一个参数是校验名字
	// 第二个参数是真正的校验规则，返回 bool 表示是否校验通过
	_ = v.RegisterValidation("is-even", func(fl validator.FieldLevel) bool {
		// fl.Field() 指向当前被娇艳的字段
		switch fl.Field().Kind() {
		// 只简单做 intXX 系列
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fl.Field().Int()%2 == 0
		// default 分支用于如果作用在不是 intXX 类型的字段上直接返回错误
		default:
			return false
		}
	})

	// is-even 注册中文翻译
	// 第一个参数必须和校验规则名一致
	// 第二个参数为注册器
	// 第三个参数函数是注册模板
	// 第四个参数函数真正生成错误信息
	_ = v.RegisterTranslation("is-even", trans, func(ut ut.Translator) error {
		// 第一个参数为翻译 key
		// 第二个参数是翻译模板，{0} 会被替换为字段名字来自于 v.RegisterTagNameFunc(...)
		// 第三个参数为是否允许覆盖已有模板
		return ut.Add("is-even", "{0}必须为偶数", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		// fe.Field() 已经中文化处理
		// 找到 is-even 对应模板替换 {0} 并返回最终字符串
		t, _ := ut.T("is-even", fe.Field())
		return t
	})

	// 注册 cn-mobile 手机号校验
	reMobile := regexp.MustCompile(`^1[3-9]\d{9}$`)
	_ = v.RegisterValidation("cn-mobile", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return false
		}
		return reMobile.MatchString(fl.Field().String())
	})

	_ = v.RegisterTranslation("cn-mobile", trans, func(ut ut.Translator) error {
		return ut.Add("cn-mobile", "{0}不符合中国手机号要求", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("cn-mobile", fe.Field())
		return t
	})

}
