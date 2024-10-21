package models

/*
	goodsAttr为商品规格和包装的保存地
*/

type GoodsAttr struct {
	Id              int
	GoodsId         int
	AttributeCateId int
	AttributeId     int
	AttributeTitle  string
	AttributeType   int
	AttributeValue  string
	Sort            int
	AddTime         int
	Status          int
}

func (GoodsAttr) TableName() string {
	return "goods_attr"
}
