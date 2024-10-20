package hjh

import (
	"github.com/gin-gonic/gin"
	. "github.com/hunterhug/go_image"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	c.String(200, "首页")

}
func (con DefaultController) Thumbnail1(c *gin.Context) {
	//宽度缩放
	fileName := "static/admin/imageResources/1631938291.png"
	savePath := "static/upload/1631938291_pro.png"
	err := ScaleF2F(fileName, savePath, 600)
	if err != nil {
		c.String(200, "生成图片失败")
		return
	}
	c.String(200, "Thumbnail1 成功")
}

func (con DefaultController) Thumbnail2(c *gin.Context) {
	//按宽度与高度比例缩放
	fileName := "static/admin/imageResources/1631938291.png"
	savePath := "static/upload/1631938291_Thumb400.png"
	//参数不等时进行剪切
	err := ThumbnailF2F(fileName, savePath, 400, 400)
	if err != nil {
		c.String(200, "生成图片失败")
		return
	}
	c.String(200, "Thumbnail2 成功")
}

func (con DefaultController) Qrcode1(c *gin.Context) {
	var png []byte
	png, err := qrcode.Encode("https://www.hjh.com", qrcode.Medium, 256)
	if err != nil {
		c.String(200, "二维码生成失败")
	}
	c.String(200, string(png))
}

func (con DefaultController) Qrcode2(c *gin.Context) {
	savePath := "static/upload/qrcode2.png"
	err := qrcode.WriteFile("https://www.hjh.com", qrcode.Medium, 556, savePath)
	if err != nil {
		c.String(200, "生成二维码失败")
	}
	file, _ := ioutil.ReadFile(savePath)
	c.String(200, string(file))
}
