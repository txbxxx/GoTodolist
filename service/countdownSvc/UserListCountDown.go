/**
 * @Author tanchang
 * @Description 列出倒计时
 * @Date 2024/8/30 16:23
 * @File:  UserListCountDown
 * @Software: GoLand
 **/

package countdownSvc

import (
	"GoToDoList/model"
	serializes "GoToDoList/serialized"
	"GoToDoList/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"strconv"
)

type UserListCountDownService struct {
}

// List 列出所有倒计时
func (svc UserListCountDownService) List(token string) gin.H {
	ctx := context.Background()
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	// 从redis中读取countdown信息
	keys, _, err := utils.Cache.Scan(ctx, 0, user.Name+":countdown:*", 300).Result()
	if err != nil {
		logrus.Error("查询redis中Countdown的数据失败", err)
		return gin.H{
			"code": -1,
			"msg":  "系统繁忙请稍后再试",
		}
	}
	// 获取倒计时个数
	result, err := utils.Cache.Get(ctx, user.Name+":countdown_num").Result()
	// 将字符串转换为int
	countdownNum, _ := strconv.Atoi(result)
	fmt.Println(len(keys), countdownNum)
	if len(keys) != countdownNum {
		fmt.Println("Mysql")
		var group singleflight.Group
		countdown, err, _ := group.Do(user.Name, func() (interface{}, error) {
			countdown, err := RefreshDayForMysql(user.Name)
			if err != nil {
				return nil, err
			}
			return countdown, nil
		})
		if err != nil {
			logrus.Error("从mysql中读取数据失败", err)
			return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
		}
		return gin.H{
			"code": 200,
			"msg":  "获取倒计时列表成功！",
			"data": serializes.CountdownSerializeListModel(countdown.([]model.CountDown)),
		}
	}
	fmt.Println("Redis")
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

// GetCountDownByCategory 根据分类获取倒计时
// Param categoryIdentity 分类标识
// Param token token密钥
func (svc UserListCountDownService) GetCountDownByCategory(token, categoryIdentity string) gin.H {
	list := svc.List(token)
	if list["code"] != 200 {
		return list
	}
	var c []serializes.CountdownSerialize
	a, ok := list["data"].([]serializes.CountdownSerialize)
	if ok {
		for _, i := range a {
			if i.Category == categoryIdentity {
				c = append(c, i)
			}
		}
		list["data"] = c
	}
	return list
}

// RefreshDayForMysql 从mysql中读取数据刷新倒计时
// 从redis读取倒计时列表
// 将倒计时列表中的数据同步至redis
func RefreshDayForMysql(userName string) ([]model.CountDown, error) {
	countdown := make([]model.CountDown, 0)
	// 使用关联查询直接查询出当前用户下面的countdown
	var user model.User
	if err := utils.DB.Where("name = ?", userName).Preload("Category.CountDown").Take(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询倒计时失败: %v", err)
		}
	}
	// 遍历category拿出countdown
	for _, category := range user.Category {
		countdown = append(countdown, category.CountDown...)
	}
	for _, count := range countdown {
		if err := isOecORFdcModel(count, userName); err != nil {
			return nil, fmt.Errorf("同步至redis失败: %v", err)
		}
	}
	return countdown, nil
}
