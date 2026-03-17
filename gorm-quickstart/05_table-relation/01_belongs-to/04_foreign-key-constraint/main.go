package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
可以通过 constraint 标签配置 OnUpdate、OnDelete 实现外键约束
使用 GORM 进行迁移时会自动创建
*/

type User struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Company struct {
	ID   int
	Name string
}

func main() {
	global.Connect()
	err := global.DB.AutoMigrate(User{}, Company{})
	/*
		CREATE TABLE `users` (
		  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
		  `created_at` datetime(3) DEFAULT NULL,
		  `updated_at` datetime(3) DEFAULT NULL,
		  `deleted_at` datetime(3) DEFAULT NULL,
		  `name` longtext COLLATE utf8mb4_general_ci,
		  `company_id` bigint DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `idx_users_deleted_at` (`deleted_at`),
		  KEY `fk_users_company` (`company_id`),
		  CONSTRAINT `fk_users_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("create table successful")
	}
}
