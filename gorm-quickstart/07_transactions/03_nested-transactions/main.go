package main

import (
	"context"
	"errors"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

/*
嵌套事务，您可以回滚较大事务内执行的一部分操作
*/

type User struct {
	ID   uint   `gorm:"primaryKey;"`
	Name string `gorm:"size:64;uniqueIndex"`
	Age  int
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{})
}

func main() {
	// JimLee 和 JamesBond 会被创建
	ctx := context.Background()
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 第一层事务
		err := gorm.G[User](tx).Create(ctx, &User{
			Name: "JimLee",
			Age:  20,
		})
		if err != nil {
			return err
		}

		// 再开一层事务（外部事务不接受错误）
		_ = tx.Transaction(func(tx2 *gorm.DB) error {
			_ = gorm.G[User](tx2).Create(ctx, &User{
				Name: "BruceLee",
				Age:  10,
			})
			// 手动触发错误，这里的用户会被回滚创建失败
			return errors.New("rollback BruceLee")
		})

		// 再开一层事务（外部事务不接受错误）
		_ = tx.Transaction(func(tx3 *gorm.DB) error {
			_ = gorm.G[User](tx3).Create(ctx, &User{
				Name: "JamesBond",
				Age:  50,
			})
			return nil
		})

		return nil

	}); err != nil {
		log.Fatal(err)
	}
}
