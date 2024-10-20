package admin

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GoodsCateController struct {
	BaseController
}

func (con GoodsCateController) Index(c *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = 0").Preload("GoodsCateItems").Find(&goodsCateList)

	c.HTML(http.StatusOK, "admin/goodsCate/index.html", gin.H{
		"goodsCateList": goodsCateList,
	})
}

func (con GoodsCateController) Add(c *gin.Context) {
	//获取所有顶级分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = 0").Find(&goodsCateList)

	c.HTML(http.StatusOK, "admin/goodsCate/add.html", gin.H{
		"goodsCateList": goodsCateList,
	})
}

func (con GoodsCateController) DoAdd(c *gin.Context) {
	title := c.PostForm("title")
	pid, err1 := models.StringToInt(c.PostForm("pid"))
	link := c.PostForm("link")
	template := c.PostForm("template")
	subTitle := c.PostForm("sub_title")
	keywords := c.PostForm("keywords")
	description := c.PostForm("description")
	sort, err2 := models.StringToInt(c.PostForm("sort"))
	status, err3 := models.StringToInt(c.PostForm("status"))
	if err1 != nil || err3 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsCate/add")
		return
	}
	if err2 != nil {
		con.Error(c, "排序值必须为整数", "/admin/goodsCate/add")
		return
	}

	cateImgDir, _ := models.UploadFile(c, "cate_img")
	goodsCate := models.GoodsCate{
		Title:       title,
		Pid:         pid,
		Link:        link,
		Template:    template,
		Keywords:    keywords,
		SubTitle:    subTitle,
		Description: description,
		Sort:        sort,
		Status:      status,
		CateImg:     cateImgDir,
		AddTime:     int(models.GetUnix()),
	}
	err := models.DB.Create(&goodsCate).Error
	if err != nil {
		con.Error(c, "增加数据失败", "/admin/goodsCate/add")
		return
	}
	con.Success(c, "增加数据成功", "/admin/goodsCate")
}

func (con GoodsCateController) Edit(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数错误", "/admin/goodsCate")
	}

	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)

	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid = 0").Find(&goodsCateList)

	c.HTML(http.StatusOK, "admin/goodsCate/edit.html", gin.H{
		"goodsCateList": goodsCateList,
		"goodsCate":     goodsCate,
	})
}

func (con GoodsCateController) DoEdit(c *gin.Context) {
	id, err1 := models.StringToInt(c.PostForm("id"))
	title := c.PostForm("title")
	pid, err2 := models.StringToInt(c.PostForm("pid"))
	link := c.PostForm("link")
	template := c.PostForm("template")
	subTitle := c.PostForm("sub_title")
	keywords := c.PostForm("keywords")
	description := c.PostForm("description")
	sort, err3 := models.StringToInt(c.PostForm("sort"))
	status, err4 := models.StringToInt(c.PostForm("status"))
	if err1 != nil || err2 != nil || err4 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsCate/add")
		return
	}
	if err3 != nil {
		con.Error(c, "排序值必须为整数", "/admin/goodsCate/add")
		return
	}
	cateImgDir, _ := models.UploadFile(c, "cate_img")

	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)
	goodsCate.Title = title
	goodsCate.Pid = pid
	goodsCate.Link = link
	goodsCate.Template = template
	goodsCate.SubTitle = subTitle
	goodsCate.Keywords = keywords
	goodsCate.Description = description
	goodsCate.Sort = sort
	goodsCate.Status = status

	if cateImgDir != "" {
		goodsCate.CateImg = cateImgDir
	}
	err := models.DB.Save(&goodsCate).Error
	if err != nil {
		con.Error(c, "修改失败", "/admin/goodsCate/edit?id="+models.IntToString(id))
		return
	}
	con.Success(c, "修改成功", "/admin/goodsCate")

}

func (con GoodsCateController) Delete(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数错误", "/admin/goodsCate")
	} else {
		goodsCate := models.GoodsCate{Id: id}
		models.DB.Find(&goodsCate)
		if goodsCate.Pid == 0 { //判断是否为顶级分类
			goodsCateList := []models.GoodsCate{}
			models.DB.Where("pid = ?", goodsCate.Id).Find(&goodsCateList)
			if len(goodsCateList) > 0 {
				con.Error(c, "顶级分类下仍有子分类，请先处理子分类", "/admin/goodsCate")
			} else {
				models.DB.Delete(&goodsCate)
				con.Success(c, "分类删除成功", "/admin/goodsCate")
			}
		} else {
			models.DB.Delete(&goodsCate)
			con.Success(c, "分类删除成功", "/admin/goodsCate")
		}
	}
}
