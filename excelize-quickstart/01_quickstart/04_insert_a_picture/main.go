package main

// Go 的 image.Decode 机制依赖格式注册
import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("files/Book1.xlsx")
	if err != nil {
		fmt.Println("Open file error:", err)
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 1. 插入图片
	if err = f.AddPicture("Sheet1", "A2", "files/image.png", nil); err != nil {
		fmt.Println("add picture error:", err)
		return
	}

	// 2. 在插入图片时设置图片的缩放比例
	_, _ = f.NewSheet("Sheet2")
	if err = f.AddPicture("Sheet2", "A2", "files/image.png", &excelize.GraphicOptions{
		// 长宽都缩放为 0.5 倍
		ScaleX: 0.5,
		ScaleY: 0.5,
	}); err != nil {
		fmt.Println(err)
		return
	}

	// 3. 插入图片并设置图片的打印属性
	enable, disable := true, false
	if err = f.AddPicture("Sheet2", "H2", "files/image.gif", &excelize.GraphicOptions{
		PrintObject:     &enable,  // 打印 excel 时是否打印该图片
		LockAspectRatio: true,     // 是否锁定宽高比例
		OffsetX:         15,       // 单元格中的 X 偏移量
		OffsetY:         10,       // 单元格中的 Y 偏移量
		Locked:          &disable, // 图片是否被锁定（是否能被编辑）
	}); err != nil {
		fmt.Println(err)
		return
	}

	// 保存文件
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}
}
