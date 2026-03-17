package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

/*
GORM 的预加载 Preload 就是在查询主表时，提前把关联表中的数据一并查询出来，避免后续反复查询数据库
解决典型的 N + 1 查询问题
*/

type User struct {
	gorm.Model
	Username string
	Role     Role
	Orders   []Order
}

type Order struct {
	gorm.Model
	UserID uint
	Price  float64
}

type Role struct {
	gorm.Model
	UserID uint
	Name   string
}

func init() {
	global.Connect()

	_ = global.DB.AutoMigrate(User{}, Order{}, Role{})

	users := []User{
		{
			Username: "alice",
			Orders: []Order{
				{Price: 100.5},
				{Price: 200.0},
			},
			Role: Role{
				Name: "Admin",
			},
		},
		{
			Username: "bob",
			Orders: []Order{
				{Price: 50.0},
			},
			Role: Role{
				Name: "normal user",
			},
		},
	}

	// 4. 一次性创建（会自动创建 Orders）
	if err := global.DB.Create(&users).Error; err != nil {
		log.Fatal(err)
	}
}

func main() {
	// GenericsApiDemo()
	TraditionalApiDemo()
}

func GenericsApiDemo() {
	// Generics API
	// 查找用户时预加载订单
	users, err := gorm.G[User](global.DB).Preload("Orders", nil).Preload("Role", nil).Find(context.Background())
	// SELECT * FROM users;
	// SELECT * FROM roles WHERE user_id IN (...);
	// SELECT * FROM orders WHERE user_id IN (...);
	if err != nil {
		log.Fatal(err)
	} else {
		for _, user := range users {
			fmt.Println(user.ID, user.Username, user.Role.Name)
			for _, order := range user.Orders {
				fmt.Println("   ", order.ID, order.Price)
			}
		}
	}
	fmt.Println("=============================")
	// 自定义预加载 SQL
	users, err = gorm.G[User](global.DB).Preload("Role", nil).Preload("Orders", func(db gorm.PreloadBuilder) error {
		db.Order("orders.price DESC")
		return nil
	}).Find(context.Background())
	// SELECT * FROM users;
	// SELECT * FROM roles WHERE user_id IN (...);
	// SELECT * FROM orders WHERE user_id IN (...) order by orders.price DESC;
	if err != nil {
		log.Fatal(err)
	} else {
		for _, user := range users {
			fmt.Println(user.ID, user.Username, user.Role.Name)
			for _, order := range user.Orders {
				fmt.Println("   ", order.ID, order.Price)
			}
		}
	}

}

func TraditionalApiDemo() {
	var users1 []User
	// SELECT * FROM users;
	// SELECT * FROM orders WHERE user_id IN (...);
	if err := global.DB.Preload("Orders").Find(&users1).Error; err != nil {
		log.Fatal(err)
	} else {
		for _, user := range users1 {
			fmt.Println(user.ID, user.Username, user.Role.Name)
			for _, order := range user.Orders {
				fmt.Println("   ", order.ID, order.Price)
			}
		}
	}
	fmt.Println("=============================")
	var users2 []User
	if err := global.DB.Preload("Orders").Preload("Role").Find(&users2).Error; err != nil {
		log.Fatal(err)
	} else {
		for _, user := range users2 {
			fmt.Println(user.ID, user.Username, user.Role.Name)
			for _, order := range user.Orders {
				fmt.Println("   ", order.ID, order.Price)
			}
		}
	}

}
