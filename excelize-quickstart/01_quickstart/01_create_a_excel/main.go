package main

/*
创建一个 Excel 文档
*/

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 1. 创建一个新的文件
	/*
		内部会是自动生成一个 excel 文件结构并且默认包含：
		一个工作表 sheet1
		当前激活表 sheet1
	*/
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 2. 创建一个新的工作表 sheet ()
	index, err := f.NewSheet("sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3. 给工作表插入值
	_ = f.SetCellValue("sheet1", "A2", "Hello World")
	_ = f.SetCellValue("sheet2", "B2", 123)

	// 3. 设置 excel 的默认工作表
	f.SetActiveSheet(index)

	// 4. 设置保存 excel 路径
	if err = f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
