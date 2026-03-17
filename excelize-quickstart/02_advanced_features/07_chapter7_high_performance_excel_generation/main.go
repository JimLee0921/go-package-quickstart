package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {

	// 主动调用内存回收
	runtime.GC()
	startTime := time.Now()
	_, mainPath, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("get main path error")
		return
	}
	filesPath := filepath.Join(filepath.Dir(mainPath), "files")

	// 1. 创建一个新的 excel 文件
	file := excelize.NewFile()

	// 2. 通过指定单元格创建流式写入器
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		log.Println("create new stream writer error: ", err)
		return
	}

	// 3. 设置列的宽度，需要在按行写入之前设置
	if err = streamWriter.SetColWidth(1, 10, 15); err != nil {
		log.Println("set col width error: ", err)
		return
	}

	// 4. 向一行写入写入数据并设置样式, 第一个参数是起始单元格坐标
	style1, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#DFEBF6"},
			Pattern: 1,
		},
	})
	if err != nil {
		log.Println("create style1 error: ", err)
		return
	}
	if err = streamWriter.SetRow("A1", []any{
		excelize.Cell{
			Value:   "商品订单数据报表",
			StyleID: style1,
		},
	}, excelize.RowOpts{Height: 30, Hidden: false}); err != nil {
		log.Println("set row error: ", err)
		return
	}

	// 5. 流式合并单元格
	if err = streamWriter.MergeCell("A1", "J1"); err != nil {
		log.Println("merge cell error: ", err)
		return
	}

	// 6. 为工作表第二行表头设置数据
	var header []any

	for _, cell := range []string{
		"订单号", "商家 ID", "买家 ID", "商品 ID", "商品单价", "交易件数", "物流公司 ID", "运单编号", "运单状态码", "签收状态码",
	} {
		header = append(header, cell)
	}
	// 设置单元格公式
	header = append(header, excelize.Cell{Formula: "SUM(F3:F1000000)"})
	err = streamWriter.SetRow("A2", header)
	if err != nil {
		log.Println("set row A2 error: ", err)
		return
	}

	// 7. 大规模写入数据, 在 Sheet1 的第三行到第一百万行写入数据
	// 这里加个计时操作对比 py
	writeStart := time.Now()
	for rowID := 3; rowID <= 1000000; rowID++ {
		row := make([]any, 10)
		for colID := 0; colID < 10; colID++ {
			// 每个单元格都是随机数
			row[colID] = rand.Intn(640000)
		}
		// 计算单元格坐标
		cell, err := excelize.CoordinatesToCellName(1, rowID)
		if err != nil {
			log.Println("get cell name error: ", err)
		}
		// 写入每行
		if err = streamWriter.SetRow(cell, row); err != nil {
			log.Println("set row data error: ", err)
			return
		}
	}
	writeCost := time.Since(writeStart)
	log.Printf("stream set row cost: %v\n", writeCost)

	// 8. 流式创建表格
	showRowStripes := false
	if err = streamWriter.AddTable(&excelize.Table{
		Range:             "A2:J1000000",
		Name:              "table",
		StyleName:         "TableStyleMedium2",
		ShowFirstColumn:   true,
		ShowLastColumn:    true,
		ShowRowStripes:    &showRowStripes,
		ShowColumnStripes: true,
	}); err != nil {
		log.Println("add table error: ", err)
		return
	}

	// 数据写完后调用 Flush 关闭流式写入器
	if err = streamWriter.Flush(); err != nil {
		log.Println("flush error: ", err)
		return
	}

	if err = file.SaveAs(filepath.Join(filesPath, "data.xlsx")); err != nil {
		log.Println("save file error: ", err)
	} else {
		log.Println("save file successful!")
	}
	printBenchmarkInfo("generate 10 columns * 100,0000 rows: ", startTime)

	defer func() {
		if err = file.Close(); err != nil {
			log.Println("close file error!")
		}
	}()
}

// Linux / Mac 版本, 里面有些东西 Windows 没有
//
//	func printBenchmarkInfo(info string, startTime time.Time) {
//		var memStats runtime.MemStats
//		var rusage syscall.Rusage
//
//		// 字节转 MB
//		bToMb := func(b uint64) uint64 {
//			return b / 1024 / 1024
//		}
//
//		// Go 运行时内存
//		runtime.ReadMemStats(&memStats)
//
//		// 系统资源（RSS）
//		syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)
//
//		fmt.Printf(
//			"%s\nRSS = %v MB\nAlloc = %v MB\nTotalAlloc = %v MB\nSys = %v MB\nNumGC = %v\nCost = %s\n",
//			info,
//			bToMb(uint64(rusage.Maxrss)),
//			bToMb(memStats.Alloc),
//			bToMb(memStats.TotalAlloc),
//			bToMb(memStats.Sys),
//			memStats.NumGC,
//			time.Since(startTime),
//		)
//	}

// Windows 简单版本
func printBenchmarkInfo(info string, startTime time.Time) {
	var memStats runtime.MemStats

	// 读取 Go 内存信息
	runtime.ReadMemStats(&memStats)

	fmt.Printf(
		"%s\nAlloc = %v MB\nTotalAlloc = %v MB\nSys = %v MB\nNumGC = %v\nCost = %s\n\n",
		info,
		memStats.Alloc/1024/1024,
		memStats.TotalAlloc/1024/1024,
		memStats.Sys/1024/1024,
		memStats.NumGC,
		time.Since(startTime),
	)
}
