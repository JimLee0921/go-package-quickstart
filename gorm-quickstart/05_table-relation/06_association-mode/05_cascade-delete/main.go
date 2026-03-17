package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"
)

/*
GORM 关联的级联删除，在删除主模型时顺带请理关联
Select + Delete
*/

type User struct {
	ID        int
	Name      string
	Emails    []Email
	Account   Account
	Languages []Language `gorm:"many2many:user_languages;"`
}

// Account 一对一
type Account struct {
	ID      int
	UserID  int
	Balance int
}

// Email 一对多
type Email struct {
	ID     int
	UserID int
	Email  string
}

// Language 多对多
type Language struct {
	ID   int
	Name string
}

func init() {
	global.Connect()

	// ---------- 清空旧表（保证可重复运行） ----------
	global.DB.Exec("DROP TABLE IF EXISTS user_languages")
	global.DB.Exec("DROP TABLE IF EXISTS emails")
	global.DB.Exec("DROP TABLE IF EXISTS accounts")
	global.DB.Exec("DROP TABLE IF EXISTS users")
	global.DB.Exec("DROP TABLE IF EXISTS languages")

	if err := global.DB.AutoMigrate(
		User{},
		Email{},
		Language{},
		Account{},
	); err != nil {
		log.Fatal(err)
	}

	// ---------- 准备 Language ----------
	langs := []Language{
		{Name: "EN"},
		{Name: "ZH"},
		{Name: "JP"},
	}
	if err := global.DB.Create(&langs).Error; err != nil {
		log.Fatal(err)
	}

	// ---------- 准备 User + Email + Language ----------
	alice := User{
		Name: "Alice",
		Account: Account{
			Balance: 100,
		},
		Emails: []Email{
			{Email: "alice-1@test.com"},
			{Email: "alice-2@test.com"},
		},
		Languages: []Language{
			langs[0], // EN
			langs[1], // ZH
		},
	}

	if err := global.DB.Create(&alice).Error; err != nil {
		log.Fatal(err)
	}

}

func main() {
	var user User
	if err := global.DB.Where("name = ?", "Alice").First(&user).Error; err != nil {
		log.Fatal(err)
	}

	//// 1. 删除1对1 User + Account（不会成功）
	//if err := global.DB.Select("Account").Delete(&user).Error; err != nil {
	//	fmt.Println(err)
	//}
	//
	//// 2. 删除1对多 User + Emails（不会成功）
	//if err := global.DB.Select("Emails").Delete(&user).Error; err != nil {
	//	fmt.Println(err)
	//}
	//
	//// 3. 删除多对多 User + Languages（不会成功）多对多删除的是关联表中的数据，不会删除 Languages 表中的数据
	//if err := global.DB.Select("Languages").Delete(&user).Error; err != nil {
	//	fmt.Println(err)
	//}

	// 4. 选择多个关联进行删除（这里全选了，所以会成功，等价于 Select(clause.Associations) 全量关联）
	if err := global.DB.Select("Languages", "Emails", "Account").Delete(&user).Error; err != nil {
		fmt.Println(err)
	}

	/*
		因为这些 Association Tags 都没设置 constraint，所以默认删除都会默认引起关联冲突
		所以这里示例可能会有问题，主要是为了演示，其实删除方式都是一样的
	*/
}
