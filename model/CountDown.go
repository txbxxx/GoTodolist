/**
 * @Author tanchang
 * @Description 倒计时数据表模型
 * @Date 2024/8/29 11:40
 * @File:  CountDown
 * @Software: GoLand
 **/

package model

import (
	"gorm.io/gorm"
)

type CountDown struct {
	gorm.Model
	Identity         string `gorm:"column:identity;type:varchar(36);unique" json:"identity"` // 倒计时唯一标识
	Name             string `gorm:"column:name;type:varchar(100);" json:"name"`              // 倒计时名
	StartTime        int64  `gorm:"column:start_time;type:bigint" json:"startTime"`          // 开始时间
	EndTime          int64  `gorm:"column:end_time;type:bigint;default:0" json:"endTime"`    // 结束时间
	Background       string `gorm:"column:background;type:varchar(255)" json:"background"`   // 倒计时背景图
	CategoryIdentity string // 倒计时的分类信息
}
