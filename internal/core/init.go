package core

import (
	logs "github.com/sirupsen/logrus"
	"time"
)

func Init() {
	initTimeLoc()
	initLog()
	InitStore()
	initWeb()
	initToHandleMessage()

	AddCommand("", Function{
		Rules: []string{"hello"},
		Admin: true,
		Handle: func(s Sender) interface{} {

			logs.Printf("获取参数：%s", s.Get(0))
			return "你好，小小 为您服务。"
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
