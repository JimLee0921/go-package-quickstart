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
如果指定的值不包括主属性，那么 GORM 会执行批量删除，它将删除所有匹配的记录
当你试图执行不带任何条件的批量删除时，GORM将不会运行并返回 ErrMissingWhereClause 错误
可以通过配置开启 AllowGlobalUpdate
*/

func main() {
	// Generic API
	ctx := context.Background()

	rowsAffected, err := gorm.G[User](global.DB).Where("name LIKE ?", "%jinzhu%").Delete(ctx)
	fmt.Println(rowsAffected, err)

	// Traditional API
	result := global.DB.Where("name LIKE ?", "%jinzhu%").Delete(&User{})
	// 等价于 global.DB.Delete(&User{}, "email LIKE ?", "%jinzhu%")
	// DELETE FROM users where mame LIKE "%jinzhu%";
	fmt.Println(result.Error, result.RowsAffected)

	// 传递切片
	var users = []User{{ID: 1}, {ID: 2}, {ID: 3}}
	global.DB.Delete(&users)
	// DELETE FROM users WHERE id IN (1,2,3);

	global.DB.Delete(&users, "name LIKE ?", "%jinzhu%")
	// DELETE FROM users WHERE name LIKE "%jinzhu%" AND id IN (1,2,3);

}
