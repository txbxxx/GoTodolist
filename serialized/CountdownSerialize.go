/**
 * @Author tanchang
 * @Description 倒计时序列化
 * @Date 2024/8/30 18:14
 * @File:  CountdownSerialize
 * @Software: GoLand
 **/

package serializes

import (
	"GoToDoList/model"
	"GoToDoList/utils"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type CountdownSerialize struct {
	Identity   string `json:"identity"`
	Name       string `json:"name"`
	Day        string `json:"day"`
	Background string `json:"background"`
}

// CountdownSerializeList 多个序列化
func CountdownSerializeList(countdowns []map[string]string, identity string) []CountdownSerialize {
	var countdownList []CountdownSerialize
	for _, countdown := range countdowns {
		countdownList = append(countdownList, CountdownSerialize{
			Identity:   identity,
			Name:       countdown["name"],
			Day:        countdown["day"],
			Background: countdown["background"],
		})
	}
	return countdownList
}

func CountdownSerializeSingle(countdown map[string]string, identity string) CountdownSerialize {
	return CountdownSerialize{
		Identity:   identity,
		Name:       countdown["name"],
		Day:        countdown["day"],
		Background: countdown["background"],
	}
}

// CountdownSerializeSingleModel 单个序列化
func CountdownSerializeSingleModel(countdown model.CountDown) CountdownSerialize {
	keyPrefix := "countdown:"
	var day float64
	var err error
	if countdown.EndTime > 0 {
		day, err = utils.OecCalculate(time.Now().Unix(), countdown.StartTime, keyPrefix+"OEC"+countdown.Identity, countdown.Background, countdown.Name)
		if err != nil {
			logrus.Error("计算日期错误:", err)
		}
	} else {
		day, err = utils.FdcCalculate(time.Now().Unix(), countdown.EndTime, keyPrefix+"OEC"+countdown.Identity, countdown.Background, countdown.Name)
		if err != nil {
			logrus.Error("计算日期错误:", err)
		}
	}
	return CountdownSerialize{
		Identity:   countdown.Identity,
		Name:       countdown.Name,
		Day:        strconv.FormatFloat(day, 'f', 2, 64),
		Background: countdown.Background,
	}
}
