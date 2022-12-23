package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

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

var functions []Function

func AddCommand(prefix string, funArray ...Function) {
	for j := range funArray {
		if funArray[j].Disable {
			continue
		}
		for i := range funArray[j].Rules {
			if strings.Contains(funArray[j].Rules[i], "raw ") {
				funArray[j].Rules[i] = strings.Replace(funArray[j].Rules[i], "raw ", "", -1)
				continue
			}
			funArray[j].Rules[i] = strings.ReplaceAll(funArray[j].Rules[i], `\r\a\w`, "raw")
			if strings.Contains(funArray[j].Rules[i], "$") {
				continue
			}
			if prefix != "" {
				funArray[j].Rules[i] = prefix + `\s+` + funArray[j].Rules[i]
			}
			funArray[j].Rules[i] = strings.Replace(funArray[j].Rules[i], "(", `[(]`, -1)
			funArray[j].Rules[i] = strings.Replace(funArray[j].Rules[i], ")", `[)]`, -1)
			funArray[j].Rules[i] = regexp.MustCompile(`\?$`).ReplaceAllString(funArray[j].Rules[i], `([\s\S]+)`)
			funArray[j].Rules[i] = strings.Replace(funArray[j].Rules[i], " ", `\s+`, -1)
			funArray[j].Rules[i] = strings.Replace(funArray[j].Rules[i], "?", `(\S+)`, -1)
			funArray[j].Rules[i] = "^" + funArray[j].Rules[i] + "$"
		}
		{
			lf := len(functions)
			for i := range functions {
				f := lf - i - 1
				if functions[f].Priority > funArray[j].Priority {
					functions = append(functions[:f+1], append([]Function{funArray[j]}, functions[f+1:]...)...)
					break
				}
			}
			if len(functions) == lf {
				if lf > 0 {
					if functions[0].Priority < funArray[j].Priority && functions[lf-1].Priority < funArray[j].Priority {
						functions = append([]Function{funArray[j]}, functions...)
					} else {
						functions = append(functions, funArray[j])
					}
				} else {
					functions = append(functions, funArray[j])
				}
			}
		}

		if funArray[j].Cron != "" {
			cmd := funArray[j]
			if _, err := C.AddFunc(funArray[j].Cron, func() {
				cmd.Handle(&Faker{})
			}); err != nil {

			} else {

			}
		}
	}
}

func parseFunction(sender Sender) {
	ct := sender.GetContent()
	content := TrimHiddenCharacter(ct)

	for _, function := range functions {
		if black(function.ImType, sender.GetImType()) || black(function.UserId, sender.GetUserID()) || black(function.GroupId, fmt.Sprint(sender.GetChatID())) {
			continue
		}
		for _, rule := range function.Rules {
			var matched bool
			if function.FindAll {
				if res := regexp.MustCompile(rule).FindAllStringSubmatch(content, -1); len(res) > 0 {
					var tmp [][]string
					for i := range res {
						tmp = append(tmp, res[i][1:])
					}
					if !function.Hidden {
						logs.Info("匹配到规则：%s", rule)
					}
					sender.SetAllMatch(tmp)
					matched = true
				}
			} else {
				if res := regexp.MustCompile(rule).FindStringSubmatch(content); len(res) > 0 {
					if !function.Hidden {
						logs.Info("匹配到规则：%s", rule)
					}
					sender.SetMatch(res[1:])
					matched = true
				}
			}
			if matched {
				if function.Admin && !sender.IsAdmin() {
					sender.Delete()
					sender.Disappear()
					return
				}
				rt := function.Handle(sender)
				if rt != nil {
					sender.Reply(rt)
				}
				if sender.IsContinue() {
					sender.ClearContinue()
					content = TrimHiddenCharacter(sender.GetContent())
					if !function.Hidden {
						logs.Info("继续去处理：%s", content)
					}
					goto next
				}
				return
			}
		}
	next:
	}
}
