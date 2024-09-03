/**
 * @Author tanchang
 * @Description 操作倒计时等其他操作
 * @Date 2024/8/30 21:57
 * @File:  Refresh
 * @Software: GoLand
 **/

package utils

import (
	"GoToDoList/model"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	FDCCountdownPrefix = "countdown:FDC:"
	OECCountdownPrefix = "countdown:OEC:"
	DELCountdownPrefix = "delete:"
)

// DelCountDownForRedis 从redis中删除一条数据
// 从redis中删除数据，并不加入回收站
func DelCountDownForRedis(identity string) error {
	fmt.Println("identity:", identity)
	keys, _, err := Cache.Scan(context.Background(), 0, "countdown:*:"+identity, 1).Result()
	if err != nil || len(keys) == 0 {
		return fmt.Errorf("要删除的数据不存在 %v", err)
	}
	// 删除
	err = Cache.Del(context.Background(), keys...).Err()
	if err != nil {
		return fmt.Errorf("删除redis数据失败: %v", err)
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

// OecCalculate 计算Oec
// FDC计算方法是当前期戳-开始日期的时间戳 最后在/86400 获得天数
// 使用Ceil向上取整 0.1 天也是1天
func OecCalculate(now, startTime int64, key, background, name string) (float64, error) {
	day := float64(now-startTime) / 86400
	// 将倒计时同步至redis，时间则向上取整
	if _, err := Cache.HSet(context.Background(), key, map[string]any{"startTime": startTime, "day": math.Ceil(day), "background": background, "name": name}).Result(); err != nil {
		return 0, fmt.Errorf("同步redis失败: %v", err)
	}
	return math.Ceil(day), nil
}

// FdcCalculate 计算Fdc
// FDC计算方法是结束日期戳-当前日期的时间戳 最后在/86400 获得天数
// 使用Ceil向上取整 0.1 天也是1天
func FdcCalculate(now, endTime int64, key, background, name string) (float64, error) {
	day := float64(endTime-now) / 86400
	// 将倒计时同步至redis，时间则向上取整
	if _, err := Cache.HSet(context.Background(), key, map[string]any{"endTime": endTime, "day": math.Ceil(day), "background": background, "name": name}).Result(); err != nil {
		return 0, fmt.Errorf("同步redis失败: %v", err)
	}
	return math.Ceil(day), nil
}

// RefreshDayForMysql 从mysql中读取数据刷新倒计时
// 从redis读取倒计时列表
// 将倒计时列表中的数据同步至redis
func RefreshDayForMysql() error {
	countdown := make([]model.CountDown, 1)
	if err := DB.Model(&model.CountDown{}).Find(&countdown).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("查询倒计时失败: %v", err)
		}
	}
	// 当前时间戳
	now := time.Now().Unix()
	for _, count := range countdown {
		key := OECCountdownPrefix + count.Identity
		if count.EndTime <= 0 {
			// 计算过去时间oec
			day, err := OecCalculate(now, count.StartTime, key, count.Background, count.Name)
			if err != nil {
				return err
			}
			logrus.Info("同步成功，剩余时间: ", math.Ceil(day))
		} else {
			key = FDCCountdownPrefix + count.Identity
			// 判断当前日期时间戳是否大于结束日期时间戳
			if now >= count.EndTime {
				// 大于则执行
				err := AddCountDownRecycle(key, count.Identity)
				if err != nil {
					return err
				}
				logrus.Info("到达的倒计时加入回收站成功")
				continue
			}
			//FDC
			// 如果没有大于，就计算还有多少天，使用结束时间减去现在时间
			day, err := FdcCalculate(now, count.EndTime, key, count.Background, count.Name)
			if err != nil {
				return err
			}
			logrus.Info("同步成功，剩余时间: ", math.Ceil(day))
		}
	}
	return nil
}

// RefFDC FDC刷新
// 从redis中读取数据计算剩余日期
func RefFDC() error {
	// 当前时间戳
	now := time.Now().Unix()
	// 查询redis中FDC的数据
	FDCKeys, _, err := Cache.Scan(context.Background(), 0, FDCCountdownPrefix+"*", 50).Result()
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
		//取出identity FDC的格式为 countdown:FDC:{{ identity }}
		split := strings.Split(FDC, "countdown:FDC:")
		// 判断当前日期时间戳是否大于结束日期时间戳
		if now >= endTime {
			//将已经到达的倒计时加入回收站
			err := AddCountDownRecycle(FDC, split[1])
			if err != nil {
				return err
			}
			logrus.Info("到达的倒计时加入回收站成功")
			continue
		}

		day, err := FdcCalculate(now, endTime, FDC, result["background"], result["name"])
		fmt.Println("keys:", FDC)
		if err != nil {
			return err
		}
		logrus.Info("同步成功，剩余时间: ", math.Ceil(day))
	}
	return nil
}

// RefOEC OEC刷新
// 从redis中读取数据计算过期日期
func RefOEC() error {
	// 当前时间戳
	now := time.Now().Unix()
	// 查询redis中OEC的数据
	FDCKeys, _, err := Cache.Scan(context.Background(), 0, OECCountdownPrefix+"*", 50).Result()
	if err != nil {
		logrus.Error("查询redis中OEC的数据失败", err)
		return fmt.Errorf(err.Error())
	}

	for _, OEC := range FDCKeys {
		// 获取当前OEC key里面的全部字段，返回一个字符串map
		result, err := Cache.HGetAll(context.Background(), OEC).Result()
		if err != nil {
			return fmt.Errorf("查询redis中OEC的数据失败: %v", err)
		}
		// 转换为int64
		startTime, _ := strconv.ParseInt(result["startTime"], 10, 64)
		// 计算过去时间
		day, err := OecCalculate(now, startTime, OEC, result["background"], result["name"])
		if err != nil {
			return err
		}
		logrus.Info("同步成功，已过去: ", math.Ceil(day))
	}
	return nil
}
