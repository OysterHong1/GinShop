package models

type GoodsColor struct {
	Id         int
	ColorName  string
	ColorValue string
	Status     int
	Checked    bool `gorm:"-"` //该语法表示操作数据库时忽略该项
}

func (GoodsColor) TableName() string {
	return "goods_color"
}
