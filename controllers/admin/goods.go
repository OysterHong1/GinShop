package admin

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GoodsController struct {
	BaseController
}

func (con GoodsController) Index(c *gin.Context) {
	goodsList := []models.Goods{}
	models.DB.Find(&goodsList)
	c.HTML(http.StatusOK, "admin/goods/index.html", gin.H{
		"goodsList": goodsList,
	})
}

func (con GoodsController) Add(c *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = 0").Preload("GoodsCateItems").Find(&goodsCateList)

	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)

	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	c.HTML(http.StatusOK, "admin/goods/add.html", gin.H{
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
	})
}

func (con GoodsController) DoAdd(c *gin.Context) {
	attrIdList := c.PostFormArray("attr_id_list")
	attrValueList := c.PostFormArray("attr_value_list")
	goodsImageList := c.PostFormArray("goods_image_list")
	c.JSON(http.StatusOK, gin.H{
		"attrIdList":     attrIdList,
		"attrValueList":  attrValueList,
		"goodsImageList": goodsImageList,
	})
}

func (con GoodsController) GoodsTypeAttribute(c *gin.Context) {
	cateId, err1 := models.StringToInt(c.Query("cateId"))
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	err2 := models.DB.Where("cate_id = ?", cateId).Find(&goodsTypeAttributeList).Error
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"result":  "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  goodsTypeAttributeList,
		})
	}
}

func (con GoodsController) ImageUpload(c *gin.Context) {
	imgDir, err := models.UploadFile(c, "file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"link": "/" + imgDir,
	})
}
