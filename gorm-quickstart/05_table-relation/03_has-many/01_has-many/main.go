package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
has many 与另一个模型建立一对多的连接，不同于 has one，拥有者可以有零个或多个关联模型
比如 user 和 credit card 模型，每个 user 可以有多张 credit card
*/

type User struct {
	ID          int
	CreditCards []CreditCard // GORM 使用，可以有 0 或多个 CreditCard
}

type CreditCard struct {
	ID     int
	Number string
	UserID uint // 自动查找 User 的 ID
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
		  KEY `fk_users_credit_card` (`user_id`),
		  CONSTRAINT `fk_users_credit_card` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

// GetAll 检索用户列表并预加载信用卡
func GetAll(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Model(&User{}).Preload("CreditCards").Find(&users).Error
	return users, err
}

func main() {
	users, err := GetAll(global.DB)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, user := range users {
			fmt.Println(user.ID, user.CreditCards)
		}
	}
}
