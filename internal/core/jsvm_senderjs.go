package core

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"time"
)

type SenderJs struct {
	Name string
	f    func(data interface{})
}

func NewSenderJs(name string) *SenderJs {
	s := &SenderJs{
		Name: name,
	}

	return s
}

func (s *SenderJs) Send(f func(data interface{})) {
	s.f = f
}

type Im struct {
	Faker
	data map[string]interface{}
	s    *SenderJs
}

func (i *Im) IsAdmin() bool {
	var b = BoltBucket(i.s.Name)
	for _, s := range b.GetArray("masters") {
		if s == fmt.Sprint(i.GetUserId()) {
			return true
		}
		//strings.Contains(wx.GetString("masters"), fmt.Sprint(i.GetUserId()))
	}
	return false
}

func (i *Im) Reply(msgs ...interface{}) (arr []string, err error) {
	logs.Info("reply message", msgs)
	arr, err = i.Faker.Reply(msgs)
	if len(msgs) > 0 {
		if i.s.f != nil {
			i.data["content"] = msgs[0]
			i.s.f(i.data)
		}

	}
	return
}

func (s *SenderJs) Receive(data map[string]interface{}) {
	if s.f != nil {
		//解析命令，修改命令
		//logs.Info("收到消息：", data)
		//logs.Info("收到消息content：", data["content"])
		//消息格式 {"user_id":"1234", "chat_id": 10000, "content": "消息内容"}
		if str, ok := data["content"]; ok {
			//进行解析，和替换内容
			go func() {
				chatId := 0

				switch data["chat_id"].(type) {
				case int:
					chatId = data["chat_id"].(int)
				}

				f := &Im{
					Faker: Faker{
						Type:     "carry",
						Platform: s.Name,
						Message:  fmt.Sprintf("%v", str),
						Admin:    true,
						UserId:   fmt.Sprintf("%v", data["user_id"]),
						ChatId:   chatId,
					},
					data: data,
					s:    s,
				}
				Senders <- f

			}()
			//go func() {
			//	f := &Faker{
			//		Type:     "terminal",
			//		Platform: "terminal",
			//		Message:  fmt.Sprintf("%v", str),
			//		Admin:    true,
			//	}
			//	Senders <- f
			//}()
		}

	}
}

type Running struct {
	data    map[string]interface{}
	handles []func(error, map[string]interface{}, interface{}) interface{}
}

var runningList []Running

func init() {
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for range ticker.C {
			for _, running := range runningList {
				JsRequest(running.data, running.handles...)
			}
		}
	}()

}

func (s *SenderJs) Request(running func() bool, data map[string]interface{}, handles ...func(error, map[string]interface{}, interface{}) interface{}) {
	if running() {
		runningList = append(runningList, Running{
			data:    data,
			handles: handles,
		})

	}
}
