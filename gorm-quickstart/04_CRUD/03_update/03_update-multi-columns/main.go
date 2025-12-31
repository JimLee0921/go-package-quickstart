package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"time"

	"gorm.io/gorm"
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
Updates 方法支持 struct 和 map[string]interface{} 参数
当使用 struct 更新时，默认情况下 GORM 只会更新非零值的字段
*/
func main() {
	// GenericsDemo()
	TraditionalDemo()
}

func GenericsDemo() {
	ctx := context.Background()

	// 使用结构体批量更新字段，只更新非零值
	rowsAffected, err := gorm.G[User](global.DB).Where("id = ?", 1).Updates(ctx, User{
		Name: "Crazy",
		Age:  222,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rowsAffected)
	}
}

func TraditionalDemo() {
	var user User
	global.DB.First(&user)
	// 调用 struct 更新，只会更新非零的字段
	err := global.DB.Model(&user).Updates(User{
		Name: "CCC",
		Age:  23,
	}).Error
	fmt.Println(err)

	// 根据 map 更新属性
	err = global.DB.Model(&user).Updates(map[string]any{
		"name": "Jim",
		"age":  222,
	}).Error
	fmt.Println(err)
}
