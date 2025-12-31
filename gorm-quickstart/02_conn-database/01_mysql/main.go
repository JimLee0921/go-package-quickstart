package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
Mysql 连接，
其他数据库连接和相关配置见
https://gorm.io/zh_CN/docs/connecting_to_the_database.html
*/
func main() {
	dsn := "root:Dayi@516@tcp(192.168.7.236:53306)/test"
	// 第二个参数是 gorm 的配置
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db)
}
