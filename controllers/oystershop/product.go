package oystershop

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"math"
)

type ProductController struct {
	BaseController
}

func (con ProductController) Category(c *gin.Context) {
	//分类id
	cateId, _ := models.StringToInt(c.Param("id"))
	//当前页
	page, _ := models.StringToInt(c.Query("page"))
	if page == 0 {
		page = 1
	}
	//每一页显示的数量
	pageSize := 3

	//获取当前分类
	currentCate := models.GoodsCate{}
	models.DB.Where("id=?", cateId).Find(&currentCate)
	subCate := []models.GoodsCate{}
	var tempSlice []int
	if currentCate.Pid == 0 {
		//获取二级分类
		models.DB.Where("pid=?", currentCate.Id).Find(&subCate)
		for i := 0; i < len(subCate); i++ {
			tempSlice = append(tempSlice, subCate[i].Id)
		}
	} else {
		//兄弟分类
		models.DB.Where("pid=?", currentCate.Pid).Find(&subCate)
	}
	tempSlice = append(tempSlice, cateId)
	where := "cate_id in ?"
	goodsList := []models.Goods{}
	models.DB.Where(where, tempSlice).Offset((page - 1) * pageSize).Limit(pageSize).Find(&goodsList)

	//获取总数量
	var count int64
	models.DB.Where(where, tempSlice).Table("goods").Count(&count)

	tpl := "oystershop/product/product_list.html"
	con.Render(c, tpl, gin.H{
		"page":        page,
		"goodsList":   goodsList,
		"subCate":     subCate,
		"currentCate": currentCate,
		"totalPages":  math.Ceil(float64(count) / float64(pageSize)),
	})
}
