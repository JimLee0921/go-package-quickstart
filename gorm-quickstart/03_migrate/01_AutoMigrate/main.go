package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
定义模型
*/

type User struct {
	gorm.Model
	Name string `gorm:"not null;unique"`
	Age  int    `gorm:"default:18"`
}

type Product struct {
	gorm.Model
	ProductName string `gorm:"not null;unique"`
}

type Order struct {
	gorm.Model
	OrderPrice float64
}

/*
使用 db.AutoMigrate 自动迁移 schema

注意 AutoMigrate 会创建表、缺失的外键、约束、列和索引但是多次调用不会删除字段
AutoMigrate 可以同时传入多个 schema 并且可以配合 db.Set() 设置参数

AutoMigrate 会自动创建数据库外键约束，可以在初始化时禁用
&gorm.Config{
  DisableForeignKeyConstraintWhenMigrating: true,
}
*/

func main() {
	global.Connect()

	// 使用 AutoMigrate 迁移 User 表
	//err := global.DB.AutoMigrate(&User{})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 使用 AutoMigrate 一次迁移多个表
	//err := global.DB.AutoMigrate(&User{}, &Product{}, &Order{})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 生成表命令时添加后缀
	err := global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Product{}, &Order{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
