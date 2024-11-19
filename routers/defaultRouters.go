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
	}
}
