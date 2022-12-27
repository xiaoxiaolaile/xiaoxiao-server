package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strings"
	"xiaoxiao/internal/jsvm"
)

var Senders chan jsvm.Sender

var isTerminal bool

func init() {
	isTerminal = false
}

func SetTerminal(b bool) {
	isTerminal = b
}

func initToHandleMessage() {
	//reply := BoltBucket("reply")
	Senders = make(chan jsvm.Sender)
	go func() {
		for {
			s := <-Senders
			go HandleMessage(s)
		}
	}()
}

func TrimHiddenCharacter(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 && c != 10 {
			continue
		}
		if c == 127 {
			continue
		}

		dstRunes = append(dstRunes, c)
	}
	return strings.ReplaceAll(string(dstRunes), "￼", "")
}

// HandleMessage 处理接受到的消息
func HandleMessage(sender jsvm.Sender) {
	defer func() {
		recover()
	}()
	ct := sender.GetContent()
	content := TrimHiddenCharacter(ct)
	defer func() {
		sender.Finish()
		if sender.IsAtLast() {
			s := sender.MessagesToSend()
			if s != "" {
				sender.Reply(s)
			}
		}
	}()
	u, g, i := fmt.Sprint(sender.GetUserId()), fmt.Sprint(sender.GetChatId()), fmt.Sprint(sender.GetImType())

	if isTerminal {
		logs.Printf("接收到消息 %v/%v@%v：%s", i, u, g, content)
	}

	mtd := false

	jsvm.WaitsRange(func(k, v interface{}) bool {
		c := v.(*jsvm.Carry)
		vs, _ := url.ParseQuery(k.(string))
		//userID := vs.Get("u")
		//chatID := vs.Get("c")
		//imType := vs.Get("i")
		//forGroup := vs.Get("f")
		logs.Info(vs)

		//if imType != i {
		//	return true
		//}
		//if chatID != g && (forGroup != "me" || g != "0") {
		//	return true
		//}
		//if userID != u && (forGroup == "" || forGroup == "me") {
		//	return true
		//}
		if m := regexp.MustCompile(c.Pattern).FindString(content); m != "" {
			mtd = true
			logs.Info(m)
			c.Chan <- m
		}
		return true
	})
	if mtd {
		return
	}

	parseFunction(sender)
}
