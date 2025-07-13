package main

import (
	"GinShop/models"
	"GinShop/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"html/template"
)

func main() {
	//创建默认路由引擎
	r := gin.Default()

	//配置gin允许跨域请求
	r.Use(cors.Default())

	r.SetFuncMap(template.FuncMap{
		"UnixToTime": models.UnixToTime, //注册模板函数
		"Str2Html":   models.Str2Html,
		"Sub":        models.Sub,
		"Mul":        models.Mul,
		"FormatImg":  models.FormatImg,
		"Substr":     models.Substr,
		"FormatAttr": models.FormatAttr,
	}) //配置自定义模板函数
	//自定义模板函数要放在模板加载之前

	//加载模板
	r.LoadHTMLGlob("templates/**/**/*")
	//配置静态web目录
	r.Static("/static", "./static")

	//session中间件配置
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret111"))
	r.Use(sessions.Sessions("mysession", store))

	routers.AdminRoutersInit(r)
	routers.DefaultRoutersInit(r)
	routers.ApiRoutersInit(r)

	r.Run()
}
