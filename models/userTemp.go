package models

type UserTemp struct {
	Id        int
	Ip        string
	Phone     string
	SendCount int
	AddDay    string
	AddTime   int
	Sign      string
}

func (UserTemp) TableName() string {
	return "user_temp"
}
