package main

import "gorm-quickstart/global"

type User struct {
	ID         int
	CreditCard CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type CreditCard struct {
	Number string
	UserID uint
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, CreditCard{})
	/*
		CREATE TABLE `credit_cards` (
		  `number` longtext COLLATE utf8mb4_general_ci,
		  `user_id` bigint DEFAULT NULL,
		  KEY `fk_users_credit_card` (`user_id`),
		  CONSTRAINT `fk_users_credit_card` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

func main() {

}
