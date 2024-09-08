/**
 * @Author tanchang
 * @Description 分类表
 * @Date 2024/9/7 16:51
 * @File:  Category
 * @Software: GoLand
 **/

package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Identity     string `gorm:"column:identity;type:varchar(36);unique" json:"identity"` // 分类唯一标识
	Name         string `gorm:"column:name;type:varchar(100);unique" json:"name"`        // 分类名
	Cover        string `gorm:"column:background;type:varchar(255)" json:"background"`   // 分类封面图
	UserIdentity string
	CountDown    []CountDown `gorm:"foreignKey:CategoryIdentity;references:Identity;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 和倒计时表关联表示一个分类有多个倒计时
}
