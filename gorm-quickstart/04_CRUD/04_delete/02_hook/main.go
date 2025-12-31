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

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Name == "JimLee" {
		return errors.New("JimLee user cant delete")
	}
	return
}

/*
对于删除操作，GORM 支持 BeforeDelete、AfterDelete Hook
在删除记录时会调用这些方法
*/
func main() {
	// ID 为 18 的 name 为 JimLee
	res := global.DB.Delete(&User{ID: 18})
	fmt.Println(res.Error)
	fmt.Println(res.RowsAffected)
}
