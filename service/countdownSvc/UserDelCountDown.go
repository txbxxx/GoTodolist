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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserDelCountDownService struct {
	Identity string `form:"identity" json:"identity" binding:"required"`
}

// Del 删除倒计时
func (svc *UserDelCountDownService) Del(token string) gin.H {
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
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
	key := user.Name + ":countdown:*:" + countdown.Identity
	keys, _ := utils.Cache.Scan(context.Background(), 0, key, 10).Val()
	err = utils.AddCountDownRecycle(keys[0], countdown.Identity)
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

// DelCountDownForRedis 从redis中删除一条数据
// 从redis中删除数据，并不加入回收站
func DelCountDownForRedis(userName, identity string) error {
	keys, _, err := utils.Cache.Scan(context.Background(), 0, userName+":countdown:*:"+identity, 30).Result()
	if err != nil {
		return err
	}
	// 如果数据不存在redis
	if len(keys) == 0 {
		return nil
	}
	// 删除
	err = utils.Cache.Del(context.Background(), keys[0]).Err()
	if err != nil {
		return fmt.Errorf("删除redis数据失败: %v", err)
	}
	return nil
}
