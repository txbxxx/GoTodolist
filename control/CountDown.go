/**
 * @Author tanchang
 * @Description 倒计时功能controller
 * @Date 2024/8/29 12:42
 * @File:  CountDown
 * @Software: GoLand
 **/

package control

import (
	"GoToDoList/service/countdownSvc"
	"github.com/gin-gonic/gin"
)

// CreateCountdown 创建
func CreateCountdown(c *gin.Context) {
	var svc countdownSvc.UserCreateCountDownService
	err := c.ShouldBind(&svc)
	if err == nil {
		create := svc.Create(c.GetHeader("Token"))
		c.JSON(200, create)
	} else {
		c.JSON(200, gin.H{
			"error": err,
		})
	}
}

// DelCountdown 删除
func DelCountdown(c *gin.Context) {
	var svc countdownSvc.UserDelCountDownService
	err := c.ShouldBind(&svc)
	if err == nil {
		c.JSON(200, svc.Del())
	} else {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	}
}

// ListCountDown 列出倒计时
func ListCountDown(c *gin.Context) {
	var svc countdownSvc.UserListCountDownService
	err := c.ShouldBind(&svc)
	if err == nil {
		c.JSON(200, svc.List())
	} else {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	}
}

// ModifyCountDown 修改倒计时
func ModifyCountDown(c *gin.Context) {
	var svc countdownSvc.UserModifyCountDownService
	err := c.ShouldBind(&svc)
	if err == nil {
		c.JSON(200, svc.Modify())
	} else {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	}
}

func SearchCountDown(c *gin.Context) {
	var svc countdownSvc.UserSearchCountDownService
	err := c.ShouldBind(&svc)
	if err == nil {
		c.JSON(200, svc.Search())
	} else {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	}
}
