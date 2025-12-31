package main

/*
可导出的字段默认在使用 GORM 进行 CRUD 时拥有全部权限
可以使用标签控制字符级别权限，可以设置一个字段为只读、只写、只创建、只更新或者被忽略
*/

type User struct {
	//Name string `gorm:"<-:create"`          // 允许读和创建
	//Name string `gorm:"<-:update"`          // 允许读和更新
	//Name string `gorm:"<-"`                 // 允许读和写（创建和更新）
	//Name string `gorm:"<-:false"`           // 允许读，禁止写
	//Name string `gorm:"->"`                 // 只读（除非有自定义配置，否则禁止写）
	//Name string `gorm:"->;<-:create"`       // 允许读和写
	//Name string `gorm:"->:false;<-:create"` // 仅创建（禁止从 db 读）
	//Name string `gorm:"-"`                  // 通过 struct 读写会忽略该字段
	//Name string `gorm:"-:all"`              // 通过 struct 读写、迁移会忽略该字段
	//Name string `gorm:"-:migration"`        // 通过 struct 迁移会忽略该字段
}

func main() {

}
