package oystershop

import (
	"GinShop/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	//1、获取顶部导航
	topNavList := []models.Nav{}
	if hasTopNavList := models.CacheDb.Get("topNavList", &topNavList); !hasTopNavList {
		/*redis中若没有数据，则从sql数据库之中查,随后再缓存到redis之中*/
		models.DB.Where("status=1 AND position=1").Find(&topNavList)
		models.CacheDb.Set("topNavList", topNavList, 60*60)
		fmt.Println("正在从Mysql读取数据")
	} else {
		fmt.Println("正在从redis读取数据")
	}

	//2、获取轮播图数据
	focusList := []models.Focus{}
	if hasFocusList := models.CacheDb.Get("focusList", &focusList); !hasFocusList {
		models.DB.Where("status=1 AND focus_type=1").Find(&focusList)
		models.CacheDb.Set("focusList", focusList, 60*60)
	}

	//3、获取分类的数据
	goodsCateList := []models.GoodsCate{}
	//https://gorm.io/zh_CN/docs/preload.html
	if hasCateList := models.CacheDb.Get("goodsCateList", &goodsCateList); !hasCateList {
		models.DB.Where("pid = 0 AND status=1").Order("sort DESC").Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
			return db.Where("goods_cate.status=1").Order("goods_cate.sort DESC")
		}).Find(&goodsCateList)
		models.CacheDb.Set("goodsCateList", goodsCateList, 60*60)
	}

	//4、获取中间导航
	middleNavList := []models.Nav{}
	if hasmiddleNavList := models.CacheDb.Get("middleNavList", &middleNavList); !hasmiddleNavList {
		models.DB.Where("status=1 AND position=2").Find(&middleNavList)
		models.CacheDb.Set("middleNavList", middleNavList, 60*60)
	}

	for i := 0; i < len(middleNavList); i++ {
		relation := strings.ReplaceAll(middleNavList[i].Relation, "，", ",") //21，22,23,24
		relationIds := strings.Split(relation, ",")
		goodsList := []models.Goods{}
		models.DB.Where("id in ?", relationIds).Select("id,title,goods_img,price").Find(&goodsList)
		middleNavList[i].GoodsItems = goodsList
	}

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

	c.HTML(http.StatusOK, "oystershop/index/index.html", gin.H{
		"topNavList":    topNavList,
		"focusList":     focusList,
		"goodsCateList": goodsCateList,
		"middleNavList": middleNavList,
		"phoneList":     phoneList,
		"otherList":     otherList,
	})
}
