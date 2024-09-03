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
	startTime, endTime := svc.StartTime.Unix(), svc.EndTime.Unix()
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
		return svc.txCreate(tx, err, newCountdown)
	})

	if err != nil {
		logrus.Error("创建倒计时错误: ", err)
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

// txCreate 事务处理创建并同比至redis
func (svc *UserCreateCountDownService) txCreate(tx *gorm.DB, err error, newCountdown model.CountDown) error {
	//插入数据库
	if err = tx.Create(&newCountdown).Error; err != nil {
		logrus.Error("创建倒计时错误: ", err)
		return err
	}

	// 同步至redis
	countdownModel := "FDC"
	// 如果没有填写endTime的就是OEC(那么endTime就是int64的最小数)模式填写了就是FDC
	if newCountdown.EndTime < 0 {
		// OEC模式
		countdownModel = "OEC"
		// key用countdown:OEC:{{ Identity }}
		// 这里需要同步初始时间即可，day表示当前时间和初始时间的差值
		key := "countdown:" + countdownModel + ":" + newCountdown.Identity
		// 计算过去时间oec
		if _, err := utils.OecCalculate(newCountdown.StartTime, newCountdown.StartTime, key, newCountdown.Background, newCountdown.Name); err != nil {
			return fmt.Errorf("同步至redis失败: %w", err)
		}
	} else {
		// 判断终止时间是否大于开始时间
		if svc.EndTime.Unix() <= svc.StartTime.Unix() {
			return fmt.Errorf("终止时间必须大于开始时间: %w", err)
		}
		key := "countdown:" + countdownModel + ":" + newCountdown.Identity
		// FDC
		if _, err := utils.FdcCalculate(newCountdown.StartTime, newCountdown.EndTime, key, newCountdown.Background, newCountdown.Name); err != nil {
			return fmt.Errorf("同步至redis失败: %w", err)
		}
	}
	return nil
}
