/**
 * @Author tanchang
 * @Description 用户controller
 * @Date 2024/7/11 21:12
 * @File:  User
 * @Software: GoLand
 **/

package control

import (
	"GoToDoList/service/userSvc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Login 登录
func Login(c *gin.Context) {
	var service userSvc.UserLoginService
	err := c.ShouldBind(&service)
	if err == nil {
		login := service.Login()
		c.JSON(200, login)
	} else {
		logrus.Error("绑定数据错误(ShouldBind): ", err.Error())
		c.JSON(200, gin.H{"code": -1, "err": err.Error()})
	}
}

// Register 注册
func Register(c *gin.Context) {
	var service userSvc.UserRegisterService
	err := c.ShouldBind(&service)
	if err == nil {
		register := service.Register()
		c.JSON(200, register)
	} else {
		logrus.Error("绑定数据错误(ShouldBind): ", err.Error())
		c.JSON(200, gin.H{"code": -1, "err": err.Error()})
	}
}
