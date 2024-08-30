/**
 * @Author tanchang
 * @Description redis链接
 * @Date 2024/7/11 16:31
 * @File:  Cache
 * @Software: GoLand
 **/

package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strconv"
)

var Cache *redis.Client

// RedisUtils redis连接
func RedisUtils(RDBAddr, RDBPwd, RDBDefaultDB string) {
	// 将字符串转换成int
	RDB, err := strconv.Atoi(RDBDefaultDB)
	if err != nil {
		fmt.Println("将string转换成int失败！", err)
	}

	//连接redis
	Cache = redis.NewClient(&redis.Options{
		Addr:     RDBAddr,
		Password: RDBPwd,
		DB:       RDB,
	})
	err = Cache.Ping(context.Background()).Err()
	if err != nil {
		logrus.Error("redis连接失败！", err)
	}
}
