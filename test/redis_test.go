/**
 * @Author tanchang
 * @Description 测试redis
 * @Date 2024/8/30 14:31
 * @File:  redis_test
 * @Software: GoLand
 **/

package test

import (
	serializes "GoToDoList/serialized"
	"GoToDoList/utils"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strings"
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
	keys, _, _ := Cache.Scan(context.Background(), 0, utils.DELCountdownPrefix+"*c5f3facf-ccf9-4d78-be76-959272fcfdf4", 50).Result()
	fmt.Println(len(keys))
	for _, countdown := range keys {
		identity := strings.Split(countdown, ":")[3]
		result, _ := Cache.HGetAll(context.Background(), countdown).Result()
		fmt.Println(serializes.CountdownSerializeSingle(result, identity))
	}
}
