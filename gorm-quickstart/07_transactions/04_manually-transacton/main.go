package main

import "gorm-quickstart/global"

/*
手动控制 GORM 事务，可以配合
SavePoint()：在当前事务中打一个检查点
RollbackTo()：回滚到这个检查点，但事务本身仍然继续
*/

type User struct {
	ID   uint
	Name string
}

func init() {
	global.Connect()
	_ = global.DB.AutoMigrate(User{})
}

func main() {
	Demo()
}

func Demo() error {
	// 1. 手动开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 2. 确保 panic 也能回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 3. 写入 Alice （一定成功）
	if err := tx.Create(&User{Name: "Alice"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. 打一个保存点
	tx.SavePoint("after_alice")

	// 5. 可选操作（这里可以失败）
	if err := tx.Create(&User{Name: "Bob"}).Error; err == nil {
		// 只回滚 bob，不影响 alice，不返回
		tx.RollbackTo("after_alice")
	}

	// 6. 继续事务
	if err := tx.Create(&User{Name: "Carol"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 所有操作完成提交事务
	return tx.Commit().Error

}
