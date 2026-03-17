package main

import "gorm-quickstart/global"

/*
通过 constraint 标签添加外键约束
*/

type User struct {
	ID          int
	CreditCards []CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type CreditCard struct {
	ID     int
	Number string
	UserID uint
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, CreditCard{})
	/*
		CREATE TABLE `users` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci


		CREATE TABLE `credit_cards` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `number` longtext COLLATE utf8mb4_general_ci,
		  `user_id` bigint DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `fk_users_credit_cards` (`user_id`),
		  CONSTRAINT `fk_users_credit_cards` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}
func main() {

}
