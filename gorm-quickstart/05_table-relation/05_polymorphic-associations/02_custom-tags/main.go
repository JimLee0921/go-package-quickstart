package main

import (
	"fmt"
	"gorm-quickstart/global"
)

/*
GORM 多态关联可以通过几个标签进行自定义字段名
在使用 `gorm:"polymorphic:Owner;"` 时 GORM 会要求子表中有 OwnerID 和 OwnerType 作为默认

"polymorphicType:Kind" 会要求不要用 OwnerType，而是用 Kind 这个字段来存父类型，要求子表中必须存在 Kind string
"polymorphicId:OwnerID" 指定 ID 字段的列名，默认规则是 <Prefix>ID，现在是显式告诉 GORM 父对象的 ID 存在 OwnerID 这个字段里面，子表中必须有 OwnerID uint
"polymorphicValue:master" 指定写入类型字段的固定值，默认情况下 GORM 会使用表名（cats/dogs），现在是无论父模型叫什们，Kind 字段一律写 master
*/

type Dog struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphicType:Kind;polymorphicId:OwnerID;polymorphicValue:cat_toy"`
}

type Cat struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphicType:Kind;polymorphicId:OwnerID;polymorphicValue:dog_toy"`
}

type Toy struct {
	ID      int
	Name    string
	OwnerID int
	Kind    string
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(Dog{}, Cat{}, Toy{}); err != nil {
		fmt.Println(err)
	}
}
func main() {
	// 创建 Cat + Toys
	//cat := Cat{
	//	Name: "Kitty",
	//	Toys: []Toy{
	//		{Name: "Ball"},
	//		{Name: "Mouse"},
	//	},
	//}
	//global.DB.Create(&cat)

	/*
		INSERT INTO cats (name) VALUES ('Kitty');

		INSERT INTO toys (name, owner_id, kind)
		VALUES
		('Ball',  cat.ID, 'master'),
		('Mouse', cat.ID, 'master');
	*/

	//dog := Dog{
	//	Name: "Buddy",
	//	Toys: []Toy{
	//		{Name: "Bone"},
	//		{Name: "Ball"},
	//	},
	//}
	//global.DB.Create(&dog)

	/*
		预加载，GORM 自动生成的过滤条件是：
		SELECT * FROM toys
		WHERE owner_id = cat.ID
		  AND kind = 'cat_toy';
	*/
	var cat Cat
	global.DB.Preload("Toys").First(&cat)
	fmt.Println(cat)
}
