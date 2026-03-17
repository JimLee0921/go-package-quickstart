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
Preload 是用多条 SQL 分别查询主表和关联表
Joins 预加载是使用一条 LEFT JOIN SQL 把关联一起查出来 主要用于 belongs-to/has-one 这种一对一和多对一关系
*/

type User struct {
	gorm.Model
	Username string

	CompanyID uint
	Company   Company

	// 新增：用于演示“带条件的 JOIN 只影响右表拼接”
	GlobexCompany Company `gorm:"foreignKey:CompanyID;references:ID"`

	ManagerID uint
	Manager   Manager

	Account Account
}

// Company 多对一 user belongs to Company
type Company struct {
	gorm.Model
	Name string
}

// Manager 多对一 user belongs to Manager，Manager 对 Company 也是多对一
type Manager struct {
	gorm.Model
	Name string

	CompanyID uint
	Company   Company
}

// Account 一对一
type Account struct {
	gorm.Model
	UserID  uint `gorm:"uniqueIndex"`
	Balance float64
}

func init() {
	global.Connect()

	if err := global.DB.AutoMigrate(Company{}, Manager{}, User{}, Account{}); err != nil {
		log.Fatal(err)
	}
	// 3. 插入 Company
	companies := []Company{
		{Name: "Acme Corp"},
		{Name: "Globex Inc"},
	}
	if err := global.DB.Create(&companies).Error; err != nil {
		log.Fatal(err)
	}

	// 4. 插入 Manager
	managers := []Manager{
		{Name: "Alice Manager", CompanyID: companies[0].ID}, // Alice 属于 Acme
		{Name: "Bob Manager", CompanyID: companies[1].ID},   // Bob 属于 Globex
	}

	if err := global.DB.Create(&managers).Error; err != nil {
		log.Fatal(err)
	}

	// 5. 插入 User（只设置外键，不带 Account）
	users := []User{
		{
			Username:  "tom",
			CompanyID: companies[0].ID,
			ManagerID: managers[0].ID,
		},
		{
			Username:  "jerry",
			CompanyID: companies[0].ID,
			ManagerID: managers[1].ID,
		},
		{
			Username:  "lucy",
			CompanyID: companies[1].ID,
			ManagerID: managers[1].ID,
		},
	}
	if err := global.DB.Create(&users).Error; err != nil {
		log.Fatal(err)
	}

	// 6. 插入 Account（1 对 1，显式指定 UserID）
	accounts := []Account{
		{UserID: users[0].ID, Balance: 1000},
		{UserID: users[1].ID, Balance: 2500},
		{UserID: users[2].ID, Balance: 500},
	}
	if err := global.DB.Create(&accounts).Error; err != nil {
		log.Fatal(err)
	}
}

func main() {
	// GenericsApiDemo()
	TraditionalApiDemo()
}

func GenericsApiDemo() {
	// Generics Api 中 joins 第一个参数得传入 clause.JoinTarget
	users, err := gorm.G[User](global.DB).
		Joins(clause.JoinTarget{Association: "Company"}, nil).
		Joins(clause.JoinTarget{Association: "Account"}, nil).
		Find(context.Background())
	if err != nil {
		log.Fatal(err)
	} else {
		// 因为没有 Joins Manager，所有的 Manager 都是空
		for _, user := range users {
			fmt.Printf(
				"ID=%-3d | User=%-10s | Company=%-12s | Balance=%8.2f | Manager=%-10s\n",
				user.ID,
				user.Username,
				user.Company.Name,
				user.Account.Balance,
				user.Manager.Name,
			)
		}
	}

	// 自定义预定义 SQL 语句
	users, err = gorm.G[User](global.DB).Joins(
		clause.JoinTarget{Association: "Account"},
		func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			// 设置条件只获取 Balance 为 1000 的数据
			db.Where(Account{Balance: 1000})
			return nil
		},
	).Find(context.Background())
	fmt.Println("==============")
	if err != nil {
		log.Fatal(err)
	} else {
		// 因为没有 Joins Manager，所有的 Manager 都是空
		for _, user := range users {
			fmt.Printf(
				"ID=%-3d | User=%-10s | Balance=%8.2f\n",
				user.ID,
				user.Username,
				user.Account.Balance,
			)
		}
	}
}

func TraditionalApiDemo() {
	var user1, user2 User
	var users1, users2, users3 []User
	if err := global.DB.Joins("Company").Joins("Manager").First(&user1, 1).Error; err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf(
			"ID=%-3d | User=%-10s | Company=%-12s | Balance=%8.2f | Manager=%-10s\n",
			user1.ID,
			user1.Username,
			user1.Company.Name,
			user1.Account.Balance, // 为空
			user1.Manager.Name,
		)
	}

	global.DB.Joins("Company").Joins("Manager").Joins("Account").First(&user2, "users.username = ?", "lucy")

	fmt.Printf(
		"ID=%-3d | User=%-10s | Company=%-12s | Balance=%8.2f | Manager=%-10s\n",
		user2.ID,
		user2.Username,
		user2.Company.Name,
		user2.Account.Balance,
		user2.Manager.Name,
	)
	fmt.Println("===================")
	global.DB.Joins("Company").Joins("Manager").Joins("Account").Find(&users1, "users.id IN ?", []int{2, 3})
	for _, user := range users1 {
		fmt.Printf(
			"ID=%-3d | User=%-10s | Company=%-12s | Balance=%8.2f | Manager=%-10s\n",
			user.ID,
			user.Username,
			user.Company.Name,
			user.Account.Balance, // 为空
			user.Manager.Name,
		)
	}
	fmt.Println("====================")
	// 右表条件查询，所有左表 users 全保留
	// 只有当 company 满足 name='Globex Inc' 时，右表 company 才能拼上
	// 不满足时：user.Company 字段就是零值
	// global.DB.Joins("Company", global.DB.Where("name = ?", "Globex Inc")).Find(&users2)
	global.DB.Joins("Company", global.DB.Where(&Company{Name: "Globex Inc"})).Find(&users2)
	for _, user := range users2 {
		fmt.Printf(
			"ID=%-3d | User=%-10s | Company=%-12s\n",
			user.ID,
			user.Username,
			user.Company.Name,
		)
	}
	fmt.Println("=============================")
	// 查询嵌套遮蔽字段(Manager 中的 Company)
	global.DB.Joins("Manager").Joins("Manager.Company").Find(&users3)
	for _, u := range users3 {
		fmt.Printf(
			"User=%-10s | UserCompany=%-10s | Manager=%-14s | ManagerCompany=%-10s\n",
			u.Username,
			u.Company.Name, // 注意：这里不会有值，因为没 Joins("Company")
			u.Manager.Name,
			u.Manager.Company.Name, // 这里会有值，因为 Joins("Manager.Company")
		)
	}
}
