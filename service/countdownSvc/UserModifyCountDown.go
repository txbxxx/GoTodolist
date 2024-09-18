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
	"fmt"
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

func (svc *UserModifyCountDownService) Modify(token string) gin.H {
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	// 查询是否存在于数据库中
	countdown := model.CountDown{}
	var count int64
	if err := utils.DB.Model(&model.CountDown{}).Where("identity = ?", svc.Identity).Take(&countdown).Count(&count).Error; err != nil {
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
	// 保存
	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		return svc.txSave(countdown, user.Name)
	})
	if err != nil {
		logrus.Error("修改倒计时失败", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	return gin.H{
		"code": 200,
		"msg":  "修改成功",
	}
}

func (svc *UserModifyCountDownService) txSave(countdown model.CountDown, userName string) error {
	if err := utils.DB.Save(countdown).Error; err != nil {
		return fmt.Errorf("保存失败:%w", err)
	}
	// 删除原本同步在redis的数据
	if err := DelCountDownForRedis(userName, svc.Identity); err != nil {
		return fmt.Errorf("删除redis数据失败:%w", err)
	}
	return isOecORFdcModel(countdown, userName)
}
