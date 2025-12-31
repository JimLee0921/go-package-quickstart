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
使用 db.Find() 获取全部对象
*/
func main() {
	var users []User
	result := global.DB.Find(&users)

	fmt.Println(result.RowsAffected)
	fmt.Println(result.Error)

	for _, user := range users {
		fmt.Println(user)
	}
}
