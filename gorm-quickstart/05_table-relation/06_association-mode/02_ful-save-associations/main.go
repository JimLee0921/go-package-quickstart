package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

/*
GORM 有个 FullSaveAssociations 配置选项
可以指定在更新主模型时是否对所有关联数据做全量保存（完全替换/同步管理）还是部分保存（只保存有变更的关联）

这个配置影响的是 UPDATE/SAVE 时关联如何处理，而不是 CREATE

开启（true）：更新主记录时，对所有有定义的关联都会重新保存（覆盖关联）

关闭（false，默认）：更新主记录时，只保存明确发生变更的关联，不会盲目覆盖
*/

type User struct {
	ID        int
	Name      string
	Emails    []Email
	Languages []Language `gorm:"many2many:user_languages;"`
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
	if err := global.DB.AutoMigrate(User{}, Email{}, Language{}); err != nil {
		log.Fatal(err)
	}

}
func main() {
	// 1. 创建用户+关联
	u := User{
		Name: "JimLee",
		Emails: []Email{
			{Email: "a@example.com"},
			{Email: "b@example.com"},
		},
		Languages: []Language{
			{Name: "EN"},
			{Name: "ZH"},
		},
	}
	if err := global.DB.Create(&u).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("== After Create ==")
	printUser(global.DB, u.ID)

	// 2. 取回并修改（保留关联对象的 ID 才能更新已有记录）
	var u1 User
	if err := global.DB.Preload("Emails").Preload("Languages").First(&u1, u.ID).Error; err != nil {
		log.Fatal(err)
	}

	u1.Name = "JimLee_V2"
	u1.Emails[0].Email = "a+changed@example.com"
	u1.Languages[0].Name = "EN-US"

	// 3. 默认 Save 只保存主表更新，关联对象字段通常不会被保存
	if err := global.DB.Save(&u1).Error; err != nil {
		fmt.Println("hahaha")
		log.Fatal(err)
	}

	fmt.Println("\n== After Save (default, FullSaveAssociations=false) ==")
	printUser(global.DB, u.ID)

	// 4. 再次修改关联字段
	var u2 User
	if err := global.DB.Preload("Emails").Preload("Languages").First(&u2, u.ID).Error; err != nil {
		log.Fatal(err)
	}
	u2.Name = "JimLee_V3"
	u2.Emails[0].Email = "a+fullsave@example.com"
	u2.Languages[0].Name = "EN-GB"

	// 5. 开启 FullSaveAssociations 会把关联对象一起保存（更新 emails / languages）
	if err := global.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&u2).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n== After Save (FullSaveAssociations=true) ==")
	printUser(global.DB, u.ID)
}

func printUser(db *gorm.DB, userID int) {
	var out User
	if err := db.Preload("Emails").Preload("Languages").First(&out, userID).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User:  ID=%d  Name=%s\n", out.ID, out.Name)

	fmt.Println("Emails:")
	for _, e := range out.Emails {
		fmt.Printf("  - (ID=%d) %s\n", e.ID, e.Email)
	}

	fmt.Println("Languages:")
	for _, l := range out.Languages {
		fmt.Printf("  - (ID=%d) %s\n", l.ID, l.Name)
	}

}
