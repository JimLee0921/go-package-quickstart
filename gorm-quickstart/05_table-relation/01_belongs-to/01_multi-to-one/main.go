package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
belongs to 会与另一个模型建立了一对一的连接，这种模型的每一个实例都属于另一个模型的一个实例
*/

// User 属于 Company，CompanyID 是外键，每个 user 能且只能被分配给一个 company
// 在 User 对象中，有一个和 Company 一样的 CompanyID
// 默认情况下， CompanyID （gorm 内部通过命名约定进行推断）被隐含地用来在 User 和 Company 之间创建一个外键关系
// 因此必须包含在 User 结构体中才能填充 Company 内部结构体
type User struct {
	gorm.Model
	Name      string
	CompanyID int     // 是数据库里的外键字段（存 ID）
	Company   Company // Go 结构体里的关联对象（存整行数据）
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
		  CONSTRAINT `fk_users_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
	if err != nil {
		fmt.Println(err)
	}
}
