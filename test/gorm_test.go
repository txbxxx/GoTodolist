/**
 * @Author tanchang
 * @Description 测试
 * @Date 2024/7/11 22:59
 * @File:  gormtest
 * @Software: GoLand
 **/

package test

import (
	"GoToDoList/model"
	"GoToDoList/utils"
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {

	var countdown model.CountDown
	utils.DB.Model(&model.CountDown{}).Where("identity = ?", "c5f3facf-ccf9-4d78-be76-959272fcfdf4").Take(&countdown).Unscoped()
	fmt.Println(countdown)
}
