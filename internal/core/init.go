package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"time"
)

func Init() {
	initToHandleMessage()
	initTimeLoc()
	initLog()
	InitStore()
	initSillyGirl()
	initPlugins()
	initWeb()
	initApi()

	NewPlugin(Function{
		Rules: []string{"hello"},
		Admin: true,
		Handle: func(s Sender) interface{} {
			return fmt.Sprintf("你好，%v 为您服务。", BoltBucket("sillyGirl").Get("name"))
		},
	})
}

func initTimeLoc() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc
}

func initLog() {
	logs.SetFormatter(&logs.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// Only log the warning severity or above.
	logs.SetLevel(logs.TraceLevel)
}
