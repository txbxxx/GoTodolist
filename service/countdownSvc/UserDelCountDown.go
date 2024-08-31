/**
 * @Author tanchang
 * @Description 手动删除倒计时
 * @Date 2024/8/30 16:09
 * @File:  UserDelCountDown
 * @Software: GoLand
 **/

package countdownSvc

import (
	"GoToDoList/model"
	"GoToDoList/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserDelCountDownService struct {
	Identity string `form:"identity" json:"identity" binding:"required"`
}

func (svc *UserDelCountDownService) Del() gin.H {
	// 查询倒计时是否存在
	var countdown model.CountDown
	if err := utils.DB.Model(&model.CountDown{}).Where("identity = ?", svc.Identity).Take(&countdown).Error; err != nil {
		logrus.Error("查询倒计时失败", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	if countdown.Identity == "" {
		return gin.H{
			"code": -1,
			"msg":  "倒计时不存在",
		}
	}
	// 存在则删除
	// 将redis中同步的此倒计时的数据加入delete回收站
	// 查询当前删除的数据
	keys, _ := utils.Cache.Scan(context.Background(), 0, "countdown:*:"+countdown.Identity, 10).Val()
	err := utils.AddCountDownRecycle(keys[0], countdown.Identity)
	if err != nil {
		logrus.Error("到达的倒计时加入回收站失败，", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	return gin.H{
		"code": 200,
		"msg":  "删除成功倒计时成功！！",
	}
}
