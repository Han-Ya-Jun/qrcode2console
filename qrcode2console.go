package qrcode2console

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/skip2/go-qrcode"
)

/*
 * @Author: Han-Ya-Jun
 * @Date: 2019/4/30 14:56
 * @Name: qrcode2console
 * @Function: 二维码输出到控制台
 */

const (
	// QRCodeSize 标准二维码尺寸
	QRCodeSize = 256
	// ShrinkQRCodeSize 缩小后的二维码尺寸
	ShrinkQRCodeSize = 35
	// Margin 二维码边距
	Margin = 29
	// Multiple 缩小倍数
	Multiple = 6
)

// QRCode2Console 二维码到控制台转换器
type QRCode2Console struct {
	url             string
	img             image.Image
	path            string
	kind            int
	genImg          bool
	points          [QRCodeSize][QRCodeSize]int
	tmpShrinkPoints [QRCodeSize][ShrinkQRCodeSize]int
	shrinkPoints    [ShrinkQRCodeSize][ShrinkQRCodeSize]int
}

// NewWithURL 从URL创建二维码转换器
func NewWithURL(url string, genImg bool) (*QRCode2Console, error) {
	qr := &QRCode2Console{
		url:    url,
		genImg: genImg,
		kind:   0,
	}
	if err := qr.genQRCode(); err != nil {
		return nil, fmt.Errorf("生成二维码失败: %w", err)
	}
	return qr, nil
}

// NewWithPath 从图片路径创建二维码转换器
func NewWithPath(path string) (*QRCode2Console, error) {
	qr := &QRCode2Console{
		path:   path,
		genImg: false,
		kind:   1,
	}
	if err := qr.genQRCode(); err != nil {
		return nil, fmt.Errorf("加载二维码图片失败: %w", err)
	}
	return qr, nil
}

// genQRCode 生成二维码矩阵
func (qr *QRCode2Console) genQRCode() error {
	var img image.Image

	switch qr.kind {
	case 0: // 从URL生成
		code, err := qrcode.Encode(qr.url, qrcode.Medium, QRCodeSize)
		if err != nil {
			return fmt.Errorf("编码二维码失败: %w", err)
		}
		buf := bytes.NewBuffer(code)
		img, err = png.Decode(buf)
		if err != nil {
			return fmt.Errorf("解码PNG失败: %w", err)
		}
	case 1: // 从图片加载
		file, err := os.Open(qr.path)
		if err != nil {
			return fmt.Errorf("打开文件失败: %w", err)
		}
		defer file.Close()

		dstImg, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("解码图片失败: %w", err)
		}
		img = imaging.Resize(dstImg, QRCodeSize, QRCodeSize, imaging.Lanczos)
	default:
		return fmt.Errorf("不支持的二维码类型: %d", qr.kind)
	}

	qr.img = img

	// 如果需要生成图片文件
	if qr.genImg {
		if err := qr.saveImage("qrcode.png"); err != nil {
			return fmt.Errorf("保存二维码图片失败: %w", err)
		}
	}

	return nil
}

// SetImage 设置自定义图片
func (qr *QRCode2Console) SetImage(img image.Image) error {
	qr.img = img

	if qr.genImg {
		if err := qr.saveImage("qrcode_custom.png"); err != nil {
			return fmt.Errorf("保存自定义二维码图片失败: %w", err)
		}
	}

	return nil
}

// saveImage 保存图片到文件
func (qr *QRCode2Console) saveImage(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, qr.img)
}

// binarization 二维码图片二值化 (0-1)
func (qr *QRCode2Console) binarization() error {
	if qr.img == nil {
		return fmt.Errorf("二维码图片未初始化")
	}

	gray := image.NewGray(image.Rect(0, 0, QRCodeSize, QRCodeSize))
	bounds := qr.img.Bounds()

	for x := bounds.Min.X; x < bounds.Max.X && x < QRCodeSize; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y && y < QRCodeSize; y++ {
			r32, g32, b32, _ := qr.img.At(x, y).RGBA()
			r, g, b := int(r32>>8), int(g32>>8), int(b32>>8)
			
			// 计算灰度值并二值化
			grayValue := (r + g + b) / 3
			if grayValue > 180 {
				qr.points[y][x] = 0
				gray.Set(x, y, color.Gray{uint8(255)})
			} else {
				qr.points[y][x] = 1
				gray.Set(x, y, color.Gray{uint8(0)})
			}
		}
	}

	if qr.genImg {
		if err := saveGrayImage("qrcode_binarization.png", gray); err != nil {
			return fmt.Errorf("保存二值化图片失败: %w", err)
		}
	}

	return nil
}

// saveGrayImage 保存灰度图片
func saveGrayImage(filename string, gray *image.Gray) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, gray)
}

// shrink 缩小二值化数组
func (qr *QRCode2Console) shrink() {
	// 水平方向缩小
	for x := 0; x < QRCodeSize; x++ {
		col := 1
		for y := Margin + 1; y < QRCodeSize-Margin && col < ShrinkQRCodeSize-1; y += Multiple {
			qr.tmpShrinkPoints[x][col] = qr.points[x][y]
			col++
		}
	}

	// 垂直方向缩小
	for y := 1; y < ShrinkQRCodeSize-1; y++ {
		row := 1
		for x := Margin + 1; x < QRCodeSize-Margin && row < ShrinkQRCodeSize-1; x += Multiple {
			qr.shrinkPoints[row][y] = qr.tmpShrinkPoints[x][y]
			row++
		}
	}
}

// Output 控制台输出二维码
func (qr *QRCode2Console) Output() error {
	if err := qr.binarization(); err != nil {
		return fmt.Errorf("二值化失败: %w", err)
	}
	
	qr.shrink()
	
	// 使用Unicode字符增强显示效果
	for x := 0; x < ShrinkQRCodeSize; x++ {
		for y := 0; y < ShrinkQRCodeSize; y++ {
			if qr.shrinkPoints[x][y] == 1 {
				// 黑色方块
				fmt.Print("██")
			} else {
				// 白色方块
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
	
	return nil
}

// OutputColored 彩色控制台输出二维码
func (qr *QRCode2Console) OutputColored() error {
	if err := qr.binarization(); err != nil {
		return fmt.Errorf("二值化失败: %w", err)
	}
	
	qr.shrink()
	
	// 使用ANSI颜色代码
	for x := 0; x < ShrinkQRCodeSize; x++ {
		for y := 0; y < ShrinkQRCodeSize; y++ {
			if qr.shrinkPoints[x][y] == 1 {
				// 黑色背景
				fmt.Print("\033[40m  \033[0m")
			} else {
				// 白色背景
				fmt.Print("\033[47m  \033[0m")
			}
		}
		fmt.Println()
	}
	
	return nil
}

// ToString 将二维码转换为字符串
func (qr *QRCode2Console) ToString() (string, error) {
	if err := qr.binarization(); err != nil {
		return "", fmt.Errorf("二值化失败: %w", err)
	}
	
	qr.shrink()
	
	var result string
	for x := 0; x < ShrinkQRCodeSize; x++ {
		for y := 0; y < ShrinkQRCodeSize; y++ {
			if qr.shrinkPoints[x][y] == 1 {
				result += "██"
			} else {
				result += "  "
			}
		}
		result += "\n"
	}
	
	return result, nil
}