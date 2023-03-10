package core

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	logs "github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
	//"github.com/Han-Ya-Jun/qrcode2console"
)

type Sender interface {
	GetUserId() string
	GetChatId() int
	GetImType() string
	GetMessageId() string
	RecallMessage(...interface{}) error
	GetUserName() string
	GetChatName() string
	IsReply() bool
	GetReplySenderUserId() int
	GetRawMessage() interface{}
	SetMatch([]string)
	SetRuleContent(string)
	SetAllMatch([][]string)
	GetMatch() []string
	GetAllMatch() [][]string
	Get(...int) string
	GetContent() string
	SetContent(string)
	IsAdmin() bool
	IsMedia() bool
	Reply(...interface{}) ([]string, error)
	Delete() error
	Disappear(lifetime ...time.Duration)
	Finish()
	Continue()
	IsContinue() bool
	ClearContinue()
	Await(Sender, func(Sender) interface{}, ...interface{}) interface{}
	Copy() Sender
	GroupKick(uid string, reject_add_request bool)
	GroupBan(uid string, duration int)
	AtLast()
	UAtLast()
	IsAtLast() bool
	MessagesToSend() string
	Param(interface{}) string
	GetPlatform() string
}

type Edit int
type Replace int
type Notify int
type Article []string

var E Edit
var R Replace
var N Notify

type ImageData []byte
type ImageBase64 string
type ImageUrl string
type ImagePath string
type VideoUrl string

type Faker struct {
	Message  string
	Type     string
	Platform string
	UserId   string
	ChatId   int
	Carry    chan string
	BaseSender
	Admin bool
	Vm    *goja.Runtime
}

type ResSender struct {
	context string
}

func (r *ResSender) GetContent() string {
	return r.context
}

func (sender *Faker) Listen(args ...interface{}) *ResSender {

	//进行判断 ，如果第一个参数是数字，进行监听，如果是数组，进行注册插件
	if len(args) > 0 {
		d := args[0]
		switch d.(type) {
		case int64:
			return listen(sender, d.(int64), "")
		case []interface{}:
			tmp := d.([]interface{})
			var rules []string
			for _, s := range tmp {
				rules = append(rules, fmt.Sprintf("%v", s))
			}
			NewPlugin(Function{
				Rules: rules,
				Admin: true,
				Handle: func(s Sender) interface{} {
					if len(args) > 1 {
						f := args[1].(func(call goja.FunctionCall) goja.Value)
						if sender.Vm != nil {
							call := goja.FunctionCall{
								Arguments: []goja.Value{
									sender.Vm.ToValue(&s),
								},
							}
							return f(call)
						}
					}
					return nil
				},
			})

		}
	}

	return nil

}
func listen(sender *Faker, data int64, pattern string) *ResSender {

	//进行判断 ，如果第一个参数是数字，进行监听，如果是数组，进行注册插件

	key := fmt.Sprintf("u=%v&c=%v&i=%v&t=%v", sender.GetUserId(), sender.GetChatId(), sender.GetImType(), time.Now().UnixNano())
	timeout := time.Millisecond * time.Duration(data)
	//if fg != nil {
	//	if *fg == "me" {
	//		key += fmt.Sprintf("&f=me")
	//	} else {
	//		key += fmt.Sprintf("&f=true")
	//	}
	//}
	c := &Carry{Sender: sender, Pattern: pattern}

	if c.Pattern == "" {
		c.Pattern = `[\s\S]*`
	}

	c.Chan = make(chan interface{}, 1)
	c.Result = make(chan interface{}, 1)

	if oc, ok := waits.LoadOrStore(key, c); ok {
		oc.(*Carry).Chan <- InterruptError
	}
	defer func() {
		waits.Delete(key)
	}()

	for {
		select {
		case result := <-c.Chan:
			logs.Info(result)
			switch result.(type) {
			case string:
				return &ResSender{context: result.(string)}
			}
		case <-time.After(timeout):
			return nil
		}
	}

}

func (sender *Faker) GetContent() string {
	if sender.Content != "" {
		return sender.Content
	}
	return sender.Message
}

func (sender *Faker) GetUserId() string {
	return sender.UserId
}

func (sender *Faker) GetChatId() int {
	return sender.ChatId
}

func (sender *Faker) GetPlatform() string {
	if sender.Platform == "" {
		return "fake"
	}
	return sender.Platform
}

func (sender *Faker) GetImType() string {
	if sender.Type == "" {
		return "fake"
	}
	return sender.Type
}

func (sender *Faker) GetMessageId() string {
	return ""
}

func (sender *Faker) GetUserName() string {
	return ""
}

func (sender *Faker) GetChatName() string {
	return ""
}

func (sender *Faker) IsReply() bool {
	return false
}

func (sender *Faker) GetReplySenderUserId() int {
	return 0
}

func (sender *Faker) GetRawMessage() interface{} {
	return sender.Message
}

func (sender *Faker) IsAdmin() bool {
	return true
}

func (sender *Faker) IsMedia() bool {
	return false
}

func (sender *Faker) Reply(msgs ...interface{}) ([]string, error) {
	rt := ""
	//var n *Notify
	for _, msg := range msgs {
		switch msg.(type) {
		case []byte:
			rt = string(msg.([]byte))
		case string:
			rt = msg.(string)
		case ImageUrl:

		case Notify:
			//v := msg.(Notify)
			//n = &v
		}
	}
	//{
	//
	//	for _, v := range regexp.MustCompile(`\[CQ:image,file=([^\[\]]+)\]`).FindAllStringSubmatch(rt, -1) {
	//		qr := qrcode2console.NewQRCode2ConsoleWithUrl(v[1], true)
	//		defer qr.Output()
	//		rt = strings.Replace(rt, fmt.Sprintf(`[CQ:image,file=%s]`, v[1]), "", -1)
	//	}
	//}
	//
	//if rt != "" && n != nil {
	//	NotifyMasters(rt)
	//}

	if rt != "" {
		if sender.Type == "carry" {
			if sender.Carry != nil {
				sender.Carry <- rt
			}
		} else if sender.Type == "terminal" {
			fmt.Printf("\x1b[%dm%s \x1b[0m\n", 31, rt)
		}
	}

	return []string{rt}, nil
}

func (sender *Faker) Delete() error {
	return nil
}

func (sender *Faker) Disappear(lifetime ...time.Duration) {

}

func (sender *Faker) Finish() {
	if sender.Carry != nil {
		close(sender.Carry)
	}
}

func (sender *Faker) Copy() Sender {
	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Faker)
	return &new
}

func (sender *Faker) GroupKick(uid string, reject_add_request bool) {

}

func (sender *Faker) GroupBan(uid string, duration int) {

}

type BaseSender struct {
	matches        [][]string
	goon           bool
	child          Sender
	Content        string
	Atlast         bool
	ToSendMessages []string
	IsFinished     bool
	Duration       *time.Duration
	RuleContent    string
}

func (sender *BaseSender) SetMatch(ss []string) {
	sender.matches = [][]string{ss}
}

func (sender *BaseSender) SetRuleContent(s string) {
	sender.RuleContent = s
}
func (sender *BaseSender) SetAllMatch(ss [][]string) {
	sender.matches = ss
}

func (sender *BaseSender) SetContent(content string) {
	sender.Content = content
}

func (sender *BaseSender) GetMatch() []string {
	return sender.matches[0]
}

func (sender *BaseSender) GetAllMatch() [][]string {
	return sender.matches
}

func (sender *BaseSender) Continue() {
	sender.goon = true
}

func (sender *BaseSender) IsContinue() bool {
	return sender.goon
}

func (sender *BaseSender) ClearContinue() {
	sender.goon = false
}

func (sender *BaseSender) Get(index ...int) string {
	i := 0
	if len(index) != 0 {
		i = index[0]
	}
	if len(sender.matches) == 0 {
		return ""
	}
	if len(sender.matches[0]) < i+1 {
		return ""
	}
	return sender.matches[0][i]
}

func (sender *BaseSender) Delete() error {
	return nil
}

func (sender *BaseSender) Disappear(lifetime ...time.Duration) {

}

func (sender *BaseSender) Finish() {
	sender.IsFinished = true
}

func (sender *BaseSender) IsMedia() bool {
	return false
}

func (sender *BaseSender) GetRawMessage() interface{} {
	return nil
}

func (sender *BaseSender) IsReply() bool {
	return false
}

func (sender *BaseSender) GetMessageId() string {
	return ""
}

func (sender *BaseSender) RecallMessage(...interface{}) error {
	return nil
}

func (sender *BaseSender) GetUserId() string {
	return ""
}
func (sender *BaseSender) GetChatId() int {
	return 0
}
func (sender *BaseSender) GetImType() string {
	return ""
}

func (sender *BaseSender) GetPlatform() string {
	return ""
}

func (sender *BaseSender) GroupKick(uid string, reject_add_request bool) {

}

func (sender *BaseSender) GroupBan(uid string, duration int) {

}

func (sender *BaseSender) GetUserName() string {
	return ""
}

func (sender *BaseSender) IsAdmin() bool {
	return false
}

func (sender *BaseSender) GetChatName() string {
	return ""
}

func (sender *BaseSender) GetReplySenderUserId() int {
	return 0
}

func (sender *BaseSender) AtLast() {
	sender.Atlast = true
}

func (sender *BaseSender) UAtLast() {
	sender.Atlast = false
}

func (sender *BaseSender) IsAtLast() bool {
	return sender.Atlast
}

func (sender *BaseSender) MessagesToSend() string {
	return strings.Join(sender.ToSendMessages, "\n")
}
func (sender *BaseSender) Param(key interface{}) string {

	switch key.(type) {
	case int64:
		return sender.Get(int(key.(int64)))
	case string:
		if len(sender.matches) > 0 {
			return param(sender.RuleContent, key.(string), sender.matches[0])
		}
	}
	return ""
}

func param(str, key string, arr []string) string {
	logs.Warn(str)
	logs.Warn(key)
	logs.Warn(arr)
	logs.Warn("=========")
	for i, s := range strings.Split(str, "\\s+") {
		if strings.Contains(s, key) {
			if i < len(arr) {
				return arr[i]
			}
		}
	}
	return ""
}

var TimeOutError = errors.New("指令超时")
var InterruptError = errors.New("被其他指令中断")

var waits sync.Map

func WaitsRange(f func(key, value any) bool) {
	waits.Range(f)
}

type Carry struct {
	Chan    chan interface{}
	Pattern string
	Result  chan interface{}
	Sender  Sender
}

type forGroup string

type again string

var Again again = ""

var GoAgain = func(str string) again {
	return again(str)
}

type YesOrNo string

var YesNo YesOrNo = "yeson"
var Yes YesOrNo = "yes"
var No YesOrNo = "no"

type Range []int

type Switch []string

var ForGroup forGroup

func (_ *BaseSender) Await(sender Sender, callback func(Sender) interface{}, params ...interface{}) interface{} {
	c := &Carry{}
	timeout := time.Hour * 999999
	var handleErr func(error)
	var fg *forGroup
	for _, param := range params {
		switch param.(type) {
		case string:
			c.Pattern = param.(string)
		case time.Duration:
			du := param.(time.Duration)
			if du != 0 {
				timeout = du
			}
		case func() string:
			callback = param.(func(Sender) interface{})
		case func(error):
			handleErr = param.(func(error))
		case forGroup:
			a := param.(forGroup)
			fg = &a
		}
	}
	// if callback == nil {
	// 	return nil
	// }
	if c.Pattern == "" {
		c.Pattern = `[\s\S]*`
	}
	c.Chan = make(chan interface{}, 1)
	c.Result = make(chan interface{}, 1)

	key := fmt.Sprintf("u=%v&c=%v&i=%v&t=%v", sender.GetUserId(), sender.GetChatId(), sender.GetImType(), time.Now().UnixNano())
	if fg != nil {
		key += fmt.Sprintf("&f=true")
	}
	if oc, ok := waits.LoadOrStore(key, c); ok {
		oc.(*Carry).Chan <- InterruptError
	}
	defer func() {
		waits.Delete(key)
	}()
	for {
		select {
		case result := <-c.Chan:
			switch result.(type) {
			case Sender:
				s := result.(Sender)
				if callback == nil {
					return s.GetContent()
				}
				result := callback(s)
				if v, ok := result.(again); ok {
					if v == "" {
						c.Result <- nil
					} else {
						c.Result <- string(v)
					}
				} else if _, ok := result.(YesOrNo); ok {
					o := strings.ToLower(regexp.MustCompile("[yYnN]").FindString(s.GetContent()))
					if o == "y" {
						return Yes
					}
					if o == "n" {
						return No
					}
					c.Result <- "Y or n ?"
				} else if vv, ok := result.(Switch); ok {
					ct := s.GetContent()
					for _, v := range vv {
						if ct == v {
							return v
						}
					}
					c.Result <- fmt.Sprintf("请从%s中选择一个。", strings.Join(vv, "、"))
				} else if vv, ok := result.(Range); ok {
					ct := s.GetContent()
					n := Int(ct)
					if fmt.Sprint(n) == ct {
						if (n >= vv[0]) && (n <= vv[1]) {

							return n
						}
					}
					c.Result <- fmt.Sprintf("请从%d~%d中选择一个整数。", vv[0], vv[1])
				} else {
					c.Result <- result
					return s.GetContent()
				}
			case error:
				if handleErr != nil {
					handleErr(result.(error))
				}
				c.Result <- nil
				return nil
			}
		case <-time.After(timeout):
			if handleErr != nil {
				handleErr(TimeOutError)
			}
			c.Result <- nil
			return nil
		}
	}
}
