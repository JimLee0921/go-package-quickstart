package main

import (
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

	sheetName := "成绩单"
	// 1. 新建 excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Close file error: ", err)
		}
	}()

	// 2. 把默认的 sheet1 名字修改为 sheetName
	_ = f.SetSheetName("Sheet1", sheetName)

	// 3. 准备成绩数据, 为二维数组
	data := [][]interface{}{
		{"考试成绩统计表"},
		{"考试名称：期中考试", nil, nil, nil, "基础科目", nil, nil, "理科科目"},
		{"序号", "学号", "姓名", "班级", "数学", "英语", "语文", "化学", "生物", "物理", "总分"},
		{1, 10001, "学生A", "1班", 93, 80, 89, 86, 57, 77},
		{2, 10002, "学生B", "1班", 65, 72, 91, 75, 66, 90},
		{3, 10003, "学生C", "2班", 92, 99, 89, 90, 79, 69},
		{4, 10004, "学生D", "1班", 72, 69, 71, 82, 75, 83},
		{5, 10005, "学生E", "2班", 81, 93, 59, 76, 66, 90},
		{6, 10006, "学生F", "2班", 92, 90, 87, 88, 92, 70},
	}

	// 4. 把数据插入 excel
	for idx, row := range data {
		// 这里计算 cell 坐标也可以使用 excelize.CoordinatesToCellName(1, idx+1)，本质上是一样的，只是传递第一个参数不一样，这里的 "A" 等价于 CoordinatesToCellName 第一个参数 1
		startCell, err := excelize.JoinCellName("A", idx+1)
		if err != nil {
			log.Println("join cell name error: ", err)
			return
		}
		err = f.SetSheetRow(sheetName, startCell, &row)
		if err != nil {
			log.Println("set sheet row error: ", err)
			return
		}
	}

	// 5. 设置公式计算总分
	// 定义公式类型（STCellFormulaTypeShared为共享公式）和作用范围
	formulaType, ref := excelize.STCellFormulaTypeShared, "K4:K9"
	if err := f.SetCellFormula(sheetName, "K4", "=SUM(E4:J4)", excelize.FormulaOpts{
		Ref:  &ref,
		Type: &formulaType, // 使用 Shared Formula 类型后 excelize 会自动扩展公式
	}); err != nil {
		log.Println("set cell formula error: ", err)
		return
	}

	// 6. 合并单元格
	mergeCellRanges := [][]string{
		{"A1", "K1"},
		{"A2", "D2"},
		{"E2", "G2"},
		{"H2", "J2"},
	}
	for _, item := range mergeCellRanges {
		if err := f.MergeCell(sheetName, item[0], item[1]); err != nil {
			log.Println("merge cell error: ", err)
			return
		}
	}

	// 7. 单元格设置样式和对齐方式等
	// 创建样式，返回该样式 ID 后续应用到指定单元格上
	style1, err := f.NewStyle(&excelize.Style{
		// 对齐方式设置为单元格内容水平居中和垂直居中且自动换行
		Alignment: &excelize.Alignment{
			Horizontal: "center", // 水平对齐方式
			Vertical:   "center", // 垂直对齐方式
			WrapText:   true,     // 文本是否自动换行
		},
		// 设置背景填充
		Fill: excelize.Fill{
			Type:    "pattern",           // 使用图案填充
			Color:   []string{"#DFE8F6"}, // 背景颜色
			Pattern: 1,                   // 填充样式， 0(无填充)、1(纯色填充)、2+(图案填充)
		},
	})
	if err != nil {
		log.Println("new style error: ", err)
		return
	}
	// 把样式应用到指定的单元格上
	if err = f.SetCellStyle(sheetName, "A1", "A1", style1); err != nil {
		log.Println("set cell style error: ", err)
		return
	}

	// 创建一个新的合并居中的样式应用到其它三个合并单元格中
	style2, err := f.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{
		Horizontal: "center", // 水平对齐方式
		Vertical:   "center", // 垂直对齐方式
		WrapText:   true,     // 文本是否自动换行
	}})
	// 应用到其它三个合并单元格中
	if err != nil {
		log.Println("new style error: ", err)
		return
	}
	for _, cell := range []string{"A2", "E2", "H2"} {
		if err = f.SetCellStyle(sheetName, cell, cell, style2); err != nil {
			log.Println("set cell style error: ", err)
			return
		}
	}

	// 8. 设置列的宽度
	err = f.SetColWidth(sheetName, "D", "K", 7)
	if err != nil {
		log.Println("set col width error: ", err)
		return
	}

	// 9. 添加表格, 这里是新版本的 AddTable API, 参数与视频教程不是太一样
	if err = f.AddTable(sheetName, &excelize.Table{
		Name:      "table",            // Excel 表对象名称
		Range:     "A3:k9",            // 表格区域
		StyleName: "TableStyleLight2", // 表格样式名
	}); err != nil {
		log.Println("add table error: ", err)
		return
	}

	// 10. 创建图表, 也给改成新版本的 API 了
	if err = f.AddChart(sheetName, "A10", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			// 这里只有一个总分系列, 也可以把其他分数也弄个系列
			{
				Name:       "成绩单!$A$2",      // 图例项（系列），在图表图例和公式栏中显示
				Categories: "成绩单!$C$4:$C$9", // X 轴数据也就是横坐标标签
				Values:     "成绩单!$K$4:$K$9",
			},
		},
		Format: excelize.GraphicOptions{
			ScaleX:  1.3,
			OffsetX: 10,
			OffsetY: 20,
		},
		Legend: excelize.ChartLegend{
			// Position: "none", // 关闭图例项（其实就是设置位置, none 表示不显示）
		},
		Title: []excelize.RichTextRun{
			{
				Text: "总分成绩表",
				Font: &excelize.Font{
					Bold:  true,
					Color: "#FF00FF",
				},
			},
		},
	}); err != nil {
		log.Println("add chart error: ", err)
		return
	}

	// 11. 添加图片（这里插入的透明的图片，类似于水印效果）
	if err = f.AddPicture(sheetName, "G8", filepath.Join(filesDir, "draft.png"), &excelize.GraphicOptions{
		ScaleX:  0.1, // 图片水平缩放比例
		ScaleY:  0.1, // 图片垂直缩放比例
		OffsetX: 15,  // 图片与插入单元格的水平偏移量，默认值为 0
		OffsetY: 15,  // 图片与插入单元格的垂直偏移量，默认值为 0
	}); err != nil {
		log.Println("add picture error: ", err)
	}

	// 12. 设置工作表试图属性, 不显示网格线
	showGridLines := false // Go 中对于 false 这种字面量不能直接 &false 取址, 所以需要先定义变量再进行取址
	// SetSheetView 第二个视图索引为 0 因为大部分场景下一个 sheet 也就一个视图, 所以传入 0 即可
	if err = f.SetSheetView(sheetName, 0, &excelize.ViewOptions{
		ShowGridLines: &showGridLines, // 工作表是否显示网格线，默认值为 true
	}); err != nil {
		log.Println("set sheet view error: ", err)
		return
	}

	// 12. 设置窗格, 冻结工作表的前几列表头, 真正关键就是 YSplit 和 TopLeftCell
	if err = f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,         // 是否冻结窗格
		Split:       false,        // 是否分割窗口
		XSplit:      0,            // 横向分割位置, 设置为 0 表示不冻结列
		YSplit:      3,            // 纵向分割位置
		TopLeftCell: "A4",         // 滚动区域嘴左上角的单元格
		ActivePane:  "bottomLeft", // 激活窗格
	}); err != nil {
		log.Println("set panes error: ", err)
		return
	}

	// 保存成绩表
	if err = f.SaveAs(filepath.Join(filesDir, "成绩表.xlsx")); err != nil {
		log.Println("save file error: ", err)
		return
	}
	log.Println("save file successful!")
}
