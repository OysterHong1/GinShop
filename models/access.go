package models

type Access struct {
	Id          int
	ModuleName  string //模块名称
	ActionName  string //操作名称
	Type        int    //节点类型 :  1、表示模块    2、表示菜单     3、操作
	Url         string //路由跳转地址
	ModuleId    int    //此module_id和当前模型的id关联       module_id= 0 表示模块
	Sort        int
	Description string
	Status      int
	AddTime     int
	AccessItem  []Access `gorm:"foreignKey:ModuleId;references:Id"`
	Checked     bool     `gorm:"-"` // 通过struct读写会忽略该字段 用于给Access添加属性

}

func (Access) TableName() string {
	return "access"
}
