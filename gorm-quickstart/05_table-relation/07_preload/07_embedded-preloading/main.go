package main

import (
	"fmt"
	"gorm-quickstart/global"
)

/*
Embedded Preloading 用于 嵌入结构体（embedded struct）尤其是在同一个结构体被多次嵌入的情况下
它的语法与 Nested Preloading 类似，使用 .（点号） 分隔字段路径

Embedded Preload 只能用于 belongs to 关系
*/

type Country struct {
	ID   uint
	Name string
}

type Address struct {
	CountryID uint
	Country   Country
}

type Org struct {
	ID      uint
	Name    string
	Address Address `gorm:"embedded"`
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(Country{}, Org{})
	china := Country{Name: "China"}
	usa := Country{Name: "USA"}

	global.DB.Create(&china)
	global.DB.Create(&usa)

	org1 := Org{
		Name: "ACME Corp",
		Address: Address{
			CountryID: china.ID,
		},
	}

	org2 := Org{
		Name: "Globex Inc",
		Address: Address{
			CountryID: usa.ID,
		},
	}

	global.DB.Create(&org1)
	global.DB.Create(&org2)
}

func main() {
	var orgs []Org
	global.DB.Preload("Address.Country").Find(&orgs)
	fmt.Println(orgs[0].Name)
	fmt.Println(orgs[0].Address.Country.Name)
	fmt.Println(orgs[1].Name)
	fmt.Println(orgs[1].Address.Country.Name)
}
