package main

import (
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

/*
GORM 中对于单条写操作（CREATE/UPDATE/DELETE）默认会在内部自动开启一个短事务来确保 ACID 特性

如果没必要开启，可以在初始化时通过 &gorm.Config{ SkipDefaultTransaction: true, } 进行关闭

也就用通过 Session 进行局部关闭，不为单条写操作自动开启事务

关闭事务会有一定的性能提升，但是在一定场景下比较危险
*/

type User struct {
	gorm.Model
	Name string
	Age  int
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(User{}); err != nil {
		log.Fatal(err)
	}
}
func main() {
	// 持续会话模式
	tx := global.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	tx.Create(&User{
		Name: "JimLee",
		Age:  20,
	})
	var user User
	var users []User
	tx.First(&user, 1)
	tx.Find(&users)

	tx.Model(&user).Update("Age", 22)
}
