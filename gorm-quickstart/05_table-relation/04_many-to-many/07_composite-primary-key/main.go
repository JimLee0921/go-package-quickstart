package main

import "gorm-quickstart/global"

/*
复合主键
如果模型使用了复合主键，GORM 会默认启用复合外键
也可以覆盖默认的外键、指定多个外键，只需用逗号分隔那些键名

使用多个字段进行 `gorm:"primaryKey"` 标签就是符合主键
*/

type Tag struct {
	ID     uint   `gorm:"primaryKey"`
	Locale string `gorm:"primaryKey"`
	Value  string
}

type Blog struct {
	ID         uint   `gorm:"primaryKey"`
	Locale     string `gorm:"primaryKey"`
	Subject    string
	Body       string
	Tags       []Tag `gorm:"many2many:blog_tags;"`
	LocaleTags []Tag `gorm:"many2many:locale_blog_tags;ForeignKey:id,locale;References:id"`
	SharedTags []Tag `gorm:"many2many:shared_blog_tags;ForeignKey:id;References:id"`
}

/*
连接表：blog_tags
  foreign key: blog_id, reference: blogs.id
  foreign key: blog_locale, reference: blogs.locale
  foreign key: tag_id, reference: tags.id
  foreign key: tag_locale, reference: tags.locale

连接表：locale_blog_tags
  foreign key: blog_id, reference: blogs.id
  foreign key: blog_locale, reference: blogs.locale
  foreign key: tag_id, reference: tags.id

连接表：shared_blog_tags
  foreign key: blog_id, reference: blogs.id
  foreign key: tag_id, reference: tags.id
*/

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(Tag{}, Blog{})
}
func main() {

}
