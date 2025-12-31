package main

import (
	"fmt"
	"gorm-quickstart/global"
)

type User struct {
	Name string
	Age  int
}

/*
还可以使用 map[string]any 和 []map[string]any{} 进行创建记录
注意当使用 map 来创建时，钩子方法不会执行，关联不会被保存且不会回写主键
*/
func init() {
	global.Connect()
	err := global.DB.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
func main() {
	// 单条记录插入
	global.DB.Model(&User{}).Create(map[string]any{
		"Name": "JimLee",
		"Age":  20,
	})

	// 批量插入
	global.DB.Model(&User{}).Create([]map[string]any{
		{"Name": "jinzhu_1", "Age": 18},
		{"Name": "jinzhu_2", "Age": 20},
	})
}
