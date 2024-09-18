/**
 * @Author tanchang
 * @Description 测试redis
 * @Date 2024/8/30 14:31
 * @File:  redis_test
 * @Software: GoLand
 **/

package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRedis(t *testing.T) {
	//连接redis
	//连接redis
	Cache := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	err := Cache.Ping(context.Background()).Err()
	if err != nil {
		logrus.Error("redis连接失败！", err)
	}
	keys, _, err := Cache.Scan(context.Background(), 0, "admin:countdown:FDC:d1bba3b7-da6a-494f-b4ac-8ba67edd5d65", 3).Result()
	logrus.Println(keys)
	val := Cache.ZScore(context.Background(), "kyy", "is").Val()
	fmt.Println(val)
}
