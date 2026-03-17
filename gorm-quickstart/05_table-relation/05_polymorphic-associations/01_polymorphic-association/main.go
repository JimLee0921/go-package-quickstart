package main

import (
	"fmt"
	"gorm-quickstart/global"
)

/*
GORM 还支持 polymorphic associations 也就是多态关联:
同一个子表中的记录可能属于多个不同的父表类型（不同的模型 struct），依赖一个类型字段来区分是哪一种
比如 Cat 和 Dog 都有很多 Toy ，但是不想为 Cat 和 Dog 单独写两个独立的关联，而是想要使用同一张表 Toy 同时关联猫和狗
可以使用 polymorphic 标签进行定义

GORM 多态关联主要支持：
has one 一对一的多态
has many 一对多的多态
也就是可以让一个模型拥有的字段不止一种模型类型

*/

type Dog struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"` // 把这个字段声明为多态关联，并约定使用 OwnerID 和 OwnerType 这两个字段表示它的父对象
}

type Cat struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
	ID        int
	Name      string
	OwnerID   int    // 保存父对象的 ID
	OwnerType string // 保存父对象的类型名（默认是表名的负数，如 dogs / cats）
}

/*
db.Create(&Dog{Name: "dog1", Toys: []Toy{{Name:"toy1"}}})

插入后：

INSERT INTO toys (name, owner_id, owner_type) VALUES ("toy1", 1, "dogs");

这样就用同一个 toys 表记录了 Dog 和 Cat 的玩具，且通过 owner_type 字段区分是属于猫还是狗
*/

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(Dog{}, Cat{}, Toy{}); err != nil {
		fmt.Println(err)
	}
}
func main() {
	/*
		GORM 的 polymorphic 会在创建 Cat / Dog 时，自动向 Toy 表插入数据，并自动填充 OwnerID 和 OwnerType
	*/
	global.DB.Create(&Dog{
		Name: "Husky",
		Toys: []Toy{{Name: "Toy1"}},
	})
	global.DB.Create(&Cat{
		Name: "bosicat",
		Toys: []Toy{{Name: "Toy1"}},
	})
}
