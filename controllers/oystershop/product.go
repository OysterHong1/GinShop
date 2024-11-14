package oystershop

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	BaseController
}

func (con ProductController) Category(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("ID的值是:" + id)

	tpl := "oystershop/product/product_list.html"
	con.Render(c, tpl, gin.H{
		"page": 20,
	})
}
