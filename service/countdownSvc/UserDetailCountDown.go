/**
 * @Author tanchang
 * @Description 单个倒计时详情
 * @Date 2024/9/7 18:58
 * @File:  UserDetailCountDown
 * @Software: GoLand
 **/

package countdownSvc

import (
	serializes "GoToDoList/serialized"
	"GoToDoList/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserDetailCountDownService struct {
}

func (svc UserDetailCountDownService) Detail(identity string) gin.H {
	ctx := context.Background()
	// 从redis中查询此倒计时数据
	keys, _, err := utils.Cache.Scan(ctx, 0, "countdown:*:"+identity, 1).Result()
	if err != nil {
		logrus.Error("查询redis中Countdown的数据失败", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
	}
	// 获取次倒计时key中的value
	countdown, err := utils.Cache.HGetAll(ctx, keys[0]).Result()
	if err != nil {
		logrus.Error("获取Countdown的数据失败", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
	}
	return gin.H{
		"code": 200,
		"msg":  "获取成功！",
		"data": serializes.CountdownSerializeSingle(countdown),
	}
}
