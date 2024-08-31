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
	result, err := Cache.HGetAll(context.Background(), "countdown:FDC:de104cf5-da79-4494-9d4b-31a6ddffd66f").Result()
	fmt.Println(result["startTime"])
}
