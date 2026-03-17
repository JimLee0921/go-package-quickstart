package main

import (
	"context"
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
has one 于另一个模型建立一对一的关联，与一对一的关系有些许不同
表明一个模型的每个实例都包含或拥有另一个模型的一个实例

比如一个 user 和 credit card 模型，每个 user 只能有一张 credit card
*/

// User 每个 user 有一张 CreditCard ，UserID 是外键
type User struct {
	ID         int
	CreditCard CreditCard // 用于在 gorm 这里逻辑上是1对1的
}

type CreditCard struct {
	ID     int
	Number string
	UserID uint `gorm:"uniqueIndex"` // 这里加上唯一索引确保在数据库层面也是1对1，默认去找 users 表的 id
}

// GetAll 检索用户列表并预加载 CreditCard
func GetAll(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error
	return users, err
}
func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, CreditCard{})
	/*
		CREATE TABLE `credit_cards` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `number` longtext COLLATE utf8mb4_general_ci,
		  `user_id` bigint DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  UNIQUE KEY `idx_credit_cards_user_id` (`user_id`),
		  CONSTRAINT `fk_users_credit_card` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

func main() {
	//users, err := GetAll(global.DB)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	for index, user := range users {
	//		fmt.Println(index, user.ID, user.CreditCard.Number)
	//	}
	//}

	// GORM 默认不会自动加载关联对象，不 Preload 预加载的情况下，查询只会取 users 表的列
	// CreditCard 这个字段在 Go 里只是一个零值结构体，不会被填充
	users, err := gorm.G[User](global.DB).Find(context.Background())
	if err != nil {
		fmt.Println(err)
	} else {
		for index, user := range users {
			fmt.Println(index, user.ID, user.CreditCard.Number)
		}
	}
}
