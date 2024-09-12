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
	"gorm.io/gorm"
	"strconv"
	"time"
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
	keys, _, err := utils.Cache.Scan(ctx, 0, user.Name+":countdown:*", 100).Result()
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
	if len(keys) != countdownNum {
		countdown, err := RefreshDayForMysql(user.Name)
		if err != nil {
			logrus.Error(err)
			return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
		}
		return gin.H{
			"code": 200,
			"msg":  "获取倒计时列表成功！",
			"data": serializes.CountdownSerializeListModel(countdown),
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
	if err := utils.DB.Preload("Category.CountDown").Take(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询倒计时失败: %v", err)
		}
	}
	// 遍历category拿出countdown
	for _, category := range user.Category {
		countdown = append(countdown, category.CountDown...)
	}
	// 当前时间戳
	now := time.Now().Unix()
	for _, count := range countdown {
		key := userName + ":" + utils.OECCountdownPrefix + count.Identity
		if count.EndTime <= 0 {
			// 计算过去时间oec
			err := utils.OecCalculate(now, count, key)
			if err != nil {
				return nil, err
			}
		} else {
			key = userName + ":" + utils.FDCCountdownPrefix + count.Identity
			// 判断当前日期时间戳是否大于结束日期时间戳
			if now >= count.EndTime {
				// 大于则执行
				err := utils.AddCountDownRecycle(key, count.Identity)
				if err != nil {
					return nil, err
				}
				logrus.Info("到达的倒计时加入回收站成功")
				continue
			}
			//FDC
			// 如果没有大于，就计算还有多少天，使用结束时间减去现在时间
			if err := utils.FdcCalculate(now, count, key); err != nil {
				return nil, err
			}
		}
	}
	return countdown, nil
}
