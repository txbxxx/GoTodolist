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
	"math"
	"strconv"
	"time"
)

type CountdownSerialize struct {
	Identity   string `json:"identity"`
	Name       string `json:"name"`
	Day        string `json:"day"`
	Category   string `json:"category"`
	Background string `json:"background"`
}

// CountdownSerializeList 多个序列化
func CountdownSerializeList(countdowns []map[string]string) []CountdownSerialize {
	var countdownList []CountdownSerialize
	for _, countdown := range countdowns {
		countdownList = append(countdownList, CountdownSerialize{
			Identity:   countdown["identity"],
			Name:       countdown["name"],
			Day:        countdown["day"],
			Category:   countdown["category"],
			Background: countdown["background"],
		})
	}
	return countdownList
}

func CountdownSerializeSingle(countdown map[string]string) CountdownSerialize {
	return CountdownSerialize{
		Identity:   countdown["identity"],
		Name:       countdown["name"],
		Day:        countdown["day"],
		Category:   countdown["category"],
		Background: countdown["background"],
	}
}

// CountdownSerializeSingleModel 单个序列化
func CountdownSerializeSingleModel(countdown model.CountDown) CountdownSerialize {
	var day float64
	now := time.Now().Unix()
	if countdown.EndTime > 0 {
		day = float64(now-countdown.StartTime) / 86400
	} else {
		day = float64(countdown.EndTime-now) / 86400
	}
	return CountdownSerialize{
		Identity:   countdown.Identity,
		Name:       countdown.Name,
		Day:        strconv.FormatFloat(math.Ceil(day), 'f', 2, 64),
		Category:   countdown.CategoryIdentity,
		Background: countdown.Background,
	}
}

// CountdownSerializeListModel 单个序列化
func CountdownSerializeListModel(countdowns []model.CountDown) []CountdownSerialize {
	var countdownList []CountdownSerialize
	for _, countdown := range countdowns {
		var day float64
		now := time.Now().Unix()
		if countdown.EndTime > 0 {
			day = float64(now-countdown.StartTime) / 86400
		} else {
			day = float64(countdown.EndTime-now) / 86400
		}
		countdownList = append(countdownList, CountdownSerialize{
			Identity:   countdown.Identity,
			Name:       countdown.Name,
			Day:        strconv.FormatFloat(math.Ceil(day), 'f', 2, 64),
			Category:   countdown.CategoryIdentity,
			Background: countdown.Background,
		})
	}
	return countdownList
}
