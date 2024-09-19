/**
 * @Author tanchang
 * @Description 添加倒计时功能
 * @Date 2024/8/29 12:54
 * @File:  UserCreateCountDown
 * @Software: GoLand
 **/

package countdownSvc

import (
	"GoToDoList/model"
	"GoToDoList/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserCreateCountDownService struct {
	Name             string    `json:"name" form:"name" binding:"required,max=10"`
	EndTime          time.Time `json:"endTime" form:"endTime" time_format:"2006-01-02 15:04:05"`
	StartTime        time.Time `json:"startTime" form:"startTime" binding:"required" time_format:"2006-01-02 15:04:05"`
	Background       string    `json:"background" form:"background"`
	CategoryIdentity string    `json:"categoryIdentity" form:"categoryIdentity" binding:"required"`
}

func (svc *UserCreateCountDownService) Create(token string) gin.H {
	data := model.CountDown{}
	// 查找当前分类是否有相同倒计时存在
	if err := utils.DB.Model(&data).Take(&data, "name = ? AND category_identity = ?", svc.Name, svc.CategoryIdentity).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return gin.H{"code": -1, "msg": "倒计时已存在"}
		}
	}
	// 创建对象前置操作
	startTime, endTime := svc.StartTime.Unix(), svc.EndTime.Unix()
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	// 不存在则创建对象
	newCountdown := model.CountDown{
		Identity:         utils.GenerateUUID(),
		Name:             svc.Name,
		StartTime:        startTime,
		EndTime:          endTime,
		Background:       "",
		CategoryIdentity: svc.CategoryIdentity,
	}
	// 开启事务
	if err := utils.DB.Transaction(func(tx *gorm.DB) error {
		return svc.txCreate(tx, newCountdown, user.Name)
	}); err != nil {
		logrus.Error("创建倒计时错误: ", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
	}
	// 创建成功后将identity添加到sorted set中
	utils.Cache.ZAdd(context.Background(), user.Name+":isMysql:countdown", &redis.Z{
		Score:  1,
		Member: newCountdown.Name,
	})
	return gin.H{"code": 200, "msg": "创建成功倒计时成功！！"}
}

// txCreate 事务处理创建并同比至redis
func (svc *UserCreateCountDownService) txCreate(tx *gorm.DB, newCountdown model.CountDown, name string) error {
	//插入数据库
	if err := tx.Create(&newCountdown).Error; err != nil {
		logrus.Error("创建倒计时错误: ", err)
		return err
	}
	if err := isOecORFdcModel(newCountdown, name); err != nil {
		return err
	}
	//添加成功则添+1
	utils.Cache.IncrBy(context.Background(), name+":countdown_num", 1)
	return nil
}

func isOecORFdcModel(countdown model.CountDown, name string) error {
	now := time.Now().Unix()
	// 如果没有填写endTime的就是OEC(那么endTime就是int64的最小数)模式填写了就是FDC
	if countdown.EndTime < 0 {
		// OEC模式
		// key用countdown:OEC:{{ Identity }}
		// 这里需要同步初始时间即可，day表示当前时间和初始时间的差值
		key := name + ":" + utils.OECCountdownPrefix + countdown.Identity
		// 计算过去时间oec
		if err := utils.OecCalculate(now, countdown, key); err != nil {
			return fmt.Errorf("同步至redis失败: %w", err)
		}
	} else {
		// 判断终止时间是否大于开始时间
		// 将int64转换为time.Time
		if countdown.EndTime <= countdown.StartTime {
			return fmt.Errorf("终止时间必须大于开始时间")
		}
		key := name + ":" + utils.FDCCountdownPrefix + countdown.Identity
		// 判断当前日期时间戳是否大于结束日期时间戳
		if now >= countdown.EndTime {
			// 大于则执行
			err := utils.AddCountDownRecycle(key, countdown.Identity)
			if err != nil {
				return fmt.Errorf("添加至回收站失败: %w", err)
			}
			logrus.Info("到达的倒计时加入回收站成功")
		}
		// FDC
		if err := utils.FdcCalculate(now, countdown, key); err != nil {
			return fmt.Errorf("同步至redis失败: %w", err)
		}
	}
	return nil
}
