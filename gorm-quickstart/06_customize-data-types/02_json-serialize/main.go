package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

/*
如果只是存入 JSON 数据，可以直接使用 gorm 的序列化，自带 json（也可以自定义序列化器）

并且有许多第三方包实现了 Scanner/Valuer 接口，可与 GORM 一起使用
*/

type User struct {
	// 这两个只能用于 Postgres 数据库，使用 "github.com/lib/pq" 实现
	//ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	//Tags.StringArray `gorm:"type:text[]"`
	ID   int
	Name string
	Info Info `gorm:"type:longtext;serializer:json" json:"info"`
}

type Info struct {
	Tags []string
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(User{}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	global.DB.Create(&User{
		Name: "JimLee",
		Info: Info{
			Tags: []string{
				"Go", "Python", "Docker",
			},
		},
	})

	user, err := gorm.G[User](global.DB).Where("name = ?", "JimLee").First(context.Background())
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(user.Name, user.Info.Tags)
	}
}
