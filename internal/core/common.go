package core

import (
	"fmt"
	"github.com/google/uuid"
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
