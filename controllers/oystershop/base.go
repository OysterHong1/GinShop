package oystershop

import (
	"GinShop/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type BaseController struct{}

func (con BaseController) Render(c *gin.Context, tpl string, data map[string]interface{}) {

	//1、获取顶部导航
	topNavList := []models.Nav{}
	if hasTopNavList := models.CacheDb.Get("topNavList", &topNavList); !hasTopNavList {
		/*redis中若没有数据，则从sql数据库之中查,随后再缓存到redis之中*/
		models.DB.Where("status=1 AND position=1").Find(&topNavList)
		models.CacheDb.Set("topNavList", topNavList, 60*60)
		fmt.Println("正在从Mysql读取数据")
	} else {
		fmt.Println("正在从redis读取数据")
	}

	//2、获取分类的数据
	goodsCateList := []models.GoodsCate{}
	//https://gorm.io/zh_CN/docs/preload.html
	if hasCateList := models.CacheDb.Get("goodsCateList", &goodsCateList); !hasCateList {
		models.DB.Where("pid = 0 AND status=1").Order("sort DESC").Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
			return db.Where("goods_cate.status=1").Order("goods_cate.sort DESC")
		}).Find(&goodsCateList)
		models.CacheDb.Set("goodsCateList", goodsCateList, 60*60)
	}

	//3、获取中间导航
	middleNavList := []models.Nav{}
	if hasmiddleNavList := models.CacheDb.Get("middleNavList", &middleNavList); !hasmiddleNavList {
		models.DB.Where("status=1 AND position=2").Find(&middleNavList)
		models.CacheDb.Set("middleNavList", middleNavList, 60*60)
	}

	for i := 0; i < len(middleNavList); i++ {
		relation := strings.ReplaceAll(middleNavList[i].Relation, "，", ",") //21，22,23,24
		relationIds := strings.Split(relation, ",")
		goodsList := []models.Goods{}
		models.DB.Where("id in ?", relationIds).Select("id,title,goods_img,price").Find(&goodsList)
		middleNavList[i].GoodsItems = goodsList
	}

	//获取cookie中储存的用户信息
	user := models.User{}
	isLogin := models.Cookie.Get(c, "userinfo", &user)
	models.Cookie.Get(c, "userinfo", &user)
	var userinfo string
	if isLogin && len(user.Phone) == 11 {
		userinfo = fmt.Sprintf(`<li class="userinfo">
			<a href="#">%v</a>		

			<i class="i"></i>
			<ol>
				<li><a href="#">个人中心</a></li>

				<li><a href="#">喜欢</a></li>

				<li><a href="/pass/loginOut">退出登录</a></li>
			</ol>								
		</li> `, user.Phone)
	} else {
		userinfo = fmt.Sprintf(`<li><a href="/pass/login">登录</a></li>
		<li>|</li>
		<li><a href="/pass/registerStep1">注册</a></li>
		<li>|</li>`)
	}

	renderData := gin.H{
		"topNavList":    topNavList,
		"goodsCateList": goodsCateList,
		"middleNavList": middleNavList,
		"userinfo":      userinfo,
	}
	for key, v := range data {
		renderData[key] = v
	}
	c.HTML(http.StatusOK, tpl, renderData)
}
