package models

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

// 默认store创建
//var store = base64Captcha.DefaultMemStore

// RedisStore的配置
var store base64Captcha.Store = RedisStore{}

// 获取验证码
func MakeCaptcha() (id string, b64s string, err error) {
	var driver base64Captcha.Driver
	driverString := base64Captcha.DriverString{
		Height:          40,
		Width:           100,
		ShowLineOptions: 2 | 4,
		NoiseCount:      0,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		//Source:  "洪佳豪文冉苗黄懿卓松童林千雅钟天榕黄嶝奕桂林",
		Length:  2,
		Fonts:   []string{"wqy-microhei.ttc"},
		BgColor: &color.RGBA{R: 3, G: 102, B: 214, A: 125},
	}
	driver = driverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err = c.Generate()
	return id, b64s, err
}

// 验证验证码
func VerifyCaptcha(id string, VerifyValue string) bool {
	fmt.Println(id, VerifyValue)
	if store.Verify(id, VerifyValue, true) {
		return true
	} else {
		return false
	}
}
