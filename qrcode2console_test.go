package qrcode2console

import "testing"

/*
* @Author:hanyajun
* @Date:2019/4/30 15:04
* @Name:qrcode2console
* @Function:test
 */

func TestNewQRCode2ConsoleWithUrl(t *testing.T) {
	qr := NewQRCode2ConsoleWithUrl("https://hanyajun.com", true)
	//qr.Debug()
	qr.Output()
}

func TestNewQRCode2ConsoleWithPath(t *testing.T) {
	qr := NewQRCode2ConsoleWithPath("qrcode.jpg")
	//qr.Debug()
	qr.Output()
}
