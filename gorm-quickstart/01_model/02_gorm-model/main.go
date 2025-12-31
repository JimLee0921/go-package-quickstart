package main

import (
	"gorm.io/gorm"
)

/*
GORM 提供了预定义的结构体 gorm.Model 包含四个常用字段：
ID			 每个记录的唯一标识符（主键）
CreatedAt	 在创建记录时自动设置为当前时间
UpdatedAt	 每当记录更新时，自动更新为当前时间
DeletedAt	 用于软删除（将记录标记为已删除，而实际上并未从数据库中删除）
可以直接拿来用嵌入自己的结构体中
*/

type User struct {
	// 自动包含四个字段
	gorm.Model
	Name string `gorm:"not null;unique"`
	Age  int    `gorm:"default:20"`
}

func main() {

}
