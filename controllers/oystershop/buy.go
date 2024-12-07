package oystershop

import (
	"github.com/gin-gonic/gin"
)

type BuyController struct {
	BaseController
}

func (con BuyController) Checkout(c *gin.Context) {

	con.Render(c, "oystershop/buy/checkout.html", gin.H{})

}
