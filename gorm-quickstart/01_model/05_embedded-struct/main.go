package main

/*
嵌入结构体

对于匿名字段，GORM 会将其字段包含在父结构体中
*/

type Author struct {
	Name  string
	Email string
}

type BlogOne struct {
	Author
	ID      int
	Upvotes int32
}

// 上面两个等价于：
/*
type BlogOne struct {
	ID      int
	Upvotes int32
	Email   string
	Upvotes int32
}
*/

/*
对于正常的结构体结构也可以使用标签 embedded 将其嵌入
*/

type BlogTwo struct {
	ID      int
	Author  Author `gorm:"embedded"`
	Upvotes int32
}

// 等价于：
/*
type BlogOne struct {
	ID      int
	Name    string
	Email   string
	Upvotes int32
}
*/

/*
还可以使用标签 embeddedPrefix 来为 db 中的字段命添加前缀
*/

type BlogThree struct {
	ID      int
	Author  Author `gorm:"embedded;embeddedPrefix:author_"`
	Upvotes int32
}

// 等价于：
/*
type BlogThree struct {
	ID      	  int
	AuthorName    string
	AuthorEmail   string
	Upvotes       int32
}
*/

func main() {

}
