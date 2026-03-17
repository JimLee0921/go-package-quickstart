package main

import "gorm-quickstart/global"

/*
many to many 会在两个 model 中添加一张连接表
比如 user 和 language 表，一个 user 可以有多个 language
多个 user 可以属于一个 language

使用 GORM 的 AutoMigrate 为 User 创建表时，GORM 会自动创建连接表

这里演示 User 单向（Unidirectional relationship） many2many（只能从 User 查到 Language）
*/

type User struct {
	ID        int
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	ID   int
	Name string
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, Language{})
}

func main() {

}
