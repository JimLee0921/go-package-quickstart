package main

import (
	"gorm-quickstart/global"
)

/*
对于 has one 关系，同样必须存在外键字段，拥有者把属于它的模型主键保存到这个字段
这个字段名称通常是由 has one 模型的类型加上其主键生成，对于 User 和 CreditCard 就是 UserID
为 user 添加 credit card 时，它会将 user 的 ID 保存到自己的 UserID 字段
如果想要使用另一个字段，可以使用 foreignKey 标签更改
*/

type User struct {
	ID         int
	CreditCard CreditCard `gorm:"foreignKey:UserName"` // 使用 CreditCard 的 UserName 作为外键
}

type CreditCard struct {
	ID       int
	Number   string
	UserName string // 对应的还是 User 的 ID 字段，只是这里叫 UserName
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, CreditCard{})
	/*
		CREATE TABLE `credit_cards` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `number` longtext COLLATE utf8mb4_general_ci,
		  `user_name` bigint DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `fk_users_credit_card` (`user_name`),
		  CONSTRAINT `fk_users_credit_card` FOREIGN KEY (`user_name`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

func main() {

}
