/**
 * @Author tanchang
 * @Description 定时任务
 * @Date 2024/8/29 21:11
 * @File:  corn
 * @Software: GoLand
 **/

package utils

import (
	"GoToDoList/model"
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// RefreshDayMysql 从mysql中读取数据刷新倒计时
func RefreshDayMysql() error {
	fmt.Println("开始刷新")
	countdown := make([]model.CountDown, 1)
	if err := DB.Model(&model.CountDown{}).Find(&countdown).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error("查询倒计时失败", err)
			return err
		}
	}
	// 当前时间戳
	now := time.Now().Unix()
	for _, count := range countdown {
		if count.EndTime <= 0 {
			// 计算剩余时间
			day := float64(now-count.StartTime) / 86400
			key := "countdown:OEC:" + count.Identity
			// 将倒计时同步至redis，时间则向上取整
			set := Cache.HMSet(context.Background(), key, "day", math.Ceil(day))
			if !set.Val() {
				logrus.Error("同步redis失败")
				return fmt.Errorf("同步redis失败")
			}
			logrus.Info("同步成功，剩余时间: ", math.Ceil(day))
		} else {
			key := "countdown:FDC:" + count.Identity
			// 判断当前日期时间戳是否大于结束日期时间戳
			if now >= count.EndTime {
				//将已经到达的倒计时加入回收站
				rename := Cache.Rename(context.Background(), key, "delete:"+key)
				if rename.Err() != nil {
					logrus.Error("到达的倒计时加入回收站失败")
					return fmt.Errorf("将已经到达的倒计时加入回收站失败")
				}
				// 删除sql数据
				err := DB.Model(&model.CountDown{}).Delete(&model.CountDown{Identity: count.Identity}).Error
				if err != nil {
					logrus.Error("删除sql数据失败")
					return fmt.Errorf("删除sql数据失败")
				}
				logrus.Info("到达的倒计时加入回收站成功")
				continue
			}
			// 如果没有大于，就计算还有多少天，使用结束时间减去现在时间
			day := float64(count.EndTime-now) / 86400
			// 将倒计时同步至redis，时间则向上取整
			set := Cache.HMSet(context.Background(), key, "day", math.Ceil(day))
			if !set.Val() {
				logrus.Error("同步redis失败")
				return fmt.Errorf("同步redis失败")
			}
			logrus.Info("同步成功，剩余时间: ", math.Ceil(day))
		}
	}
	return nil
}

// RefFDC FDC刷新
func RefFDC() error {
	// 当前时间戳
	now := time.Now().Unix()
	// 查询redis中FDC和OEC的数据
	FDCKeys, _, err := Cache.Scan(context.Background(), 0, "countdown:FDC:*", 50).Result()
	if err != nil {
		logrus.Error("查询redis中FDC和OEC的数据失败", err)
		return fmt.Errorf(err.Error())
	}

	for _, FDC := range FDCKeys {
		result := Cache.HMGet(context.Background(), FDC, "endTime").Val()
		// 转换为int64
		s := result[0].(string)
		//取出identity
		split := strings.Split(s, "countdown:FDC:")
		endTime, _ := strconv.ParseInt(s, 10, 64)
		// 判断当前日期时间戳是否大于结束日期时间戳
		if now >= endTime {
			//将已经到达的倒计时加入回收站
			if rename := Cache.Rename(context.Background(), FDC, "delete:"+FDC); rename.Err() != nil {
				logrus.Error("到达的倒计时加入回收站失败")
				return fmt.Errorf("将已经到达的倒计时加入回收站失败")
			}
			// 删除sql数据
			if err := DB.Model(&model.CountDown{}).Delete(&model.CountDown{Identity: split[0]}).Error; err != nil {
				logrus.Error("删除sql数据失败")
				return fmt.Errorf("删除sql数据失败")
			}
			logrus.Info("到达的倒计时加入回收站成功")
			continue
		}
		// 如果没有大于，就计算还有多少天，使用结束时间减去现在时间
		day := float64(endTime-now) / 86400
		// 将倒计时同步至redis，时间则向上取整
		if set := Cache.HMSet(context.Background(), FDC, "day", math.Ceil(day)); !set.Val() {
			logrus.Error("同步redis失败")
			return fmt.Errorf("同步redis失败")
		}
		logrus.Info("同步成功，剩余时间: ", math.Ceil(day))
	}
	return nil
}

// RefOEC OEC刷新
func RefOEC() error {
	// 当前时间戳
	now := time.Now().Unix()
	// 查询redis中OEC的数据
	FDCKeys, _, err := Cache.Scan(context.Background(), 0, "countdown:OEC:*", 50).Result()
	if err != nil {
		logrus.Error("查询redis中OEC的数据失败", err)
		return fmt.Errorf(err.Error())
	}

	for _, OEC := range FDCKeys {
		result := Cache.HMGet(context.Background(), OEC, "startTime").Val()
		// 转换为int64
		s := result[0].(string)
		startTime, _ := strconv.ParseInt(s, 10, 64)
		day := float64(now-startTime) / 86400
		// 将倒计时同步至redis，时间则向上取整
		if set := Cache.HMSet(context.Background(), OEC, "day", math.Ceil(day)); !set.Val() {
			logrus.Error("同步redis失败")
			return fmt.Errorf("同步redis失败")
		}
		logrus.Info("同步成功，已过去: ", math.Ceil(day))
	}
	return nil
}

// Run 运行
func Run(job func() error) {
	from := time.Now().UnixNano()
	err := job()
	to := time.Now().UnixNano()
	jobName := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()
	if err != nil {
		fmt.Printf("%s error: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	} else {
		fmt.Printf("%s success: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	}
}

func CronJob() {
	c := cron.New()
	defer c.Stop()
	// 每分钟执行一次
	if _, err := c.AddFunc("*/1 * * * *", func() { Run(RefFDC) }); err != nil {
		logrus.Error("将倒计时同步至redis错误: ", err)
	}
	if _, err := c.AddFunc("*/1 * * * *", func() { Run(RefOEC) }); err != nil {
		logrus.Error("将倒计时同步至redis错误: ", err)
	}
	c.Start()
	fmt.Println("定时任务启动成功")
	// 阻塞协程
	select {}
}
