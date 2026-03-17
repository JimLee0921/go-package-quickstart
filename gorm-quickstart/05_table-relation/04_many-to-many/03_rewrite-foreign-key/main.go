package main

import (
	"gorm-quickstart/global"

	"gorm.io/gorm"
)

/*
many2many 关系中，连接表会同时拥有两个模型的外键

可以使用标签 foreignKey、references、joinForeignKey、joinReferences 进行重写
也可以只进行部份重写
*/

type User struct {
	gorm.Model
	Profiles []Profile `gorm:"many2many:user_profiles;foreignKey:Refer;joinForeignKey:UserReferID;References:UserRefer;joinReferences:ProfileRefer"`
	Refer    uint      `gorm:"index:,unique"`
}

type Profile struct {
	gorm.Model
	Name      string
	UserRefer uint `gorm:"index:,unique"`
}

/*
会创建连接表：user_profiles

	foreign key: user_refer_id, reference: users.refer
	foreign key: profile_refer, reference: profiles.user_refer
*/
func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{}, Profile{})
}
func main() {

}
