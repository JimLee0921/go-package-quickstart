package main

import "time"

/*
GORM 已经自定义了 CreatedAt 和 UpdatedAt 进行追踪创建/更新时间
如果定义了这两个字段 GORM 内部会自动在创建/更新时进行处理
如果需要自定义名称是由 `gorm:"autoCreateTime"` 和 `gorm:"autoUpdateTime"`
默认是 time.Time 级别，可以指定时间戳单位
*/

type User struct {
	CreatedAt time.Time // 当前时间填充
	UpdatedAt int       // 时间戳秒数
	Updated   int64     `gorm:"autoUpdateTime:nano"` // 时间戳纳秒
	//Updated   int64     `gorm:"autoUpdateTime:milli"` // 使用时间戳毫秒数填充更新时间
	Created int64 `gorm:"autoCreateTime"` // 使用时间戳秒数填充创建时间
}

func main() {

}
