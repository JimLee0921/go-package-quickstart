package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

type User struct {
	ID      int
	Name    string
	Account Account
	Emails  []Email
	Tags    []Tag `gorm:"many2many:user_tags;"`
}

type Account struct {
	ID        int
	UserID    int
	Balance   int
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除字段
}

// Email has many：私有数据（适合 Unscoped）
type Email struct {
	ID        int
	UserID    *int
	Email     string
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除字段
}

// Tag many2many：共享数据（慎用 Unscoped，软删除和硬删除都不推荐）
type Tag struct {
	ID   int
	Name string
}

func init() {
	global.Connect()

	global.DB.Exec("DROP TABLE IF EXISTS user_tags")
	global.DB.Exec("DROP TABLE IF EXISTS emails")
	global.DB.Exec("DROP TABLE IF EXISTS tags")
	global.DB.Exec("DROP TABLE IF EXISTS accounts")
	global.DB.Exec("DROP TABLE IF EXISTS users")

	if err := global.DB.AutoMigrate(&User{}, &Email{}, &Tag{}, Account{}); err != nil {
		log.Fatal(err)
	}

	tags := []Tag{{Name: "go"}, {Name: "db"}}
	global.DB.Create(&tags)

	u := User{
		Name: "Alice",
		Emails: []Email{
			{Email: "a1@test.com"},
			{Email: "a2@test.com"},
		},
		Tags: []Tag{tags[0], tags[1]},
		Account: Account{
			Balance: 500,
		},
	}
	global.DB.Create(&u)
}

func main() {
	/*
		默认情况
			默认情况下 Association Mode 的 Replace / Delete / Clear 默认只会断开关系，而不会删除记录
			hse one / has many：把子表的外键设为 NULL
			many2many：删除 join table 里的关联行
			不会删除关联表里的记录
			可以使用 Unscoped 开启越权开关，不使用默认的安全策略，软删除 = 可恢复、保留历史；硬删除 = 不可逆、风险高
		Unscoped
			使用 Unscoped 可以开启软删除或硬删除
			软删除并不真正从数据库中删除记录，而是通过一个标志字段（通常是 deleted_at）来表示已删除，只要模型中包含 gorm.DeletedAt 字段，就会自动启用软删除机制
			db.Model(&user).Association("Languages").Unscoped().Clear()
			硬删除是真正把数据从数据库中删除
			db.Unscoped().Model(&user).Association("Languages").Unscoped().Clear()
		使用原则
			一对一可以使用软删除，也可以使用硬删除
			一对多建议使用软删除，也可以使用硬删除但不太建议
			多对多不建议使用 Unscope 进行操作，因为比如 languages 中的数据也在被其它数据使用着，join table 默认就是硬删除了
	*/
	// SoftDeleteOneToMany()
	// PermanentDeleteOneToMany()
	// SoftDeleteOneToOne()
	// PermanentDeleteOneToOne()

}

func SoftDeleteOneToMany() {
	fmt.Println("-------------SoftDeleteOneToMany-------------")
	/*
		has many 软删除 Association + Unscoped = 触发软删除（DeletedAt 被设置）
			emails 表：记录仍在
			emails.deleted_at：被设置
			普通查询 db.Find(&emails) 看不到这些记录
			db.Unscoped().Find(&emails) 能看到

	*/
	var u User
	global.DB.Where("name = ?", "Alice").First(&u)

	if err := global.DB.Model(&u).Association("Emails").Unscoped().Clear(); err != nil {
		log.Fatal(err)
	}

	var emails []Email
	if err := global.DB.Unscoped().Find(&emails).Error; err != nil {
		log.Fatal(err)
	}
	for _, email := range emails {
		fmt.Println(email.ID, email.UserID, email.Email, email.DeletedAt)
	}
}

func PermanentDeleteOneToMany() {
	/*
		一对多硬删除，这里使用 Delete 删除 email 为 a1@test.com 的
	*/
	fmt.Println("-------------PermanentDeleteOneToMany-------------")

	var u User
	global.DB.Where("name = ?", "Alice").First(&u)

	var emails []Email
	if err := global.DB.Model(&u).Association("Emails").Find(&emails); err != nil {
		log.Fatal(err)
	}
	// 前面的 db.Unscoped()：允许操作软删 user
	// 后面的 Association().Unscoped()：强制对关联记录的删除不走软删（若 Email 有 DeletedAt）
	if err := global.DB.Unscoped().Model(&u).Association("Emails").Unscoped().Delete(&emails); err != nil {
		log.Fatal(err)
	}
}

func SoftDeleteOneToOne() {
	fmt.Println("-------------SoftDeleteOneToOne-------------")
	/*
		软链接删除 Account
	*/
	var u User
	global.DB.Where("name = ?", "Alice").First(&u)

	var account Account
	if err := global.DB.Model(&u).Association("Account").Find(&account); err != nil {
		log.Fatal(err)
	}

	if err := global.DB.Model(&u).Association("Account").Unscoped().Delete(&account); err != nil {
		log.Fatal(err)
	}

}

func PermanentDeleteOneToOne() {
	fmt.Println("-------------PermanentDeleteOneToOne-------------")
	/*
		硬链接删除 Account
	*/
	var u User
	global.DB.Where("name = ?", "Alice").First(&u)

	var account Account
	if err := global.DB.Model(&u).Association("Account").Find(&account); err != nil {
		log.Fatal(err)
	}

	if err := global.DB.Unscoped().Model(&u).Association("Account").Unscoped().Delete(&account); err != nil {
		log.Fatal(err)
	}

}
