package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
如果是需要多个写操作放在一个事务还是需要手动使用 db.Transaction(...) 进行开启自动事务
*/

type Account struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"size:255;uniqueIndex"`
	Balance int64
}

type TransferLog struct {
	ID     uint `gorm:"primaryKey"`
	FromID uint
	ToID   uint
	Amount int64
	Status string // SUCCESS / FAILED
	Reason string // 失败原因
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(Account{}, TransferLog{}); err != nil {
		log.Fatal(err)
	}

	// 初始化账号
	//if err := global.DB.Create([]Account{
	//	{Name: "JimLee", Balance: 1000},
	//	{Name: "JamesBond", Balance: 5000},
	//}).Error; err != nil {
	//	log.Fatal(err)
	//}
}

func main() {
	//_ = TransferTraditional("JimLee", "JamesBond", 10000)
	//_ = TransferTraditional("JamesBond", "JimLee", 200)

	_ = TransferGenerics("JimLee", "JamesBond", 1000)
	_ = TransferGenerics("JamesBond", "JimLee", 200)
}

// TransferTraditional
// 进行转帐，如果发起转帐方余额不足会进行回滚
func TransferTraditional(fromName, toName string, amount int64) error {
	// 1. 查询双方账户 ID （也可以放在事务中进行）
	var from, to Account
	if err := global.DB.Where("name = ?", fromName).First(&from).Error; err != nil {
		log.Fatal(err)
	}
	if err := global.DB.Where("name = ?", toName).First(&to).Error; err != nil {
		log.Fatal(err)
	}
	// 2. 主事务（扣款、加款、记录日志）
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 锁定付款方，防止并发超扣（Clauses(clause.Locking{Strength: "UPDATE"})意味着在事务提交或回滚前其它事务不能修改它）
		// 收款方按理也需要锁定，但是这里更新的时候使用的 gorm.Expr 已经保证了收款方 balance 的原子性，这里可以不锁定
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&from, from.ID).Error; err != nil {
			return err
		}

		// 2. 判断是否可以扣款成功
		if from.Balance < amount {
			return fmt.Errorf("insufficient balance: %d < %d", from.Balance, amount)
		}

		// 3. 付款方扣款
		if err := tx.Model(&Account{}).Where("id = ?", from.ID).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// 4. 收款方加钱
		if err := tx.Model(&Account{}).Where("id = ?", to.ID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// 5. 成功日志（属于主事务的一部分）
		return tx.Create(&TransferLog{
			FromID: from.ID,
			ToID:   to.ID,
			Amount: amount,
			Status: "SUCCESS",
		}).Error
	})

	// 3. 失败日志，必须要在主事务外进行写入，否则会随着回滚消失
	if err != nil {
		_ = global.DB.Create(&TransferLog{
			FromID: from.ID,
			ToID:   to.ID,
			Amount: amount,
			Status: "FAILED",
			Reason: err.Error(),
		})
	}
	return nil
}

// TransferGenerics 泛型写法
func TransferGenerics(fromName, toName string, amount int64) error {
	// 1. 查询双方账户 ID （也可以放在事务中进行）
	ctx := context.Background()
	from, err := gorm.G[Account](global.DB).Where("name = ?", fromName).First(ctx)
	if err != nil {
		log.Fatal(err)
	}
	to, err := gorm.G[Account](global.DB).Where("name = ?", toName).First(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 开启主事务
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 给付款方加锁防止多并发场景发生错误时
		if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&from, from.ID).Error; err != nil {
			return err
		}

		// 2. 判断是否可以扣款成功
		if from.Balance < amount {
			return fmt.Errorf("insufficient balance: %d < %d", from.Balance, amount)
		}

		// 3. 付款方扣款
		if _, err = gorm.G[Account](tx).Where("id = ?", from.ID).Update(ctx, "balance", gorm.Expr("balance - ?", amount)); err != nil {
			return err
		}

		// 4. 收款方收款
		if _, err = gorm.G[Account](tx).Where("id = ?", to.ID).Update(ctx, "balance", gorm.Expr("balance + ?", amount)); err != nil {
			return err
		}

		// 5. 成功日志（属于主事务的一部分）
		return gorm.G[TransferLog](tx).Create(ctx, &TransferLog{
			FromID: from.ID,
			ToID:   to.ID,
			Amount: amount,
			Status: "SUCCESS",
		})

	})

	// 3. 失败日志
	if err != nil {
		return gorm.G[TransferLog](global.DB).Create(ctx, &TransferLog{
			FromID: from.ID,
			ToID:   to.ID,
			Amount: amount,
			Status: "FAILED",
			Reason: err.Error(),
		})
	}
	return nil
}
