/**
 * @Author tanchang
 * @Description  从删除的数据中恢复倒计时
 * @Date 2024/9/3 15:49
 * @File:  recoverForMysql
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
)

type UserRecoverCountDownService struct {
	Identity string `json:"identity" form:"identity"`
}

// RecoverCountDown 恢复倒计时数据
func (svc *UserRecoverCountDownService) RecoverCountDown() gin.H {
	// 先查找被删除的数据
	if err := utils.DB.Unscoped().Model(&model.CountDown{}).Where("identity = ?", svc.Identity).Update("DeletedAt", nil).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error("RecoverCountDown: 查找被删除的数据失败", err)
			return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}

		}
		logrus.Error("RecoverCountDown: 恢复失败", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
	}
	// 从回收站删除
	if err := utils.DeleteForRecycle(svc.Identity); err != nil {
		logrus.Error(err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
	}
	// 从数据中同步至redis
	if err := utils.RefreshDayForMysql(); err != nil {
		logrus.Error("RecoverCountDown: 从数据中同步至redis失败", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
	}
	return gin.H{
		"code": 200,
		"msg":  "恢复倒计时成功",
	}
}
