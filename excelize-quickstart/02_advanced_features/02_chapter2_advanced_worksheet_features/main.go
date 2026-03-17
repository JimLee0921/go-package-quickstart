package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/xuri/excelize/v2"
)

func main() {
	sheet1Name := "成绩单"
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("get current file path error")
		return
	}
	filesDir := filepath.Join(filepath.Dir(currentFile), "files")
	if err := os.Chdir(filepath.Dir(currentFile)); err != nil {
		log.Println("change working directory error: ", err)
		return
	}

	f, err := excelize.OpenFile(filepath.Join(filesDir, "成绩表.xlsx"))
	if err != nil {
		log.Println("open file error", err)
		return
	}
	// 1. 创建图表工作表
	if err = f.AddChartSheet("物理成绩统计图", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			// 这里只有一个总分系列, 也可以把其他分数也弄个系列
			{
				Name:       "成绩单!$A$2",      // 图例项（系列），在图表图例和公式栏中显示
				Categories: "成绩单!$C$4:$C$9", // X 轴数据也就是横坐标标签
				Values:     "成绩单!$J$4:$J$9",
			},
		},
		Title: []excelize.RichTextRun{
			{
				Text: "总分成绩表",
			},
		},
	}); err != nil {
		log.Println("add chart sheet error: ", err)
		return
	}

	// 2. 设置 sheet 背景图片
	if err = f.SetSheetBackground(sheet1Name, filepath.Join(filesDir, "watermark.png")); err != nil {
		log.Println("set sheet background image error: ", err)
		return
	}

	// 3. 创建条件样式
	red, err := f.NewConditionalStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#9A0511",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FEC7CE"},
			Pattern: 1,
		},
	})
	if err != nil {
		log.Println("create new conditional style error: ", err)
		return
	}
	green, err := f.NewConditionalStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#00FF00",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#006400"},
			Pattern: 1,
		},
	})
	if err != nil {
		log.Println("create new conditional style error: ", err)
		return
	}

	// 4. 设置条件格式（新版本，老版本使用字符串拼接）
	for _, col := range []string{"E", "F", "G", "H", "I", "J"} {
		if err = f.SetConditionalFormat(sheet1Name, fmt.Sprintf("%s4:%s9", col, col), []excelize.ConditionalFormatOptions{
			{ // 最低分数标为红色
				Type:     "bottom", // 条件类型，bottom 为最小值
				Criteria: "=",      // 设置单元格数据的条件格式运算符
				Value:    "1",      // 与 Criteria 参数一起使用，可以用确定的值作为设置单元格条件格式的条件参数
				Format:   &red,
			}, { // 最高分标为绿色
				Type:     "top",
				Criteria: "=",
				Value:    "1",
				Format:   &green,
			},
		}); err != nil {
			log.Println("set conditional format error: ", err)
			return
		}
	}

	// 5. 为单元格添加批注
	if err = f.AddComment(sheet1Name, excelize.Comment{
		Author: "老师",
		Cell:   "K9",
		Text:   "真实优秀啊你小子",
	}); err != nil {
		log.Println("add comment error: ", err)
		return
	}

	// 6. 添加数据校验（这里设置为班级列只能为这三个值中的一个）
	dvRange := excelize.NewDataValidation(true)
	dvRange.Sqref = "D4:D9"
	if err = dvRange.SetDropList([]string{"1班", "2班", "3班"}); err != nil {
		log.Println("create data validation error: ", err)
		return
	}
	if err = f.AddDataValidation(sheet1Name, dvRange); err != nil {
		log.Println("data validation error: ", err)
		return
	}

	// 保存
	if err = f.Save(); err != nil {
		log.Println("save file error: ", err)
		return
	}
	log.Println("save file successful!")
}
