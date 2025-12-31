package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

type User struct {
	Name string `gorm:"not null;unique"`
}

func main() {
	global.Connect()

	migrator := global.DB.Migrator()

	// 获取当前数据库名
	fmt.Println(migrator.CurrentDatabase())

	//TableOperation(migrator)

	ColumnOperation(migrator)

}

func TableOperation(migrator gorm.Migrator) {
	/*
		表相关
	*/
	// 为 User 重建表
	err := migrator.AutoMigrate(&User{})
	// db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&User{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// 检查 	User 表是否存在
	ok := migrator.HasTable(&User{})
	// 等价于 migrator.HasTable("users")
	if ok {
		fmt.Println("Table User exists")
	} else {
		fmt.Println("Table User dont exists")
	}

	// 如果表存在则删除（删除时会忽略、删除外键约束）
	err = migrator.DropTable(&User{})
	// 等价于 migrator.DropTable("users")
	if err != nil {
		return
	}

	// 重命名表
	err = migrator.RenameTable("users", "user_infos")
	// 等价于 migrator.RenameTable(&User{}, &UserInfo{}) 但是推荐上面这种
	// 不然还得定义一个只是临时使用的结构体，UserInfo 结构体结构不关心，只使用名字 user_infos
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ColumnOperation(migrator gorm.Migrator) {
	/*
		字段 Column 列相关操作
	*/

	// 检查 name 字段是否存在
	ok := migrator.HasColumn(&User{}, "Name")
	if ok {
		fmt.Println("Name field exists")
	} else {
		fmt.Println("Name field dont exists")
	}

	// 删除 name 字段
	err := migrator.DropColumn(&User{}, "Name")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// 添加 name 字段
	err = migrator.AddColumn(&User{}, "Name")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 修改 name 字段 Alter 比较危险，默认不会自动执行，必须显式调用
	err = migrator.AlterColumn(&User{}, "Name")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 字段重命名
	err = migrator.RenameColumn(&User{}, "name", "new_name")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 字段类型
	types, err := migrator.ColumnTypes(&User{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range types {
		fmt.Println(t)
	}

}
