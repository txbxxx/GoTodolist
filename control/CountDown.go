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
