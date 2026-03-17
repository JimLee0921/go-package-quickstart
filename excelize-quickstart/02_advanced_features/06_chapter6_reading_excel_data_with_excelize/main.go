package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/xuri/excelize/v2"
)

func main() {
	_, mainPath, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("get main path error")
		return
	}
	filesPath := filepath.Join(filepath.Dir(mainPath), "files")

	// 1. 使用密码打开 excel
	file, err := excelize.OpenFile(filepath.Join(filesPath, "成绩表.xlsx"), excelize.Options{
		Password: "password",
	})
	if err != nil {
		log.Println("open file error: ", err)
		return
	}

	// 2. 读取工作表列表
	for sheetIdx, sheetName := range file.GetSheetList() {
		log.Printf("sheetIdx: %d, sheetName: %s\n", sheetIdx, sheetName)
	}

	// 3. 获取活动工作表名称
	activeSheetName := file.GetSheetList()[file.GetActiveSheetIndex()]
	fmt.Println("active sheet name: ", activeSheetName)

	// 4. 读取活动工作表的合并单元格
	mergeCells, err := file.GetMergeCells(activeSheetName)
	if err != nil {
		log.Println("get merge cells error: ", err)
		return
	}

	for _, cell := range mergeCells {
		log.Printf("merge cell range: %s:%s, merge cell value: %s", cell.GetStartAxis(), cell.GetEndAxis(), cell.GetCellValue())
	}

	// 5. 指定工作表通过指定值搜索单元格坐标切片
	searchCells, err := file.SearchSheet(activeSheetName, "90")
	if err != nil {
		log.Println("search sheet error: ", err)
		return
	}
	log.Println(searchCells)

	// 6. 通过正则表达式进行搜索
	searchCells, err = file.SearchSheet(activeSheetName, "^7", true)
	if err != nil {
		log.Println("search sheet by reg error: ", err)
		return
	}
	log.Println(searchCells)

	defer func() {
		if err = file.Close(); err != nil {
			log.Println("close file error: ", err)
		}
	}()

	// 7. 获取每一行然后按行获取每一个坐标的值
	rows, err := file.GetRows(activeSheetName)
	if err != nil {
		log.Println("get rows error: ", err)
		return
	}
	for _, row := range rows {
		for _, cell := range row {
			fmt.Printf("%s\t: ", cell)
		}
		fmt.Println()
	}

	// 8. 获取工作表中的批注(老版 API 是直接获取工作簿中所有工作表批注，这个版本是传入工作表获取指定工作表的批注)
	comments, err := file.GetComments(activeSheetName)
	if err != nil {
		log.Println("get comments error: ", err)
		return
	}
	for _, comment := range comments {
		log.Println("comment: ", comment)
	}

	// 9. 获取列宽
	colAWidth, err := file.GetColWidth(activeSheetName, "A")
	if err != nil {
		log.Println("get col A width error: ", err)
		return
	}
	fmt.Println("col A width: ", colAWidth)

	colGWidth, err := file.GetColWidth(activeSheetName, "G")
	if err != nil {
		log.Println("get col G width error: ", err)
		return
	}
	fmt.Println("col G width: ", colGWidth)

	// 10 获取行高
	row1Height, err := file.GetRowHeight(activeSheetName, 1)
	if err != nil {
		log.Println("get row 1 height error: ", err)
		return
	}
	fmt.Println("row 1 height: ", row1Height)

	row10Height, err := file.GetRowHeight(activeSheetName, 10)
	if err != nil {
		log.Println("get row 10 height error: ", err)
		return
	}
	fmt.Println("row 10 height: ", row10Height)

	// 11. 获取图片 新版本的 Picture 结构体没有图片名
	pictures, err := file.GetPictures(activeSheetName, "G8")
	if err != nil {
		log.Println("get pictures error: ", err)
		return
	}
	for idx, picture := range pictures {
		name := fmt.Sprintf("image%d%s", idx+1, picture.Extension)
		size := len(picture.File)

		fmt.Printf("picture name: %s, size: %d bytes\n", name, size)
	}

	// 12. 获取工作簿默认字体设置
	defaultFont, err := file.GetDefaultFont()
	if err != nil {
		log.Println("get default font error: ", err)
		return
	}
	fmt.Println("default font: ", defaultFont)

	if err = file.SaveAs(filepath.Join(filesPath, "成绩表_new.xlsx")); err != nil {
		log.Println("save file error: ", err)
	} else {
		log.Println("save file successful!")
	}

}
