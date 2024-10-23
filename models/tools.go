package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

// 时间戳转换成日期
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

// 日期转换成时间戳 2020-05-02 15:04:05
func DateToUnix(str string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// 获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// 获取当前的日期
func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// 获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// Md5加密
func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// String转换为Int
func StringToInt(str string) (int, error) {
	n, err := strconv.Atoi(str)
	return n, err
}

// string转float64
func Float(str string) (float64, error) {
	n, err := strconv.ParseFloat(str, 64)
	return n, err
}

func IntToString(num int) string {
	return strconv.Itoa(num)
}

// 上传文件
func UploadFile(c *gin.Context, picName string) (string, error) {
	file, err := c.FormFile(picName)
	if err != nil {
		return "", err
	}

	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	if _, ok := allowExtMap[extName]; !ok {
		c.String(200, "上传的文件不合法")
		return "", errors.New("文件后缀名不合法")
	}

	day := GetDay()
	dir := "./static/upload/" + day

	err1 := os.MkdirAll(dir, 0666)
	if err1 != nil {
		c.String(200, "MkdirAll失败")
		return "", err1
	}

	fileName := strconv.FormatInt(GetUnixNano(), 10) + extName
	dst := path.Join(dir, fileName)
	err2 := c.SaveUploadedFile(file, dst)
	if err2 != nil {
		return "", err2
	}
	return dst, nil
}

// 字符串解析为html
func Str2Html(str string) template.HTML {
	return template.HTML(str)
}
