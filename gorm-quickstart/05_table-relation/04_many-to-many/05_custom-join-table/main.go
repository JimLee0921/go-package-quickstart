package main

import (
	"fmt"
	"gorm-quickstart/global"
	"time"

	"gorm.io/gorm"
)

/*
可以自定义连接表，支持软删除，钩子函数等功能，并且可以具有更多字段
通过 SetupJoinTable 进行设置
*/

type Person struct {
	ID        int
	Name      string
	Addresses []Address `gorm:"many2many:person_addresses;"`
}
type Address struct {
	ID   int
	Name string
}

// PersonAddress 中间表，PersonID + AddressID 作为主键
type PersonAddress struct {
	PersonID  int `gorm:"primaryKey"`
	AddressID int `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (PersonAddress) BeforeCreate(db *gorm.DB) error {
	// ...
	return nil
}

func init() {
	global.Connect()
	// 先替换 join table
	// 修改 Person 的 Addresses 字段连接表为 PersonAddress
	// PersonAddress 必须定义好所需的外键，否则会报错
	if err := global.DB.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{}); err != nil {
		fmt.Println("SetupJoinTable error:", err)
		return
	}

	// 再迁移（可以把 join table 也手动指定，不指定也会自动迁移进去）

	if err := global.DB.AutoMigrate(&Person{}, &Address{}, &PersonAddress{}); err != nil {
		fmt.Println("AutoMigrate error:", err)
		return
	}
	/*
		CREATE TABLE `people` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `name` longtext COLLATE utf8mb4_general_ci,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci


		CREATE TABLE `addresses` (
		  `id` bigint NOT NULL AUTO_INCREMENT,
		  `name` longtext COLLATE utf8mb4_general_ci,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci


		CREATE TABLE `person_addresses` (
		  `person_id` bigint NOT NULL,
		  `address_id` bigint NOT NULL,
		  `created_at` datetime(3) DEFAULT NULL,
		  `deleted_at` datetime(3) DEFAULT NULL,
		  PRIMARY KEY (`person_id`,`address_id`),
		  KEY `fk_person_addresses_address` (`address_id`),
		  CONSTRAINT `fk_person_addresses_address` FOREIGN KEY (`address_id`) REFERENCES `addresses` (`id`),
		  CONSTRAINT `fk_person_addresses_person` FOREIGN KEY (`person_id`) REFERENCES `people` (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	*/
}

func main() {

}
