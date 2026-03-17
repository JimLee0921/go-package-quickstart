package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"
)

/*
预加载支持条件查询，类似于行内条件
*/

type User struct {
	ID     int
	Name   string
	State  string
	Orders []Order
}

type Order struct {
	ID          int
	OrderNumber string
	UserID      int
	State       string
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(User{}, Order{}); err != nil {
		log.Fatal(err)
	}

	users := []User{
		{ID: 1, Name: "Alice", State: "active"},
		{ID: 2, Name: "Bob", State: "inactive"},
		{ID: 3, Name: "Charlie", State: "active"},
	}
	if err := global.DB.Create(&users).Error; err != nil {
		log.Fatal(err)
	}

	orders := []Order{
		{ID: 1, OrderNumber: "ORD-001", UserID: 1, State: "paid"},
		{ID: 2, OrderNumber: "ORD-002", UserID: 1, State: "pending"},
		{ID: 3, OrderNumber: "ORD-003", UserID: 2, State: "cancelled"},
		{ID: 4, OrderNumber: "ORD-004", UserID: 3, State: "paid"},
		{ID: 5, OrderNumber: "ORD-005", UserID: 3, State: "cancelled"},
	}
	if err := global.DB.Create(&orders).Error; err != nil {
		log.Fatal(err)
	}
}

func main() {
	var users1, users2 []User
	global.DB.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users1)
	for _, u := range users1 {
		fmt.Printf(
			"User=%s State=%s\n",
			u.Name,
			u.State,
		)
		for _, o := range u.Orders {
			fmt.Printf("		orderNumber=%s State=%s\n",
				o.OrderNumber,
				o.State,
			)
		}
	}
	fmt.Println("============================")
	global.DB.Where("state = ?", "active").Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users2)
	for _, u := range users1 {
		fmt.Printf(
			"User=%s State=%s\n",
			u.Name,
			u.State,
		)
		for _, o := range u.Orders {
			fmt.Printf("		orderNumber=%s State=%s\n",
				o.OrderNumber,
				o.State,
			)
		}
	}
}
