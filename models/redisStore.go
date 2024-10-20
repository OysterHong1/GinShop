package models

import (
	"context"
	"fmt"
	"time"
)

var Storectx = context.Background()

const CAPTCHA = "captcha:"

type RedisStore struct {
}

// 实现设置captcha的方法
func (r RedisStore) Set(id string, value string) error {
	key := CAPTCHA + id
	err := RedisDb.Set(ctx, key, value, time.Minute*2).Err()

	return err
}

// 实现获取captcha的方法
func (r RedisStore) Get(id string, clear bool) string {
	key := CAPTCHA + id
	val, err := RedisDb.Get(Storectx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear {
		err := RedisDb.Del(Storectx, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// 实现验证captcha的方法
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := RedisStore{}.Get(id, clear)
	//fmt.Println("key:"+id+";value:"+v+";answer:"+answer)
	return v == answer
}
