package oystershop

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/objcoding/wxpay"
	qrcode "github.com/skip2/go-qrcode"
)

type WxpayController struct{}

func (con WxpayController) Wxpay(c *gin.Context) {

	//1、获取订单号 判断此订单号是否值当前用户的
	//2、获取订单里面的支付信息

	//1、配置基本信息
	account := wxpay.NewAccount(
		"wx7bf3787c783116e4",
		"1502539541",
		"zhongyuantengitying6666666666666",
		false,
	)
	client := wxpay.NewClient(account)

	//2、获取ip地址   订单号等信息
	ip := c.ClientIP()
	template := "200601021504"
	tradeNo := time.Now().Format(template)
	//3、调用统一下单
	params := make(wxpay.Params)
	params.SetString("body", "Gin微信支付11").
		SetString("out_trade_no", tradeNo).
		SetInt64("total_fee", 1). //1分
		SetString("spbill_create_ip", ip).
		SetString("notify_url", "http://pay.apiying.com/wxpay/notify"). //必须在商户平台的Native支付回调链接里面配置
		// SetString("trade_type", "APP")
		SetString("trade_type", "NATIVE") //网站支付需要改为NATIVE

	p, err1 := client.UnifiedOrder(params)
	if err1 != nil {
		c.String(http.StatusOK, "生成二维码失败")
		return
	}
	//4、获取code_url 生成支付二维码

	png, err := qrcode.Encode(p["code_url"], qrcode.Medium, 256)
	if err != nil {
		c.String(http.StatusOK, "生成二维码失败")
		return
	}
	c.String(http.StatusOK, string(png))
}

/*
1、发布到服务器测试
2、必须在商户平台的Native支付回调链接里面配置
3、如何接收XML的数据  c.GetRawData()
4、如何获取数据
5、如何验证数据
6、更新数据
*/
func (con WxpayController) WxpayNotify(c *gin.Context) {
	//1、获取表单传过来的xml数据
	xmlByte, _ := c.GetRawData()
	xmlStr := string(xmlByte)

	postParams := wxpay.XmlToMap(xmlStr)

	//2、校验签名
	account := wxpay.NewAccount(
		"wx7bf3787c783116e4",
		"1502539541",
		"zhongyuantengitying6666666666666",
		false,
	)
	client := wxpay.NewClient(account)
	isValidate := client.ValidSign(postParams)

	fmt.Println(isValidate)
	fmt.Println("-----更新订单-----")
	fmt.Println(postParams)

	//3、更新订单
	c.String(200, "ok")
}
