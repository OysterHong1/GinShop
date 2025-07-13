package api

import "github.com/gin-gonic/gin"

type V2Controller struct{}

func (con V2Controller) Index(c *gin.Context) {
	c.String(200, "我是一个api接口")
}
func (con V2Controller) Userlist(c *gin.Context) {
	c.String(200, "我是一个api接口-Userlist")
}
func (con V2Controller) Plist(c *gin.Context) {
	c.String(200, "我是一个api接口-Plist")
}
