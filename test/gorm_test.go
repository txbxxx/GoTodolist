/**
 * @Author tanchang
 * @Description 测试
 * @Date 2024/7/11 22:59
 * @File:  gorm_test
 * @Software: GoLand
 **/

package test

import (
	"GoToDoList/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func TestCreateUser(t *testing.T) {
	databases := "root:000000@tcp(127.0.0.1:3306)/todolist_test?charset=utf8mb4&parseTime=True&loc=Local"
	//配置数据库
	db, err := gorm.Open(mysql.Open(databases), &gorm.Config{
		SkipDefaultTransaction: false, //禁用事务
		NamingStrategy: schema.NamingStrategy{ //命名策略
			SingularTable: true, //禁用复数名称
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	sqlDB, err := db.DB()

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)

	//  设置最大打开连接数
	sqlDB.SetMaxOpenConns(20)

	if err = sqlDB.Ping(); err != nil {
		logrus.Println("链接失败")
	}
	var user model.User
	db.Preload("Category.CountDown").First(&user)
	for _, category := range user.Category {
		fmt.Println(category.CountDown)
	}
}
