package main

import (
	"gorm-quickstart/global"
	"log"
)

/*
Select/Omit 还可以直接指定关联的字段使用

在GORM中创建或者更新记录时，可以使用Select和Omit方法来指定是否包含某个关联的字段
Omit则能够排除关联模型的特定字段
*/

type User struct {
	ID   int
	Name string

	BillingAddressID *int
	BillingAddress   Address `gorm:"foreignKey:BillingAddressID;references:ID"`

	ShippingAddressID *int
	ShippingAddress   Address `gorm:"foreignKey:ShippingAddressID;references:ID"`
}

type Address struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	Address1 string
	Address2 string
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(User{}, Address{}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Select
	user := User{
		Name: "JimLee",
		BillingAddress: Address{
			Address1: "Billing Addr 1",
			Address2: "Billing Addr 2",
		},
		ShippingAddress: Address{
			Address1: "Shipping Addr 1",
			Address2: "Shipping Addr 2",
		},
	}

	// 只允许 写入 Name，BillingAddressID 和 BillingAddress 的 Address2
	// 可能因为关联关系、某些版本原因，这里测试的时候，如果单独指定 BillingAddress 就会把 Address1 和 Address2 都写入
	// 如果只是使用 BillingAddress.Address1 的话都不会写入
	// 必须是指定 "BillingAddress" 后再指定 "BillingAddress.Address2" 才会单独写入 "BillingAddress.Address2"
	global.DB.Select(
		"name",
		"BillingAddressID",
		"BillingAddress",
		"BillingAddress.Address2",
	).Create(&user)

	// Omit
	user2 := User{
		Name: "BruceLee",
		BillingAddress: Address{
			Address1: "user2 Billing Addr 1",
			Address2: "user2 Billing Addr 2",
		},
		ShippingAddress: Address{
			Address1: "user2 Shipping Addr 1",
			Address2: "user2 Shipping Addr 2",
		},
	}
	// 省略 BillingAddress.Address1
	global.DB.Omit(
		"BillingAddress.Address1",
	).Create(&user2)

}
