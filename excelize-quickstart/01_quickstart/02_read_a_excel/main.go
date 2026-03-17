package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 1. 打开一个 excel
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 2. 获取制定工作表中指定单元格的内容, GetCellValue 默认取出来都是字符串, 需要自己手动再转化
	cell1, err := f.GetCellValue("sheet1", "B1")
	cell2, err := f.GetCellValue("sheet1", "B2")
	fmt.Printf("cell b1 value: %v type: %T\n", cell1, cell1)
	fmt.Printf("cell b2 value: %v type: %T\n", cell2, cell2)

	// 3. 获取 sheet1 中所有单元格
	rows, err := f.GetRows("sheet1")
	if err != nil {
		fmt.Println(err)
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Printf("value: %v type: %T\n", colCell, colCell)
		}
		fmt.Println()
	}
}
