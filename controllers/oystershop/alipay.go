package oystershop

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
)

type AlipayController struct{}

func (con AlipayController) Alipay(c *gin.Context) {
	//1、获取订单号 判断此订单号是否值当前用户的
	//2、获取订单里面的支付信息

	var privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCPUi2yj4RJcqAKfjV/AxjqIjjUq007pGhLClbrwEoCVhgnQU9bAVFnPlDaVdO2xBfu8D/gHCwV9czZKJh51yg1kwE6Yv/hSd966PnGS9NszfxxWfTbWeC6DIHZj66nTECK+vYWX36tKxIG+juzXyfoAuyL1h58oFgF/Zwa1FzpKDTcMUIE+npf3DpMS0Uatf6TsREzVQQs3i0WVxYY3lv1Dmualr6Q3GAFE1j/xt/STE993uh8MVeLS+RcTrrjilPSVH5Z0DLAjkDSH7XUK1lIHASpgOrddEJ8MeT8L7bUJ8nDs5qQ8zfbPVJaXfsF2NYS162HkIW2bl0r5c5Mudt7AgMBAAECggEAWUHVweHNgiyH7WECkhJsvswHVrNEi0NtzGYpEfOUY/YYXsI22Lduaf0OP5u6GZXwTdeEAF+rORX2uLumkiLkINFnr2QedcEbFCHqBIwOpTF36WQbsUw9P8EwUT1BiWFcxPFctzxL2S78sCnBaol1gfHoPYJhRD5b84cpZDAjmPSJk1XAtAtKChUIskLBAsCvwlGHbx/6UQwM9eKgwo0Y67MCPW9wBjE9bRFWBfaeLszEVu3nKyOKLwcUGDXrbmBS8bj9YtyqTG7RjKZmIGuJtPKehEjlNn2ALYYMHUA5VSdVP4LrBiLVLQE5tTDedAzr3uHuWOgLaMeDnfwQg8Hj6QKBgQDGaCoZ3JRg4PbGetduHpGfWOstwZKNwqkDeYvpA03GcRYhzhqG98qdTGEYaLGkszQrh0ZcYk7fs8jsR5Lu2WTzImtiAYDCDEjSDHj5N5PfWR0r2pxrUvX8shun4QKpz4QQ0RUjeujZ2hJHkeCviF2+k5IBDvz6YtBo+H3IgR9QjwKBgQC47IVtuvlYZ0/2TAcAt40YhEgLZOr3NuR4eSnx76zf/8vRmHDEfIqvUdIzzJy3RVYT3uXiB7DYwaHr76ouP3lMOhgLgUlK4Dt8L3UbMP6Asr+6D/uggVmlIHKK76HZIdBL5nGQvTOvwE+fmEivlR5QV+cczqQdCYyNZXETmgAkVQKBgF32JMIcqZR71cLHmFDJX1Okq7P+sWY7YwmHPZA7hVDOa5nU3tE+dpEqA+2oX0DNsY5PwS2tTQc6QJRNjTNadymCCnLenVjIso/vYjc8b+ZdcKg9HsjhACgNPXWy5S0AXt4L9sPXyICreu60EkFvBl5jysh/jaUSuPqNfBxBsk/XAoGBAKYxkzzp4/wKZXfSHh0L2VemUuVCnlTtVWncYtEXeQObbX8CBJ7h2vXzj/mTs2iWfOTA11NLXCmB5FcZfpWv4ACc2U1FtSwA2BUkxZdZcfESNHMwuBEpDvrzbV3mPUvaMsxz366YC+Kw8B5biz+ZwbOtPHzMTfv2wAW3nGdkaSo9AoGBALW8kQ6QwXu3gB7q8Es1zCvSg07tnJks57nJDTtOOg6oc52RvKjndWHT6aFwp/YrESeHNvoudRhPm5nKUTbKDiLixDXcrvEU3b70dJGwWHkBdYU/s5J2o5oaclwiyiUptnIp2eA+Cu+wH1eREGmrOAdeBSFOfdjcfJTFmOCoalJR" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var client, err = alipay.New("2021001186696588", privateKey, true)
	client.LoadAppPublicCertFromFile("crt/appCertPublicKey_2021001186696588.crt") // 加载应用公钥证书
	client.LoadAliPayRootCertFromFile("crt/alipayRootCert.crt")                   // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("crt/alipayCertPublicKey_RSA2.crt")       // 加载支付宝公钥证书

	// 将 key 的验证调整到初始化阶段
	if err != nil {
		fmt.Println(err)
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://localhost:8080/v3/alipayNotify"
	p.ReturnURL = "http://localhost:8080/v3/alipayReturn"
	p.Subject = "测试 公钥证书模式-这是一个gin订单"
	template := "2006-01-02 15:04:05"
	p.OutTradeNo = time.Now().Format(template)
	p.TotalAmount = "0.1"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err4 = client.TradePagePay(p)
	if err4 != nil {
		fmt.Println(err4)
	}

	var payURL = url.String()
	fmt.Println(payURL)

	c.Redirect(302, payURL)

}
func (con AlipayController) AlipayNotify(c *gin.Context) {
	fmt.Println("AlipayNotify")

	var privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCPUi2yj4RJcqAKfjV/AxjqIjjUq007pGhLClbrwEoCVhgnQU9bAVFnPlDaVdO2xBfu8D/gHCwV9czZKJh51yg1kwE6Yv/hSd966PnGS9NszfxxWfTbWeC6DIHZj66nTECK+vYWX36tKxIG+juzXyfoAuyL1h58oFgF/Zwa1FzpKDTcMUIE+npf3DpMS0Uatf6TsREzVQQs3i0WVxYY3lv1Dmualr6Q3GAFE1j/xt/STE993uh8MVeLS+RcTrrjilPSVH5Z0DLAjkDSH7XUK1lIHASpgOrddEJ8MeT8L7bUJ8nDs5qQ8zfbPVJaXfsF2NYS162HkIW2bl0r5c5Mudt7AgMBAAECggEAWUHVweHNgiyH7WECkhJsvswHVrNEi0NtzGYpEfOUY/YYXsI22Lduaf0OP5u6GZXwTdeEAF+rORX2uLumkiLkINFnr2QedcEbFCHqBIwOpTF36WQbsUw9P8EwUT1BiWFcxPFctzxL2S78sCnBaol1gfHoPYJhRD5b84cpZDAjmPSJk1XAtAtKChUIskLBAsCvwlGHbx/6UQwM9eKgwo0Y67MCPW9wBjE9bRFWBfaeLszEVu3nKyOKLwcUGDXrbmBS8bj9YtyqTG7RjKZmIGuJtPKehEjlNn2ALYYMHUA5VSdVP4LrBiLVLQE5tTDedAzr3uHuWOgLaMeDnfwQg8Hj6QKBgQDGaCoZ3JRg4PbGetduHpGfWOstwZKNwqkDeYvpA03GcRYhzhqG98qdTGEYaLGkszQrh0ZcYk7fs8jsR5Lu2WTzImtiAYDCDEjSDHj5N5PfWR0r2pxrUvX8shun4QKpz4QQ0RUjeujZ2hJHkeCviF2+k5IBDvz6YtBo+H3IgR9QjwKBgQC47IVtuvlYZ0/2TAcAt40YhEgLZOr3NuR4eSnx76zf/8vRmHDEfIqvUdIzzJy3RVYT3uXiB7DYwaHr76ouP3lMOhgLgUlK4Dt8L3UbMP6Asr+6D/uggVmlIHKK76HZIdBL5nGQvTOvwE+fmEivlR5QV+cczqQdCYyNZXETmgAkVQKBgF32JMIcqZR71cLHmFDJX1Okq7P+sWY7YwmHPZA7hVDOa5nU3tE+dpEqA+2oX0DNsY5PwS2tTQc6QJRNjTNadymCCnLenVjIso/vYjc8b+ZdcKg9HsjhACgNPXWy5S0AXt4L9sPXyICreu60EkFvBl5jysh/jaUSuPqNfBxBsk/XAoGBAKYxkzzp4/wKZXfSHh0L2VemUuVCnlTtVWncYtEXeQObbX8CBJ7h2vXzj/mTs2iWfOTA11NLXCmB5FcZfpWv4ACc2U1FtSwA2BUkxZdZcfESNHMwuBEpDvrzbV3mPUvaMsxz366YC+Kw8B5biz+ZwbOtPHzMTfv2wAW3nGdkaSo9AoGBALW8kQ6QwXu3gB7q8Es1zCvSg07tnJks57nJDTtOOg6oc52RvKjndWHT6aFwp/YrESeHNvoudRhPm5nKUTbKDiLixDXcrvEU3b70dJGwWHkBdYU/s5J2o5oaclwiyiUptnIp2eA+Cu+wH1eREGmrOAdeBSFOfdjcfJTFmOCoalJR" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var client, err = alipay.New("2021001186696588", privateKey, true)
	client.LoadAppPublicCertFromFile("crt/appCertPublicKey_2021001186696588.crt") // 加载应用公钥证书
	client.LoadAliPayRootCertFromFile("crt/alipayRootCert.crt")                   // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("crt/alipayCertPublicKey_RSA2.crt")       // 加载支付宝公钥证书

	if err != nil {
		fmt.Println(err)
		return
	}

	req := c.Request
	req.ParseForm()
	ok := client.VerifySign(req.Form)

	fmt.Println(ok)

	fmt.Println(req.Form)

	c.String(200, "ok")
}
func (con AlipayController) AlipayReturn(c *gin.Context) {
	c.String(200, "支付成功")
}
