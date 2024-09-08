/**
 * @Author tanchang
 * @Description 创建倒计时
 * @Date 2024/9/7 19:22
 * @File:  UserCreateCategory
 * @Software: GoLand
 **/

package categorySvc

import (
	"GoToDoList/model"
	"GoToDoList/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserCreateCategoryService struct {
	Name  string `json:"name" form:"name" binding:"required,max=10"`
	Cover string `json:"cover" form:"cover"`
}

func (svc UserCreateCategoryService) Create(token string) gin.H {
	var num int64
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	// 查询是否当前用户已有分类
	if err := utils.DB.Model(&model.Category{}).Where("name = ? AND user_identity = ?", svc.Name, user.Identity).Take(&model.Category{}).Count(&num).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("查询分类失败:", err)
		return gin.H{"code": -1, "msg": "请求繁忙请稍后再试"}
	}
	if num > 0 {
		return gin.H{"code": -1, "msg": "分类已存在"}
	}
	// 不存在则创建
	category := model.Category{
		Identity:     utils.GenerateUUID(),
		Name:         svc.Name,
		Cover:        svc.Cover,
		UserIdentity: user.Identity,
	}
	// 开启事务
	if err := utils.DB.Transaction(func(tx *gorm.DB) error { return txCreate(tx, category, user.Name) }); err != nil {
		logrus.Error(err)
		return gin.H{"code": -1, "msg": "请求繁忙请稍后再试"}
	}
	return gin.H{"code": 200, "msg": "创建分类成功!"}
}

func txCreate(tx *gorm.DB, category model.Category, name string) error {
	if err := tx.Create(&category).Error; err != nil {
		return fmt.Errorf("创建分类错误: %v", err)
	}
	// 同步至redis 使用set
	ctx := context.Background()
	key := name + ":category:" + category.Identity
	if _, err := utils.Cache.HSet(ctx, key, map[string]any{"name": category.Name, "identity": category.Identity, "cover": category.Cover}).Result(); err != nil {
		return fmt.Errorf("同步至redis错误: %v", err)
	}
	// 设置过期时间30分钟(初始)
	duration := time.Hour / 2
	utils.Cache.Expire(ctx, key, duration)
	//添加成功则添+1
	utils.Cache.IncrBy(ctx, name+":category_num", 1)
	return nil
}

// TODO 还没定时同步或者访问则同步
