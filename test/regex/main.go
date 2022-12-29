package main

import (
	logs "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func main() {
	str := "^[存储操作:get,delete]\\s+[桶]\\s+[键]$"
	str1 := "^[存储操作:test]\\s+[桶:hello]\\s+[键]$"
	test(str, str1)
	rule := `^(get)\s+(.*)\s+(.*)$`
	content := "get SillyGirl name"
	var arr []string
	if res := regexp.MustCompile(rule).FindStringSubmatch(content); len(res) > 0 {

		for i, s := range res {
			logs.Info(i, " = ", s)
		}
		arr = res[1:]
	}
	logs.Info(param(str, "键", arr))
	logs.Info(param(str, "存储操作", arr))

	if res := regexp.MustCompile(`\[.*\]`).FindStringSubmatch(str); len(res) > 0 {
		logs.Info("是要解析的类型")
	}

}

func param(str, key string, arr []string) string {
	for i, s := range strings.Split(str, "\\s+") {
		if strings.Contains(s, key) {
			if i < len(arr) {
				return arr[i]
			}
		}
	}
	return ""
}

func test(arr ...string) {
	for _, str := range arr {
		str = strings.Replace(str, "^", "", -1)
		str = strings.Replace(str, "$", "", -1)
		logs.Info(str)
		arr := strings.Split(str, "\\s+")
		logs.Info(arr)
		var myArr []string
		for _, s := range arr {
			str := ""
			reg := ""
			if strings.Contains(s, ":") {
				reg = ":.*]"
				str = regexp.MustCompile(reg).FindString(s)
				str = strings.Replace(str, ":", "", -1)
				str = strings.Replace(str, "]", "", -1)
				str = strings.Replace(str, ",", "|", -1)
				str = "(" + str + ")"
			} else {
				reg = `\[.*\]`
				str = regexp.MustCompile(reg).ReplaceAllString(s, "(.*)")

			}
			logs.Info(reg, str)
			myArr = append(myArr, str)
		}

		s := "^" + strings.Join(myArr, `\s+`) + "$"
		logs.Info(s)
		logs.Info()
	}
}

type Test struct {
	march string
}
