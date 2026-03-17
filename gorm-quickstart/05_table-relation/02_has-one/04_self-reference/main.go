package main

import "gorm-quickstart/global"

/*
一张表中自关联 self-referential association 模型

每个 User 有一个 Manager 上级，而 Manager 也是 User
*/

type User struct {
	ID        int
	Name      string
	ManagerID *uint // 类型是指针，可以为空
	Manager   *User // GORM 判定自关联
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
		  KEY `fk_users_manager` (`manager_id`),
		  CONSTRAINT `fk_users_manager` FOREIGN KEY (`manager_id`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

func main() {

}
