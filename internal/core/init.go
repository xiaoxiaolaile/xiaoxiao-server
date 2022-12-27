package core

import (
	logs "github.com/sirupsen/logrus"
	"time"
	"xiaoxiao/internal/jsvm"
)

func Init() {
	initTimeLoc()
	initLog()
	jsvm.InitStore()
	initPlugins()
	initWeb()
	initToHandleMessage()

	NewPlugin(Function{
		Rules: []string{"hello"},
		Admin: true,
		Handle: func(s jsvm.Sender) interface{} {

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
