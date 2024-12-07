package middlewares

import (
	"GinShop/models"

	"github.com/gin-gonic/gin"
)

func InitUserAuthMiddleware(c *gin.Context) {
	//判断用户有没有登录
	user := models.User{}
	isLogin := models.Cookie.Get(c, "userinfo", &user)
	if !isLogin || len(user.Phone) != 11 {
		c.Redirect(302, "/pass/login")
		return
	}
}
