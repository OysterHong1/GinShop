package routers

import (
	"GinShop/controllers/oystershop"
	"GinShop/middlewares"
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
		defaultRouters.POST("/pass/doLogin", oystershop.PassController{}.DoLogin)
		defaultRouters.GET("/pass/loginOut", oystershop.PassController{}.LoginOut)
		//判断用户权限
		defaultRouters.GET("/buy/checkout", middlewares.InitUserAuthMiddleware, oystershop.BuyController{}.Checkout)
		defaultRouters.POST("/buy/doCheckout", middlewares.InitUserAuthMiddleware, oystershop.BuyController{}.DoCheckout)
		defaultRouters.GET("/buy/pay", middlewares.InitUserAuthMiddleware, oystershop.BuyController{}.Pay)
		defaultRouters.GET("/buy/orderPayStatus", middlewares.InitUserAuthMiddleware, oystershop.BuyController{}.OrderPayStatus)

		defaultRouters.POST("/address/addAddress", middlewares.InitUserAuthMiddleware, oystershop.AddressController{}.AddAddress)
		defaultRouters.POST("/address/editAddress", middlewares.InitUserAuthMiddleware, oystershop.AddressController{}.EditAddress)
		defaultRouters.GET("/address/changeDefaultAddress", middlewares.InitUserAuthMiddleware, oystershop.AddressController{}.ChangeDefaultAddress)
		defaultRouters.GET("/address/getOneAddressList", middlewares.InitUserAuthMiddleware, oystershop.AddressController{}.GetOneAddressList)

		defaultRouters.GET("/alipay", middlewares.InitUserAuthMiddleware, oystershop.AlipayController{}.Alipay)
		defaultRouters.POST("/alipayNotify", oystershop.AlipayController{}.AlipayNotify)
		defaultRouters.GET("/alipayReturn", middlewares.InitUserAuthMiddleware, oystershop.AlipayController{}.AlipayReturn)

		defaultRouters.GET("/wxpay", middlewares.InitUserAuthMiddleware, oystershop.WxpayController{}.Wxpay)
		defaultRouters.POST("/wxpay/notify", oystershop.WxpayController{}.WxpayNotify)

		defaultRouters.GET("/user", middlewares.InitUserAuthMiddleware, oystershop.UserController{}.Index)
		defaultRouters.GET("/user/order", middlewares.InitUserAuthMiddleware, oystershop.UserController{}.OrderList)
		defaultRouters.GET("/user/orderinfo", middlewares.InitUserAuthMiddleware, oystershop.UserController{}.OrderInfo)
	}
}
