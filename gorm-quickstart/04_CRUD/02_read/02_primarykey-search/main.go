package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"time"

	"gorm.io/gorm"
)

/*
如果主键是数字类型，可以使用内联条件来检索对象
当使用字符串时，需要额外的注意来避免SQL注入
*/

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
func main() {
	GenericsDemo()
	TraditionalDemo()
}

func GenericsDemo() {
	ctx := context.Background()

	// 数值类型的主键，查找主键为 10 的元素
	user, err := gorm.G[User](global.DB).Where("id = ?", 10).First(ctx)
	// SELECT * FROM users WHERE id = 10;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}

	// 字符串类型的主键（数值也能用 "10"）
	user, err = gorm.G[User](global.DB).Where("id = ?", "10").First(ctx)
	// SELECT * FROM users WHERE id = 10;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}

	// 查找多个主键
	users, err := gorm.G[User](global.DB).Where("id IN ?", []int{6, 7, 9}).Find(ctx)
	// SELECT * FROM users WHERE id IN (1,2,3);
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// UUID 查找示例
	user, err = gorm.G[User](global.DB).Where("id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a").First(ctx)
	// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}
}

func TraditionalDemo() {
	/*
		Traditional API
	*/
	//var user User
	//global.DB.First(&user, 10)
	//// SELECT * FROM users WHERE id = 10;
	//fmt.Println(user)

	//var user User
	//global.DB.First(&user, "10")
	//// SELECT * FROM users WHERE id = 10;
	//fmt.Println(user)

	//var users []User
	//global.DB.Find(&users, []int{1, 2, 3})
	//// SELECT * FROM users WHERE id IN (1,2,3);
	//fmt.Println(users)

	var user []User
	global.DB.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")
	// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";

}
