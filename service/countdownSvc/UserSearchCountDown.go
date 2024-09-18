/**
 * @Author tanchang
 * @Description 查询倒计时
 * @Date 2024/9/3 21:30
 * @File:  UserSearchCountDown
 * @Software: GoLand
 **/

package countdownSvc

import (
	serializes "GoToDoList/serialized"
	"GoToDoList/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type UserSearchCountDownService struct {
	Name string `json:"name" form:"name" binding:"max=10"`
	Day  int    `json:"day" form:"day"`
}

// TODO 从数据库中查询并同步到redis中

// Search 从redis中搜索
func (svc UserSearchCountDownService) Search(token string) gin.H {
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	ctx := context.Background()
	// 使用Scan从redis里面读取倒计时中的全部信息
	keys, _, err := utils.Cache.Scan(ctx, 0, user.Name+":countdown:*", 50).Result()
	if err != nil {
		logrus.Error("UserSearchCountDownService: 从redis查找所有倒计时数据失败", err)
		return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
	}
	// 判断是否查询到消息
	if len(keys) == 0 {

	}
	//顺序便利消息
	countdownList := make([]map[string]string, 0)
	for _, countdown := range keys {
		// 从redis中读取数据
		result, err := utils.Cache.HGetAll(ctx, countdown).Result()
		if err != nil {
			logrus.Error("UserSearchCountDownService: 从redis查找数据失败 ", err)
			return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
		}
		if result == nil {
			continue
		}
		// 判断map中是否存在此key
		if _, ok := result["day"]; !ok {
			continue
		}
		if _, ok := result["name"]; !ok {
			continue
		}
		// 将day转换为int
		day, err := strconv.Atoi(result["day"])
		if err != nil {
			logrus.Error("UserSearchCountDownService: 字符串转换为int失败", err)
			return gin.H{"code": -1, "msg": "系统繁忙请稍后在试"}
		}
		// 判断是否为空
		if len(svc.Name) != 0 {
			// 判断是否包含svc.Name
			if strings.Contains(result["name"], svc.Name) || day == svc.Day {
				countdownList = append(countdownList, result)
				continue
			}
		} else if svc.Day == day {
			// 满足天数
			countdownList = append(countdownList, result)
			continue
		}
	}
	return gin.H{
		"code": 200,
		"msg":  "查询成功！",
		"data": serializes.CountdownSerializeList(countdownList),
	}
}
