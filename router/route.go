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
	// 倒计时
	countdown := httpServer.Group("/countdown")
	{
		countdown.POST("/create", control.CreateCountdown)
		countdown.DELETE("/del", control.DelCountdown)
		countdown.GET("/list", control.ListCountDown)
		countdown.PUT("/modify", control.ModifyCountDown)
		countdown.GET("/search", control.SearchCountDown)
		countdown.POST("/upload", control.UploadBackground)
		countdown.GET("/detail/:identity", control.DetailCountDown)
	}
	// 回收站
	recycle := httpServer.Group("/recycle")
	{
		recycle.GET("/listCountDown", control.RecycleListCountDown)
		recycle.POST("/recoverCountDown", control.RecoverCountDown)
		recycle.GET("/recoverCountDown", control.RecoverCountDown)
	}
	// 分类
	category := httpServer.Group("/category")
	{
		category.POST("/create", control.CreateCategory)
		category.GET("/list", control.ListCategory)
	}
	return httpServer
}
