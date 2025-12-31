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
使用 Save 进行更新，会保存所有字段，即使字段是零值（upsert = update + insert）
注意事项：Save 会先尝试全字段 UPDATE，如果发现更新不到任何行，会再自动执行一次 INSERT
1. 如果没有主键值，则直接走新增操作，行为等价于 Create
2. 如果有一个主键，则先执行全字段 UPDATE(所有字段都会被更新，零值也都会被写入数据库)
3. 如果传入了主键但是主键不存在，则会重新走插入操作，也就是 Create
*/
func main() {
	user, err := gorm.G[User](global.DB).First(context.Background())
	if err != nil {
		fmt.Println(err)
	} else {
		user.Name = "Fuck"
		user.Age = 666
		// 只有 Traditional API 有 Save 方法，因为挺危险的
		global.DB.Save(&user)
		// UPDATE users SET name = "Fuck", age - 20 WHERE id=1;
	}
	// 因为 Birthday 设置了 autoCreateTime 所以会进行设置的
	global.DB.Save(&User{
		Name: "PunkJerry",
		Age:  100,
	})
	// INSERT INTO `users` (`name`,`age`,`birthday`) VALUES ("jinzhu",100,"0000-00-00 00:00:00")

	// global.DB.Save(&User{ID: 1, Name: "jiangshan", Age: 100})
	// 如果不指定 birthday 等价于 UPDATE `users` SET `name`="jinzhu",`age`=100,`birthday`="0000-00-00 00:00:00" WHERE `id` = 1
	// 但是这样零值会报错

	global.DB.Save(&User{ID: 1, Name: "jiangshan", Age: 100, Birthday: time.Now()})

}
