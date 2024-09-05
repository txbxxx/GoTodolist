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
	// 从redis的回收站中获取删除倒计时信息
	keys, _, err := utils.Cache.Scan(context.Background(), 0, utils.DELCountdownPrefix+"countdown:*", 50).Result()
	if err != nil {
		logrus.Error("获取回收站数据失败: ", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
	}
	var deleteList []map[string]string
	for _, countdown := range keys {
		// 从redis读取数据
		result, err := utils.Cache.HGetAll(context.Background(), countdown).Result()
		if err != nil {
			logrus.Error("获取redis Key数据失败: ", err)
			return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
		}
		deleteList = append(deleteList, result)
	}
	return gin.H{
		"code": 200,
		"msg":  "获取回收站倒计时数据成功",
		"data": serializes.CountdownSerializeList(deleteList),
	}
}
