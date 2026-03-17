package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"time"

	"gorm.io/gorm"
)

type Result struct {
	ID   int
	Name string
	Age  int
}
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
	QueryDemo()
	ExecDemo()
}

func QueryDemo() {
	// 原生查询 SQL / Scan

	// Scan 到自定义 result 类型
	result, err := gorm.G[Result](global.DB).Raw("SELECT id, name, age FROM users WHERE id = ?", 1).Find(context.Background())

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	// Scan 进原始类型
	age, err := gorm.G[int](global.DB).Raw("SELECT SUM(age) FROM users WHERE id IN (?, ?, ?)", 0, 1, 2).Find(context.Background())

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(age)
	}

	// Scan 进 ORM 模型切片
	users, err := gorm.G[User](global.DB).Raw("SELECT * FROM users WHERE id > ?", 100).Find(context.Background())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}
}

func ExecDemo() {
	// Exec 原生 SQL
	res := gorm.WithResult()

	err := gorm.G[any](global.DB, res).Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders(
    	id int AUTO_INCREMENT primary key,
    	order_number int,
    	price DECIMAL(10,2),
    	user int
	)`)
	fmt.Println(err, res)

	// 更新数据传入参数
	err = gorm.G[any](global.DB, res).Exec(context.Background(), "INSERT INTO orders(order_number, price, user) VALUES (?, ?, ?), (?, ?, ?)", 101, 50.20, 1, "102", "100.23", "2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.RowsAffected)
	}

	// 执行 SQL 表达式
	err = gorm.G[any](global.DB).Exec(context.Background(), "UPDATE users SET age = ? WHERE name = ?", gorm.Expr("age * ? + ?", 1, 12), "Jinx")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.RowsAffected)
	}
}
