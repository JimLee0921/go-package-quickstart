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
func main() {
	// 通过 offset + limit 进行分页获取全部数据
	limit := 5 // limit 不能为 0
	page := 1
	for {
		offset := limit * (page - 1)
		find, err := gorm.G[User](global.DB.Debug()).Limit(limit).Offset(offset).Find(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("batch %d:\n", page)
		for index, item := range find {
			fmt.Println(index, item)
		}
		if len(find) != limit {
			fmt.Println("no more data")
			break
		} else {
			page += 1
		}
		fmt.Println()
	}
}
