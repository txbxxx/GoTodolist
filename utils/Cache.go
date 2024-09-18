/**
 * @Author tanchang
 * @Description redis链接
 * @Date 2024/7/11 16:31
 * @File:  Cache
 * @Software: GoLand
 **/

package utils

import (
	"GoToDoList/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	FDCCountdownPrefix = "countdown:FDC:"
	OECCountdownPrefix = "countdown:OEC:"
	DELCountdownPrefix = "delete:"
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

// DeleteForRecycle
// 从回收站删除一条数据
func DeleteForRecycle(identity string) error {
	keys, _, err := Cache.Scan(context.Background(), 0, "*"+DELCountdownPrefix+"*"+identity, 1).Result()
	if err != nil {
		return fmt.Errorf("查询redis中回收站数据失败: %v", err)
	}
	// 检查是否有待删除的键
	if len(keys) == 0 {
		return nil
	}
	deleteCount, err := Cache.Del(context.Background(), keys...).Result()
	if err != nil {
		return fmt.Errorf("删除redis数据失败: %v", err)
	}
	if deleteCount == 0 {
		return fmt.Errorf("没有找到需要删除的数据")
	}
	return nil
}

// ListFormRedis 从redis中获取数据并添加至列表
// Param keys redis中的多个key
func ListFormRedis(ctx context.Context, keys []string) ([]map[string]string, error) {
	list := make([]map[string]string, 0)
	for _, key := range keys {
		result := Cache.HGetAll(ctx, key)
		if err := result.Err(); err != nil {
			return nil, err
		}
		list = append(list, result.Val())
	}
	return list, nil
}

// OecCalculate 计算Oec
// FDC计算方法是当前期戳-开始日期的时间戳 最后在/86400 获得天数
// 使用Ceil向上取整 0.1 天也是1天
func OecCalculate(now int64, countdown model.CountDown, key string) error {
	ctx := context.Background()
	day := float64(now-countdown.StartTime) / 86400
	// 将倒计时同步至redis，时间则向上取整
	if _, err := Cache.HSet(ctx, key,
		map[string]any{
			"startTime":        countdown.StartTime,
			"day":              math.Ceil(day),
			"background":       countdown.Background,
			"name":             countdown.Name,
			"identity":         countdown.Identity,
			"categoryIdentity": countdown.CategoryIdentity,
		}).Result(); err != nil {
		return fmt.Errorf("同步redis失败: %v", err)
	}
	// 设置过期时间
	duration := time.Hour / 2
	// 随机增加啊0-60分钟
	random := time.Duration(rand.Intn(60)) * time.Minute
	if err := Cache.Expire(ctx, key, duration+random).Err(); err != nil {
		return fmt.Errorf("设置过期时间失败: %v", err)
	}
	return nil
}

// FdcCalculate 计算Fdc
// FDC计算方法是结束日期戳-当前日期的时间戳 最后在/86400 获得天数
// 使用Ceil向上取整 0.1 天也是1天
func FdcCalculate(now int64, countdown model.CountDown, key string) error {
	ctx := context.Background()
	day := float64(countdown.EndTime-now) / 86400
	// 获取username
	// 将倒计时同步至redis，时间则向上取整
	if _, err := Cache.HSet(ctx, key,
		map[string]any{
			"endTime":          countdown.EndTime,
			"starTime":         countdown.StartTime,
			"day":              math.Ceil(day),
			"background":       countdown.Background,
			"name":             countdown.Name,
			"identity":         countdown.Identity,
			"categoryIdentity": countdown.CategoryIdentity,
		}).Result(); err != nil {
		return fmt.Errorf("同步redis失败: %v", err)
	}
	// 设置过期时间
	duration := time.Hour / 2
	// 随机增加啊0-60分钟
	random := time.Duration(rand.Intn(60)) * time.Minute
	if err := Cache.Expire(ctx, key, duration+random).Err(); err != nil {
		return fmt.Errorf("设置过期时间失败: %v", err)
	}
	return nil
}

// RefFDC FDC刷新
// 从redis中读取数据计算剩余日期
func RefFDC() error {
	// 当前时间戳
	now := time.Now().Unix()
	// 查询redis中FDC的数据
	FDCKeys, _, err := Cache.Scan(context.Background(), 0, "*:"+FDCCountdownPrefix+"*", 50).Result()
	if err != nil {
		logrus.Error("查询redis中FDC的数据失败", err)
		return fmt.Errorf(err.Error())
	}

	for _, FDC := range FDCKeys {
		// 获取当前OEC key里面的全部字段，返回一个字符串map
		result, err := Cache.HGetAll(context.Background(), FDC).Result()
		if err != nil {
			return fmt.Errorf("查询redis中OEC的数据失败: %v", err)
		}
		// 转换为int64
		endTime, _ := strconv.ParseInt(result["endTime"], 10, 64)
		startTime, _ := strconv.ParseInt(result["start"], 10, 64)
		//取出identity FDC的格式为 countdown:FDC:{{ identity }}
		identity := strings.Split(FDC, ":")[3]
		// 判断当前日期时间戳是否大于结束日期时间戳
		if now >= endTime {
			//将已经到达的倒计时加入回收站
			err := AddCountDownRecycle(FDC, identity)
			if err != nil {
				return err
			}
			logrus.Info("到达的倒计时加入回收站成功")
			continue
		}
		// 创建倒计时对象
		count := model.CountDown{
			Identity:         identity,
			EndTime:          endTime,
			StartTime:        startTime,
			Name:             result["name"],
			Background:       result["background"],
			CategoryIdentity: result["categoryIdentity"],
		}
		if err := FdcCalculate(now, count, FDC); err != nil {
			return err
		}
	}
	return nil
}

// RefOEC OEC刷新
// 从redis中读取数据计算过期日期
func RefOEC() error {
	// 当前时间戳
	now := time.Now().Unix()
	// 查询redis中OEC的数据
	FDCKeys, _, err := Cache.Scan(context.Background(), 0, "*:"+OECCountdownPrefix+"*", 50).Result()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	for _, OEC := range FDCKeys {
		// 获取identity值
		identity := strings.Split(OEC, ":")[3]
		// 获取当前OEC key里面的全部字段，返回一个字符串map
		result, err := Cache.HGetAll(context.Background(), OEC).Result()
		if err != nil {
			return fmt.Errorf("查询redis中OEC的数据失败: %v", err)
		}
		// 转换为int64
		startTime, _ := strconv.ParseInt(result["startTime"], 10, 64)
		// 创建对象
		count := model.CountDown{
			Identity:         identity,
			StartTime:        startTime,
			Name:             result["name"],
			Background:       result["background"],
			CategoryIdentity: result["categoryIdentity"],
		}
		// 计算过去时间
		if err := OecCalculate(now, count, OEC); err != nil {
			return err
		}
	}
	return nil
}

// AddCountDownRecycle 添加至回收站
// 将倒计时加上前缀delete: 表示加入回收站了
func AddCountDownRecycle(key string, identity string) error {
	//将已经到达的倒计时加入回收站
	rename := Cache.Rename(context.Background(), key, DELCountdownPrefix+key)
	if rename.Err() != nil {
		return fmt.Errorf("将已经到达的倒计时加入回收站失败: %v", rename.Err())
	}
	// 删除sql数据
	err := DB.Model(&model.CountDown{}).Where("identity = ?", identity).Delete(&model.CountDown{}).Error
	if err != nil {
		return fmt.Errorf("删除sql数据失败: %v", err)
	}
	return nil
}
