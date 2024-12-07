package oystershop

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
)

type DefaultController struct {
	BaseController
}

func (con DefaultController) Index(c *gin.Context) {
	//1、获取顶部导航 挪入了base.go

	//2、获取轮播图数据
	focusList := []models.Focus{}
	if hasFocusList := models.CacheDb.Get("focusList", &focusList); !hasFocusList {
		models.DB.Where("status=1 AND focus_type=1").Find(&focusList)
		models.CacheDb.Set("focusList", focusList, 60*60)
	}

	//3、获取分类的数据 挪入了base.go
	//4、获取中间导航 挪入了base.go

	//手机
	phoneList := []models.Goods{}
	if hasPhoneList := models.CacheDb.Get("phoneList", &phoneList); !hasPhoneList {
		phoneList := models.GetGoodsByCategory(1, "best", 8)
		models.CacheDb.Set("phoneList", phoneList, 60*60)
	}

	//配件
	otherList := []models.Goods{}
	if hasOtherList := models.CacheDb.Get("otherList", &otherList); !hasOtherList {
		otherList := models.GetGoodsByCategory(9, "all", 8)
		models.CacheDb.Set("otherList", otherList, 60*60)
	}

	con.Render(c, "oystershop/index/index.html", gin.H{
		"focusList": focusList,
		"phoneList": phoneList,
		"otherList": otherList,
	})
}
