package main

import (
	"context"
	"errors"
	"fmt"
	"gorm-quickstart/global"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
GORM允许用户通过实现这些接口
BeforeSave, BeforeCreate, AfterSave, AfterCreate来自定义钩子
这些钩子方法会在创建一条记录时被调用，可以用于字段校验、日志记录等

也可以使用 SkipHooks 会话模式跳过 Hooks 方法

DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Create(&users)
DB.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)
*/

type Role uint

const (
	RoleAdmin Role = iota
	RoleUser
	RoleMaster
)

type User struct {
	UUID     uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name     string    `gorm:"not null"`
	Age      int       `gorm:"default:18"`
	Role     Role      `gorm:"type:tinyint;not null;default:2"`
	Birthday int64     `gorm:"autoCreateTime"`
}

// BeforeCreate 钩子，用于设置 UUID 和检查 Role
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	// 1. 设置 uuid
	u.UUID = uuid.New()

	// 2. 检查 Role
	if u.Role == RoleAdmin || u.Role == RoleMaster {
		return errors.New("invalid role")
	}
	return

}

func init() {
	global.Connect()
	err := global.DB.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
func main() {

	user := User{
		Name: "JimLee",
		Age:  202,
		Role: RoleUser,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := gorm.G[User](global.DB).Create(ctx, &user)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 这里跳过钩子如果不设置 UUID 的话 uuid 就会变成了 00000000-0000-0000-0000-000000000000
	result := global.DB.Session(&gorm.Session{SkipHooks: true}).Create(&User{
		UUID: uuid.New(),
		Name: "JimLee",
		Age:  202,
		Role: RoleAdmin,
	})
	fmt.Println(result.Error)
}
