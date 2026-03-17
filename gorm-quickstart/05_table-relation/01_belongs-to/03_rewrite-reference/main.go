package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
对于 belongs to 关系，GORM 通常使用数据库表，主表（拥有者）的主键值作为外键参考
比如 User 和 Company 就是使用了主表 Company 的主键字段 ID 作为外键的参考值
如果设置了 User 实体属于 Company 实体，GORM 就会自动把 Company 中的 ID 属性保存到 User 的 CompanyID 属性中
同样也可以使用 references 标签来修改引用
*/

type User struct {
	gorm.Model
	Name      string
	CompanyID string  `gorm:"type:varchar(64);index"` // 加上索引
	Company   Company `gorm:"references:Code"`        // 使用 Company 的 Code 作为引用
}

type Company struct {
	ID   int
	Code string `gorm:"type:varchar(64);uniqueIndex"` // 必须为唯一索引
	Name string
}

/*
如果外键名恰好在拥有者类型中存在，GORM 通常会错误的认为它是 has one 关系
需要在 belongs to 关系中指定 references
*/

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
		  `company_id` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `idx_users_deleted_at` (`deleted_at`),
		  KEY `idx_users_company_id` (`company_id`),
		  CONSTRAINT `fk_users_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`code`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
	if err != nil {
		fmt.Println(err)
	}
}
