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
	file, err := excelize.OpenFile(filepath.Join(filesPath, "data.xlsx"))
	if err != nil {
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println("close file error: ", err)
		}
	}()

	// 1. 设置工作簿的默认字体为楷体, 在 WPS 中会被默认字体样式覆盖, 但是 Microsoft office excel 中是正常的
	if err = file.SetDefaultFont("KaiTi"); err != nil {
		log.Println("set default font error: ", err)
		return
	}
	sheet1Name, sheet2Name, sheet3Name, sheet4Name := "东部", "西部", "南部", "北部"

	// 2. 为每类咖啡销售额后创建迷你图
	// 三个 string 类型切片用作迷你图参数
	var sparklineLocation, sparklineRange1, sparklineRange2 []string
	for row := 2; row <= 10; row++ {
		// 插入的坐标
		location, err := excelize.JoinCellName("N", row)
		if err != nil {
			log.Println("get location cell error: ", err)
			return
		}
		// 数据起始坐标
		start, err := excelize.JoinCellName("B", row)
		if err != nil {
			log.Println("get start cell error: ", err)
			return
		}
		// 数据结束坐标
		end, err := excelize.JoinCellName("M", row)
		if err != nil {
			log.Println("get start cell error: ", err)
			return
		}
		sparklineLocation = append(sparklineLocation, location)
		sparklineRange1 = append(sparklineRange1, fmt.Sprintf("%s!%s:%s", sheet1Name, start, end))
		sparklineRange2 = append(sparklineRange2, fmt.Sprintf("%s!%s:%s", sheet2Name, start, end))
	}
	// 创建迷你图
	if err = file.AddSparkline(sheet1Name, &excelize.SparklineOptions{
		Location: sparklineLocation,
		Range:    sparklineRange1,
		Markers:  true, // 控制迷你图标记是否显示
	}); err != nil {
		log.Println(sheet1Name, "add spark line error: ", err)
		return
	}
	if err = file.AddSparkline(sheet2Name, &excelize.SparklineOptions{
		Location: sparklineLocation,
		Range:    sparklineRange2,
		Markers:  true,     // 控制迷你图标记是否显示
		Type:     "column", // 改变迷你图类型
		Style:    18,       // 设置迷你图样式配色
	}); err != nil {
		log.Println(sheet2Name, "add spark line error: ", err)
		return
	}

	for col := 2; col <= 7; col++ {
		// 3. 索引转列名(将数据类型为整型的索引转换为列名, 比如 2 就是第二列, 37 就是第三列第七行)
		column, err := excelize.ColumnNumberToName(col)
		if err != nil {
			log.Println("column number to name error: ", err)
			return
		}

		// 4. 设置列的分级显示(这里列传入的是字符串)
		if err = file.SetColOutlineLevel(sheet1Name, column, 1); err != nil {
			log.Println("set column outline level error: ", err)
			return
		}
	}

	// 5. 创建行的分级显示(这里行传入的是整数)
	for row := 7; row <= 9; row++ {
		if err = file.SetRowOutlineLevel(sheet2Name, row, 1); err != nil {
			log.Println("set row outline level error: ", err)
			return
		}
	}

	// 6. 设置页眉页脚
	if err = file.SetHeaderFooter(sheet2Name, &excelize.HeaderFooterOptions{
		OddHeader:  "&C&\"Microsoft YaHei,Bold Italic\"&U&KAB7942咖啡&K000000销售数据统计表",
		OddFooter:  "&C&T",
		EvenHeader: "&C&\"Microsoft YaHei,Bold Italic\"&U&KAB7942咖啡&K000000销售数据统计表",
		EvenFooter: "&C&T",
	}); err != nil {
		log.Println("set header footer error: ", err)
		return
	}

	// 7. 插入分页符
	if err = file.InsertPageBreak(sheet3Name, "G1"); err != nil {
		log.Println("insert page break error: ", err)
		return
	}

	// 8. 保护工作表
	if err = file.ProtectSheet(sheet1Name, &excelize.SheetProtectionOptions{
		Password:      "password", // 设置密码
		EditScenarios: false,
	}); err != nil {
		log.Println("protect sheet error: ", err)
		return
	}

	// 9. 工作表分组把东部和南部两张工作表设置为一组(注意工作表分组中必须有一个active worksheet)
	if err = file.GroupSheets([]string{sheet1Name, sheet4Name}); err != nil {
		log.Println("group sheets error: ", err)
		return
	}

	// 10 取消工作表分组
	if err = file.UngroupSheets(); err != nil {
		log.Println("ungroup sheets error: ", err)
		return
	}

	// 11. 设置工作表可见性, false 表示隐藏工作表
	if err = file.SetSheetVisible(sheet2Name, false); err != nil {
		log.Println("set sheet visible error: ", err)
	}

	if err = file.SaveAs(filepath.Join(filesPath, "new_data.xlsx")); err != nil {
		log.Println("save file error: ", err)
	} else {
		log.Println("save file successful!")
	}

}
