package models

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

//定义结构体  缓存结构体 私有
type ginCookie struct{}

//写入数据的方法
func (cookie ginCookie) Set(c *gin.Context, key string, value interface{}) {

	bytes, _ := json.Marshal(value)
	c.SetCookie(key, string(bytes), 3600, "/", "localhost", false, true)
}

//获取数据的方法
func (cookie ginCookie) Get(c *gin.Context, key string, obj interface{}) bool {

	valueStr, err1 := c.Cookie(key)
	if err1 == nil && valueStr != "" && valueStr != "[]" {
		err2 := json.Unmarshal([]byte(valueStr), obj)
		return err2 == nil
	}
	return false
}

//实例化结构体
var Cookie = &ginCookie{}
