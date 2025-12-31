package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       int       `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Age      int       `gorm:"default:18"`
	Birthday time.Time `gorm:"autoCreateTime"`
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(&User{})
}

/*
配合使用 db.Where() 添加条件
*/
func main() {
	// WhereStringDemo()
	// WhereStructAndMapDemo()
	// InlineDemo()
	// NotDemo()
	// OrDemo()
	// SelectDemo()
	// OrderDemo()
	// LimitAndOffsetDemo()
	// GroupAndHavingDemo()
	DistinctDemo()
}

// WhereStringDemo 字符串条件
func WhereStringDemo() {
	ctx := context.Background()
	// 获取第一条匹配到的记录
	user, err := gorm.G[User](global.DB).Where("name = ?", "Gogo").First(ctx)
	// SELECT * FROM users WHERE name = 'Gogo' ORDER BY id LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}

	// 匹配所有记录 != JimLee 的记录
	users, err := gorm.G[User](global.DB).Where("name <> ?", "JimLee").Find(ctx)
	// SELECT * FROM users WHERE name <> 'JimLee';
	// SELECT * FROM users WHERE name != 'JimLee';
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// IN 查找条件查找所有 name 为 JimLee, Dsb
	users, err = gorm.G[User](global.DB).Where("name IN ?", []string{"JimLee", "Dsb"}).Find(ctx)
	// SELECT * FROM users WHERE name IN ('JimLee','Dsb');
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// LIKE 模糊查询
	users, err = gorm.G[User](global.DB).Where("name LIKE ?", "%imL%").Find(ctx)
	// SELECT * FROM users WHERE name LIKE '%imL%';
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// AND 并集查询
	users, err = gorm.G[User](global.DB).Where("name = ? AND age >= ?", "JimLee", 25).Find(ctx)
	// SELECT * FROM users WHERE name = 'JimLee' AND age >= 25;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// Time 时间过滤
	users, err = gorm.G[User](global.DB).Where("birthday > ?", time.Now().Add(-1*time.Second*3600)).Find(ctx)
	// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// BETWEEN 条件查询
	// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';
	users, err = gorm.G[User](global.DB).Where("birthday BETWEEN ? AND ?", "2025-12-30 08:29:43.522", "2025-12-30 08:31:56.440").Find(ctx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// 如果对象设置了主键，条件查询将不会覆盖主键的值，而是用 And 连接条件
	var existUser = User{ID: 10}
	global.DB.Where("id = ?", 20).First(&existUser)
	// SELECT * FROM users WHERE id = 10 and id = 20 ORDER BY id ASC LIMIT 1

}

// WhereStructAndMapDemo Struct & Map & Slice 条件
func WhereStructAndMapDemo() {
	ctx := context.Background()

	// Struct 注意使用结构体作为查询条件时 GORM 只会使用非零值字段来构建查询条件
	// 如果某个字段的值是 0、""、false 或其他零值，那这个字段不会被用于生成查询条件
	first, err := gorm.G[User](global.DB).Where(&User{
		Name: "JimLee",
		Age:  22,
	}).First(ctx)
	// SELECT * FROM users WHERE name = JimLee AND age = 22 ORDER BY id LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(first)
	}

	// Struct 作为查询条件时强制让零值也参与查询
	// 可以指定字段来构建查询条件，无论是不是这些字段是不是零值
	users, err := gorm.G[User](global.DB).Where(&User{Name: "wtf"}, "name", "age").Find(ctx)
	// SELECT * FROM users WHERE name = "wtf" AND age = 0;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// 这里 name 不指定会被忽略不参与构建语句 &User{Name: "wtf"} 可以直接写为空值也就是 &User{}
	users, err = gorm.G[User](global.DB).Where(&User{Name: "wtf"}, "age").Find(ctx)
	// SELECT * FROM users WHERE age = 0;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// Map
	users, err = gorm.G[User](global.DB).Where(map[string]any{
		"name": "JimLee",
		"age":  22,
	}).Find(ctx)
	// SELECT * FROM users WHERE name = "JimLee" AND age = 22;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// 主键 Slice
	users, err = gorm.G[User](global.DB).Where([]int{20, 21, 22}).Find(ctx)
	// SELECT * FROM users WHERE id IN (20, 21, 22);
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

}

// InlineDemo 内联条件查询，把条件直接写入 First、Find、Take 等查询方法，写法和 Where 一样
// 这里不能使用泛型了
func InlineDemo() {
	var user1 User
	global.DB.First(&user1, "id = ?", 10)
	// SELECT * FROM users WHERE id = 10
	fmt.Println(user1)

	var users1 []User
	global.DB.Find(&users1, "name = ?", "JimLee")
	// SELECT * FROM users WHERE name = "JimLee"
	fmt.Println(users1)

	var users2 []User
	global.DB.Find(&users2, "name != ? AND age > ?", "JimLee", 20)
	// SELECT * FROM users WHERE name != "JimLee" AND age > 20
	fmt.Println(users2)

	// 配合 Struct
	var users3 []User
	global.DB.Find(&users3, User{Age: 22})
	// SELECT * FROM users WHERE age = 22;
	fmt.Println(users3)

	// 配合 Map
	var users4 []User
	global.DB.Find(&users4, map[string]any{
		"age": 22,
	})
	// SELECT * FROM users WHERE age = 22;
	fmt.Println(users4)
}

// NotDemo Not 条件，使用方式类似于 Where
func NotDemo() {
	ctx := context.Background()

	// 字符串
	first, err := gorm.G[User](global.DB).Not("name = ?", "JimLee").First(ctx)
	// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(first)
	}

	// Not In
	users, err := gorm.G[User](global.DB).Not(map[string]any{
		"name": []string{"JimLee", "Dsb"},
	}).Find(ctx)
	// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// Struct
	first, err = gorm.G[User](global.DB).Not(User{Name: "JimLee", Age: 20}).First(ctx)
	// SELECT * FROM users WHERE name != "JimLee" AND age != 18 ORDER BY id LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(first)
	}

	// Not In slice of primary keys
	first, err = gorm.G[User](global.DB).Not([]int{1, 2, 3}).First(ctx)
	// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(first)
	}
}

// OrDemo Or 条件查询
func OrDemo() {
	ctx := context.Background()

	// 字符串
	users, err := gorm.G[User](global.DB).Where("name = ?", "JimLee").Or("age = ?", 18).Find(ctx)
	// SELECT * FROM users WHERE name = "JimLee" OR age = 18;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// 结构体
	users, err = gorm.G[User](global.DB).Where("name = ?", "JimLee").Or(User{Name: "Dsb", Age: 123}).Find(ctx)
	// SELECT * FROM users WHERE name = "JimLee" OR (name = 'jinzhu 2' AND age = 18);
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// Map
	users, err = gorm.G[User](global.DB).Where("name = ?", "JimLee").Or(map[string]any{
		"name": "DSB",
		"age":  18,
	}).Find(ctx)
	// SELECT * FROM users WHERE name = 'JimLee' OR (name = 'DSB' AND age = 18);
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}
}

// SelectDemo Select 选择特定字段，GORM 默认是选择所有查找字段也就是 SELECT *
// 更多智能选择字段高级查询这里先不学习了
func SelectDemo() {
	ctx := context.Background()

	// 注意如果还打印其它字段就都是 0 值
	users, err := gorm.G[User](global.DB).Select("name", "age").Find(ctx)
	// 等价于 .Select([]string{"name", "age"})
	// SELECT name, age FROM users
	if err != nil {
		fmt.Println(err)
	} else {
		for index, user := range users {
			fmt.Println(index, user.Name, user.Age)
		}
	}

	//rows, err = gorm.G[User](global.DB).Select("COALESCE(age,?)", 20).Rows(ctx)
	//// SELECT COALESCE(age, '20') FROM users;
	//if err != nil {
	//	fmt.Println(err)
	//}
}

// OrderDemo order by 排序，可以使用多个
func OrderDemo() {
	ctx := context.Background()
	users, err := gorm.G[User](global.DB).Where("age > ?", 8).Order("age desc, name").Find(ctx)
	// 等价于：.Order("age desc").Order("name")
	// SELECT * FROM users ORDER BY age desc, name;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}
}

// LimitAndOffsetDemo 限制返回条数和偏移跳过
func LimitAndOffsetDemo() {
	ctx := context.Background()

	users, err := gorm.G[User](global.DB).Limit(3).Find(ctx)
	// SELECT * FROM users LIMIT 3;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// -1 为 取消限制
	users, err = gorm.G[User](global.DB).Limit(10).Limit(-1).Find(ctx)
	// SELECT * FROM users LIMIT 10;
	// SELECT * FROM users;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// 使用 Generics API 调用 Offset 必须配合 LIMIT(Traditional API 可以但是全表扫描很危险)
	users, err = gorm.G[User](global.DB).Limit(2).Offset(3).Find(ctx)
	// SELECT * FROM users OFFSET 3 LIMIT 2;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	// Offset 也能传入 -1 表示取消
	users, err = gorm.G[User](global.DB).Limit(2).Offset(-1).Find(ctx)
	// SELECT * FROM users LIMIT 2;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}
}

// GroupAndHavingDemo 分组过滤
func GroupAndHavingDemo() {

	// 定义聚合结果专用结构体
	type result struct {
		Name  string
		Total int
	}
	var results1 []result
	err := global.DB.Model(&User{}).Select("name, sum(age) as total").Where("name LIKE ?", "B%").Group("name").Find(&results1).Error
	// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "B%" GROUP BY `name` LIMIT 1;
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(results1)
	}

	var results2 []result
	err = global.DB.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "JimLee").Find(&results2).Error
	// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "JimLee";
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(results2)
	}

	rows1, err := global.DB.Table("users").Select("date(birthday) as date, sum(age) as total").Group("date(birthday)").Rows()
	if err != nil {
		fmt.Println(err)
	} else {
		defer rows1.Close()
		for rows1.Next() {
			var (
				date  time.Time
				total int64
			)
			if err := rows1.Scan(&date, &total); err != nil {
				fmt.Println("scan error:", err)
				return
			}

			fmt.Println("date:", date, "total:", total)
		}
	}

	rows2, err := global.DB.Table("users").Select("date(birthday) as date, sum(age) as total").Group("date(birthday)").Having("sum(age) > ?", 100).Rows()
	if err != nil {
		fmt.Println(err)
	} else {
		defer rows2.Close()
		for rows2.Next() {
			var (
				date  time.Time
				total int64
			)
			if err := rows2.Scan(&date, &total); err != nil {
				fmt.Println("scan error:", err)
				return
			}

			fmt.Println("date:", date, "total:", total)
		}
	}
}

// DistinctDemo 用于去重
func DistinctDemo() {
	users, err := gorm.G[User](global.DB).Distinct("name", "age").Order("name, age desc").Find(context.Background())
	// SELECT DISTINCT name, age FROM users ORDER BY name ASC, age DESC;
	if err != nil {
		fmt.Println(err)
	} else {
		for index, user := range users {
			fmt.Println(index, user.Name, user.Age)
		}
	}
}
