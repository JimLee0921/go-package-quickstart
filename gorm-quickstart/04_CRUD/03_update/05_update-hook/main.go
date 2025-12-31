package main

import (
	"errors"
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

// BeforeUpdate 使用 hook 控制更新
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Name == "JimLee" {
		return errors.New("JimLee cant update")
	}
	return
}

func main() {
	var user User
	global.DB.First(&user, 4)
	fmt.Println(user) // name 为 JimLee
	res := global.DB.Model(&user).Updates(map[string]any{
		"name": "Bruce",
	})
	fmt.Println(res.Error, res.RowsAffected)

}
