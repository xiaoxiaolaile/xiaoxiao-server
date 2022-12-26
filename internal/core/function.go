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
type Functions []*Function

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

func getModules() Functions {
	return getFunctions(func(d Function) bool {
		return d.Module
	})
}

func getRules() Functions {
	return getFunctions(func(d Function) bool {
		return d.Rules != nil
	})
}
func getCron() Functions {
	return getFunctions(func(d Function) bool {
		return len(d.Cron) > 0
	})
}
func getServers() Functions {
	return getFunctions(func(d Function) bool {
		return d.OnStart
	})
}

func getFunctions(f func(d Function) bool) []*Function {
	sort.Sort(functions)
	var data []*Function
	for _, function := range functions {
		if f(*function) {
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

	AddCommand("", functions...)
}

func NewPlugin(f Function) {
	functions = append(functions, &f)
	AddCommand("", &f)
}

func createPlugin(str string) *Function {
	reg := "/\\*(.|\\r\\n|\\n)*?\\*/"
	if res := regexp.MustCompile(reg).FindStringSubmatch(str); len(res) != 0 {
		data := res[0]
		//fmt.Println(data)
		var rules []string
		for _, res := range regexp.MustCompile(`@rule\s+(.+)`).FindAllStringSubmatch(data, -1) {
			rules = append(rules, strings.Trim(res[1], " "))
		}

		return &Function{
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
			Content:     str,
			FindAll:     true,
		}

	}
	return &Function{
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

func AddCommand(prefix string, funArray ...*Function) {
	for j := range funArray {

		fun := funArray[j]

		if fun.Disable {
			continue
		}

		addRules(prefix, fun)

		if fun.Cron != "" {
			cmd := fun
			if _, err := C.AddFunc(fun.Cron, func() {
				cmd.Handle(&Faker{})
			}); err != nil {

			} else {

			}
		}
	}
}

func addRules(prefix string, function *Function) {

	rules := function.Rules

	if len(rules) > 0 {
		for i := range rules {
			if strings.Contains(rules[i], "raw ") {
				rules[i] = strings.Replace(rules[i], "raw ", "", -1)
				continue
			}
			rules[i] = strings.ReplaceAll(rules[i], `\r\a\w`, "raw")
			if strings.Contains(rules[i], "$") {
				continue
			}
			if prefix != "" {
				rules[i] = prefix + `\s+` + rules[i]
			}
			rules[i] = strings.Replace(rules[i], "(", `[(]`, -1)
			rules[i] = strings.Replace(rules[i], ")", `[)]`, -1)
			rules[i] = regexp.MustCompile(`\?$`).ReplaceAllString(rules[i], `([\s\S]+)`)
			rules[i] = strings.Replace(rules[i], " ", `\s+`, -1)
			rules[i] = strings.Replace(rules[i], "?", `(\S+)`, -1)
			rules[i] = "^" + rules[i] + "$"
		}
		f := function.Handle
		if f == nil {
			function.Handle = func(s Sender) interface{} {
				logs.Printf(function.Content)
				return "nil"
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
						logs.Info("1:匹配到规则：%s", rule)
					}
					sender.SetAllMatch(tmp)
					matched = true
				}
			} else {
				if res := regexp.MustCompile(rule).FindStringSubmatch(content); len(res) > 0 {
					if !function.Hidden {
						logs.Info("2:匹配到规则：%s", rule)
					}
					sender.SetMatch(res[1:])
					matched = true
				}
			}
			if matched {
				//if function.Admin && !sender.IsAdmin() {
				//	sender.Delete()
				//	sender.Disappear()
				//	return
				//}
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
