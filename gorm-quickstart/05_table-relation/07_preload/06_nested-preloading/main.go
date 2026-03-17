package main

import (
	"fmt"
	"gorm-quickstart/global"
	"log"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string

	// has-one
	CreditCard CreditCard

	// has-many
	Orders []Order
}

type CreditCard struct {
	gorm.Model
	UserID uint
	Number string
}

type Order struct {
	gorm.Model
	UserID uint
	State  string

	// has-many
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Qty       int

	// belongs-to
	Product Product
}

type Product struct {
	gorm.Model
	Name  string
	Price float64
}

func init() {
	global.Connect()
	if err := global.DB.AutoMigrate(&User{}, &CreditCard{}, &Order{}, &OrderItem{}, &Product{}); err != nil {
		log.Fatal(err)
	}

	// 防止重复插入（内存库可省略，这里给模板）
	var cnt int64
	global.DB.Model(&User{}).Count(&cnt)
	if cnt > 0 {
		return
	}

	// Products
	products := []Product{
		{Name: "Keyboard", Price: 80},
		{Name: "Mouse", Price: 40},
		{Name: "Monitor", Price: 300},
	}
	if err := global.DB.Create(&products).Error; err != nil {
		log.Fatal(err)
	}

	// Users
	users := []User{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	if err := global.DB.Create(&users).Error; err != nil {
		log.Fatal(err)
	}

	// CreditCards (has-one)
	cards := []CreditCard{
		{UserID: users[0].ID, Number: "4111-1111-1111-1111"},
		{UserID: users[1].ID, Number: "5555-5555-5555-4444"},
	}
	if err := global.DB.Create(&cards).Error; err != nil {
		log.Fatal(err)
	}

	// Orders (has-many)
	orders := []Order{
		{UserID: users[0].ID, State: "paid"},      // Alice paid
		{UserID: users[0].ID, State: "pending"},   // Alice pending
		{UserID: users[1].ID, State: "paid"},      // Bob paid
		{UserID: users[1].ID, State: "cancelled"}, // Bob cancelled
	}
	if err := global.DB.Create(&orders).Error; err != nil {
		log.Fatal(err)
	}

	// OrderItems (has-many) + Product (belongs-to)
	items := []OrderItem{
		// Alice paid order -> 2 items
		{OrderID: orders[0].ID, ProductID: products[0].ID, Qty: 1}, // Keyboard
		{OrderID: orders[0].ID, ProductID: products[1].ID, Qty: 2}, // Mouse x2

		// Alice pending order -> 1 item
		{OrderID: orders[1].ID, ProductID: products[2].ID, Qty: 1}, // Monitor

		// Bob paid order -> 1 item
		{OrderID: orders[2].ID, ProductID: products[1].ID, Qty: 1}, // Mouse

		// Bob cancelled order -> 1 item
		{OrderID: orders[3].ID, ProductID: products[0].ID, Qty: 1}, // Keyboard
	}
	if err := global.DB.Create(&items).Error; err != nil {
		log.Fatal(err)
	}
}

func main() {
	// -------------------------
	// 1) 嵌套预加载 + has-one 预加载
	// -------------------------
	var users1 []User
	if err := global.DB.Preload("Orders.OrderItems.Product").
		Preload("CreditCard").
		Find(&users1).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n=== Query #1: Preload Orders.OrderItems.Product + CreditCard ===")
	printUsers(users1)

	// -------------------------
	// 2) 条件性预加载状态为已付款的订单
	//  然后预加载 Orders.OrderItems（仅针对匹配的订单）
	// -------------------------
	var users2 []User
	if err := global.DB.
		Preload("Orders", "state = ?", "paid").
		Preload("Orders.OrderItems").
		Find(&users2).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Query #2: Preload Orders(state=paid) + Orders.OrderItems ===")
	fmt.Println("Note: Unmatched orders (pending/cancelled) are not loaded, so their OrderItems also won't be loaded.")
	printUsers(users2)
}
func printUsers(users []User) {
	for _, u := range users {
		cc := u.CreditCard.Number
		if cc == "" {
			cc = "<none>"
		}

		fmt.Printf("User=%s | CreditCard=%s | Orders=%d\n", u.Name, cc, len(u.Orders))
		for _, o := range u.Orders {
			fmt.Printf("  OrderID=%d State=%s Items=%d\n", o.ID, o.State, len(o.OrderItems))
			for _, it := range o.OrderItems {
				// 在 Query #2 里没有 preload Product，所以 Product.Name 可能是空
				pname := it.Product.Name
				if pname == "" {
					pname = "<not preloaded>"
				}
				fmt.Printf("    ItemID=%d Qty=%d Product=%s\n", it.ID, it.Qty, pname)
			}
		}
	}
}
