/**
 * @Author tanchang
 * @Description 处理回收站
 * @Date 2024/9/3 18:29
 * @File:  Recycle
 * @Software: GoLand
 **/

package control

import (
	"GoToDoList/service/countdownSvc"
	"github.com/gin-gonic/gin"
)

func RecycleListCountDown(c *gin.Context) {
	var svc countdownSvc.UserRecycleCountDownListService
	if err := c.ShouldBind(&svc); err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
	} else {
		c.JSON(200, svc.List())
	}
}

func RecoverCountDown(c *gin.Context) {
	var svc countdownSvc.UserRecoverCountDownService
	if err := c.ShouldBind(&svc); err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
	} else {
		c.JSON(200, svc.RecoverCountDown(c.GetHeader("Authorization")))
	}
}
