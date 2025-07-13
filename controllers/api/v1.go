package api

import (
	"GinShop/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type V1Controller struct{}

func (con V1Controller) Index(c *gin.Context) {
	c.String(200, "我是一个api接口")
}
func (con V1Controller) Navlist(c *gin.Context) {
	navList := []models.Nav{}
	models.DB.Find(&navList)
	c.JSON(200, gin.H{
		"navList": navList,
	})
}

//后台通过c.PostForm获取数据
//Content-Type: application/json; 发过来的数据需要通过c.GetRawData() 获取

type UserInfo struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func (con V1Controller) DoLogin(c *gin.Context) {
	var userInfo UserInfo
	b, _ := c.GetRawData() //从 c.Request.Body 读取请求数据
	err := json.Unmarshal(b, &userInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"userInfo": userInfo,
		})
	}
}

type Article struct {
	Title   string `form:"title" json:"title"`
	Content string `form:"content" json:"content"`
}

func (con V1Controller) EditArticle(c *gin.Context) {

	var article Article
	b, _ := c.GetRawData() //从 c.Request.Body 读取请求数据
	err := json.Unmarshal(b, &article)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"article": article,
		})
	}
}

func (con V1Controller) DeleteNav(c *gin.Context) {

	id := c.Query("id")

	c.JSON(200, gin.H{
		"message": "删除数据成功",
		"id":      id,
	})
}
