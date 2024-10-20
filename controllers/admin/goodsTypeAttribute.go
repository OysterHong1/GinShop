package admin

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type GoodsTypeAttributeController struct {
	BaseController
}

func (con GoodsTypeAttributeController) Index(c *gin.Context) {
	cateId, err1 := models.StringToInt(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsType")
		return
	}

	//获取商品类型属性
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	models.DB.Where("cate_id=?", cateId).Find(&goodsTypeAttributeList)
	//获取商品类型属性对应的类型

	goodsType := models.GoodsType{}
	models.DB.Where("id=?", cateId).Find(&goodsType)

	c.HTML(http.StatusOK, "admin/goodsTypeAttribute/index.html", gin.H{
		"cateId":                 cateId,
		"goodsTypeAttributeList": goodsTypeAttributeList,
		"goodsType":              goodsType,
	})
}

func (con GoodsTypeAttributeController) Add(c *gin.Context) {
	//获取当前商品类型属性对应的id
	cateId, err1 := models.StringToInt(c.Query("cate_id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsType")
		return
	}

	//获取所有商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	c.HTML(http.StatusOK, "admin/goodsTypeAttribute/add.html", gin.H{
		"goodsTypeList": goodsTypeList,
		"cateId":        cateId,
	})
}

func (con GoodsTypeAttributeController) DoAdd(c *gin.Context) {
	title := strings.Trim(c.PostForm("title"), " ")
	cateId, err1 := models.StringToInt(c.PostForm("cate_id"))
	attrType, err2 := models.StringToInt(c.PostForm("attr_type"))
	attrValue := c.PostForm("attr_value")
	sort, err3 := models.StringToInt(c.PostForm("sort"))

	if err1 != nil || err2 != nil {
		con.Error(c, "非法请求", "/admin/goodsType")
		return
	}
	if title == "" {
		con.Error(c, "商品类型属性名称不能为空", "/admin/goodsTypeAttribute/add?cate_id="+models.IntToString(cateId))
		return
	}

	if err3 != nil {
		con.Error(c, "排序值不对", "/admin/goodsTypeAttribute/add?cate_id="+models.IntToString(cateId))
		return
	}

	goodsTypeAttr := models.GoodsTypeAttribute{
		Title:     title,
		CateId:    cateId,
		AttrType:  attrType,
		AttrValue: attrValue,
		Status:    1,
		Sort:      sort,
		AddTime:   int(models.GetUnix()),
	}
	err := models.DB.Create(&goodsTypeAttr).Error
	if err != nil {
		con.Error(c, "增加商品类型属性失败 请重试", "/admin/goodsTypeAttribute/add?cate_id="+models.IntToString(cateId))
	} else {
		con.Success(c, "增加商品类型属性成功", "/admin/goodsTypeAttribute?id="+models.IntToString(cateId))
	}
}

func (con GoodsTypeAttributeController) Edit(c *gin.Context) {
	//查询接收int类型，需要将原string类型进行转换
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		con.Error(c, "增加角色失败，请重试", "/admin/goodsType")
		return
	}
	goodsTypeAttribute := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTypeAttribute)

	//获取所有商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	c.HTML(http.StatusOK, "admin/goodsTypeAttribute/edit.html", gin.H{
		"goodsTypeList":      goodsTypeList,
		"goodsTypeAttribute": goodsTypeAttribute,
	})

}

func (con GoodsTypeAttributeController) DoEdit(c *gin.Context) {
	id, err1 := models.StringToInt(c.PostForm("id"))
	title := strings.Trim(c.PostForm("title"), " ")
	cateId, err2 := models.StringToInt(c.PostForm("cate_id"))
	attrType, err3 := models.StringToInt(c.PostForm("attr_type"))
	attrValue := c.PostForm("attr_value")
	sort, err4 := models.StringToInt(c.PostForm("sort"))

	if err1 != nil || err2 != nil || err3 != nil {
		con.Error(c, "非法请求", "/admin/goodsType")
		return
	}
	if title == "" {
		con.Error(c, "商品类型属性名称不能为空", "/admin/goodsTypeAttribute/edit?id="+models.IntToString(id))
		return
	}

	if err4 != nil {
		con.Error(c, "排序值不对", "/admin/goodsTypeAttribute/edit?id="+models.IntToString(id))
		return
	}

	goodsTypeAttr := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTypeAttr)
	goodsTypeAttr.Title = title
	goodsTypeAttr.CateId = cateId
	goodsTypeAttr.AttrType = attrType
	goodsTypeAttr.AttrValue = attrValue
	goodsTypeAttr.Sort = sort

	err := models.DB.Save(&goodsTypeAttr).Error
	if err != nil {
		con.Error(c, "修改商品类型属性失败 请重试", "/admin/goodsTypeAttribute/add?id="+models.IntToString(id))
	} else {
		con.Success(c, "修改商品类型属性成功", "/admin/goodsType")
	}
}

func (con GoodsTypeAttributeController) Delete(c *gin.Context) {
	id, err1 := models.StringToInt(c.Query("id"))
	cateId, err2 := models.StringToInt(c.Query("cate_id"))
	if err1 != nil || err2 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/goodsType")
	} else {
		goodsTypeAttr := models.GoodsTypeAttribute{Id: id}
		models.DB.Delete(&goodsTypeAttr)
		con.Success(c, "删除数据成功", "/admin/goodsTypeAttribute?id="+models.IntToString(cateId))
	}
}
