package admin

import (
	"GinShop/models"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	//验证MD5是否正确
	fmt.Println("MD5: " + models.Md5("123456"))
	c.HTML(http.StatusOK, "admin/login/login.html", gin.H{})
}

func (con LoginController) DoLogin(c *gin.Context) {
	captchaId := c.PostForm("captchaId")
	verifyValue := c.PostForm("verifyValue")

	username := c.PostForm("username")
	password := c.PostForm("password")
	//验证验证码是否正确
	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag {
		//若验证成功，查询数据库，判断用户与密码是否存在
		userinfoList := []models.Manager{}
		password = models.Md5(password) //密码加密
		models.DB.Where("username=? AND password=?", username, password).Find(&userinfoList)

		if len(userinfoList) > 0 { //若找到了数据
			//执行登录 保存用户信息 跳转
			// cookie/session 前者保存在浏览器，后者保存在服务器 为安全考虑选择session
			session := sessions.Default(c)

			//session.Set没法直接保存结构体对应的切片，故转化为字符串保存
			userinfoSlice, _ := json.Marshal(userinfoList)
			session.Set("userinfo", string(userinfoSlice))
			session.Save()
			con.Success(c, "登录成功", "/admin")
		} else { //若未找到数据
			con.Error(c, "用户名/密码错误", "/admin/login")
		}
	} else {
		con.Success(c, "验证码验证失败", "/admin/login")
	}

}

func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con LoginController) LogOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userinfo")
	session.Save()
	con.Success(c, "退出登录成功", "/admin/login")
}
