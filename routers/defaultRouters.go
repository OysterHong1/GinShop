package routers

import (
	"GinShop/controllers/oystershop"
	"github.com/gin-gonic/gin"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", oystershop.DefaultController{}.Index)
		defaultRouters.GET("/category:id", oystershop.ProductController{}.Category)
		defaultRouters.GET("/detail", oystershop.ProductController{}.Detail)
		defaultRouters.GET("/product/getImgList", oystershop.ProductController{}.GetImgList)

		defaultRouters.GET("/cart", oystershop.CartController{}.Get)
		defaultRouters.GET("/cart/addCart", oystershop.CartController{}.AddCart)

		defaultRouters.GET("/cart/successTip", oystershop.CartController{}.AddCartSuccess)

		defaultRouters.GET("/cart/decCart", oystershop.CartController{}.DecCart)
		defaultRouters.GET("/cart/incCart", oystershop.CartController{}.IncCart)
		defaultRouters.GET("/cart/changeOneCart", oystershop.CartController{}.ChangeOneCart)
		defaultRouters.GET("/cart/changeAllCart", oystershop.CartController{}.ChangeAllCart)
		defaultRouters.GET("/cart/delCart", oystershop.CartController{}.DelCart)

		defaultRouters.GET("/pass/login", oystershop.PassController{}.Login)
		defaultRouters.GET("/pass/captcha", oystershop.PassController{}.Captcha)

		defaultRouters.GET("/pass/registerStep1", oystershop.PassController{}.RegisterStep1)
		defaultRouters.GET("/pass/registerStep2", oystershop.PassController{}.RegisterStep2)
		defaultRouters.GET("/pass/registerStep3", oystershop.PassController{}.RegisterStep3)
		defaultRouters.GET("/pass/sendCode", oystershop.PassController{}.SendCode)
		defaultRouters.GET("/pass/validateSmsCode", oystershop.PassController{}.ValidateSmsCode)
		defaultRouters.POST("/pass/doRegister", oystershop.PassController{}.DoRegister)
	}
}
