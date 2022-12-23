package core

import logs "github.com/sirupsen/logrus"

func Init() {
	initLog()
	initStore()
	initToHandleMessage()
}

func initLog() {
	logs.SetFormatter(&logs.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// Only log the warning severity or above.
	logs.SetLevel(logs.TraceLevel)
}
