package middlewares

import (
	"GinShop/models"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"net/http"
	"os"
	"strings"
)

// 权限判断中间件
func InitAdminAuthMiddleware(c *gin.Context) {
	//获取url，取得访问的地址
	pathname := strings.Split(c.Request.URL.String(), "?")[0]
	fmt.Println("当前地址:" + pathname)

	//获取session中保存的用户信息
	session := sessions.Default(c)
	userinfo := session.Get("userinfo") //空接口类型 要进行类型断言

	//判断userinfo是否是一个string
	userinfoStr, ok := userinfo.(string)
	if ok { //判断session中的用户信息是否存在，若不存在则跳转回登录页面
		//若userinfo存在，需要获取其中信息
		var userinfoStruct []models.Manager
		err := json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		if err != nil || !(len(userinfoStruct) > 0 && userinfoStruct[0].Username != "") {
			if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
				c.Redirect(302, "/admin/login")
			}
		} else { //用户登录成功后需要做权限判断(超级管理员无需判断)
			urlPath := strings.Replace(pathname, "/admin/", "", 1)
			if userinfoStruct[0].IsSuper == 0 && !excludeAuthPath("/"+urlPath) {
				//获取权限列表
				accessList := []models.Access{}
				models.DB.Where("module_id = ?", 0).Preload("AccessItem").Find(&accessList)

				//获取当前角色已拥有的权限,放到map中
				roleAccess := []models.RoleAccess{}
				models.DB.Where("role_id = ?", userinfoStruct[0].RoleId).Find(&roleAccess)
				roleAccessMap := make(map[int]int)
				for _, access := range roleAccess {
					roleAccessMap[access.AccessId] = access.AccessId
				}

				//获取当前访问的url对应的权限id，判断权限id是否在角色对应权限中
				access := models.Access{}
				models.DB.Where("url = ?", urlPath).Find(&access)

				if _, ok := roleAccessMap[access.Id]; !ok {
					c.String(http.StatusOK, "您没有权限访问，访问中止")
					c.Abort()
				}
			}
		}
	} else {
		if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
			c.Redirect(302, "/admin/login")
		}
	}
}

// 排除权限判断的方法
func excludeAuthPath(urlpath string) bool {
	//加载app.ini配置文件
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Println(iniErr)
		os.Exit(1)
	}
	excludeAuthPath := config.Section("").Key("excludeAuthPath").String()
	excludeAuthPathSlice := strings.Split(excludeAuthPath, ",")
	for _, v := range excludeAuthPathSlice {
		if v == urlpath {
			return true
		}
	}
	return false
}
