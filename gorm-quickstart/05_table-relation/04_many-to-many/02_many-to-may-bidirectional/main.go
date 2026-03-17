package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
双向 many to many 绑定，更常用，表迁移结构上没有区别，只是可以在代码层面上进行双向查询

单向关系可以用值，因为对象图时树形，直走一个方向
双向关系必须用指针，因为对象图时图结构（存在环）
*/

type User struct {
	ID        int
	Name      string
	Languages []*Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	ID    int
	Name  string
	Users []*User `gorm:"many2many:user_languages;"`
}

// 检索

// GetAllUsers 检索 User 列表并预加载 Language
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Model(&User{}).Preload("Languages").Find(&users).Error
	return users, err
}

// GetAllLanguages 检索 Language 列表并预加载 User
func GetAllLanguages(db *gorm.DB) ([]Language, error) {
	var languages []Language
	err := db.Model(&Language{}).Preload("Users").Find(&languages).Error
	return languages, err
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, Language{})
	/*
		CREATE TABLE `users` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `name` longtext COLLATE utf8mb4_general_ci,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci

		CREATE TABLE `languages` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `name` longtext COLLATE utf8mb4_general_ci,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci

		CREATE TABLE `user_languages` (
		  `language_id` bigint NOT NULL,
		  `user_id` bigint NOT NULL,
		  PRIMARY KEY (`language_id`,`user_id`),
		  KEY `fk_user_languages_user` (`user_id`),
		  CONSTRAINT `fk_user_languages_language` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`),
		  CONSTRAINT `fk_user_languages_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}
func main() {
	users, err := GetAllUsers(global.DB)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, user := range users {
			fmt.Println(user.Name, user.Languages)
		}
	}

	languages, err := GetAllLanguages(global.DB)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, language := range languages {
			fmt.Println(language.Name, language.Users)
		}
	}
}
