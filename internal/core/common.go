package core

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

func black(filter *Filter, str string) bool {
	if filter != nil {
		if filter.BlackMode {
			if Contains(filter.Items, str) {
				return true
			}
		} else {
			if !Contains(filter.Items, str) {
				return true
			}
		}
	}
	return false
}

func Contains(strs []string, str string) bool {
	for _, o := range strs {
		if str == o {
			return true
		}
	}
	return false
}

func GenUUID() string {
	u2, _ := uuid.NewUUID()
	return u2.String()
}

func refreshPlugins() {
	InitWatch()
	functions = []*Function{}
	runningList = []Running{}
	initPlugins()
	keyMap = initServerPlugin(getServers()...)
}

func getMessage(msgs ...interface{}) string {
	message := ""
	for _, msg := range msgs {
		fmt.Println("test -> ", msg)
		message += fmt.Sprintf("%v", msg)
	}

	return message
}

func unicode2utf8(source string) string {
	var res = []string{""}
	sUnicode := strings.Split(source, "\\u")
	var context = ""
	for _, v := range sUnicode {
		var additional = ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}
