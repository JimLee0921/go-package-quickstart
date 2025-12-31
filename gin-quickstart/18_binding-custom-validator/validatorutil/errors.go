package validatorutil

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateError(err error, req any) (res map[string]string) {
	res = make(map[string]string)

	var errs validator.ValidationErrors
	// 非校验错误，直接返回错误信息字符串
	if !errors.As(err, &errs) {
		res["error"] = err.Error()
		return
	}

	// 反射获取结构体拼接 label 和 json tag 最后返回错误信息进行拼接更规则：
	/*
		{
			"errors" : {
				"name": "性名最长为xxx",
				"age": "age年龄必须大于或等于18"
			}
		}
	*/

	// 1. 拿到结构体类型
	rt := reflect.TypeOf(req)

	// 2. 判断是否为指针类型
	// 因为这里传入的 &User，而 TypeOf(&User).Kind() 是 Ptr
	// 但是 tag 属性定义在 User 上，而不是 *User 上
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	// 3. 处理校验失败的情况
	for _, fe := range errs {
		// 1. 获取字段名（StructField 返回的是 Go 结构体字段名）
		// 比如 UserName string `json:"user_name"` 返回的就是 UserName
		jsonKey := fe.StructField()
		// 2. 兜底判断，正常传入这里都是 struct 结构体类型
		if rt.Kind() == reflect.Struct {
			// 3. 通过字段名匹配对应字段
			if f, ok := rt.FieldByName(fe.StructField()); ok {
				// 4. 获取 json 标签的值
				if tag := f.Tag.Get("json"); tag != "" && tag != "-" {
					// 5. json 标签的值可能为 `json:"user_name,omitempty"` 需要根据 , 裁剪获取真正的 tag name
					jsonKey = strings.Split(tag, ",")[0]
				}
			}
		}
		// 2. 进行翻译
		msg := fe.Translate(translator)

		// 3. 可选：去除消息中的字段名，只保留规则描述
		// fe.Field() 经过 RegisterTagNameFunc，可能是中文 label
		//fieldName := fe.Field()
		//msg = strings.TrimSpace(strings.TrimPrefix(msg, fieldName))

		res[jsonKey] = msg
	}
	return
}
