/**
 * @Author tanchang
 * @Description 单个倒计时详情
 * @Date 2024/9/7 18:58
 * @File:  UserDetailCountDown
 * @Software: GoLand
 **/

package countdownSvc

import (
	"GoToDoList/model"
	serializes "GoToDoList/serialized"
	"GoToDoList/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

type UserDetailCountDownService struct {
}

// Detail 获取单个倒计时详情
// @param identity
// @param token
func (svc UserDetailCountDownService) Detail(identity, token string) gin.H {
	ctx := context.Background()
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	// 从redis中查询此倒计时数据
	keys, _, err := utils.Cache.Scan(ctx, 0, user.Name+":countdown:*:"+identity, 1).Result()
	if err != nil {
		logrus.Error("查询redis中Countdown的数据失败", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
	}
	if len(keys) != 0 {
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
	// 从mysql获取
	var group singleflight.Group
	countdown, err, _ := group.Do(identity, func() (interface{}, error) {
		var countdown model.CountDown
		if err = utils.DB.Model(&model.User{}).Preload("Category.CountDown").Where("identity = ?", identity).Take(&countdown).Error; err != nil {
			return nil, fmt.Errorf("获取Countdown的数据失败: %v", err)
		}
		// 同步至redis
		if err = isOecORFdcModel(countdown, user.Name); err != nil {
			return nil, fmt.Errorf("同步Countdown的数据失败: %v", err)
		}
		return countdown, nil
	})
	return gin.H{
		"code": 200,
		"msg":  "获取成功！",
		"data": serializes.CountdownSerializeSingleModel(countdown.(model.CountDown)),
	}
}
