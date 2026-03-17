package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := "root:Dayi@516@tcp(192.168.7.236:53306)/test?parseTime=True"
	// 第二个参数是 gorm 的配置
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info), // 控制 SQL 日志级别、慢 SQL、输出格式
		SkipDefaultTransaction: true,                                // 是否跳过关闭默认事务
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:         "t_", // 统一前缀
			SingularTable:       true, // 使用单数表名
			NameReplacer:        nil,
			NoLowerCase:         false, // 不推荐修改，把 UserName -> User_Name
			IdentifierMaxLength: 0,
		}, // 自定义命名策略，需要实现 schema.Namer 接口
		PrepareStmt:                              true, // 开启 Prepared Statement 缓存，性能优化
		DisableForeignKeyConstraintWhenMigrating: true, // AutoMigrate 时不自动创建外键
		NowFunc: func() time.Time {
			return time.Now().UTC() // 自动时间字段统一来源
		},
		QueryFields:          true, // 字段精确查询 * 替换为 id, name, age, ...
		CreateBatchSize:      1000, // 批量插入，db.Create() 时自动分批
		DryRun:               true, // 只生成 SQL 不执行
		DisableAutomaticPing: true, // 关闭启动时的 Ping 操作，特殊链接场景
	})
	if err != nil {
		panic(err)
	}
}
