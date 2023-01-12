package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var Senders chan Sender

var isTerminal bool

func init() {
	isTerminal = false
}

func SetTerminal(b bool) {
	isTerminal = b
}

func initToHandleMessage() {
	//reply := BoltBucket("reply")
	Senders = make(chan Sender)
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
func HandleMessage(sender Sender) {
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

	key := fmt.Sprintf("u=%v&c=%v&i=%v", sender.GetUserId(), sender.GetChatId(), sender.GetImType())
	if v, ok := waits.Load(key); ok {
		c := v.(*Carry)
		if m := regexp.MustCompile(c.Pattern).FindString(content); m != "" {
			logs.Info(m)
			c.Chan <- m
		}
		return
	}

	parseFunction(sender)
}
