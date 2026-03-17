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
		log.Println("get mainPath error: ")
		return
	}
	filesPath := filepath.Join(filepath.Dir(mainPath), "files")
	f, err := excelize.OpenFile(filepath.Join(filesPath, "data.xlsx"))
	if err != nil {
		log.Println("open file error: ", err)
		return
	}

	// 1. 为单元格设置测试数据
	sheetName := "Sheet1"
	for row := 2; row <= 73; row++ {
		cell, err := excelize.JoinCellName("E", row)
		if err != nil {
			log.Println("get cell name error: ", err)
			return
		}
		if err = f.SetCellFormula(sheetName, cell, "=RANDBETWEEN(2000, 6000)"); err != nil {
			log.Println("set cell formula error", err)
			return
		}
	}
	// TIP 这里有个问题，如果 Scope 传入的话
	// 2. 设置自定义名称
	if err = f.SetDefinedName(&excelize.DefinedName{
		Name:     "源数据",
		RefersTo: "Sheet1!$A$1:$E$73",
		Comment:  "自定义名称",
		Scope:    sheetName,
	}); err != nil {
		log.Println("set default name error: ", err)
		return
	}

	// 3. 创建按区域以年份为主数据透视表
	if err = f.AddPivotTable(&excelize.PivotTableOptions{
		DataRange:       "源数据",
		PivotTableRange: "Sheet1!$G$2:$L$13",
		Rows:            []excelize.PivotTableField{{Data: "年", DefaultSubtotal: true}, {Data: "月"}},
		Filter:          []excelize.PivotTableField{{Data: "区域"}},
		Columns:         []excelize.PivotTableField{{Data: "类型"}},
		Data:            []excelize.PivotTableField{{Data: "销售额", Name: "累计销售额", Subtotal: "Sum"}},
		RowGrandTotals:  true,
		ColGrandTotals:  true,
		ShowDrill:       true,
		ShowRowHeaders:  true,
		ShowColHeaders:  true,
		ShowLastColumn:  true,
	}); err != nil {
		fmt.Println("add pivot table error: ", err)
		return
	}

	// 3. 创建按区域以月份为主数据透视表
	if err = f.AddPivotTable(&excelize.PivotTableOptions{
		DataRange:           "源数据",
		PivotTableRange:     "Sheet1!$G$18:$T$25",
		Rows:                []excelize.PivotTableField{{Data: "区域", DefaultSubtotal: true}},
		Filter:              []excelize.PivotTableField{{Data: "年"}},
		Columns:             []excelize.PivotTableField{{Data: "月"}, {Data: "类型"}},
		Data:                []excelize.PivotTableField{{Data: "销售额", Name: "累计销售额", Subtotal: "Sum"}},
		RowGrandTotals:      true,
		ColGrandTotals:      true,
		ShowDrill:           true,
		ShowRowHeaders:      true,
		ShowColHeaders:      true,
		ShowLastColumn:      true,
		PivotTableStyleName: "PivotStyleLight21",
	}); err != nil {
		fmt.Println("add pivot table error: ", err)
		return
	}

	// 4. 按月份以区域为维度数据透视表
	if err = f.AddPivotTable(&excelize.PivotTableOptions{
		DataRange:           "源数据",
		PivotTableRange:     "Sheet1!$G$30:$X$36",
		Rows:                []excelize.PivotTableField{{Data: "月", DefaultSubtotal: true}},
		Filter:              []excelize.PivotTableField{{Data: "年"}},
		Columns:             []excelize.PivotTableField{{Data: "区域"}, {Data: "类型"}},
		Data:                []excelize.PivotTableField{{Data: "销售额", Name: "累计销售额", Subtotal: "Sum"}},
		RowGrandTotals:      true,
		ColGrandTotals:      true,
		ShowDrill:           true,
		ShowRowHeaders:      true,
		ShowColHeaders:      true,
		ShowLastColumn:      true,
		PivotTableStyleName: "PivotStyleLight18",
	}); err != nil {
		fmt.Println("add pivot table error: ", err)
		return
	}

	// 5. 设置工作表视图属性
	var zoomScale float64 = 120 // 放大工作表试图比例为 120
	var topLeftCell = "G1"      // 设置打开工作表最左边的行为 G 行
	if err = f.SetSheetView(sheetName, 0, &excelize.ViewOptions{ZoomScale: &zoomScale, TopLeftCell: &topLeftCell}); err != nil {
		log.Println("set sheet view error: ", err)
		return
	}

	// 6. 在指定单元格添加形状
	lineWidth := 1.5
	if err = f.AddShape(sheetName, &excelize.Shape{
		Cell: "N3",
		Type: "rect",
		Fill: excelize.Fill{Color: []string{"#8EB9FF"}},
		Line: excelize.ShapeLine{Color: "#4286F4", Width: &lineWidth},
		Paragraph: []excelize.RichTextRun{
			{
				Text: "数据透视表",
				Font: &excelize.Font{
					Bold:      true,
					Italic:    true,
					Family:    "Times New Roman",
					Size:      36,
					Color:     "#777777",
					Underline: "sng",
				},
			},
		},
		Width:  270,
		Height: 70,
	}); err != nil {
		log.Println("add shape error: ", err)
		return
	}

	// 7. 设置工作簿属性
	if err = f.SetDocProps(&excelize.DocProperties{
		Category:    "category",
		Creator:     "Excelize",
		Description: "文档描述",
		Keywords:    "SpreadSheet",
		Subject:     "各区域销售数据",
		Title:       "销售额报表",
	}); err != nil {
		log.Println("set doc props error: ", err)
		return
	}

	if err = f.SaveAs(filepath.Join(filesPath, "new_data.xlsx")); err != nil {
		log.Println("save new_data error: ", err)
	} else {
		log.Println("save new_data successful!")
	}

	if err = f.Close(); err != nil {
		fmt.Println("close file error: ", err)
		return
	}
}
