package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

/*
GORM 大多数在存数据的时候可能有需要存入 json 数据的需求，这时候就需要用到自定义数据类型

自定义的数据类型必须实现两个方法：driver.Valuer 写入数据库和 sql.Scanner 从数据库读取
*/

// StringSlice 自定义数据类型
type StringSlice []string

// Value 实现 driver.Valuer 方法，返回 json value，这里必须是值接收者
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 方法，将 value 扫描为 json，这里必须是指针接收者
func (s *StringSlice) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}
	// 这里断言 []byte 是因为 Value 使用的 Marshal 返回 []byte
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringSlice")
	}
	return json.Unmarshal(bytes, s)
}

type User struct {
	ID   uint
	Name string `gorm:"size:32"`
	// Tags StringSlice // 自定义数据类型，这里在数据库中如何存储取决于数据库
	// Tags StringSlice `gorm:"type:json"` // 数据库中使用 json 存储
	Tags StringSlice `gorm:"type:text" json:"info"` // 数据库中使用 text 存储
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(User{}); err != nil {
		log.Fatal(err)
	}
}
func main() {
	//global.DB.Create(&User{
	//	Name: "JimLee",
	//	Tags: []string{
	//		"Golang",
	//		"Python",
	//		"Docker",
	//	},
	//})

	user, err := gorm.G[User](global.DB).Where("name = ?", "JimLee").First(context.Background())
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(user.Name, user.Tags)
	}

}
