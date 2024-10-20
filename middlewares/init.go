package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func InitMiddleware(c *gin.Context) {
	//判断用户是否登录
	fmt.Println(time.Now())
	fmt.Println(c.Request.URL)

	c.Set("username", "hjh")

	/*
		中间件使用gotoutine时，不能使用原始的上下文
		必须使用其只读副本(c.Copy())
	*/
	CPc := c.Copy()
	go func() {
		time.Sleep(4 * time.Second)
		fmt.Println("Done! in path" + CPc.Request.URL.Path)
	}()

}
