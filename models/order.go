package models

type Order struct {
	Id          int
	OrderId     string
	Uid         int
	AllPrice    float64
	Phone       string
	Name        string
	Address     string
	PayStatus   int // 支付状态： 0 表示未支付     1 已经支付
	PayType     int // 支付类型： 0 alipay    1 wechat
	OrderStatus int // 订单状态： 0 已下单  1 已付款  2 已配货  3、发货   4、交易成功   5、退货   6、取消
	AddTime     int
	/*	PayTime          int //支付时间
		DistributionTime int //配货时间
		ExwarehouseTime  int //出库时间
		SuccessfulTime   int //交易成功时间
		CancelTime       int //取消时间
		ReturnTime       int //退款时间
		LogisticsCompany int //物流公司
		WaybillNo        int //运单号*/
	//其他的字段
	OrderItem []OrderItem `gorm:"foreignKey:OrderId;references:Id"`
}

func (Order) TableName() string {
	return "order"
}
