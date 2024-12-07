package oystershop

import (
	"GinShop/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type PassController struct {
	BaseController
}

// 获取验证码
func (con PassController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha(50, 120, 2)

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con PassController) Login(c *gin.Context) {
	prevPage := c.Request.Referer()

	c.HTML(http.StatusOK, "oystershop/pass/login.html", gin.H{
		"prevPage": prevPage,
	})

}

func (con PassController) RegisterStep1(c *gin.Context) {
	c.HTML(http.StatusOK, "oystershop/pass/register_step1.html", gin.H{})
}

func (con PassController) RegisterStep2(c *gin.Context) {
	sign := c.Query("sign")
	verifyCode := c.Query("verifyCode")

	//1、验证图形验证码是否正确

	session := sessions.Default(c)
	sessionVerifyCode := session.Get("verifyCode")
	sessionVerifyCodeStr, ok := sessionVerifyCode.(string)
	fmt.Println("---------")
	fmt.Println(verifyCode, sessionVerifyCodeStr)
	if !ok || verifyCode != sessionVerifyCodeStr {
		c.Redirect(302, "/pass/registerStep1")
	}

	//2、获取sign 判断sign是否合法
	userTemp := []models.UserTemp{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.HTML(http.StatusOK, "oystershop/pass/register_step2.html", gin.H{
			"phone":      userTemp[0].Phone,
			"verifyCode": verifyCode,
			"sign":       sign,
		})
	} else {
		c.Redirect(302, "/pass/registerStep1")
	}

}

func (con PassController) RegisterStep3(c *gin.Context) {

	sign := c.Query("sign")
	smsCode := c.Query("smsCode")

	//1、验证短信验证码是否正确
	session := sessions.Default(c)
	sessionSmsCode := session.Get("smsCode")
	sessionSmsCodeStr, ok := sessionSmsCode.(string)

	if !ok || smsCode != sessionSmsCodeStr {
		c.Redirect(302, "/pass/registerStep1")
	}

	//2、获取sign 判断sign是否合法

	userTemp := []models.UserTemp{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.HTML(http.StatusOK, "oystershop/pass/register_step3.html", gin.H{
			"smsCode": smsCode,
			"sign":    sign,
		})
	} else {
		c.Redirect(302, "/pass/registerStep1")
	}

}

func (con PassController) DoRegister(c *gin.Context) {

	//1、获取表单传过来的数据

	sign := c.PostForm("sign")
	smsCode := c.PostForm("smsCode")
	password := c.PostForm("password")
	rpassword := c.PostForm("rpassword")

	//2、验证smsCode是否合法
	session := sessions.Default(c)
	sessionSmsCode := session.Get("smsCode")
	sessionSmsCodeStr, ok := sessionSmsCode.(string)

	if !ok || smsCode != sessionSmsCodeStr {
		c.Redirect(302, "/")
	}
	//3、验证密码是否合法
	if len(password) < 6 || password != rpassword {
		c.Redirect(302, "/")
	}
	//4、验证签名是否合法
	userTemp := []models.UserTemp{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		//4、完成注册
		user := models.User{
			Phone:    userTemp[0].Phone,
			Password: models.Md5(password), //密码要加密
			LastIp:   userTemp[0].Ip,
			AddTime:  int(models.GetUnix()),
			Status:   1,
		}
		models.DB.Create(&user)

		//5、执行登录
		models.Cookie.Set(c, "userinfo", user)

		c.Redirect(302, "/")

	} else {
		c.Redirect(302, "/")
	}

}

func (con PassController) SendCode(c *gin.Context) {

	phone := c.Query("phone")
	verifyCode := c.Query("verifyCode")
	captchaId := c.Query("captchaId")

	if captchaId == "resend" {
		// 1、注册第二个页面发送验证码的时候需要验证图形验证码
		sessionDefault := sessions.Default(c)
		sessionVerifyCode := sessionDefault.Get("verifyCode")
		sessionVerifyCodeStr, ok := sessionVerifyCode.(string)
		if !ok || verifyCode != sessionVerifyCodeStr {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "非法请求",
			})
			return
		}
	} else {
		// 1、验证图形验证码是否正确 保存图形验证码
		if flag := models.VerifyCaptcha(captchaId, verifyCode); !flag {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "验证码输入错误，请重试",
			})
			return
		}
		//保存图形验证码
		sessionDefault := sessions.Default(c)
		sessionDefault.Set("verifyCode", verifyCode)
		sessionDefault.Save()
	}
	/*
		2、判断手机格式是否合法
				pattern := `^[\d]{11}$`
				reg := regexp.MustCompile(pattern)
				reg.MatchString(phone)
	*/
	pattern := `^[\d]{11}$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(phone) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号格式不合法",
		})
		return
	}

	//3、验证手机号是否注册过
	userList := []models.User{}
	models.DB.Where("phone = ?", phone).Find(&userList)
	if len(userList) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号已经注册，请直接登录",
		})
		return
	}
	//4、判断当前ip地址今天发送短信的次数

	ip := c.ClientIP()
	currentDay := models.GetDay() //20211211
	var sendCount int64
	models.DB.Table("user_temp").Where("ip=? AND add_day=?", ip, currentDay).Count(&sendCount)
	if sendCount > 4 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "此ip今天发送短信的次数已经达到上限，请明天再试",
		})
		return
	}
	//5、验证当前手机号今天发送的次数是否合法
	userTemp := []models.UserTemp{}
	smsCode := models.GetRandomNum()
	sign := models.Md5(phone + currentDay) //签名：主要用于页面跳转传值
	models.DB.Where("phone = ? AND add_day=?", phone, currentDay).Find(&userTemp)
	if len(userTemp) > 0 {
		if userTemp[0].SendCount > 2 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "此手机号今天发送短信的次数已经达到上限，请明天再试",
			})
			return
		} else {
			//1、生成短信验证码  发送验证码  调用前面课程的接口
			fmt.Println("----------自己集成发送短信的接口--------")
			fmt.Println(smsCode)
			//2、服务器保存验证码
			session := sessions.Default(c)
			session.Set("smsCode", smsCode)
			session.Save()

			//3、更新发送短信的次数
			oneUserTemp := models.UserTemp{}
			models.DB.Where("id=?", userTemp[0].Id).Find(&oneUserTemp)
			oneUserTemp.SendCount = oneUserTemp.SendCount + 1
			oneUserTemp.AddTime = int(models.GetUnix())

			models.DB.Save(&oneUserTemp)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "发送短信成功",
				"sign":    sign,
			})
			return
		}

	} else {
		//1、生成短信验证码  发送验证码  调用前面课程的接口
		fmt.Println("----------自己集成发送短信的接口--------")
		fmt.Println(smsCode)
		//2、服务器保存验证码
		session := sessions.Default(c)
		session.Set("smsCode", smsCode)
		session.Save()

		//3、记录发送短信的次数

		oneUserTemp := models.UserTemp{
			Ip:        ip,
			Phone:     phone,
			SendCount: 1,
			AddDay:    currentDay,
			AddTime:   int(models.GetUnix()),
			Sign:      sign,
		}
		models.DB.Create(&oneUserTemp)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "发送短信成功",
			"sign":    sign,
		})
		return

	}
}

// 验证验证码
func (con PassController) ValidateSmsCode(c *gin.Context) {

	sign := c.Query("sign")
	smsCode := c.Query("smsCode")
	//1、验证数据是否合法

	userTemp := []models.UserTemp{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "非法请求",
		})
		return
	}

	//2、验证短信验证码是否正确
	sessionDefault := sessions.Default(c)
	sessionSmsCode := sessionDefault.Get("smsCode")
	sessionSmsCodeStr, ok := sessionSmsCode.(string)
	if !ok || smsCode != sessionSmsCodeStr {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "短信验证码输入错误",
		})
		return
	}
	//3、判断验证码有没有过期   15分

	nowTime := models.GetUnix()
	if (nowTime-int64(userTemp[0].AddTime))/1000/60 > 15 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "短信验证码已过期",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "验证码输入正确",
	})
}

func (con PassController) DoLogin(c *gin.Context) {

	phone := strings.Trim(c.PostForm("phone"), " ")
	password := c.PostForm("password")
	captchaId := c.PostForm("captchaId")
	captchaVal := c.PostForm("captchaVal")

	//1、验证图形验证码是否合法
	if flag := models.VerifyCaptcha(captchaId, captchaVal); !flag {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "图形验证码不正确",
		})
		return
	}

	//2、验证用户名密码是否正确
	password = models.Md5(strings.Trim(password, " "))
	userList := []models.User{}
	models.DB.Where("phone = ? AND password = ?", phone, password).Find(&userList)
	if len(userList) > 0 {
		//执行登录
		models.Cookie.Set(c, "userinfo", userList[0])
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "用户登录成功",
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户名或者密码错误",
		})
		return
	}

}

func (con PassController) LoginOut(c *gin.Context) {
	//删除cookie里面的userinfo执行跳转
	models.Cookie.Remove(c, "userinfo")
	prevPage := c.Request.Referer()
	if len(prevPage) > 0 {
		c.Redirect(302, prevPage)
	} else {
		c.Redirect(302, "/")
	}
}
