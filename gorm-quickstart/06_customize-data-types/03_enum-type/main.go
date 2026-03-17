package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/*
枚举类型在返回前端时通常使用 JSON 序列化进行返回，但是因为在数据库中这种枚举存入的是0，1，2这些数字

如果想要 JSON 返回时自动转为字符串需要枚举类型实现 MarshalJSON
*/

type Level int8

func (l Level) MarshalJSON() (data []byte, err error) {
	var str string
	switch l {
	case InfoLevel:
		str = "info"
	case WarningLevel:
		str = "warning"
	case ErrorLevel:
		str = "error"
	}
	return json.Marshal(map[string]any{
		// 把对应的数字和字符串都返回给前端(最简单的写法)
		"value": int8(l),
		"label": str,
	})
}

const (
	InfoLevel Level = iota
	WarningLevel
	ErrorLevel
)

type MyLog struct {
	gorm.Model
	LogID datatypes.UUID
	Title string `gorm:"size:256"`
	Level Level
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(MyLog{}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	logs := []MyLog{
		{
			LogID: datatypes.NewUUIDv4(),
			Title: "system start",
			Level: InfoLevel,
		},
		{
			LogID: datatypes.NewUUIDv4(),
			Title: "disk space warning",
			Level: WarningLevel,
		},
		{
			LogID: datatypes.NewUUIDv4(),
			Title: "db connection failed",
			Level: ErrorLevel,
		},
	}

	if err := global.DB.Create(&logs).Error; err != nil {
		log.Fatal(err)
	}

	logsFind, err := gorm.G[MyLog](global.DB).Find(context.Background())
	if err != nil {
		log.Fatal(err)
	} else {
		// 模拟 JSON 序列化返回给前端
		data, err := json.Marshal(logsFind)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
	}
}
