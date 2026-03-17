package main

import "gorm-quickstart/global"

/*
要定义 has many 关系，同样必须存在外键
默认的外键名是拥有者的类型名加上其主键字段名
例如，要定义一个属于 User 的模型，则其外键应该是 UserID
可以改写外键来使用另一个字段名作为外键，使用 foreignKey 标签自定义
*/

type User struct {
	ID          int
	CreditCards []CreditCard `gorm:"foreignKey:UserRefer"` // 改外键字段名为 UserRefer
}

type CreditCard struct {
	ID        int
	Number    string
	UserRefer uint // 指向 User.ID
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, CreditCard{})
	/*
		CREATE TABLE `credit_cards` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `number` longtext COLLATE utf8mb4_general_ci,
		  `user_refer` bigint DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `fk_users_credit_cards` (`user_refer`),
		  CONSTRAINT `fk_users_credit_cards` FOREIGN KEY (`user_refer`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

func main() {

}
