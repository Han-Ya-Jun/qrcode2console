package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/Han-Ya-Jun/qrcode2console"
)

func main() {
	url := flag.String("url", "", "要生成二维码的URL")
	path := flag.String("path", "", "二维码图片路径")
	colored := flag.Bool("color", false, "使用彩色输出")
	output := flag.String("output", "", "输出到文件（可选）")
	genImg := flag.Bool("gen-img", false, "生成二维码图片文件")
	flag.Parse()
	
	var qr *qrcode2console.QRCode2Console
	var err error
	
	if *url != "" {
		qr, err = qrcode2console.NewWithURL(*url, *genImg)
	} else if *path != "" {
		qr, err = qrcode2console.NewWithPath(*path)
	} else {
		fmt.Println("错误: 请提供 -url 或 -path 参数")
		fmt.Println("用法:")
		flag.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  qrcode2console -url https://example.com")
		fmt.Println("  qrcode2console -path qrcode.png -color")
		fmt.Println("  qrcode2console -url https://github.com -gen-img -output qr.txt")
		os.Exit(1)
	}
	
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
	
	// 如果指定输出文件
	if *output != "" {
		qrString, err := qr.ToString()
		if err != nil {
			fmt.Printf("错误: 转换二维码失败: %v\n", err)
			os.Exit(1)
		}
		
		file, err := os.Create(*output)
		if err != nil {
			fmt.Printf("错误: 创建文件失败: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		
		_, err = file.WriteString(qrString)
		if err != nil {
			fmt.Printf("错误: 写入文件失败: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("二维码已保存到: %s\n", *output)
		return
	}
	
	// 控制台输出
	if *colored {
		if err := qr.OutputColored(); err != nil {
			fmt.Printf("错误: 输出二维码失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := qr.Output(); err != nil {
			fmt.Printf("错误: 输出二维码失败: %v\n", err)
			os.Exit(1)
		}
	}
}