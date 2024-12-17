package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
)

var ctx = context.Background()
var rdbClient *redis.Client
var redisEnable bool

func init() {

	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}

	ip := config.Section("redis").Key("ip").String()
	port := config.Section("redis").Key("port").String()
	redisEnable, _ = config.Section("redis").Key("redisEnable").Bool()

	if redisEnable {
		//连接redis数据库
		rdbClient = redis.NewClient(&redis.Options{
			Addr:     ip + ":" + port,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		_, err := rdbClient.Ping(ctx).Result()
		if err != nil {
			fmt.Println("redis数据库连接失败")
		} else {
			fmt.Println("redis数据库连接成功...")
		}
	}

}

type cacheDb struct{}

func (c cacheDb) Set(key string, value interface{}, expiration int) {
	if redisEnable {
		v, err := json.Marshal(value)
		if err == nil {
			rdbClient.Set(ctx, key, string(v), time.Second*time.Duration(expiration))
		}
	}
}

func (c cacheDb) Get(key string, obj interface{}) bool {
	if redisEnable {
		valueStr, err1 := rdbClient.Get(ctx, key).Result()
		//如果数据库连不上的话 也会缓存到redis中
		if err1 == nil && valueStr != "" && valueStr != "[]" {
			err2 := json.Unmarshal([]byte(valueStr), obj)
			return err2 == nil
		}
		return false
	}
	return false
}

//清除缓存
func (c cacheDb) FlushAll() {
	if redisEnable {
		rdbClient.FlushAll(ctx)
	}
}

var CacheDb = &cacheDb{}
