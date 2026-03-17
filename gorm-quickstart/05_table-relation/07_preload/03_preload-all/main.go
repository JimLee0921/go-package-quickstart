package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
可以使用 clause.Associations 进行全部预加载
Preload(clause.Associations) = 自动预加载当前模型的所有一层关联
1. 只加载直接关联（第一层）
2. 不会递归加载嵌套关联
*/

type User struct {
	gorm.Model
	Name      string
	CompanyID uint
	Company   Company
	RoleID    uint
	Role      Role
	Orders    []Order
}

type Company struct {
	gorm.Model
	Name string
}

type Role struct {
	gorm.Model
	Name string
}

type Order struct {
	gorm.Model
	UserID uint
	Price  float64
}

func init() {
	global.Connect()
	db := global.DB

	// 1. 自动建表
	if err := db.AutoMigrate(
		&Company{},
		&Role{},
		&User{},
		&Order{},
	); err != nil {
		log.Fatal(err)
	}

	// 2. 防止重复插入
	var cnt int64
	db.Model(&User{}).Count(&cnt)
	if cnt > 0 {
		return
	}

	// 3. Company
	companies := []Company{
		{Name: "Acme Corp"},
		{Name: "Globex Inc"},
	}
	if err := db.Create(&companies).Error; err != nil {
		log.Fatal(err)
	}

	// 4. Role
	roles := []Role{
		{Name: "Admin"},
		{Name: "User"},
	}
	if err := db.Create(&roles).Error; err != nil {
		log.Fatal(err)
	}

	// 5. User（belongs-to Company + Role）
	users := []User{
		{
			Name:      "Tom",
			CompanyID: companies[0].ID,
			Role:      roles[0], // Admin
		},
		{
			Name:      "Jerry",
			CompanyID: companies[0].ID,
			Role:      roles[1], // User
		},
		{
			Name:      "Lucy",
			CompanyID: companies[1].ID,
			Role:      roles[1], // User
		},
	}
	if err := db.Create(&users).Error; err != nil {
		log.Fatal(err)
	}

	// 6. Orders（has-many）
	orders := []Order{
		{Price: 100},
		{Price: 200},
		{Price: 300},
	}

	// Tom：2 个订单
	if err := db.Model(&users[0]).
		Association("Orders").
		Append(&orders[0], &orders[1]); err != nil {
		log.Fatal(err)
	}

	// Jerry：1 个订单
	if err := db.Model(&users[1]).
		Association("Orders").
		Append(&orders[2]); err != nil {
		log.Fatal(err)
	}

	// Lucy：0 个订单（用于演示 has-many 为空）
}

func main() {
	var users []User
	global.DB.
		Preload(clause.Associations).
		Find(&users)

	for _, u := range users {
		fmt.Printf(
			"User=%-5s | Company=%-10s | Role=%-5s | Orders=%d\n",
			u.Name,
			u.Company.Name,
			u.Role.Name,
			len(u.Orders),
		)
	}
}
