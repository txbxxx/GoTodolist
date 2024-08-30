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
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserCreateCountDownService struct {
	Name       string    `json:"name" form:"name" binding:"required,max=10"`
	EndTime    time.Time `json:"endTime" form:"endTime" time_format:"2006-01-02 15:04:05"`
	StartTime  time.Time `json:"startTime" form:"startTime" binding:"required" time_format:"2006-01-02 15:04:05"`
	Background string    `json:"background" form:"background"`
}

func (svc *UserCreateCountDownService) Create(token string) gin.H {
	data := model.CountDown{}
	//查找是否有相同倒计时存在
	if err := utils.DB.Model(&data).Take(&data, "name = ?", svc.Name).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return gin.H{
				"code": -1,
				"msg":  "系统繁忙请稍后再试",
			}
		}
	}

	if data.Name != "" {
		return gin.H{
			"code": -1,
			"msg":  "倒计时已存在",
		}
	}

	// 创建对象前置操作
	startTime := svc.StartTime.Unix()
	endTime := svc.EndTime.Unix()
	// 解析Token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{
			"code": -1,
			"msg":  "登录错误",
		}
	}

	// 不存在则创建对象
	newCountdown := model.CountDown{
		Identity:     utils.GenerateUUID(),
		Name:         svc.Name,
		StartTime:    startTime,
		EndTime:      endTime,
		Background:   "",
		UserIdentity: user.Identity,
	}

	// 开启事务
	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		//插入数据库
		if err = tx.Create(&newCountdown).Error; err != nil {
			logrus.Error("创建倒计时错误: ", err)
			return err
		}

		// 同步至redis
		countdownModel := "FDC"
		// 如果没有填写endTime的就是OEC(那么endTime就是int64的最小数)模式填写了就是FDC
		if newCountdown.EndTime <= 0 {
			// OEC模式
			countdownModel = "OEC"
			// key用countdown:OEC:{{ Identity }}
			// 这里需要同步初始时间即可，day表示当前时间和初始时间的差值
			key := "countdown:" + countdownModel + ":" + newCountdown.Identity
			if is := utils.Cache.HMSet(context.Background(), key, map[string]any{"startTime": newCountdown.StartTime, "day": 0}); !is.Val() {
				return fmt.Errorf("同步至redis失败")
			}
		} else {
			// FDC
			key := "countdown:" + countdownModel + ":" + newCountdown.Identity
			if is := utils.Cache.HMSet(context.Background(), key, map[string]any{"endTime": newCountdown.EndTime, "day": 0}); !is.Val() {
				return fmt.Errorf("同步至redis失败")
			}
		}
		return nil
	})

	if err != nil {
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}

	return gin.H{
		"code": 200,
		"msg":  "创建成功倒计时成功！！",
	}
}
