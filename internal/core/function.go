package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

type Function struct {
	Rules   []string
	ImType  *Filter
	UserId  *Filter
	GroupId *Filter
	FindAll bool
	Handle  func(s Sender) interface{} `json:"-"`
	Show    string
	Hidden  bool

	Author      string
	Origin      string
	CreateAt    string
	Description string
	Version     string
	Title       string
	Platform    string
	Priority    int
	Cron        string
	Admin       bool
	Public      bool
	Icon        string
	Encrypt     bool
	Disable     bool
	Content     string
}

type Filter struct {
	BlackMode bool
	Items     []string
}

var functions []Function

func getPlugins() []Function {
	return functions
}

func createPlugins() {
	db := BoltBucket("plugins")
	db.Foreach(func(k, v []byte) error {
		functions = append(functions, createPlugin(string(v)))
		return nil
	})
}

func createPlugin(str string) Function {
	reg := "/\\*(.|\\r\\n|\\n)*?\\*/"
	if res := regexp.MustCompile(reg).FindStringSubmatch(str); len(res) != 0 {
		data := res[0]
		//fmt.Println(data)
		var rules []string
		for _, res := range regexp.MustCompile(`@rule\s+(.+)`).FindAllStringSubmatch(data, -1) {
			rules = append(rules, strings.Trim(res[1], " "))
		}

		return Function{
			Rules:       rules,
			Author:      getString("author", data),
			Origin:      getString("origin", data),
			CreateAt:    getString("create_at", data),
			Description: getString("description", data),
			Version:     getString("version", data),
			Title:       getString("title", data),
			Platform:    getString("platform", data),
			Priority:    getInt("priority", data),
			Cron:        getString("cron", data),
			Admin:       getBool("admin", data),
			Public:      getBool("public", data),
			Icon:        getString("icon", data),
			Encrypt:     getBool("encrypt", data),
			Disable:     getBool("disable", data),
			FindAll:     true,
		}

	}
	return Function{
		Content: str,
		FindAll: true,
	}
}

func getString(key, data string) string {
	r := ""
	for _, res := range regexp.MustCompile("@"+key+`\s+(.+)`).FindAllStringSubmatch(data, -1) {
		r = strings.Trim(res[1], " ")
	}
	return r
}

func getInt(key, data string) int {
	r := 0
	for _, res := range regexp.MustCompile("@"+key+`\s+(.+)`).FindAllStringSubmatch(data, -1) {
		s := strings.Trim(res[1], " ")
		r, _ = strconv.Atoi(fmt.Sprint(s))
	}
	return r
}

func getBool(key, data string) bool {
	r := false
	for _, res := range regexp.MustCompile("@"+key+`\s+(.+)`).FindAllStringSubmatch(data, -1) {
		r = strings.Trim(res[1], " ") == "true"
	}
	return r
}

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
