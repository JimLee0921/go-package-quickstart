package main

import (
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
	var result1 Result
	err1 := global.DB.Raw("SELECT id, name, age FROM users WHERE id = ?", 1).Scan(&result1).Error
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(result1)
	}

	var result2 Result
	err2 := global.DB.Raw("SELECT id, name, age FROM users WHERE name = ?", "JimLee").Scan(&result2).Error
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println(result2)
	}

	var age int
	err3 := global.DB.Raw("SELECT SUM(age) FROM users WHERE age = ?", 22).Scan(&age).Error
	if err2 != nil {
		fmt.Println(err3)
	} else {
		fmt.Println(age)
	}

	var users []User
	err4 := global.DB.Raw("SELECT * FROM users WHERE id > ?", 100).Scan(&users).Error
	if err4 != nil {
		fmt.Println(err4)
	} else {
		fmt.Println(users)
	}
}

func ExecDemo() {
	global.DB.Exec(`CREATE TABLE IF NOT EXISTS orders(
    	id int AUTO_INCREMENT primary key,
    	order_number int,
    	price DECIMAL(10,2),
    	user int
	)`)

	res1 := global.DB.Exec("INSERT INTO orders(order_number, price, user) VALUES (?, ?, ?), (?, ?, ?)", 101, 50.20, 1, "102", "100.23", "2")
	fmt.Println(res1.Error, res1.RowsAffected)

	res2 := global.DB.Exec("UPDATE orders SET price = ? WHERE user = ?", gorm.Expr("price * ? + ?", 10000, 1), 1)
	fmt.Println(res2.Error, res2.RowsAffected)
}
