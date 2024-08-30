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
	fmt.Println(Cache)
	set := Cache.HMSet(context.Background(), "countdown:FDC:0aef70aa-8ef4-4781-924e-88a6e4e037d9", "day", 2)
	fmt.Println(set)
}
