package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"
)

/*
GORM 在创建或更新记录时会自动保存其关联和引用，主要使用 upsert 技术更新现有关联的外键引用
当创建一条新的记录时，GORM 会自动保存它的关联数据（这个过程包括向关联表插入数据以及维护外键引用）
*/

type User struct {
	ID                int
	Name              string
	BillingAddressID  *uint
	BillingAddress    Address `gorm:"foreignKey:BillingAddressID;references:ID"`
	ShippingAddressID *uint
	ShippingAddress   Address `gorm:"foreignKey:ShippingAddressID;references:ID"`

	Emails    []Email
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Address struct {
	ID       int
	Address1 string
}

type Email struct {
	ID     int
	UserID uint
	Email  string
}

type Language struct {
	ID   int
	Name string `gorm:"type:varchar(64);uniqueIndex"`
}

func init() {
	global.Connect()
	// 迁移：主表+关联表+many2many 连接表（GORM 自动创建）
	if err := global.DB.AutoMigrate(User{}, Address{}, Email{}, Language{}); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1. 创建 User
	user := &User{
		Name: "JimLee",
		BillingAddress: Address{
			Address1: "Billing Address - Address 1",
		},
		ShippingAddress: Address{
			Address1: "Shipping Address - Address 1",
		},
		Emails: []Email{
			{Email: "JimLee@example.com"},
			{Email: "JimLee-2@example.com"},
		},
		Languages: []Language{
			{Name: "ZH"},
			{Name: "EN"},
		},
	}
	/*
		[4.037ms] [rows:1] INSERT INTO `addresses` (`address1`) VALUES ('Billing Address - Address 1') ON DUPLICATE KEY UPDATE `id`=`id`
		[10.052ms] [rows:1] INSERT INTO `addresses` (`address1`) VALUES ('Shipping Address - Address 1') ON DUPLICATE KEY UPDATE `id`=`id`
		[3.516ms] [rows:2] INSERT INTO `emails` (`user_id`,`email`) VALUES (1,'JimLee@example.com'),(1,'JimLee-2@example.com') ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`)
		[4.376ms] [rows:2] INSERT INTO `languages` (`name`) VALUES ('ZH'),('EN') ON DUPLICATE KEY UPDATE `id`=`id`
		[3.123ms] [rows:2] INSERT INTO `user_languages` (`user_id`,`language_id`) VALUES (1,1),(1,2) ON DUPLICATE KEY UPDATE `user_id`=`user_id`
		[51.380ms] [rows:1] INSERT INTO `users` (`name`,`billing_address_id`,`shipping_address_id`) VALUES ('JimLee',1,2)
	*/

	// 2. Create 会自动保存关联+维护外键+写入 user_languages 连接表
	// 因为 Language 设置了 uniqueIndex 第二次运行会失败
	if err := global.DB.Debug().Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user: ID=%d, BillingAddressID=%v, ShippingAddressID=%v\n",
		user.ID, user.BillingAddressID, user.ShippingAddressID)

	// 3. 预加载查回所有关联进行验证
	var got User
	if err := global.DB.Preload("BillingAddress").
		Preload("ShippingAddress").
		Preload("Emails").
		Preload("Languages").
		First(&got, user.ID).Error; err != nil {
		log.Fatal(err)

	}

	fmt.Println("---- Loaded result ----")
	fmt.Println("User:", got.ID, got.Name)
	fmt.Println("BillingAddress:", got.BillingAddress.ID, got.BillingAddress.Address1)
	fmt.Println("ShippingAddress:", got.ShippingAddress.ID, got.ShippingAddress.Address1)
	fmt.Println("Emails:", len(got.Emails))
	for _, e := range got.Emails {
		fmt.Println(" -", e.Email)
	}
	fmt.Println("Languages:", len(got.Languages))
	for _, l := range got.Languages {
		fmt.Println(" -", l.Name)
	}
}
