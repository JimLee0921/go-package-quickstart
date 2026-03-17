package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
GORM 在做关联保存时会默认把所有可保存的字段和关联都参与进来，但是可能遇到下面一些问题

1. 想只保存主表指定字段，不想自动保存所有关联
2. 想保存主表+某些关联字段但不想触发关联自动创建/更新
3. 在 many2many 场景想控制哪些关联参语插入/更新
4. 想跳过某些字段的更新（例如有默认值、由 DB 触发的字段等）

比如 db.Create(&user) 且 user 由关联字段：
Emails    []Email
Languages []Language `gorm:"many2many:user_languages;"`
GORM 默认行为时：
1. 先插入主表 user
2. 再插入 emails
3. 再插入 languages
4. 再插入 user_language

这时可能已在 DB 里有现存的 Email/Language 不想重复插入或只想插入主表和特定关联或不想写关联，GORM 提供了
Select("Field1", "Field2", ...)：只指定要保存的字段/关联，只对指定字段或关联执行插入/更新
Omit("FieldX", "AssociationY", ...)：对比 Select 的只保留，Omit 是剔除不想保存的字段/关联
作为对关联/字段范围控制的标准做法
*/

type User struct {
	ID        int
	Name      string
	Emails    []Email
	Languages []Language `gorm:"many2many:user_languages"`
}

type Email struct {
	ID     int
	UserID int
	Email  string
}

type Language struct {
	ID   int
	Name string `type:"varchar(60)"`
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(User{}, Email{}, Language{}); err != nil {
		log.Fatal(err)
	}
}

func main() {

	user := User{
		Name:      "Bob",
		Emails:    []Email{{Email: "bob@x.com"}},
		Languages: []Language{{Name: "EN"}, {Name: "DE"}},
	}
	fmt.Println("=== Select(Name, Emails) ===")
	global.DB.Select("Name", "Emails").Create(&user)
	printAll(global.DB)

	user2 := User{
		Name:      "JimLee",
		Emails:    []Email{{Email: "JimLee@x.com"}},
		Languages: []Language{{Name: "ZH-CN"}, {Name: "ZH-EB"}},
	}

	fmt.Println("=== Omit(Emails) ===")
	global.DB.Omit("Emails").Create(&user2)
	printAll(global.DB)

	// 使用 Omit(clause.Associations) 跳过所有关联
	user3 := User{
		Name:      "Crazy",
		Emails:    []Email{{Email: "Crazy@x.com"}},
		Languages: []Language{{Name: "ZH"}, {Name: "ZH"}},
	}
	fmt.Println("=== Omit(ALL Associations) ===")
	global.DB.Omit(clause.Associations).Create(&user3)
	printAll(global.DB)

}

func printAll(db *gorm.DB) {
	var users []User
	db.Preload("Emails").Preload("Languages").Find(&users)

	fmt.Println("Users:")
	for _, u := range users {
		fmt.Printf("  %d: %s\n", u.ID, u.Name)
		fmt.Println("    Emails:", len(u.Emails))
		for _, e := range u.Emails {
			fmt.Println("      -", e.Email)
		}
		fmt.Println("    Languages:", len(u.Languages))
		for _, l := range u.Languages {
			fmt.Println("      -", l.Name)
		}
	}
	fmt.Println()
}
