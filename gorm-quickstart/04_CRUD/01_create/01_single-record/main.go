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

/*
创建并插入单条记录
*/
func main() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{})

	user := User{
		Name:     "JimLee",
		Age:      20,
		Birthday: time.Now(),
	}

	// 传统 API
	resultTraditional := global.DB.Create(&user)
	fmt.Println(user.ID)                        // 插入数据的主键
	fmt.Println(resultTraditional.Error)        // 是否有错误
	fmt.Println(resultTraditional.RowsAffected) // 返回插入记录的总数

	// 泛型 API
	resultGenerics := gorm.WithResult()
	// gorm.G[models.User](global.DB, resultGenerics) 是针对 User 类型的 ORM 操作器
	// 后面的 opts 参数切片传入了个 result 是一个操作结果收集器
	err := gorm.G[User](global.DB, resultGenerics).Create(context.Background(), &user)
	fmt.Println(user.ID)                     // 插入数据的主键
	fmt.Println(err)                         // 是否有错误
	fmt.Println(resultGenerics.RowsAffected) // 返回插入记录的总数
}
