package main

import "gorm-quickstart/global"

/*
自引用 many2many 关系
*/

type User struct {
	ID      int
	Friends []*User `gorm:"many2many:user_friends"`
}

/*
会创建连接表：user_friends

foreign key: user_id, reference: users.id
foreign key: friend_id, reference: users.id
*/

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{})
}
func main() {

}
