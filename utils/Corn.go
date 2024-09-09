/**
 * @Author tanchang
 * @Description 定时任务
 * @Date 2024/8/29 21:11
 * @File:  corn
 * @Software: GoLand
 **/

package utils

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"time"
)

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
	// 每小时执行一次,从数据库中读取数据
	//if _, err := c.AddFunc("*/1 * * * *", func() { Run(RefreshDayForMysql()) }); err != nil {
	//	logrus.Error("将倒计时从同步至redis错误: ", err)
	//}
	c.Start()
	fmt.Println("定时任务启动成功")
	// 阻塞协程
	select {}
}
