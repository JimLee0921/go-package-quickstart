package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Age      int    `gorm:"default:18"`
	Birthday time.Time
}

func main() {
	/*
		传入切片插入多条数据
		泛型不支持传入切片直接批量插入，但是可以使用 gorm gen
	*/
	global.Connect()
	_ = global.DB.AutoMigrate(&User{})
	users := []User{
		{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
		{Name: "Jackson", Age: 19, Birthday: time.Now()},
	}
	// 配合 context
	// result := global.DB.WithContext(context.Background()).Create(users)
	// 使用 CreateInBathes 指定批量插入的批次大小，也可以在创建 db 是直接进行设置 batchSize
	result := global.DB.WithContext(context.Background()).CreateInBatches(users, 100)
	fmt.Println(result.Error)
	fmt.Println(result.RowsAffected)
}
