package jsvm

import (
	"fmt"
	"math/rand"
)

type SillyGirlJs struct {
	//BucketGet  func(bucket, key string) string                 `json:"bucketGet"`
	//BucketSet  func(bucket, key, value string)                 `json:"bucketSet"`
	//BucketKeys func(bucket string) []string                    `json:"bucketKeys"`
	//Push       func(obj map[string]interface{})                `json:"push"`
	//Session    func(wt interface{}) func(...int) SessionResult `json:"session"`
	//Call       func(key string) interface{}                    `json:"call"`

	defaultUserId string
}

func NewSillyGirl() *SillyGirlJs {
	defaultUserId := fmt.Sprintf("carry_%d", rand.Int63())
	return &SillyGirlJs{
		defaultUserId: defaultUserId,
	}
}

type SessionResult struct {
	HasNext bool
	Message string
}

func getBucket(name string) *BucketJs {
	return &BucketJs{Bucket: BoltBucket(name)}
}

func (s *SillyGirlJs) BucketGet(bucket, key string) string {
	return getBucket(bucket).Get(key, "")
}
func (s *SillyGirlJs) BucketSet(bucket, key, value string) {
	getBucket(bucket).Set(key, value)
}
func (s *SillyGirlJs) BucketKeys(bucket string) []string {
	return getBucket(bucket).Keys()
}
func (s *SillyGirlJs) Push(obj map[string]interface{}) {
	imType := obj["imType"].(string)
	groupCode := 0
	var userID interface{}
	if _, ok := obj["groupCode"]; ok {
		groupCode = Int(obj["groupCode"])
	} else {
		userID = obj["userID"]
	}
	content := obj["content"].(string)
	if groupCode != 0 {
		if push, ok := GroupPushs[imType]; ok {
			push(groupCode, userID, content, "")
		}
	} else {
		if push, ok := Pushs[imType]; ok {
			push(userID, content, nil, "")
		}
	}
}
func (s *SillyGirlJs) Session(info interface{}) func(...int) SessionResult {
	//userId := s.defaultUserId
	//msg := ""
	//imTpye := "carry"
	//chatId := 0
	//switch info.(type) {
	//case string:
	//	msg = info.(string)
	//default:
	//	props := info.(map[string]interface{})
	//	for i := range props {
	//		switch strings.ToLower(i) {
	//		case "imtype":
	//			imTpye = props[i].(string)
	//		case "msg":
	//			msg = props[i].(string)
	//		case "chatid":
	//			chatId = Int(props[i])
	//		case "userid":
	//			userId = props[i].(string)
	//		}
	//	}
	//}
	//if msg == "" {
	//	return nil
	//}
	//c := &Faker{
	//	Type:    imTpye,
	//	Message: msg,
	//	Carry:   make(chan string),
	//	UserId:  userId,
	//	ChatId:  chatId,
	//}
	//Senders <- c
	//var f = func(i ...int) SessionResult {
	//	timeOut := 1000 * 100
	//	if len(i) > 0 {
	//		timeOut = i[0]
	//	}
	//	select {
	//	case v, ok := <-c.Listen():
	//		return SessionResult{
	//			HasNext: ok,
	//			Message: v,
	//		}
	//	case <-time.After(time.Millisecond * time.Duration(timeOut)):
	//		return SessionResult{
	//			HasNext: false,
	//			Message: "已超时",
	//		}
	//	}
	//}
	//return f
	return nil
}
func (s *SillyGirlJs) Call(key string) interface{} {
	if f, ok := OttoFuncs[key]; ok {
		return f
	}
	return nil
}

var OttoFuncs = map[string]interface{}{
	//"machineId": func(_ string) string {
	//	return GetMachineID()
	//},
	//"uuid": func(_ string) string {
	//	return utils.GenUUID()
	//},
	//"md5": utils.Md5,
	//"timeFormat": func(str string) string {
	//	return time.Now().Format(str)
	//},
	//"now": func() string {
	//	return time.Now().Format("2000-01-01 00:00:00")
	//},
	//"timeFormater": func(time time.Time, format string) string {
	//	return time.Format(format)
	//},
}
