package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func main() {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("get current file path error")
		return
	}
	filesDir := filepath.Join(filepath.Dir(currentFile), "files")

	// 1. 读取 csv
	csvFile, err := os.Open(filepath.Join(filesDir, "MSFT.csv"))
	if err != nil {
		log.Println("read csv file error: ", err)
		return
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)

	// 2. 把 csv 中的文件存入到 excel 中
	f := excelize.NewFile()
	sheetName := "Sheet1"
	row := 1
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("read csv file line error", err)
			break
		}
		// 计算单元格坐标
		cell, err := excelize.CoordinatesToCellName(1, row)
		if err != nil {
			log.Println("get cell error: ", err)
			return
		}
		// 第一行表头都是字符串直接写入
		if row == 1 {
			if err = f.SetSheetRow(sheetName, cell, &line); err != nil {
				log.Println("set first row error: ", err)
				break
			}
			row++
			continue
		}
		// 第二行开始把数值转为浮点数后写入
		number, err := convertSlice(line)
		if err != nil {
			log.Println("convert slice error: ", err)
			break
		}
		if err = f.SetSheetRow(sheetName, cell, &number); err != nil {
			log.Println("set numbers error: ", err)
			break
		}
		row++
	}

	// 3. 把首航单元格中的英文字符串改为中文的
	// A1 直接使用 SetCellValue
	if err = f.SetCellValue(sheetName, "A1", "日期"); err != nil {
		log.Println("set cell value error: ", err)
		return
	}

	// A2 使用 SetCellStr
	if err = f.SetCellStr(sheetName, "B1", "开盘价"); err != nil {
		log.Println("set cell str error: ", err)
		return
	}

	// 其它几个单元格使用 SetSheetRow 进行设置
	if err = f.SetSheetRow(sheetName, "C1", &[]string{"最高价", "最低价", "收盘价", "收盘调价", "成交量"}); err != nil {
		log.Println("set sheet first row error: ", err)
		return
	}

	// 4. 设置列的样式
	// 为从 B:F 列的全部单元格设置保留两位小数的数字格式
	style1, err := f.NewStyle(&excelize.Style{NumFmt: 2})
	if err != nil {
		log.Println("create style1 error: ", err)
		return
	}
	if err = f.SetColStyle(sheetName, "B:F", style1); err != nil {
		log.Println("set col style error: ", err)
		return
	}

	// 把 G 列设置为千分位以 , 分割的数字格式
	style2, err := f.NewStyle(&excelize.Style{NumFmt: 3})
	if err != nil {
		log.Println("create style2 error: ", err)
		return
	}
	if err = f.SetColStyle(sheetName, "G", style2); err != nil {
		log.Println("set col style error: ", err)
		return
	}

	// 设置列宽
	if err = f.SetColWidth(sheetName, "A", "G", 11); err != nil {
		log.Println("set col width error: ", err)
		return
	}

	// 隐藏收盘调价列
	if err = f.SetColVisible(sheetName, "F", false); err != nil {
		log.Println("set col visible error: ", err)
		return
	}

	// 5. 插入行并合并单元格后插入数据(使用的新API，视频教程里面的 InsertRow 应该已经移除了)
	if err = f.InsertRows(sheetName, 1, 1); err != nil {
		log.Println("insert rows error: ", err)
		return
	}
	if err = f.MergeCell(sheetName, "A1", "G1"); err != nil {
		log.Println("merge cell error: ", err)
		return
	}

	// 6. 创建富文本格式并设置单元格样式
	if err = f.SetCellRichText(sheetName, "A1", []excelize.RichTextRun{
		{
			Text: "MSFT\r\n",
			Font: &excelize.Font{
				Bold:   true,
				Color:  "2354e8",
				Size:   20,
				Family: "Times New Roman",
			},
		}, {
			Text: "近五年数据",
			Font: &excelize.Font{
				Family: "Microsoft YaHei",
			},
		},
	}); err != nil {
		log.Println("set cell rich text error: ", err)
		return
	}

	style3, err := f.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{WrapText: true, Horizontal: "center", Vertical: "center"}})
	if err = f.SetCellStyle(sheetName, "A1", "A1", style3); err != nil {
		log.Println("set cell style error: ", err)
		return
	}

	// 7. 设置 A1 的行高
	if err = f.SetRowHeight(sheetName, 1, 60); err != nil {
		log.Println("set row height error: ", err)
		return
	}

	// 8. 设置超链接
	if err = f.SetCellValue(sheetName, "H1", "数据来源: finance.yahoo.com"); err != nil {
		log.Println("set hyperlink text error: ", err)
		return
	}
	style4, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color:     "#008000",
			Underline: "single",
		},
	})
	if err = f.SetCellHyperLink(sheetName, "H1", "https://finance.yahoo.com", "External"); err != nil {
		log.Println("set hyperlink error: ", err)
		return
	}
	if err = f.SetCellStyle(sheetName, "H1", "H1", style4); err != nil {
		log.Println("set cell style4: ", err)
		return
	}

	// 9. 创建走势表 sheet 并设置为默认工作表
	newSheet, err := f.NewSheet("走势表")
	if err != nil {
		log.Println("create new sheet error: ", err)
		return
	}
	f.SetActiveSheet(newSheet)

	// 10. 创建收盘价图表
	if err = f.AddChart("走势表", "A1", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$E$2",
				Categories: "Sheet1!$A$3:$A$1189",
				Values:     "Sheet1!$E$3:$E$1189",
				Marker: excelize.ChartMarker{
					Symbol: "none", // 关闭数据标记形状
				},
			},
		},
		Format: excelize.GraphicOptions{
			ScaleX: 1.6,
			ScaleY: 1.5,
		},
		Title: []excelize.RichTextRun{
			{
				Text: "收盘价",
				Font: &excelize.Font{
					Color: "#228B22",
					Size:  30,
					Bold:  true,
				},
			},
		},
		Legend: excelize.ChartLegend{
			Position: "none", // 关闭图例项
		},
		XAxis: excelize.ChartAxis{
			TickLabelSkip: 60,
		},
	}); err != nil {
		log.Println("add chart error: ", err)
		return
	}

	// 11. 创建成交量图表
	if err = f.AddChart("走势表", "A24", &excelize.Chart{
		Type: excelize.Area,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$G$2",
				Categories: "Sheet1!$A$3:$A$1189",
				Values:     "Sheet1!$G$3:$G$1189",
				Marker: excelize.ChartMarker{
					Symbol: "none", // 关闭数据标记形状
				},
			},
		},
		Format: excelize.GraphicOptions{
			ScaleX: 1.6,
			ScaleY: 1.5,
		},
		Title: []excelize.RichTextRun{
			{
				Text: "成交量",
				Font: &excelize.Font{
					Color: "#228B22",
					Size:  30,
					Bold:  true,
				},
			},
		},
		Legend: excelize.ChartLegend{
			Position: "none", // 关闭图例项
		},
		XAxis: excelize.ChartAxis{
			TickLabelSkip: 60,
		},
	}); err != nil {
		log.Println("add chart error: ", err)
		return
	}

	if err = f.SaveAs(filepath.Join(filesDir, "stock.xlsx")); err != nil {
		log.Println("save stock.xlsx error: ", err)
	} else {
		log.Println("save sotck.xlsx successful!")
	}
}

func convertSlice(record []string) (numbers []any, err error) {
	for _, arg := range record {
		if value, err := strconv.ParseFloat(arg, 64); err == nil {
			numbers = append(numbers, value)
			continue
		}
		numbers = append(numbers, arg)
	}
	return
}
