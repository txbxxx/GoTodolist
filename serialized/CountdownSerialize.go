/**
 * @Author tanchang
 * @Description //TODO
 * @Date 2024/8/30 18:14
 * @File:  CountdownSerialize
 * @Software: GoLand
 **/

package serializes

type CountdownSerialize struct {
	Name       string `json:"name"`
	Day        string `json:"day"`
	Background string `json:"background"`
}

func CountdownSerializeList(countdowns []map[string]string) []CountdownSerialize {
	var countdownList []CountdownSerialize
	for _, countdown := range countdowns {
		countdownList = append(countdownList, CountdownSerialize{
			Name:       countdown["name"],
			Day:        countdown["day"],
			Background: countdown["background"],
		})
	}
	return countdownList
}
