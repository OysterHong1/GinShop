package oystershop

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	. "github.com/hunterhug/go_image"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	//1、获取顶部导航
	topNavList := []models.Nav{}
	models.DB.Where("status=1 AND position=1").Find(&topNavList)
	//2、获取轮播图数据
	focusList := []models.Focus{}
	models.DB.Where("status=1 AND focus_type=1").Find(&focusList)
	//3、获取分类的数据
	goodsCateList := []models.GoodsCate{}
	//https://gorm.io/zh_CN/docs/preload.html
	models.DB.Where("pid = 0 AND status=1").Order("sort DESC").Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
		return db.Where("goods_cate.status=1").Order("goods_cate.sort DESC")
	}).Find(&goodsCateList)

	//4、获取中间导航
	middleNavList := []models.Nav{}
	models.DB.Where("status=1 AND position=2").Find(&middleNavList)

	for i := 0; i < len(middleNavList); i++ {
		relation := strings.ReplaceAll(middleNavList[i].Relation, "，", ",") //21，22,23,24
		relationIds := strings.Split(relation, ",")
		goodsList := []models.Goods{}
		models.DB.Where("id in ?", relationIds).Select("id,title,goods_img,price").Find(&goodsList)
		middleNavList[i].GoodsItems = goodsList
	}

	//手机

	phoneList := models.GetGoodsByCategory(1, "best", 8)

	//配件

	otherList := models.GetGoodsByCategory(9, "all", 1)

	c.HTML(http.StatusOK, "oystershop/index/index.html", gin.H{
		"topNavList":    topNavList,
		"focusList":     focusList,
		"goodsCateList": goodsCateList,
		"middleNavList": middleNavList,
		"phoneList":     phoneList,
		"otherList":     otherList,
	})

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
