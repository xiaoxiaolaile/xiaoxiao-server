package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
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

type Function struct {
	Rules    []string
	ImType   *Filter
	UserId   *Filter
	GroupId  *Filter
	FindAll  bool
	Admin    bool
	Handle   func(s Sender) interface{}
	Cron     string
	Show     string
	Priority int
	Disable  bool
	Hash     string
	Hidden   bool
}
type Filter struct {
	BlackMode bool
	Items     []string
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
	u, g, i := fmt.Sprint(sender.GetUserID()), fmt.Sprint(sender.GetChatID()), fmt.Sprint(sender.GetImType())

	if isTerminal {
		logs.Printf("接收到消息 %v/%v@%v：%s", i, u, g, content)
	}

	//con := true
	//mtd := false
	//waits.Range(func(k, v interface{}) bool {
	//	c := v.(*Carry)
	//	vs, _ := url.ParseQuery(k.(string))
	//	userID := vs.Get("u")
	//	chatID := vs.Get("c")
	//	imType := vs.Get("i")
	//	forGroup := vs.Get("f")
	//	if imType != i {
	//		return true
	//	}
	//	if chatID != g && (forGroup != "me" || g != "0") {
	//		return true
	//	}
	//	if userID != u && (forGroup == "" || forGroup == "me") {
	//		return true
	//	}
	//	if m := regexp.MustCompile(c.Pattern).FindString(content); m != "" {
	//		mtd = true
	//		if f, ok := c.Sender.(*Faker); ok && f.Carry != nil {
	//			if s1, o := sender.(*Faker); o && s1.Carry != nil {
	//				f.Carry = s1.Carry
	//				c := make(chan string)
	//				oc := s1.Carry
	//				s1.Carry = c
	//				go func() {
	//					for {
	//						r, o := <-c
	//						if !o {
	//							break
	//						}
	//						oc <- r
	//					}
	//				}()
	//			}
	//		}
	//		c.Chan <- sender
	//		sender.Reply(<-c.Result)
	//		if !sender.IsContinue() {
	//			con = false
	//			return false
	//		}
	//		content = TrimHiddenCharacter(sender.GetContent())
	//	}
	//	return true
	//})
	//if mtd && !con {
	//	return
	//}
	//replied := false
	//MakeBucket(fmt.Sprintf("reply%s%d", sender.GetImType(), sender.GetChatID())).Foreach(func(k, v []byte) error {
	//	if string(v) == "" {
	//		return nil
	//	}
	//	reg, err := regexp.Compile(string(k))
	//	if err == nil {
	//		if reg.FindString(content) != "" {
	//			replied = true
	//			r := string(v)
	//			if strings.Contains(r, "$") {
	//				sender.Reply(reg.ReplaceAllString(content, r))
	//			} else {
	//				sender.Reply(r)
	//			}
	//		}
	//	}
	//	return nil
	//})
	//
	//if !replied {
	//	reply.Foreach(func(k, v []byte) error {
	//		if string(v) == "" {
	//			return nil
	//		}
	//		reg, err := regexp.Compile(string(k))
	//		if err == nil {
	//			if reg.FindString(content) != "" {
	//				replied = true
	//				r := string(v)
	//				if strings.Contains(r, "$") {
	//					sender.Reply(reg.ReplaceAllString(content, r))
	//				} else {
	//					sender.Reply(r)
	//				}
	//			}
	//		}
	//		return nil
	//	})
	//}
	//
	//for _, function := range Functions {
	//	if black(function.ImType, sender.GetImType()) || black(function.UserId, sender.GetUserID()) || black(function.GroupId, fmt.Sprint(sender.GetChatID())) {
	//		continue
	//	}
	//	for _, rule := range function.Rules {
	//		var matched bool
	//		if function.FindAll {
	//			if res := regexp.MustCompile(rule).FindAllStringSubmatch(content, -1); len(res) > 0 {
	//				var tmp [][]string
	//				for i := range res {
	//					tmp = append(tmp, res[i][1:])
	//				}
	//				if !function.Hidden {
	//					logs.Info("匹配到规则：%s", rule)
	//				}
	//				sender.SetAllMatch(tmp)
	//				matched = true
	//			}
	//		} else {
	//			if res := regexp.MustCompile(rule).FindStringSubmatch(content); len(res) > 0 {
	//				if !function.Hidden {
	//					logs.Info("匹配到规则：%s", rule)
	//				}
	//				sender.SetMatch(res[1:])
	//				matched = true
	//			}
	//		}
	//		if matched {
	//			if function.Admin && !sender.IsAdmin() {
	//				sender.Delete()
	//				sender.Disappear()
	//				return
	//			}
	//			rt := function.Handle(sender)
	//			if rt != nil {
	//				sender.Reply(rt)
	//			}
	//			if sender.IsContinue() {
	//				sender.ClearContinue()
	//				content = utils.TrimHiddenCharacter(sender.GetContent())
	//				if !function.Hidden {
	//					logs.Info("继续去处理：%s", content)
	//				}
	//				goto next
	//			}
	//			return
	//		}
	//	}
	//next:
	//}
	//
	//recall := sillyGirl.GetString("recall")
	//if recall != "" {
	//	recalled := false
	//	for _, v := range strings.Split(recall, "&") {
	//		reg, err := regexp.Compile(v)
	//		if err == nil {
	//			if reg.FindString(content) != "" {
	//				if !sender.IsAdmin() && sender.GetImType() != "wx" {
	//					sender.Delete()
	//					recalled = true
	//					break
	//				}
	//			}
	//		}
	//	}
	//	if recalled {
	//		return
	//	}
	//}
}
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
