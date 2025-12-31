package main

import (
	"database/sql"
	"time"
)

/*
GORM 模型本质就是普通 Go struct + 反射映射，通过反射读取导出字段来决定建表映射规则

主键：GORM 使用一个名为ID 的字段作为每个模型的默认主键
表名：默认情况下，GORM 将结构体名称转换为 snake_case 并为表名加上复数形式。 For instance, a User struct becomes users in the database, and a GormUserName becomes gorm_user_names
列名：GORM 自动将结构体字段名称转换为 snake_case 作为数据库中的列名
时间戳字段：GORM使用字段 CreatedAt 和 UpdatedAt 来自动跟踪记录的创建和更新时间
*/

// User 示例
type User struct {
	// ID 默认就是主键
	ID uint
	// 基本类型，默认 not null 不能为空
	Name string
	Age  uint8
	// 指针类型，可可用
	Email    *string
	Birthday *time.Time
	// 可以为空，参考 sql.NULL
	MemberNumber sql.Null[string]
	ActivateAt   sql.Null[time.Time]
	// 自动时间字段 这两个字段在 Create/Update时 会被 GORM 自动赋值/更新（只要字段存在且导出）
	// 配套 tag: gorm:"autoCreateTime" / gorm:"autoUpdateTime"
	// 即使不写 tag 只要字段名符合规范会自动生效
	CreatedAt time.Time
	UpdatedAt time.Time
	// GORM 直接忽略，不建列也不读写
	ignored string
}

func main() {

}
