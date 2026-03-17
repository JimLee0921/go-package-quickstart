package main

import (
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {

	f := excelize.NewFile()

	// 创建数据
	f.SetCellValue("Sheet1", "A1", "姓名")
	f.SetCellValue("Sheet1", "B1", "分数")

	f.SetCellValue("Sheet1", "A2", "张三")
	f.SetCellValue("Sheet1", "A3", "李四")
	f.SetCellValue("Sheet1", "A4", "王五")

	f.SetCellValue("Sheet1", "B2", 90)
	f.SetCellValue("Sheet1", "B3", 85)
	f.SetCellValue("Sheet1", "B4", 95)

	// 创建图表 sheet
	err := f.AddChartSheet("统计图", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$B$1",
				Categories: "Sheet1!$A$2:$A$4",
				Values:     "Sheet1!$B$2:$B$4",
			},
		},
		Title: []excelize.RichTextRun{
			{
				Text: "成绩统计",
			},
		},
		Legend: excelize.ChartLegend{
			Position: "right",
		},
		Dimension: excelize.ChartDimension{
			Width:  800,
			Height: 400,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	// 设置默认打开 sheet
	index, _ := f.GetSheetIndex("统计图")
	f.SetActiveSheet(index)

	err = f.SaveAs("chart_sheet_demo.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("生成成功")
}
