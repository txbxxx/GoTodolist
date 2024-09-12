/**
 * @Author tanchang
 * @Description 列出回收站内倒计时数据
 * @Date 2024/9/3 16:57
 * @File:  UserRecycleCountDownList
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

type UserRecycleCountDownListService struct{}

func (svc UserRecycleCountDownListService) List() gin.H {
	ctx := context.Background()
	// 从redis的回收站中获取删除倒计时信息
	keys, _, err := utils.Cache.Scan(ctx, 0, "*"+utils.DELCountdownPrefix+"countdown:*", 50).Result()
	if err != nil {
		logrus.Error("获取回收站数据失败: ", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
	}
	deleteList, err := utils.ListFormRedis(ctx, keys)
	if err != nil {
		logrus.Error("获取回收站数据失败: ", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
	}
	return gin.H{
		"code": 200,
		"msg":  "获取回收站倒计时数据成功",
		"data": serializes.CountdownSerializeList(deleteList),
	}
}
