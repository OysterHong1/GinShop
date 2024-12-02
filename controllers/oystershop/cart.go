package oystershop

import (
	"GinShop/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	BaseController
}

func (con CartController) Get(c *gin.Context) {
	//获取购物车数据 显示购物车数据
	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)

	var allPrice float64

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	var tpl = "oystershop/cart/cart.html"
	con.Render(c, tpl, gin.H{
		"cartList": cartList,
		"allPrice": allPrice,
	})

}

func (con CartController) AddCart(c *gin.Context) {

	/*
	   购物车数据保持到哪里？：

	           1、购物车数据保存在本地    （cookie）

	           2、购物车数据保存到服务器(mysql)   （必须登录）

	           3、没有登录 购物车数据保存到本地 ， 登录成功后购物车数据保存到服务器  （用的最多）


	   增加购物车的实现逻辑：

	           1、获取增加购物车的数据  （把哪一个商品加入到购物车）

	           2、判断购物车有没有数据   （cookie）

	           3、如果购物车没有任何数据  直接把当前数据写入cookie

	           4、如果购物车有数据

	              1、判断购物车有没有当前数据

	                       有当前数据让当前数据的数量加1，然后写入到cookie

	              2、如果没有当前数据直接写入cookie
	*/

	// 1、获取增加购物车的数据,放在结构体里面  （把哪一个商品加入到购物车）
	colorId, _ := models.StringToInt(c.Query("color_id"))
	goodsId, err := models.StringToInt(c.Query("goods_id"))
	if err != nil {
		c.Redirect(302, "/")
	}

	goods := models.Goods{}
	goodsColor := models.GoodsColor{}
	models.DB.Where("id=?", goodsId).Find(&goods)
	models.DB.Where("id=?", colorId).Find(&goodsColor)

	currentData := models.Cart{
		Id:           goodsId,
		Title:        goods.Title,
		Price:        goods.Price,
		GoodsVersion: goods.GoodsVersion,
		Num:          1,
		GoodsColor:   goodsColor.ColorName,
		GoodsImg:     goods.GoodsImg,
		GoodsGift:    goods.GoodsGift, /*赠品*/
		GoodsAttr:    "",              //根据自己的需求拓展
		Checked:      true,            /*默认选中*/
	}

	// 2、判断购物车有没有数据   （cookie）
	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)

	if len(cartList) > 0 {
		//4、购物车有数据  判断购物车有没有当前数据
		if models.HasCartData(cartList, currentData) {
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == currentData.Id && cartList[i].GoodsColor == currentData.GoodsColor && cartList[i].GoodsAttr == currentData.GoodsAttr {
					cartList[i].Num = cartList[i].Num + 1
				}
			}
		} else {
			cartList = append(cartList, currentData)
		}

		models.Cookie.Set(c, "cartList", cartList)

	} else {
		// 3、如果购物车没有任何数据  直接把当前数据写入cookie
		cartList = append(cartList, currentData)
		models.Cookie.Set(c, "cartList", cartList)
	}
	c.Redirect(302, "/cart/successTip?goods_id="+models.IntToString(goodsId))
}

func (con CartController) AddCartSuccess(c *gin.Context) {

	goodsId, err := models.StringToInt(c.Query("goods_id"))
	if err != nil {
		c.Redirect(302, "/")
	}

	goods := models.Goods{}
	models.DB.Where("id=?", goodsId).Find(&goods)

	var tpl = "oystershop/cart/addcart_success.html"
	con.Render(c, tpl, gin.H{
		"goods": goods,
	})
}

// 增加购物车数量
func (con CartController) IncCart(c *gin.Context) {
	//1、获取客户端穿过来的数据
	goodsId, err := models.StringToInt(c.Query("goods_id"))
	goodsColor := c.Query("goods_color")
	GoodsAttr := ""

	//定义返回的数据
	var allPrice float64
	var currentPrice float64
	var num int

	var response gin.H
	//2、判断数据是否合法
	if err != nil {
		response = gin.H{
			"success": false,
			"message": "传入参数错误",
		}
	} else {
		cartList := []models.Cart{}
		models.Cookie.Get(c, "cartList", &cartList)
		if len(cartList) > 0 {
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == goodsId && cartList[i].GoodsColor == goodsColor && cartList[i].GoodsAttr == GoodsAttr {
					cartList[i].Num = cartList[i].Num + 1
					currentPrice = float64(cartList[i].Num) * cartList[i].Price
					num = cartList[i].Num
				}

				if cartList[i].Checked {
					allPrice += cartList[i].Price * float64(cartList[i].Num)
				}

			}
			//重新写入数据
			models.Cookie.Set(c, "cartList", cartList)

			response = gin.H{
				"success":      true,
				"message":      "更新数据成功",
				"allPrice":     allPrice,
				"num":          num,
				"currentPrice": currentPrice,
			}
		} else {
			response = gin.H{
				"success": false,
				"message": "传入参数错误",
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

// 减少购物车数量
func (con CartController) DecCart(c *gin.Context) {
	//1、获取客户端穿过来的数据
	goodsId, err := models.StringToInt(c.Query("goods_id"))
	goodsColor := c.Query("goods_color")
	GoodsAttr := ""

	//定义返回的数据
	var allPrice float64
	var currentPrice float64
	var num int

	var response gin.H
	//2、判断数据是否合法
	if err != nil {
		response = gin.H{
			"success": false,
			"message": "传入参数错误",
		}
	} else {
		cartList := []models.Cart{}
		models.Cookie.Get(c, "cartList", &cartList)
		if len(cartList) > 0 {
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == goodsId && cartList[i].GoodsColor == goodsColor && cartList[i].GoodsAttr == GoodsAttr {
					if cartList[i].Num > 1 {
						cartList[i].Num = cartList[i].Num - 1
					}
					currentPrice = float64(cartList[i].Num) * cartList[i].Price
					num = cartList[i].Num
				}

				if cartList[i].Checked {
					allPrice += cartList[i].Price * float64(cartList[i].Num)
				}

			}
			//重新写入数据
			models.Cookie.Set(c, "cartList", cartList)

			response = gin.H{
				"success":      true,
				"message":      "更新数据成功",
				"allPrice":     allPrice,
				"num":          num,
				"currentPrice": currentPrice,
			}
		} else {
			response = gin.H{
				"success": false,
				"message": "传入参数错误",
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

func (con CartController) ChangeOneCart(c *gin.Context) {
	//1、获取客户端穿过来的数据
	goodsId, err := models.StringToInt(c.Query("goods_id"))
	goodsColor := c.Query("goods_color")
	GoodsAttr := ""

	//定义返回的数据
	var allPrice float64

	var response gin.H
	//2、判断数据是否合法
	if err != nil {
		response = gin.H{
			"success": false,
			"message": "传入参数错误",
		}
	} else {
		cartList := []models.Cart{}
		models.Cookie.Get(c, "cartList", &cartList)
		if len(cartList) > 0 {
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == goodsId && cartList[i].GoodsColor == goodsColor && cartList[i].GoodsAttr == GoodsAttr {
					cartList[i].Checked = !cartList[i].Checked

				}
				if cartList[i].Checked {
					allPrice += cartList[i].Price * float64(cartList[i].Num)
				}
			}
			//重新写入数据
			models.Cookie.Set(c, "cartList", cartList)

			response = gin.H{
				"success":  true,
				"message":  "更新数据成功",
				"allPrice": allPrice,
			}
		} else {
			response = gin.H{
				"success": false,
				"message": "传入参数错误",
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

func (con CartController) ChangeAllCart(c *gin.Context) {
	flag, _ := models.StringToInt(c.Query("flag"))

	var allPrice float64
	var response gin.H

	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)
	if len(cartList) > 0 {
		for i := 0; i < len(cartList); i++ {
			if flag == 1 {
				cartList[i].Checked = true
			} else {
				cartList[i].Checked = false
			}

			if cartList[i].Checked {
				allPrice += cartList[i].Price * float64(cartList[i].Num)
			}
		}
		//重新写入数据
		models.Cookie.Set(c, "cartList", cartList)

		response = gin.H{
			"success":  true,
			"message":  "更新数据成功",
			"allPrice": allPrice,
		}
	} else {
		response = gin.H{
			"success": false,
			"message": "传入参数错误",
		}
	}
	c.JSON(http.StatusOK, response)
}

func (con CartController) DelCart(c *gin.Context) {

	goodsId, _ := models.StringToInt(c.Query("goods_id"))
	goodsColor := c.Query("goods_color")
	GoodsAttr := ""

	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == goodsId && cartList[i].GoodsColor == goodsColor && cartList[i].GoodsAttr == GoodsAttr {
			cartList = append(cartList[:i], cartList[(i+1):]...)
		}
	}
	//重新写入数据
	models.Cookie.Set(c, "cartList", cartList)
	c.Redirect(http.StatusFound, "/cart")
}
