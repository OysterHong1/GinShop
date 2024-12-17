package oystershop

import (
	"GinShop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BuyController struct {
	BaseController
}

func (con BuyController) Checkout(c *gin.Context) {
	//1、获取购物车中选择的商品

	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)

	orderList := []models.Cart{}
	var allPrice float64
	var allNum int

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
			allNum += cartList[i].Num
		}
	}

	//2、获取当前用户的收货地址

	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	addressList := []models.Address{}
	models.DB.Where("uid = ?", user.Id).Order("id desc").Find(&addressList)

	//3.生成签名
	orderSign := models.Md5(models.GetRandomNum())
	session := sessions.Default(c)
	session.Set("orderSign", orderSign)
	session.Save()

	//判断orderList数据是否存在
	if len(orderList) == 0 {
		c.Redirect(302, "/")
		return
	}

	con.Render(c, "oystershop/buy/checkout.html", gin.H{
		"orderList":   orderList,
		"allPrice":    allPrice,
		"allNum":      allNum,
		"addressList": addressList,
		"orderSign":   orderSign,
	})

}

/*
提交订单执行结算

	1、获取用户信息 获取用户的收货地址信息
	2、获取购买商品的信息
	3、把订单信息放在订单表，把商品信息放在商品表
	4、删除购物车里面的选中数据
	5、跳转到支付页面
*/
func (con BuyController) DoCheckout(c *gin.Context) {

	//防止订单重复提交
	orderSignClient := c.PostForm("orderSign")
	session := sessions.Default(c)
	orderSignSession := session.Get("orderSign")
	orderSignServer, ok := orderSignSession.(string)
	if !ok {
		c.Redirect(302, "/")
		return
	}
	if orderSignServer != orderSignClient {
		c.Redirect(302, "/")
		return
	}

	session.Delete("orderSign") //完成后删除订单标签
	session.Save()

	// 1、获取用户信息 获取用户的收货地址信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	addressResult := []models.Address{}
	models.DB.Where("uid = ? AND default_address=1", user.Id).Find(&addressResult)
	if len(addressResult) == 0 {
		c.Redirect(302, "/buy/checkout")
		return
	}

	// 2、获取购买商品的信息
	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)
	orderList := []models.Cart{}
	var allPrice float64
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
		}
	}
	// 3、把订单信息放在订单表，把商品信息放在商品表
	order := models.Order{
		OrderId:     models.GetOrderId(),
		Uid:         user.Id,
		AllPrice:    allPrice,
		Phone:       addressResult[0].Phone,
		Name:        addressResult[0].Name,
		Address:     addressResult[0].Address,
		PayStatus:   0,
		PayType:     0,
		OrderStatus: 0,
		AddTime:     int(models.GetUnix()),
	}

	err := models.DB.Create(&order).Error
	//增加数据成功以后可以通过  order.Id
	if err == nil {
		// 把商品信息放在商品对应的订单表
		for i := 0; i < len(orderList); i++ {
			orderItem := models.OrderItem{
				OrderId:      order.Id,
				Uid:          user.Id,
				ProductTitle: orderList[i].Title,
				ProductId:    orderList[i].Id,
				ProductImg:   orderList[i].GoodsImg,
				ProductPrice: orderList[i].Price,
				ProductNum:   orderList[i].Num,
				GoodsVersion: orderList[i].GoodsVersion,
				GoodsColor:   orderList[i].GoodsColor,
			}
			models.DB.Create(&orderItem)
		}
	}

	// 4、删除购物车里面的选中数据
	noSelectCartList := []models.Cart{}
	for i := 0; i < len(cartList); i++ {
		if !cartList[i].Checked {
			noSelectCartList = append(noSelectCartList, cartList[i])
		}
	}
	models.Cookie.Set(c, "cartList", noSelectCartList)

	c.Redirect(302, "/buy/pay?orderId="+models.IntToString(order.Id))
}

// 支付
func (con BuyController) Pay(c *gin.Context) {
	orderId, err := models.StringToInt(c.Query("orderId"))
	if err != nil {
		c.Redirect(302, "/")
	}
	//获取用户信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	//获取订单信息
	order := models.Order{}
	models.DB.Where("id = ?", orderId).Find(&order)
	if order.Uid != user.Id {
		c.Redirect(302, "/")
		return
	}
	//获取订单对应的商品

	orderItems := []models.OrderItem{}
	models.DB.Where("order_id = ?", orderId).Find(&orderItems)

	con.Render(c, "oystershop/buy/pay.html", gin.H{
		"order":      order,
		"orderItems": orderItems,
	})
}

func (con BuyController) OrderPayStatus(c *gin.Context) {

	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入参数错误",
		})
		return
	}
	//获取用户信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	//获取主订单信息
	order := models.Order{}
	models.DB.Where("id=?", id).Find(&order)

	//判断当前数据是否合法
	if user.Id != order.Uid {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "非法请求",
		})
		return
	}

	//判断是否支付
	if order.PayStatus == 1 && order.OrderStatus == 1 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "支付成功",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "支付成功",
		})
		return
	}

}
