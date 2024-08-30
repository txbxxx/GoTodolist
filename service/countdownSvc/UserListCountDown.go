/**
 * @Author tanchang
 * @Description 列出倒计时
 * @Date 2024/8/30 16:23
 * @File:  UserListCountDown
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

type UserListCountDownService struct {
}

func (svc UserListCountDownService) List() gin.H {
	// 从redis中读取countdown信息
	keys, _, err := utils.Cache.Scan(context.Background(), 0, "countdown:FDC:*", 50).Result()
	if err != nil {
		logrus.Error("查询redis中Countdown的数据失败", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	countdownList := make([]map[string]string, len(keys))
	for _, countdown := range keys {
		result := utils.Cache.HGetAll(context.Background(), countdown)
		if err := result.Err(); err != nil {
			logrus.Error("查询redis中Countdown的数据失败", err)
			return gin.H{
				"code": -1,
				"msg":  "系统繁忙请稍后再试",
			}
		}
		countdownList = append(countdownList, result.Val())
	}
	return gin.H{
		"code": 200,
		"msg":  "获取倒计时列表成功！",
		"data": serializes.CountdownSerializeList(countdownList),
	}
}
