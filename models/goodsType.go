package models

// goodstype为商品类型
// 如手机，电脑，路由器......
type GoodsType struct {
	Id          int
	Title       string
	Description string
	Status      int
	AddTime     int
}

func (GoodsType) TableName() string {
	return "goods_type"
}
