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
如果您想要在更新时选择、忽略某些字段，可以使用 Select、Omit
*/
func main() {
	GenericsDemo()
	// TraditionalDemo()
}

func GenericsDemo() {
	ctx := context.Background()

	// 使用结构体批量更新字段，只更新非零值
	rowsAffected, err := gorm.G[User](global.DB).Where("id = ?", 4).Select("name").Updates(ctx, User{
		Name: "Cra2zy",
		Age:  123,
	})
	// UPDATE users SET name = "Crazy" WHERE id = 4;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rowsAffected)
	}

	rowsAffected, err = gorm.G[User](global.DB).Where("id = ?", 24).Omit("Age").Updates(ctx, User{
		Name: "Fuck",
		Age:  1,
	})
	// UPDATE users SET name = "Fuck" WHERE id = 24;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rowsAffected)
	}

	// 更新零值
	rowsAffected, err = gorm.G[User](global.DB).Where("id = ?", 25).Select("Name", "Age").Updates(ctx, User{
		Name: "BruceTT",
		Age:  0,
	})
	// UPDATE users SET name = "BruceTT" AND age = 0 WHERE id = 25;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rowsAffected)
	}

	// 选择所有字段没有传入设置零值（这里需要忽略主键 ID 和一些不能为 0 的 值，不然也会自动设置 id = 0）
	rowsAffected, err = gorm.G[User](global.DB).Where("id = ?", 26).Select("*").Omit("id", "birthday").Updates(ctx, User{
		Name: "CrazyLeeJim",
	})
	// UPDATE users SET name = "CrazyLeeJim" AND age = 0 WHERE id = 26;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rowsAffected)
	}
}

func TraditionalDemo() {
	// 调用 struct 更新，只会更新非零的字段
	err := global.DB.Model(&User{ID: 1}).Updates(User{
		Name: "CCC",
		Age:  23,
	}).Error
	// UPDATE users SET name = "CCC" ABD age = 23 WHERE id = 1;
	fmt.Println(err)

	// 根据 map 更新属性
	err = global.DB.Model(&User{ID: 8}).Updates(map[string]any{
		"name": "Jim",
		"age":  222,
	}).Error
	// UPDATE users SET name = "hello" AND age = 222 WHERE id = 8;
	fmt.Println(err)

	// 选中零值字段
	err = global.DB.Model(&User{ID: 2}).Select("name", "age").Updates(User{
		Name: "TTT",
		Age:  0,
	}).Error
	// UPDATE users SET name = "TTT" AND age = 0 WHERE id = 2;
	fmt.Println(err)

	// 不更新指定字段
	err = global.DB.Model(User{ID: 32}).Omit("name").Updates(map[string]any{
		"name": "hello",
		"age":  55,
	}).Error
	// UPDATE users SET age = 55 WHERE id = 32;
	fmt.Println(err)

	// 选择所有字段(包括零值)
	err = global.DB.Model(User{ID: 15}).Select("name", "Age").Updates(User{
		Name: "NewName",
		Age:  0,
	}).Error
	fmt.Println(err)

	// 选择除指定字段外所有字段（这里需要忽略主键 ID 和一些不能为 0 的 值，不然也会自动设置 id = 0）
	err = global.DB.Model(&User{ID: 1}).Select("*").Omit("id", "birthday").Updates(User{
		Name: "Jinx",
		Age:  123,
	}).Error
	fmt.Println(err)
}
