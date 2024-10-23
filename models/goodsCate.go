package models

type GoodsCate struct {
	Id             int
	Title          string
	CateImg        string
	Link           string
	Template       string
	Pid            int
	SubTitle       string
	Keywords       string
	Description    string
	Sort           int
	Status         int
	AddTime        int
	GoodsCateItems []GoodsCate `gorm:"foreignKey:pid;references:Id"`
}

func (GoodsCate) TableName() string {
	return "goods_cate"
}
