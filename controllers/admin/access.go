package admin

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AccessController struct {
	BaseController
}

func (con AccessController) Index(c *gin.Context) {
	accessList := []models.Access{}
	//获取module_id=0（即模块后），自关联查询顶级模块下的数据
	models.DB.Where("module_id = ?", 0).Preload("AccessItem").Find(&accessList)

	c.HTML(http.StatusOK, "admin/access/index.html", gin.H{
		"accessList": accessList,
	})
}

func (con AccessController) Add(c *gin.Context) {
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Find(&accessList)
	c.HTML(http.StatusOK, "admin/access/add.html", gin.H{
		"accessList": accessList,
	})
}

func (con AccessController) DoAdd(c *gin.Context) {
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	actionName := strings.Trim(c.PostForm("action_name"), " ")
	accessType, err1 := models.StringToInt(c.PostForm("type"))
	url := c.PostForm("url")
	moduleId, err2 := models.StringToInt(c.PostForm("module_id"))
	sort, err3 := models.StringToInt(c.PostForm("sort"))
	status, err4 := models.StringToInt(c.PostForm("status"))
	description := c.PostForm("description")

	if moduleName == "" {
		con.Error(c, "模块名不可为空", "/admin/access/add")
	}
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "传入参数错误", "/admin/access/add")
		return
	}
	access := models.Access{
		ModuleName:  moduleName,
		ActionName:  actionName,
		Type:        accessType,
		Url:         url,
		ModuleId:    moduleId,
		Sort:        sort,
		Status:      status,
		Description: description,
	}
	err5 := models.DB.Create(&access).Error
	if err5 != nil {
		con.Error(c, "增加权限失败", "/admin/access/add")
		return
	}
	con.Success(c, "增加权限成功", "/admin/access/add")
}

func (con AccessController) Edit(c *gin.Context) {
	id, err1 := models.StringToInt(c.Query("id"))
	if err1 != nil {
		con.Error(c, "id查询失败,参数错误", "/admin/access")
	}
	access := models.Access{Id: id}
	models.DB.Find(&access)

	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Find(&accessList)
	c.HTML(http.StatusOK, "admin/access/edit.html", gin.H{
		"access":     access,
		"accessList": accessList,
	})
}

func (con AccessController) DoEdit(c *gin.Context) {
	id, err1 := models.StringToInt(c.PostForm("id"))
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	actionName := strings.Trim(c.PostForm("action_name"), " ")
	accessType, err2 := models.StringToInt(c.PostForm("type"))
	url := c.PostForm("url")
	moduleId, err3 := models.StringToInt(c.PostForm("module_id"))
	sort, err4 := models.StringToInt(c.PostForm("sort"))
	status, err5 := models.StringToInt(c.PostForm("status"))
	description := c.PostForm("description")

	if moduleName == "" {
		con.Error(c, "模块名不可为空", "/admin/access/edit?id="+models.IntToString(id))
	}
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		con.Error(c, "传入参数错误", "/admin/access")
		return
	}

	access := models.Access{Id: id}
	models.DB.Find(&access)
	access.ModuleName = moduleName
	access.ActionName = actionName
	access.Type = accessType
	access.Url = url
	access.ModuleId = moduleId
	access.Sort = sort
	access.Status = status
	access.Description = description

	err6 := models.DB.Save(&access).Error
	if err6 != nil {
		con.Error(c, "权限修改失败", "/admin/access")
		return
	}
	con.Success(c, "权限修改成功", "/admin/access")
}

func (con AccessController) Delete(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数错误", "/admin/access")
	} else {
		access := models.Access{Id: id}
		models.DB.Find(&access)
		if access.ModuleId == 0 { //顶级模块需额外判断
			accessList := []models.Access{}
			models.DB.Where("module_id=?", access.Id).Find(&accessList)
			if len(accessList) > 0 {
				con.Error(c, "顶级模块下仍有数据，请先处理低级模块的权限", "/admin/access")
			} else {
				models.DB.Delete(&access)
				con.Success(c, "权限删除成功", "/admin/access")
			}
		} else {
			models.DB.Delete(&access)
			con.Success(c, "权限删除成功", "/admin/access")
		}
	}
}
