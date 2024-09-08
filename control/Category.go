/**
 * @Author tanchang
 * @Description 处理分类
 * @Date 2024/9/7 20:40
 * @File:  Category
 * @Software: GoLand
 **/

package control

import (
	"GoToDoList/service/categorySvc"
	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var svc categorySvc.UserCreateCategoryService
	if err := c.ShouldBind(&svc); err == nil {
		c.JSON(200, svc.Create(c.GetHeader("token")))
	} else {
		c.JSON(200, err.Error())
	}
}

func ListCategory(c *gin.Context) {
	var svc categorySvc.UserListCategoryService
	if err := c.ShouldBind(&svc); err == nil {
		c.JSON(200, svc.List(c.GetHeader("token")))
	} else {
		c.JSON(200, err.Error())
	}
}
