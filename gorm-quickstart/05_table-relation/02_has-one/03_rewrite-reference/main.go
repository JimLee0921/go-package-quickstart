package main

import (
	"gorm-quickstart/global"
)

/*
默认情况下，拥有者实体会将 has one 对于模型的逐渐作为外键，也可以修改用另外一个字段保存
使用 references 标签进行更改
*/

type User struct {
	ID         int
	Name       string     `gorm:"type:varchar(64);index"` // 在 User 表创建一个普通索引
	CreditCard CreditCard `gorm:"foreignKey:UserName;references:Name"`
	// foreignKey 表示外键字段在 CreditCard 表中，字段名是 user_name
	// references 表示 credit_cards.user_name 引用的是用户表的 users.name 字段
}

type CreditCard struct {
	ID       int
	Number   string
	UserName string
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, CreditCard{})
	/*
		CREATE TABLE `credit_cards` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `number` longtext COLLATE utf8mb4_general_ci,
		  `user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `fk_users_credit_card` (`user_name`),
		  CONSTRAINT `fk_users_credit_card` FOREIGN KEY (`user_name`) REFERENCES `users` (`name`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

func main() {

}
