package oystershop

import (
	"GinShop/models"
	"math"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
}

func (con UserController) Index(c *gin.Context) {
	// c.String(http.StatusOK, "首页")
	var tpl = "oystershop/user/welcome.html"
	con.Render(c, tpl, gin.H{})
}

func (con UserController) OrderList(c *gin.Context) {
	// 当前页
	page, _ := models.StringToInt(c.Query("page"))
	if page == 0 {
		page = 1
	}
	//pageSize

	pageSize := 2

	//获取当前用户
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	//模糊查询
	where := "uid =" + models.IntToString(user.Id)
	keywords := c.Query("keywords")
	if keywords != "" {
		//查询
		orderItemList := []models.OrderItem{}
		models.DB.Where("product_title like ?", "%"+keywords+"%").Find(&orderItemList)
		var str string
		// 字符串：   12,12,22
		for i := 0; i < len(orderItemList); i++ {
			if i == 0 {
				str += models.IntToString(orderItemList[i].OrderId)
			} else {
				str += "," + models.IntToString(orderItemList[i].OrderId)
			}
		}
		where += " AND id in (" + str + ")"
	}
	//按照状态筛选订单
	orderStatus, statusErr := models.StringToInt(c.Query("orderStatus"))
	if statusErr == nil && orderStatus >= 0 {
		where += " AND order_status=" + models.IntToString(orderStatus)
	} else {
		orderStatus = -1
	}

	//获取当前用户下面订单信息
	orderList := []models.Order{}
	models.DB.Where(where).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&orderList)

	//获取总数量
	var count int64
	models.DB.Where(where).Table("order").Count(&count)

	var tpl = "oystershop/user/order.html"
	con.Render(c, tpl, gin.H{
		"order":       orderList,
		"page":        page,
		"keywords":    keywords,
		"orderStatus": orderStatus,
		"totalPages":  math.Ceil(float64(count) / float64(pageSize)),
	})
}

func (con UserController) OrderInfo(c *gin.Context) {
	id, err := models.StringToInt(c.Query("id"))
	if err != nil {
		c.Redirect(302, "/user/order")
	}
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	order := []models.Order{}
	models.DB.Where("id=? And uid=?", id, user.Id).Preload("OrderItem").Find(&order)

	if len(order) == 0 {
		c.Redirect(302, "/user/order")
		return
	}
	var tpl = "oystershop/user/order_info.html"
	con.Render(c, tpl, gin.H{
		"order": order[0],
	})
}
