package global

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect 统一连接数据库
func Connect() {
	dsn := "root:Dayi@516@tcp(192.168.7.236:53306)/test?parseTime=True"
	// 第二个参数是 gorm 的配置
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info), // 日志中开启 SQL 语句详情信息
	})
	if err != nil {
		panic(err)
	}
	DB = db
}
