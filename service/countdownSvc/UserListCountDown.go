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
	ctx := context.Background()
	// 从redis中读取countdown信息
	keys, _, err := utils.Cache.Scan(ctx, 0, "countdown:*", 50).Result()
	if err != nil {
		logrus.Error("查询redis中Countdown的数据失败", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	// 遍历keys,HGetAll返回map[string]string
	countdownList, err := utils.ListFormRedis(ctx, keys)
	if err != nil {
		return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
	}
	return gin.H{
		"code": 200,
		"msg":  "获取倒计时列表成功！",
		"data": serializes.CountdownSerializeList(countdownList),
	}
}
