package main

import (
	"GinShop/models"
	"GinShop/routers"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

func main() {
	// 检查数据库连接
	if err := checkDatabaseConnections(); err != nil {
		log.Fatal("数据库连接检查失败:", err)
		os.Exit(1)
	}

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

	// 配置session中间件 - 使用配置文件中的Redis配置
	store, err := createRedisStore()
	if err != nil {
		log.Printf("Redis session store 创建失败: %v，使用内存存储", err)
		// 如果Redis连接失败，使用cookie存储作为备选方案
		store = cookie.NewStore([]byte("secret111"))
	}
	r.Use(sessions.Sessions("mysession", store))

	routers.AdminRoutersInit(r)
	routers.DefaultRoutersInit(r)
	routers.ApiRoutersInit(r)

	fmt.Println("🚀 服务器启动成功，监听端口 :9000")
	r.Run(":9000")
}

// 检查数据库连接
func checkDatabaseConnections() error {
	fmt.Println("📊 正在检查数据库连接...")

	// 检查MySQL连接
	if models.DB == nil {
		return fmt.Errorf("MySQL 数据库连接失败")
	}

	sqlDB, err := models.DB.DB()
	if err != nil {
		return fmt.Errorf("获取MySQL数据库实例失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("MySQL 数据库连接测试失败: %v", err)
	}

	fmt.Println("✅ MySQL 数据库连接正常")

	// Redis连接检查已在models/redislCache.go中处理
	fmt.Println("✅ 数据库连接检查完成")
	return nil
}

// 创建Redis存储
func createRedisStore() (sessions.Store, error) {
	config, err := ini.Load("./conf/app.ini")
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	redisIP := config.Section("redis").Key("ip").String()
	redisPort := config.Section("redis").Key("port").String()
	redisAddr := redisIP + ":" + redisPort

	fmt.Printf("🔗 连接Redis: %s\n", redisAddr)

	store, err := redis.NewStore(10, "tcp", redisAddr, "", []byte("secret111"))
	if err != nil {
		return nil, fmt.Errorf("创建Redis store失败: %v", err)
	}

	return store, nil
}
