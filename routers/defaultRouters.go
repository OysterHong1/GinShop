package routers

import (
	"GinShop/controllers/hjh"
	"github.com/gin-gonic/gin"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", hjh.DefaultController{}.Index)
		defaultRouters.GET("/thumbnail1", hjh.DefaultController{}.Thumbnail1)
		defaultRouters.GET("/thumbnail2", hjh.DefaultController{}.Thumbnail2)
		defaultRouters.GET("/qrcode1", hjh.DefaultController{}.Qrcode1)
		defaultRouters.GET("/qrcode2", hjh.DefaultController{}.Qrcode2)
	}
}
