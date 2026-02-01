package qrcode2console

import (
	"image"
	"os"
	"testing"
)

/*
 * @Author: Han-Ya-Jun
 * @Date: 2019/4/30 15:04
 * @Name: qrcode2console
 * @Function: 测试文件
 */

func TestNewWithURL_Success(t *testing.T) {
	qr, err := NewWithURL("https://hanyajun.com", false)
	if err != nil {
		t.Fatalf("NewWithURL 失败: %v", err)
	}
	
	if qr == nil {
		t.Fatal("返回的 QRCode2Console 为 nil")
	}
}

func TestNewWithURL_WithImageGeneration(t *testing.T) {
	qr, err := NewWithURL("https://github.com/Han-Ya-Jun", true)
	if err != nil {
		t.Fatalf("NewWithURL 失败: %v", err)
	}
	
	// 检查是否生成了图片文件
	if _, err := os.Stat("qrcode.png"); os.IsNotExist(err) {
		t.Error("未生成二维码图片文件")
	} else {
		// 清理测试文件
		os.Remove("qrcode.png")
	}
	
	if qr == nil {
		t.Fatal("返回的 QRCode2Console 为 nil")
	}
}

func TestNewWithPath_Success(t *testing.T) {
	// 使用项目中的测试图片
	qr, err := NewWithPath("qrcode.jpg")
	if err != nil {
		t.Fatalf("NewWithPath 失败: %v", err)
	}
	
	if qr == nil {
		t.Fatal("返回的 QRCode2Console 为 nil")
	}
}

func TestNewWithPath_FileNotFound(t *testing.T) {
	_, err := NewWithPath("nonexistent.png")
	if err == nil {
		t.Fatal("预期错误但未返回错误")
	}
	
	t.Logf("预期错误: %v", err)
}

func TestOutput(t *testing.T) {
	qr, err := NewWithURL("https://example.com", false)
	if err != nil {
		t.Fatalf("NewWithURL 失败: %v", err)
	}
	
	if err := qr.Output(); err != nil {
		t.Fatalf("Output 失败: %v", err)
	}
}

func TestOutputColored(t *testing.T) {
	qr, err := NewWithURL("https://example.com", false)
	if err != nil {
		t.Fatalf("NewWithURL 失败: %v", err)
	}
	
	if err := qr.OutputColored(); err != nil {
		t.Fatalf("OutputColored 失败: %v", err)
	}
}

func TestToString(t *testing.T) {
	qr, err := NewWithURL("https://example.com", false)
	if err != nil {
		t.Fatalf("NewWithURL 失败: %v", err)
	}
	
	qrString, err := qr.ToString()
	if err != nil {
		t.Fatalf("ToString 失败: %v", err)
	}
	
	if qrString == "" {
		t.Fatal("返回的字符串为空")
	}
	
	// 验证字符串格式
	if len(qrString) < 100 {
		t.Fatal("二维码字符串太短")
	}
	
	t.Logf("生成的二维码字符串长度: %d", len(qrString))
}

func TestSetImage(t *testing.T) {
	qr, err := NewWithURL("https://example.com", false)
	if err != nil {
		t.Fatalf("NewWithURL 失败: %v", err)
	}
	
	// 创建一个简单的测试图片
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	err = qr.SetImage(testImg)
	if err != nil {
		t.Fatalf("SetImage 失败: %v", err)
	}
}

func TestBinarization_Error(t *testing.T) {
	qr := &QRCode2Console{}
	
	// 尝试二值化未初始化的图片
	err := qr.binarization()
	if err == nil {
		t.Fatal("预期错误但未返回错误")
	}
	
	t.Logf("预期错误: %v", err)
}

func TestIntegration(t *testing.T) {
	// 集成测试：完整的流程
	urls := []string{
		"https://github.com",
		"https://google.com",
		"https://openai.com",
	}
	
	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			qr, err := NewWithURL(url, false)
			if err != nil {
				t.Fatalf("NewWithURL 失败: %v", err)
			}
			
			// 测试黑白输出
			if err := qr.Output(); err != nil {
				t.Fatalf("Output 失败: %v", err)
			}
			
			// 测试彩色输出
			if err := qr.OutputColored(); err != nil {
				t.Fatalf("OutputColored 失败: %v", err)
			}
			
			// 测试字符串输出
			qrString, err := qr.ToString()
			if err != nil {
				t.Fatalf("ToString 失败: %v", err)
			}
			
			if qrString == "" {
				t.Fatal("返回的字符串为空")
			}
		})
	}
}

func BenchmarkNewWithURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewWithURL("https://example.com", false)
		if err != nil {
			b.Fatalf("NewWithURL 失败: %v", err)
		}
	}
}

func BenchmarkOutput(b *testing.B) {
	qr, err := NewWithURL("https://example.com", false)
	if err != nil {
		b.Fatalf("NewWithURL 失败: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qr.Output()
	}
}

func BenchmarkToString(b *testing.B) {
	qr, err := NewWithURL("https://example.com", false)
	if err != nil {
		b.Fatalf("NewWithURL 失败: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qr.ToString()
	}
}