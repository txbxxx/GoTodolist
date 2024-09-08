/**
 * @Author tanchang
 * @Description 列出当前用户的类别
 * @Date 2024/9/7 21:27
 * @File:  UserListCategory
 * @Software: GoLand
 **/

package categorySvc

import (
	"GoToDoList/model"
	serializes "GoToDoList/serialized"
	"GoToDoList/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type UserListCategoryService struct {
}

// List 列出类别
// 如果redis中查询不到，则从mysql中查询
func (svc UserListCategoryService) List(token string) gin.H {
	ctx := context.Background()
	// 解析token
	user, err := utils.AnalyseToken(token)
	if err != nil {
		logrus.Error("Token 解析错误：", err.Error())
		return gin.H{"code": -1, "msg": "登录错误"}
	}
	// 查询当前用户的类别列表
	// 先从redis中查取，如果redis中查询不到，在从mysql中查并同步到redis中
	keys, err := GetCategoryFormRedis(ctx, user)
	if err != nil {
		logrus.Error("从redis查询类别列表失败", err.Error())
		return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
	}
	// 获取倒计时个数
	result, err := utils.Cache.Get(ctx, user.Name+":category_num").Result()
	// 将字符串转换为int
	CategoryNum, _ := strconv.Atoi(result)
	if len(keys) != CategoryNum {
		//从mysql中读取
		categoryMap, err := GetCategoryFormMysql(user)
		if err != nil {
			logrus.Error("从mysql中查询类别列表失败", err.Error())
			return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
		}
		return gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": serializes.CategorySerializeListFromModel(categoryMap),
		}
	}
	categoryMap, err := utils.ListFormRedis(ctx, keys)
	if err != nil {
		logrus.Error("获取分类的数据失败", err.Error())
		return gin.H{"code": -1, "msg": "系统繁忙请稍后再试"}
	}
	return gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": serializes.CategorySerializeList(categoryMap),
	}
}

// GetCategoryFormMysql 从mysql中获取分类数据并且同步至redis
func GetCategoryFormMysql(user *utils.UserClaims) ([]model.Category, error) {
	// 从mysql中查询数据后直接放回从mysql查询的数据
	countdownList, err := RefCategoryForMysql(user.Identity, user.Name)
	if err != nil {
		return nil, fmt.Errorf("mysql: %v", err)
	}
	// 查询到了则返回
	return countdownList, nil
}

// GetCategoryFormRedis 从redis中获取分类信息
func GetCategoryFormRedis(ctx context.Context, user *utils.UserClaims) ([]string, error) {
	keys, _, err := utils.Cache.Scan(ctx, 0, user.Name+":category:*", 100).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// RefCategoryForMysql 同步全量分类数据至redis
func RefCategoryForMysql(userIdentity, name string) ([]model.Category, error) {
	ctx := context.Background()
	// 从mysql查询数据
	categoryList := make([]model.Category, 1)
	utils.DB.Model(&model.Category{}).Where("user_identity = ?", userIdentity).Find(&categoryList)
	for i, category := range categoryList {
		key := name + ":category:" + category.Identity
		// 同步至redis
		if _, err := utils.Cache.HSet(ctx, key, map[string]any{"name": category.Name, "identity": category.Identity, "cover": category.Cover}).Result(); err != nil {
			return nil, err
		}
		// 设置过期时间30 + 遍历索引分钟数分钟
		duration := time.Hour/2 + time.Duration(i)*time.Minute
		utils.Cache.Expire(ctx, key, duration)
	}
	return categoryList, nil
}
