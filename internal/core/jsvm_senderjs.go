package core

import (
	"time"
)

type SenderJs struct {
	Name string
	f    func(data interface{})
}

func (s *SenderJs) Send(f func(data interface{})) {
	s.f = f
}

func (s *SenderJs) Receive(data map[string]interface{}) {
	if s.f != nil {
		//解析命令，修改命令
		//if str, ok := data["content"]; ok {
		//	//进行解析，和替换内容
		//
		//}
		s.f(data)
	}
}

func (s *SenderJs) Request(running func() bool, data map[string]interface{}, handles ...func(error, map[string]interface{}, interface{}) interface{}) {
	if running() {
		JsRequest(data, handles...)
		ticker := time.NewTicker(time.Millisecond * 10000)
		for range ticker.C {
			JsRequest(data, handles...)
		}

	}
}
