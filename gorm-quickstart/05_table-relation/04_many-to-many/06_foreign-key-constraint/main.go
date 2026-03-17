package main

import "gorm-quickstart/global"

/*
可以通过为标签 constraint 配置 OnUpdate、OnDelete 实现外键约束，在使用 GORM 进行迁移时会被自动创建
*/

type User struct {
	ID        int
	Languages []Language `gorm:"many2many:user_speaks;"`
}

type Language struct {
	Code string `gorm:"primaryKey"`
	Name string
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, Language{})
	/*
		CREATE TABLE `user_speaks` (
		  `user_id` bigint NOT NULL,
		  `language_code` varchar(191) COLLATE utf8mb4_general_ci NOT NULL,
		  PRIMARY KEY (`user_id`,`language_code`),
		  KEY `fk_user_speaks_language` (`language_code`),
		  CONSTRAINT `fk_user_speaks_language` FOREIGN KEY (`language_code`) REFERENCES `languages` (`code`),
		  CONSTRAINT `fk_user_speaks_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}
func main() {

}
