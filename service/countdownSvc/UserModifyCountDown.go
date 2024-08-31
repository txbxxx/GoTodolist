/**
 * @Author tanchang
 * @Description 修改倒计时
 * @Date 2024/8/30 16:08
 * @File:  UserModifyCountDown
 * @Software: GoLand
 **/

package countdownSvc

import "time"

type UserModifyCountDownService struct {
	Name       string    `json:"name" form:"name" binding:"required,max=10"`
	EndTime    time.Time `json:"endTime" form:"endTime" time_format:"2006-01-02 15:04:05"`
	StartTime  time.Time `json:"startTime" form:"startTime" binding:"required" time_format:"2006-01-02 15:04:05"`
	Background string    `json:"background" form:"background"`
}
