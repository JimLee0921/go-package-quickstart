package main

import "gorm-quickstart/global"

/*
GORM 通常使用拥有者的主键作为外键的值
对于上面的例子，它是 User 的 ID 字段
为 user 添加 credit card 时，GORM 会将 user 的 ID 字段保存到 credit card 的 UserID 字段
同样也可以使用标签 references 来自定义更改
*/

type User struct {
	ID           int
	MemberNumber string       `gorm:"type:varchar(64);uniqueIndex"`
	CreditCards  []CreditCard `gorm:"foreignKey:UserNumber;references:MemberNumber"`
}
type CreditCard struct {
	Number     string
	UserNumber string `gorm:"type:varchar(64)"`
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, CreditCard{})
	/*
		CREATE TABLE `credit_cards` (
		  `number` longtext COLLATE utf8mb4_general_ci,
		  `user_number` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
		  KEY `fk_users_credit_cards` (`user_number`),
		  CONSTRAINT `fk_users_credit_cards` FOREIGN KEY (`user_number`) REFERENCES `users` (`member_number`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci


		CREATE TABLE `users` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `member_number` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  UNIQUE KEY `idx_users_member_number` (`member_number`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}
func main() {

}
