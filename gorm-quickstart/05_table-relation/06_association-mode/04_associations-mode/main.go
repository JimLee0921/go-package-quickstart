package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"
)

/*
GORM 中的关联模式提供了多种辅助方法来处理模型之间的关系，为管理关联数据提供了高效的方式

要启动关联模式，需要指定源模型和关系的字段名称
源模型必须包含主键并且关系的字段名称应与现有的关联字段相匹配
*/
type User struct {
	ID        int
	Name      string
	Emails    []Email
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Email struct {
	ID     int
	UserID int
	Email  string
}

type Language struct {
	ID   int
	Name string
}

func init() {
	global.Connect()

	if err := global.DB.AutoMigrate(
		User{},
		Email{},
		Language{},
	); err != nil {
		log.Fatal(err)
	}

	// ---------- 清空旧数据（保证可重复运行） ----------
	global.DB.Exec("DELETE FROM user_languages")
	global.DB.Exec("DELETE FROM emails")
	global.DB.Exec("DELETE FROM users")
	global.DB.Exec("DELETE FROM languages")

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
		Emails: []Email{
			{Email: "alice-1@example.com"},
			{Email: "alice-2@example.com"},
		},
		Languages: []Language{
			langs[0], // EN
			langs[1], // ZH
		},
	}

	bob := User{
		Name: "Bob",
		Emails: []Email{
			{Email: "bob-1@example.com"},
		},
		Languages: []Language{
			langs[1], // ZH
		},
	}

	if err := global.DB.Create(&alice).Error; err != nil {
		log.Fatal(err)
	}
	if err := global.DB.Create(&bob).Error; err != nil {
		log.Fatal(err)
	}
}

func main() {
	// StartAssociation()
	// QueryAssociationsDemo()
	// AppendAssociationsDemo()
	// ReplaceAssociationsDemo()
	// DeleteAssociationsDemo()
	// ClearAssociationsDemo()
	// CountAssociationsDemo()
	// ConditionalFindAssociationsDemo()
}

func StartAssociation() {
	var user User
	global.DB.Model(&user).Association("Languages")
	// 检查错误
	if err := global.DB.Model(&user).Association("Languages").Error; err != nil {
		fmt.Println(err)
	}
}

// --------- Demo 1: Query Associations ---------

func QueryAssociationsDemo() {
	fmt.Println("\n===== QueryAssociationsDemo =====")
	var alice User
	// 1. 先拿到 alice（Association Mode 仍然需要一个 source model）
	if err := global.DB.Where("name = ?", "Alice").First(&alice).Error; err != nil {
		log.Fatal(err)
	}

	// 2. 查找 Alice 的 Languages
	var languages []Language
	if err := global.DB.Model(&alice).Association("Languages").Find(&languages); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Alice Languages:")
	for _, l := range languages {
		fmt.Printf("  - ID=%d Name=%s\n", l.ID, l.Name)
	}

	// 3. 查看 Alice 的 Emails
	var emails []Email
	if err := global.DB.Model(&alice).Association("Emails").Find(&emails); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Alice Emails:")
	for _, e := range emails {
		fmt.Printf("  - ID=%d Email=%s\n", e.ID, e.Email)
	}
}

// --------- Demo 2: Append Associations ---------

func AppendAssociationsDemo() {
	fmt.Println("\n===== AppendAssociationsDemo =====")

	// 1. 先找到 user
	var bob User
	if err := global.DB.Where("name = ?", "Bob").First(&bob).Error; err != nil {
		log.Fatal(err)
	}
	// 2. 追加前查看语言
	var languagesBeforeAppend []Language
	if err := global.DB.Model(&bob).Association("Languages").Find(&languagesBeforeAppend); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Before Append Languages:")
	for _, l := range languagesBeforeAppend {
		fmt.Printf("  - ID=%d Language=%s\n", l.ID, l.Name)
	}

	// 3. 找到需要追加的语言
	var jp Language
	if err := global.DB.Where("name = ?", "JP").Find(&jp).Error; err != nil {
		log.Fatal(err)
	}

	// 4. 追加关联
	if err := global.DB.Model(&bob).Association("Languages").Append(&jp); err != nil {
		log.Fatal(err)
	}

	// 5. 再次查看结果
	var languagesAfterAppend []Language
	if err := global.DB.Model(&bob).Association("Languages").Find(&languagesAfterAppend); err != nil {
		log.Fatal(err)
	}
	fmt.Println("After Append Languages:")
	for _, l := range languagesAfterAppend {
		fmt.Printf("  - ID=%d Language=%s\n", l.ID, l.Name)
	}
}

// --------- Demo 3: Replace Associations ---------

func ReplaceAssociationsDemo() {
	// 把 Bob 的 Languages 从 [ZH] 替换成 [EN, JP]。
	fmt.Println("\n===== ReplaceAssociationsDemo =====")

	// 1. 查找 user
	var bob User
	if err := global.DB.Where("name = ?", "Bob").First(&bob).Error; err != nil {
		log.Fatal(err)
	}

	// 2. 替换前打印语言
	var languagesBeforeReplace []Language
	if err := global.DB.Model(&bob).Association("Languages").Find(&languagesBeforeReplace); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Before Replace Languages:")
	for _, l := range languagesBeforeReplace {
		fmt.Printf("  - ID=%d Language=%s\n", l.ID, l.Name)
	}

	// 3. 查找 EN / JP 两种语言
	var en, jp Language
	if err := global.DB.Where("name = ?", "EN").First(&en).Error; err != nil {
		log.Fatal(err)
	}
	if err := global.DB.Where("name = ?", "JP").First(&jp).Error; err != nil {
		log.Fatal(err)
	}

	// 4. 替换关联
	if err := global.DB.Model(&bob).Association("Languages").Replace(&en, &jp); err != nil {
		log.Fatal(err)
	}

	// 5. 替换后查询
	var languagesAfterReplace []Language
	if err := global.DB.Model(&bob).Association("Languages").Find(&languagesAfterReplace); err != nil {
		log.Fatal(err)
	}
	fmt.Println("After Replace Languages:")
	for _, l := range languagesAfterReplace {
		fmt.Printf("  - ID=%d Language=%s\n", l.ID, l.Name)
	}
}

// --------- Demo 4: Delete Associations ---------

func DeleteAssociationsDemo() {
	// 删除 Alice 的一个 Language (ZH)
	fmt.Println("\n===== DeleteAssociationsDemo =====")

	// 1. 查找用户
	var alice User
	if err := global.DB.Where("name = ?", "Alice").First(&alice).Error; err != nil {
		log.Fatal(err)
	}

	// 2. 查找 ZH 语言
	var zh Language
	if err := global.DB.Where("name = ?", "ZH").First(&zh).Error; err != nil {
		log.Fatal(err)
	}

	// 3. 删除关系（不会删除 Languages 表中的 ZH）
	if err := global.DB.Model(&alice).Association("Languages").Delete(&zh); err != nil {
		log.Fatal(err)
	}

	// 4. 查询结果
	var languages []Language
	if err := global.DB.Model(&alice).Association("Languages").Find(&languages); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Alice Languages After Delete ZH:")
	for _, l := range languages {
		fmt.Printf("  - %s\n", l.Name)
	}
}

// --------- Demo 5: Clear Associations ---------

func ClearAssociationsDemo() {
	// 清空 Alice 的所有 Languages 和 Emails，清空当前对象的所有关联，对于 many2many 只清空连接表
	fmt.Println("\n===== ClearAssociationsDemo =====")

	// 1. 查找 Alice
	var alice User
	if err := global.DB.Where("name = ?", "Alice").First(&alice).Error; err != nil {
		log.Fatal(err)
	}

	// 2. 清空关联
	if err := global.DB.Model(&alice).Association("Languages").Clear(); err != nil {
		log.Fatal(err)
	}
	if err := global.DB.Model(&alice).Association("Emails").Clear(); err != nil {
		log.Fatal(err)
	}

	// 3. 查询结果
	var (
		emails    []Email
		languages []Language
	)

	if err := global.DB.Model(&alice).Association("Emails").Find(&emails); err != nil {
		log.Fatal(err)
	}
	if err := global.DB.Model(&alice).Association("Languages").Find(&languages); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Alice Languages After Clear: count=%d\n", len(languages))
	fmt.Printf("Alice Emails After Clear: count=%d\n", len(emails))

}

// --------- Demo 6: Associations Count ---------

func CountAssociationsDemo() {
	// 统计 Alice 的 Emails 和 Languages 数量
	fmt.Println("\n===== CountAssociationsDemo =====")

	// 1. 查找 alice
	var alice User
	if err := global.DB.Where("name = ?", "Alice").First(&alice).Error; err != nil {
		log.Fatal(err)
	}

	// 2. 进行统计
	emailCount := global.DB.Model(&alice).Association("Emails").Count()
	languageCount := global.DB.Model(&alice).Association("Languages").Count()

	fmt.Printf("Alice Emails Count: %d\n", emailCount)
	fmt.Printf("Alice Languages Count: %d\n", languageCount)

}

// --------- Demo 7: Conditional Find Association ---------

func ConditionalFindAssociationsDemo() {
	// 只查询 Alice 的 EN / JP Languages
	fmt.Println("\n===== ConditionalFindAssociationsDemo =====")

	// 1. 查找 alice
	var alice User
	if err := global.DB.Where("name = ?", "Alice").First(&alice).Error; err != nil {
		log.Fatal(err)
	}

	// 2. 先查找全部 languages
	var languages []Language
	if err := global.DB.Model(&alice).Association("Languages").Find(&languages); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Alice Languages:")
	for index, language := range languages {
		fmt.Printf("  - index=%d Name=%s\n", index, language.Name)
	}

	// 3. 根据条件查询
	names := []string{"EN", "JP"}
	var languagesFiltered []Language
	if err := global.DB.Model(&alice).Where("name IN ?", names).Association("Languages").Find(&languagesFiltered); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Alice Languages Filtered(EN/JP):")
	for index, language := range languagesFiltered {
		fmt.Printf("  - index=%d Name=%s\n", index, language.Name)
	}
}
