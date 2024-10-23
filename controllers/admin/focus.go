package admin

import (
	"GinShop/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FocusController struct {
	BaseController
}

func (con FocusController) Index(c *gin.Context) {
	focusList := []models.Focus{}
	models.DB.Find(&focusList)
	c.HTML(http.StatusOK, "admin/focus/index.html", gin.H{
		"focusList": focusList,
	})
}

func (con FocusController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/focus/add.html", gin.H{})
}

func (con FocusController) DoAdd(c *gin.Context) {
	title := c.PostForm("title")
	focusType, err1 := models.StringToInt(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort, err2 := models.StringToInt(c.PostForm("sort"))
	status, err3 := models.StringToInt(c.PostForm("status"))

	if err1 != nil || err3 != nil {
		con.Error(c, "非法请求", "/admin/focus/add.html")
	}
	if err2 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/focus/add.html")
	}

	//上传文件
	focusImgSrc, err4 := models.UploadImg(c, "focus_img")
	if err4 != nil {
		fmt.Println(err4)
	}

	focus := models.Focus{
		Title:     title,
		FocusType: focusType,
		FocusImg:  focusImgSrc,
		Link:      link,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}
	err5 := models.DB.Create(&focus).Error
	if err5 != nil {
		con.Error(c, "增加轮播图失败", "/admin/focus/add")
	} else {
		con.Success(c, "增加轮播图成功", "/admin/focus")
	}
}

func (con FocusController) Edit(c *gin.Context) {
	id, err1 := models.StringToInt(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/focus")
		return
	}

	focus := models.Focus{Id: id}
	models.DB.Find(&focus)

	c.HTML(http.StatusOK, "admin/focus/edit.html", gin.H{
		"focus": focus,
	})
}

func (con FocusController) DoEdit(c *gin.Context) {
	id, err1 := models.StringToInt(c.PostForm("id"))
	title := c.PostForm("title")
	focusType, err2 := models.StringToInt(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort, err3 := models.StringToInt(c.PostForm("sort"))
	status, err4 := models.StringToInt(c.PostForm("status"))

	if err1 != nil || err2 != nil || err4 != nil {
		con.Error(c, "非法请求", "/admin/focus")
	}
	if err3 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/focus/edit?id="+models.IntToString(id))
	}

	focusImg, _ := models.UploadImg(c, "focus_img")
	focus := models.Focus{Id: id}
	models.DB.Find(&focus)
	focus.Title = title
	focus.FocusType = focusType
	focus.Link = link
	focus.Sort = sort
	focus.Status = status
	if focusImg != "" {
		focus.FocusImg = focusImg
	}
	err5 := models.DB.Save(&focus).Error
	if err5 != nil {
		con.Error(c, "增加轮播图失败", "/admin/focus/edit?id="+models.IntToString(id))
	} else {
		con.Success(c, "增加轮播图成功", "admin/focus")
	}
}

func (con FocusController) Delete(c *gin.Context) {
	id, err1 := models.StringToInt(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/focus")
	} else {
		focus := models.Focus{Id: id}
		models.DB.Delete(&focus)
		con.Success(c, "删除数据成功", "/admin/focus")
	}
}
