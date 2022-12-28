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
		if str, ok := data["content"]; ok {
			//进行解析，和替换内容
			go func() {
				f := &Im{
					Faker: Faker{
						Type:     "carry",
						Platform: s.Name,
						Message:  fmt.Sprintf("%v", str),
						Admin:    true,
					},
					data: data,
					s:    s,
				}
				Senders <- f

			}()
			go func() {
				f := &Faker{
					Type:     "terminal",
					Platform: "terminal",
					Message:  fmt.Sprintf("%v", str),
					Admin:    true,
				}
				Senders <- f
			}()
		}

	}
}

func (s *SenderJs) Request(running func() bool, data map[string]interface{}, handles ...func(error, map[string]interface{}, interface{}) interface{}) {
	go func() {
		if running() {
			JsRequest(data, handles...)
			ticker := time.NewTicker(time.Millisecond * 5000)
			for range ticker.C {
				JsRequest(data, handles...)
			}

		}
	}()
}
