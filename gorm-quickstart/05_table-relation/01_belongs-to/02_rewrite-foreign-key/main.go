package main

import (
	"fmt"
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
定义一对多 belongs to 关系，数据库表中必须存在外键
默认情况外键的名字拥有拥有者的类型名称加上表的主键的字段名字
例如 User 实体属于 Company 实体，那么外键名称一般默认使用 CompanyID(到了数据库默认命名规范自动改为 company_id)
可以通过 tag 标签自定义外键名
*/

type User struct {
	gorm.Model
	Name         string
	CompanyRefer int
	Company      Company `gorm:"foreignKey:CompanyRefer"` // 指定使用 CompanyRefer 作为外键代替默认的 CompanyID

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
		  `company_refer` bigint DEFAULT NULL,
		  PRIMARY KEY (`id`),
		  KEY `idx_users_deleted_at` (`deleted_at`),
		  KEY `fk_users_company` (`company_refer`),
		  CONSTRAINT `fk_users_company` FOREIGN KEY (`company_refer`) REFERENCES `companies` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
	if err != nil {
		fmt.Println(err)
	}
}
