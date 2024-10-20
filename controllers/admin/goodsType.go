package admin

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type GoodsTypeController struct {
	BaseController
}

func (con GoodsTypeController) Index(c *gin.Context) {
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	c.HTML(http.StatusOK, "admin/goodsType/index.html", gin.H{
		"goodsTypeList": goodsTypeList,
	})
}

func (con GoodsTypeController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/goodsType/add.html", gin.H{})
}

func (con GoodsTypeController) DoAdd(c *gin.Context) {
	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	status, err1 := models.StringToInt(c.PostForm("status"))

	if err1 != nil {
		con.Error(c, "传入的参数不正确", "/admin/goodsType/add")
	}

	if title == "" {
		con.Error(c, "标题不能为空", "/admin/goodsType/add")
		return
	}

	goodsType := models.GoodsType{
		Title:       title,
		Description: description,
		Status:      status,
		AddTime:     int(models.GetUnix()),
	}

	err := models.DB.Create(&goodsType).Error
	if err != nil {
		con.Error(c, "增加商品类型失败，请重试", "/admin/goodsType/add")
	} else {
		con.Success(c, "增加商品类型成功！", "/admin/goodsType")
	}
}

func (con GoodsTypeController) Edit(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))

	if err != nil {
		con.Error(c, "传入参数错误", "/admin/goodsType")
	} else {
		goodsType := models.GoodsType{Id: id}
		models.DB.Find(&goodsType)
		c.HTML(http.StatusOK, "admin/goodsType/edit.html", gin.H{
			"goodsType": goodsType,
		})
	}

}

func (con GoodsTypeController) DoEdit(c *gin.Context) {

	id, err1 := models.StringToInt(c.PostForm("id"))

	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	status, err2 := models.StringToInt(c.PostForm("description"))
	if err1 != nil || err2 != nil {
		con.Error(c, "传入的参数不正确", "/admin/goodsType/edit")
	}
	if title == "" {
		con.Error(c, "角色标题不能为空", "/admin/goodsType/edit?id=?"+models.IntToString(id))
	}

	//查询要修改的数据后执行修改
	goodsType := models.GoodsType{Id: id}
	models.DB.Find(&goodsType)
	goodsType.Title = title
	goodsType.Description = description
	goodsType.Status = status

	err3 := models.DB.Save(&goodsType).Error
	if err3 != nil {
		con.Error(c, "数据保存失败", "/admin/goodsType/edit?id="+models.IntToString(id))
	} else {
		con.Success(c, "数据修改成功", "/admin/goodsType")
	}
}

func (con GoodsTypeController) Delete(c *gin.Context) {
	id, err1 := models.StringToInt(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/goodsType")
	} else {
		goodsType := models.GoodsType{Id: id}
		models.DB.Delete(&goodsType)
		con.Success(c, "删除数据成功", "/admin/goodsType")
	}
}
