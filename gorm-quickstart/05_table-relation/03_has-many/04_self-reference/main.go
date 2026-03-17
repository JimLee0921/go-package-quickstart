package main

import "gorm-quickstart/global"

// 自引用 has many，一个用户可能有多个 Manager，每个 Manager 也是 User

type User struct {
	ID        int
	Name      string
	ManagerID *uint
	Team      []User `gorm:"foreignKey:ManagerID"`
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{})
	/*
		CREATE TABLE `users` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `name` longtext COLLATE utf8mb4_general_ci,
		  `manager_id` bigint DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `fk_users_team` (`manager_id`),
		  CONSTRAINT `fk_users_team` FOREIGN KEY (`manager_id`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}
func main() {

}
