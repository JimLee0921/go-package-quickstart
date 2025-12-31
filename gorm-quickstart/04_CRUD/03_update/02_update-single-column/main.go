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
当使用 Update 更新单列时，需要有一些条件
否则将会引起 ErrMissingWhereClause 错误
*/
func main() {
	// TraditionalDemo()
	GenericsDemo()
}

func TraditionalDemo() {

	// 条件更新
	err := global.DB.Model(&User{}).Where("name = ?", "Dsb").Update("birthday", time.Now()).Error
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("update successful")
	}

	// 指定 ID
	user1 := User{ID: 111}
	err = global.DB.Model(&user1).Update("name", "Hello").Error
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("update successful")
	}

	// 根据条件和 model 的值进行更新
	user2 := User{ID: 111}
	err = global.DB.Model(&user2).Where("name = ?", "Hello").Update("name", "HHHH").Error
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("update successful")
	}
}

func GenericsDemo() {
	ctx := context.Background()

	// 条件更新
	rowsAffected, err := gorm.G[User](global.DB).Where("name = ?", "Dsb").Update(ctx, "name", "BruceLee")
	fmt.Println(rowsAffected, err)

	// 传入 ID
	rowsAffected, err = gorm.G[User](global.DB).Where("id = ?", 90).Update(ctx, "name", "World")
	fmt.Println(rowsAffected, err)

	// 多个条件
	rowsAffected, err = gorm.G[User](global.DB).Where("name = ? AND age = ?", "BruceLee", 18).Update(ctx, "age", 22)
	fmt.Println(rowsAffected, err)
}
