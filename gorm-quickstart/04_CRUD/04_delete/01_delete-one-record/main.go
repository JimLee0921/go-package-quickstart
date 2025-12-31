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
删除一条记录时，删除对象需要指定主键，否则会触发 批量删除
*/
func main() {
	// TraditionDemo()
	GenericsDemo()
}

func TraditionDemo() {
	res := global.DB.Delete(&User{ID: 20})
	// DELETE from users WHERE id = 20;
	// 等价于 global.DB.Delete(&User{ID: 20}) 和 global.DB.Delete(&User{ID: "20"})
	fmt.Println(res.Error, res.RowsAffected)

	// 带额外条件
	res = global.DB.Where("name = ? AND id = ?", "BruceLee", 9).Delete(&User{})
	fmt.Println(res.Error, res.RowsAffected)
}

func GenericsDemo() {
	ctx := context.Background()

	// 通过 ID 删除
	rowsAffected, err := gorm.G[User](global.DB).Where("id = ?", 10).Delete(ctx)
	fmt.Println(rowsAffected, err)

	rowsAffected, err = gorm.G[User](global.DB).Where("id = ? AND name = ?", 20, "JimLee").Delete(ctx)
	// DELETE from emails where id = 10 AND name = "jinzhu";
	fmt.Println(rowsAffected, err)
}
