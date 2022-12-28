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

func (s *SenderJs) Send(f func(data interface{})) {
	s.f = f
}

type Im struct {
	Faker
}

func (s *SenderJs) Receive(data map[string]interface{}) {
	if s.f != nil {
		//解析命令，修改命令
		logs.Info("收到消息：", data)
		logs.Info("收到消息content：", data["content"])
		if str, ok := data["content"]; ok {
			//进行解析，和替换内容
			go func() {

				carry := make(chan string, 1)
				f := &Im{
					Faker{
						Type:     "carry",
						Platform: s.Name,
						Carry:    carry,
						Message:  fmt.Sprintf("%v", str),
						Admin:    true,
					},
				}
				Senders <- f
				timeout := time.Millisecond * 2
				for {
					select {
					case result := <-carry:
						data["content"] = result
						s.f(data)
						return
					case <-time.After(timeout):
						return
					}
				}

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
