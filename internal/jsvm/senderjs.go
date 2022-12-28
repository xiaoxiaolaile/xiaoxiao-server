package jsvm

type SenderJs struct {
	Name string
}

func (s *SenderJs) Send(f func(message string)) {

}

func (s *SenderJs) Receive(map[string]interface{}) {

}
