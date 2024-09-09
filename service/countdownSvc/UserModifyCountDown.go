/**
 * @Author tanchang
 * @Description 修改倒计时
 * @Date 2024/8/30 16:08
 * @File:  UserModifyCountDown
 * @Software: GoLand
 **/

package countdownSvc

import (
	"GoToDoList/model"
	"GoToDoList/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

// UserModifyCountDownService  修改倒计时
type UserModifyCountDownService struct {
	Identity   string    `json:"identity" form:"identity" binding:"required"`
	Name       string    `json:"name" form:"name" binding:"max=10"`
	EndTime    time.Time `json:"endTime" form:"endTime" time_format:"2006-01-02 15:04:05"`
	StartTime  time.Time `json:"startTime" form:"startTime" time_format:"2006-01-02 15:04:05"`
	Background string    `json:"background" form:"background"`
}

// TODO 修改数据时如果同步到redis，但是刚好在准备改的时候有节点可能会读取，这样就读取的是旧数据了

func (svc *UserModifyCountDownService) Modify(token string) gin.H {
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	// 查询是否存在于数据库中
	countdown := &model.CountDown{}
	var count int64
	if err := utils.DB.Model(&model.CountDown{}).Where("identity = ?", svc.Identity).Take(countdown).Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gin.H{
				"code": -1,
				"msg":  "倒计时不存在",
			}
		}
		logrus.Error("查询倒计时失败", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	// 修改倒计时
	countdown.Name, countdown.EndTime, countdown.StartTime, countdown.Background = svc.Name, svc.EndTime.Unix(), svc.StartTime.Unix(), svc.Background
	// 判断终止时间是否大于开始时间
	if svc.EndTime.Unix() <= svc.StartTime.Unix() {
		logrus.Error("终止时间必须大于开始时间", err)
		return gin.H{
			"code": -1,
			"msg":  "终止时间必须大于开始时间",
		}
	}
	// 保存
	if err := utils.DB.Save(countdown).Error; err != nil {
		logrus.Error("保存倒计时失败", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	// 删除原本同步在redis的数据
	if err := DelCountDownForRedis(svc.Identity); err != nil {
		logrus.Error(err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	// 同步至redis
	countdownModel := "FDC"
	// 如果没有填写endTime的就是OEC(那么endTime就是int64的最小数)模式填写了就是FDC
	if countdown.EndTime < 0 {
		// OEC模式
		countdownModel = "OEC"
		// key用countdown:OEC:{{ Identity }}
		// 这里需要同步初始时间即可，day表示当前时间和初始时间的差值
		key := user.Name + ":countdown:" + countdownModel + ":" + countdown.Identity
		// 计算过去时间oec
		if err := utils.OecCalculate(countdown.StartTime, countdown.StartTime, key, countdown.Background, countdown.Name, countdown.Identity); err != nil {
			logrus.Error("同步至redis失败", err)
			return gin.H{
				"code": -1,
				"msg":  "系统繁忙请稍后再试",
			}
		}
	} else {
		key := user.Name + ":countdown:" + countdownModel + ":" + countdown.Identity
		// FDC
		if err := utils.FdcCalculate(countdown.StartTime, countdown.StartTime, countdown.EndTime, key, countdown.Background, countdown.Name, countdown.Identity); err != nil {
			logrus.Error("同步至redis失败", err)
			return gin.H{
				"code": -1,
				"msg":  "系统繁忙请稍后再试",
			}
		}
	}
	return gin.H{
		"code": 200,
		"msg":  "修改成功",
	}
}
