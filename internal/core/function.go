package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"regexp"
	"sort"
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
	Module      bool
	OnStart     bool
}

type Filter struct {
	BlackMode bool
	Items     []string
}

var functions Functions

// Functions 将[]string定义为MyStringList类型
type Functions []Function

// Len 实现sort.Interface接口的获取元素数量方法
func (m Functions) Len() int {
	return len(m)
}

// Less 实现sort.Interface接口的比较元素方法
func (m Functions) Less(i, j int) bool {
	return m[i].Priority > m[j].Priority
}

// Swap 实现sort.Interface接口的交换元素方法
func (m Functions) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func getPlugins() Functions {
	return getFunctions(func(d Function) bool {
		return true
	})
}

func getModules() []Function {
	return getFunctions(func(d Function) bool {
		return d.Module
	})
}

func getRules() []Function {
	return getFunctions(func(d Function) bool {
		return d.Rules != nil
	})
}
func getCron() []Function {
	return getFunctions(func(d Function) bool {
		return len(d.Cron) > 0
	})
}
func getServers() []Function {
	return getFunctions(func(d Function) bool {
		return d.OnStart
	})
}

func getFunctions(f func(d Function) bool) []Function {
	sort.Sort(functions)
	var data []Function
	for _, function := range functions {
		if f(function) {
			data = append(data, function)
		}
	}
	return data
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
			Module:      getBool("module", data),
			OnStart:     getBool("on_start", data),
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
