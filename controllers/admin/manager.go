package admin

import (
	"GinShop/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ManagerController struct {
	BaseController
}

func (con ManagerController) Index(c *gin.Context) {
	managerList := []models.Manager{}
	models.DB.Preload("Role").Find(&managerList)

	c.HTML(http.StatusOK, "admin/manager/index.html", gin.H{
		"managerList": managerList,
	})
}

func (con ManagerController) Add(c *gin.Context) {
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	c.HTML(http.StatusOK, "admin/manager/add.html", gin.H{
		"roleList": roleList,
	})
}

func (con ManagerController) DoAdd(c *gin.Context) {
	roleId, err1 := models.StringToInt(c.PostForm("role_id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/manager")
	}
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")
	if len(username) < 2 || len(password) < 6 {
		con.Error(c, "用户名或者密码的长度不合法", "/admin/manager/add")
		return
	}

	//判断管理员是否存在
	managerList := []models.Manager{}
	models.DB.Where("username=?", username).Find(&managerList)
	if len(managerList) > 0 {
		con.Error(c, "此管理员已存在", "/admin/manager/add")
		return
	}

	manager := models.Manager{
		Username: username,
		Password: models.Md5(password),
		Email:    email,
		Mobile:   mobile,
		RoleId:   roleId,
		Status:   1,
		AddTime:  int(models.GetUnix()),
	}
	err2 := models.DB.Create(&manager).Error
	if err2 != nil {
		con.Error(c, "增加管理员失败", "/admin/manager/add")
		return
	}
	con.Success(c, "增加管理员成功", "/admin/manager/add")
}

func (con ManagerController) Edit(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数错误，请重试", "/admin/manager")
		return
	}

	manager := models.Manager{Id: id}
	models.DB.Find(&manager)

	//获取所有角色
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	c.HTML(http.StatusOK, "admin/manager/edit.html", gin.H{
		"manager":  manager,
		"roleList": roleList,
	})
}

func (con ManagerController) DoEdit(c *gin.Context) {
	id, err1 := models.StringToInt(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/manager")
	}
	roleId, err2 := models.StringToInt(c.PostForm("role_id"))
	if err2 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/manager")
	}

	password := strings.Trim(c.PostForm("password"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	//执行修改
	//找到修改对象
	manager := models.Manager{Id: id}
	models.DB.Find(&manager)

	//检查密码是否为空
	if password != "" {
		if len(password) < 6 {
			con.Error(c, "密码长度不合法", "/admin/manager/edit?id="+models.IntToString(id))
			return
		} else {
			manager.Password = models.Md5(password)
		}
	}

	//检查手机号码长度
	if len(mobile) > 11 {
		con.Error(c, "mobile长度不合法", "/admin/manager/edit?id="+models.IntToString(id))
		return
	} else {
		manager.Mobile = mobile
	}
	manager.Email = email
	manager.RoleId = roleId

	err3 := models.DB.Save(&manager).Error
	if err3 != nil {
		con.Error(c, "修改数据失败", "/admin/manager/edit?id="+models.IntToString(id))
		return
	}
	con.Success(c, "修改数据成功", "/admin/manager/edit?id="+models.IntToString(id))
}

func (con ManagerController) Delete(c *gin.Context) {
	id, err1 := models.StringToInt(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误，请重试", "/admin/manager")
	} else {
		manager := models.Manager{Id: id}
		models.DB.Delete(&manager)
		con.Success(c, "删除成功", "/admin/manager")
	}
}
