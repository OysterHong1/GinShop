package models

type GoodsTypeAttribute struct {
	Id        int
	CateId    int
	Title     string
	AttrType  int
	AttrValue string
	Status    int
	Sort      int
	AddTime   int
}

func (GoodsTypeAttribute) TableName() string {
	return "goods_type_attribute"
}
