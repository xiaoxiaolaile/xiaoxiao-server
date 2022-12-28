package jsvm

import (
	"time"
)

type SenderJs struct {
	Name string
}

func (s *SenderJs) Send(f func(message string)) {

}

func (s *SenderJs) Receive(data map[string]interface{}) {

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
