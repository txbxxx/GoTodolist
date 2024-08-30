/**
 * @Author tanchang
 * @Description 路由设置
 * @Date 2024/7/11 15:28
 * @File:  route
 * @Software: GoLand
 **/

package router

import (
	"GoToDoList/control"
	"GoToDoList/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	httpServer := gin.Default()
	//跨域
	httpServer.Use(middleware.Cors())

	user := httpServer.Group("/user")
	{
		user.POST("/login", control.Login)
		user.POST("/register", control.Register)
	}
	countdown := httpServer.Group("/countdown")
	{
		countdown.POST("/createCountDown", control.CreateCountdown)
		countdown.POST("/delCountDown", control.DelCountdown)
	}
	return httpServer
}
