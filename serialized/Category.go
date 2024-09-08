/**
 * @Author tanchang
 * @Description 分类列表序列化
 * @Date 2024/9/8 20:56
 * @File:  Category
 * @Software: GoLand
 **/

package serializes

import "GoToDoList/model"

type CategorySerialize struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	Cover    string `json:"Cover"`
}

// CategorySerializeList 多个序列化
func CategorySerializeList(category []map[string]string) []CategorySerialize {
	var countdownList []CategorySerialize
	for _, key := range category {
		countdownList = append(countdownList, CategorySerialize{
			Identity: key["identity"],
			Name:     key["name"],
			Cover:    key["background"],
		})
	}
	return countdownList
}

func CategorySerializeSingle(category map[string]string) CategorySerialize {
	return CategorySerialize{
		Identity: category["identity"],
		Name:     category["name"],
		Cover:    category["background"],
	}
}

func CategorySerializeSingleFromModel(category model.Category) CategorySerialize {
	return CategorySerialize{
		Identity: category.Identity,
		Name:     category.Name,
		Cover:    category.Cover,
	}
}

func CategorySerializeListFromModel(category []model.Category) []CategorySerialize {
	var countdownList []CategorySerialize
	for _, key := range category {
		countdownList = append(countdownList, CategorySerialize{
			Identity: key.Identity,
			Name:     key.Name,
			Cover:    key.Cover,
		})
	}
	return countdownList
}
