package routers

import (
	"GinShop/controllers/api"

	"github.com/gin-gonic/gin"
)

func ApiRoutersInit(r *gin.Engine) {
	apiV1Routers := r.Group("/api/v1")
	{
		apiV1Routers.GET("/", api.V1Controller{}.Index)
		apiV1Routers.GET("/navList", api.V1Controller{}.Navlist)
		apiV1Routers.POST("/doLogin", api.V1Controller{}.DoLogin)
		apiV1Routers.PUT("/editArticle", api.V1Controller{}.EditArticle)
		apiV1Routers.DELETE("/deleteNav", api.V1Controller{}.DeleteNav)
	}

	apiV2Routers := r.Group("/api/v2")
	{
		apiV2Routers.GET("/", api.V2Controller{}.Index)
		apiV2Routers.GET("/userlist", api.V2Controller{}.Userlist)
		apiV2Routers.GET("/plist", api.V2Controller{}.Plist)
	}

}
