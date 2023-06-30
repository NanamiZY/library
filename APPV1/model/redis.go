package model

import (
	"fmt"
	"github.com/go-redis/redis"
)

var RedisConn *redis.Client
var RedisConn1 *redis.Client

func init() {
	// 存储验证码
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Redis 未设置密码时为空
		DB:       1,  // 使用默认数据库
	})
	//存储用户路径
	RedisConn1 = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Redis 未设置密码时为空
		DB:       0,  // 使用默认数据库
	})
	// 测试连接是否成功
	_, err := RedisConn.Ping().Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v", err)
		return
	}
	fmt.Println("Connected to Redis")
}
