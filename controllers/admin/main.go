package admin

import (
	"GinShop/models"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type MainController struct{}

func (con MainController) Index(c *gin.Context) {
	//获取userinfo对应的session用于权限判断
	session := sessions.Default(c)
	userinfo := session.Get("userinfo") //空接口类型 要进行类型断言

	userinfoStr, ok := userinfo.(string)
	if ok {
		//获取用户信息
		var userinfoStruct []models.Manager
		json.Unmarshal([]byte(userinfoStr), &userinfoStruct)

		//获取权限列表
		accessList := []models.Access{}
		//子模块排序使用回调函数 详情参考gorm文档
		models.DB.Where("module_id = ?", 0).Preload("AccessItem", func(db *gorm.DB) *gorm.DB {
			return db.Order("access.sort DESC")
		}).Order("sort DESC").Find(&accessList)

		//获取当前角色拥有的权限
		roleAccess := []models.RoleAccess{}
		models.DB.Where("role_id = ?", userinfoStruct[0].RoleId).Find(&roleAccess)
		roleAccessMap := make(map[int]int)
		for _, access := range roleAccess {
			roleAccessMap[access.AccessId] = access.AccessId
		}
		//遍历所有的权限数据，判断当前权限的id是否在角色id的map中
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

		c.HTML(http.StatusOK, "admin/main/index.html", gin.H{
			"username":   userinfoStruct[0].Username,
			"accessList": accessList,
			"is_super":   userinfoStruct[0].IsSuper,
		})
	} else {
		c.Redirect(http.StatusFound, "/admin/login")
	}
}

func (con MainController) Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/main/welcome.html", gin.H{})
}

// 公共方法，用于状态修改
func (con MainController) ChangeStatus(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "传入的参数错误",
		})
		return
	}

	table := c.Query("table")
	field := c.Query("field")
	err1 := models.DB.Exec("update "+table+" set "+field+"=ABS("+field+"-1) where id = ?", id).Error
	if err1 != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "修改失败，请重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改数据成功",
	})
}

// 公共方法，用于数字修改
func (con MainController) ChangeNum(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "传入的参数错误",
		})
		return
	}

	table := c.Query("table")
	field := c.Query("field")
	num := c.Query("num")
	err1 := models.DB.Exec("update "+table+" set "+field+"="+num+" where id = ?", id).Error
	if err1 != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "修改失败，请重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改数据成功",
	})
}
