package admin

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RoleController struct {
	BaseController
}

func (con RoleController) Index(c *gin.Context) {
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	c.HTML(http.StatusOK, "admin/role/index.html", gin.H{
		"roleList": roleList,
	})
}

func (con RoleController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/role/add.html", gin.H{})
}

func (con RoleController) DoAdd(c *gin.Context) {
	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	if title == "" {
		con.Error(c, "角色标题不能为空", "/admin/role/add")
		return
	}
	role := models.Role{}
	role.Title = title
	role.Description = description
	role.Status = 1
	role.AddTime = (int)(models.GetUnix())

	err := models.DB.Create(&role).Error
	if err != nil {
		con.Error(c, "增加角色失败，请重试", "/admin/role/add")
	} else {
		con.Success(c, "增加角色成功！", "/admin/role/add")
	}

	c.String(http.StatusOK, "DoAdd")
}

func (con RoleController) Edit(c *gin.Context) {
	//查询接收int类型，需要将原string类型进行转换
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		con.Error(c, "增加角色失败，请重试", "/admin/role")
	} else {
		role := models.Role{Id: id}
		models.DB.Find(&role)
		c.HTML(http.StatusOK, "admin/role/edit.html", gin.H{
			"role": role,
		})
	}

}

func (con RoleController) DoEdit(c *gin.Context) {
	id, err1 := models.StringToInt(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/role")
	}

	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")
	if title == "" {
		con.Error(c, "角色标题不能为空", "/admin/role/edit")
		return
	}
	//查询要修改的数据后执行修改
	role := models.Role{Id: id}
	models.DB.Find(&role)
	role.Title = title
	role.Description = description
	err2 := models.DB.Save(&role).Error
	if err2 != nil {
		con.Error(c, "数据保存失败", "/admin/role/edit?id="+models.IntToString(role.Id))
	} else {
		con.Success(c, "数据修改成功", "/admin/role")
	}
}

func (con RoleController) Delete(c *gin.Context) {
	id, err1 := models.StringToInt(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/role")
	} else {
		role := models.Role{Id: id}
		models.DB.Delete(&role)
		con.Success(c, "删除数据成功", "/admin/role")
	}
}

func (con RoleController) Auth(c *gin.Context) {
	roleId, err := models.StringToInt(c.Query("id"))
	if err != nil {
		con.Error(c, "参数传入错误", "/admin/role")
		return
	}

	//获取权限列表
	accessList := []models.Access{}
	models.DB.Where("module_id = ?", 0).Preload("AccessItem").Find(&accessList)

	//获取当前角色已拥有的权限,放到map中
	roleAccess := []models.RoleAccess{}
	models.DB.Where("role_id = ?", roleId).Find(&roleAccess)
	roleAccessMap := make(map[int]int)
	for _, access := range roleAccess {
		roleAccessMap[access.AccessId] = access.AccessId
	}

	//遍历权限数据，若数据在角色权限map中，则使其check=true，进行选中
	for i := 0; i < len(accessList); i++ {
		if _, ok := roleAccessMap[accessList[i].Id]; ok {
			accessList[i].Checked = true
		}
		//二级模块
		for j := 0; j < len(accessList[i].AccessItem); j++ {
			if _, ok := roleAccessMap[accessList[i].AccessItem[j].Id]; ok {
				accessList[i].AccessItem[j].Checked = true
			}
		}
	}

	//将数据传入模板
	c.HTML(http.StatusOK, "admin/role/auth.html", gin.H{
		"roleId":        roleId,
		"roleAccessMap": roleAccessMap,
		"accessList":    accessList,
	})
}

func (con RoleController) DoAuth(c *gin.Context) {
	//获取角色ID
	roleId, err := models.StringToInt(c.PostForm("role_id"))
	if err != nil {
		con.Error(c, "参数传入错误", "/admin/role")
		return
	}
	//获取权限ID
	accessIds := c.PostFormArray("access_node[]")

	//删除当前角色对应的权限
	roleAccess := models.RoleAccess{}
	models.DB.Where("role_id = ?", roleId).Delete(&roleAccess)

	//循环添加权限
	for _, v := range accessIds {
		roleAccess.RoleId = roleId
		accessId, _ := models.StringToInt(v)
		roleAccess.AccessId = accessId
		models.DB.Create(&roleAccess)
	}
	con.Success(c, "权限添加成功", "admin/role/auth")
}
