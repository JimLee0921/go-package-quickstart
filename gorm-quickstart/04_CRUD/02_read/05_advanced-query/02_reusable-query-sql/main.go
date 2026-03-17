package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
Scopes 方法允许将常用的查询条件定义为可重用的函数
这些作用域再可以在其它查询中进行引用
*/

type Order struct {
	ID          int     `gorm:"primaryKey"`
	OrderNumber int     `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	User        int
}

// 定义 Scopes

// AmountGreaterThan1000 amount 大于 1000
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("price > ?", 1000)
}

// OrderNumberLessThan105 筛选 OrderNumber 小于 105 的订单 Scope
func OrderNumberLessThan105(db *gorm.DB) *gorm.DB {
	return db.Where("order_number < ?", 105)
}

// OrderPayer 用于按用户筛选订单的 Scope
func OrderPayer(users []int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user IN (?)", users)
	}
}

func main() {
	global.Connect()
	find, err := gorm.G[Order](global.DB).Find(context.Background())
	if err != nil {
		return
	}
	fmt.Println(find)

	// 查询中使用 Scopes 方法可以连接一个或多个 Scope 进行查询
	var orders []Order
	global.DB.Scopes(AmountGreaterThan1000, OrderNumberLessThan105, OrderPayer([]int{1})).Find(&orders)
	fmt.Println(orders)
}
