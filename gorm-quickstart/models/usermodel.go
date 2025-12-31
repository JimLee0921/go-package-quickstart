package models

import (
	"time"

	"gorm.io/gorm"
)

// User 可以使用 Model 也就是 UserModel 结尾表示这个结构体是表结构
type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Age      int    `gorm:"default:18"`
	Birthday time.Time
}
