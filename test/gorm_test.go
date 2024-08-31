/**
 * @Author tanchang
 * @Description //TODO
 * @Date 2024/7/11 22:59
 * @File:  gormtest
 * @Software: GoLand
 **/

package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	s := "countdown:OEC:20c67cbf-2678-4e46-ad46-786f9e4cc62e"
	split := strings.Split(s, ":")
	fmt.Println(split[2])
}
