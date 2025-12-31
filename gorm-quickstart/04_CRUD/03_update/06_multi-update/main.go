package main

import (
	"fmt"
	"gorm-quickstart/global"
	"time"
)

type User struct {
	ID       int       `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Age      int       `gorm:"default:18"`
	Birthday time.Time `gorm:"autoCreateTime"`
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(&User{})
}

/*
如果没有通过 Model 指定一个含有主键的记录，GORM 会执行批量更新

全局更新：
如果执行一个没有任何条件的批量更新，GORM 默认不会运行
并且会返回 ErrMissingWhereClause 错误

可以通过配置启用 AllowGlobalUpdate 模式
*/
func main() {
	res := global.DB.Model(User{}).Where("name = ?", "JimLee").Updates(User{
		Name: "Bruce",
		Age:  24,
	})
	fmt.Println(res.Error, res.RowsAffected)

	res = global.DB.Model(User{}).Where("id in ?", []int{10, 16, 12, 55}).Updates(map[string]any{
		"name": "JimLee",
		"age":  18,
	})
	fmt.Println(res.Error, res.RowsAffected)
}
