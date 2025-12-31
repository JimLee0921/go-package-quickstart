package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"time"

	"gorm.io/gorm"
)

/*
GORM 提供了 First Take Last 方法以便从数据库中检索单个对象
当查询数据库时它添加了 LIMIT 1 条件，且没有找到记录时，它会返回 ErrRecordNotFound 错误

First Take 使用的主键排序查找第一条和最后一条
只有在目标 struct 是指针或者通过 db.Model() 指定 model 时，该方法才有效
此外，如果相关 model 没有定义主键，那么将按 model 的第一个字段进行排序

如果使用 gorm 的特定字段类型（例如 gorm.DeletedAt），它将运行不同的查询来检索对象
db.First(&user)
//  SELECT * FROM `users` WHERE `users`.`id` = '15' AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
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
	// 插入一些记录
	global.DB.Create([]*User{
		{Name: "JimLee", Age: 22},
		{Name: "Dsb"},
		{Name: "Bruce", Age: 10},
		{Name: "Bond", Age: 55},
		{Name: "Gogo", Age: 17},
		{Name: "SSSS", Age: 13},
		{Name: "wtf", Age: 8},
	})
}
func main() {
	GenericsDemo()
	TraditionalDemo()
}

func GenericsDemo() {
	/*
		Generics API
	*/
	ctx := context.Background()

	// 根据主键获取第一条记录
	user, err := gorm.G[User](global.DB).First(ctx)
	// 等价于 SELECT * FROM users ORDER BY id LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}

	// 获取第一条记录，没有特殊顺序
	user, err = gorm.G[User](global.DB).Take(ctx)
	// SELECT * FROM users LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}

	// 按照主键获取最后一条记录
	user, err = gorm.G[User](global.DB).Last(ctx)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}
}

func TraditionalDemo() {
	/*
		Traditional API 这几个分开运行，不然结果一样
	*/
	var user User
	//global.DB.First(&user)
	//fmt.Println(user)
	//global.DB.Take(&user)
	//fmt.Println(user)
	global.DB.Last(&user)
	fmt.Println(user)
}
