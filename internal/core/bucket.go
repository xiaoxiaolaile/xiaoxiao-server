package core

//var Zero Bucket

//func MakeBucket(name string) Bucket {
//	if Zero == nil {
//		logs.Error("找不到存储器，开发者自行实现接口。")
//	}
//	return Zero.Copy(name)
//}

type Bucket interface {
	//Copy(string) Bucket
	Set(interface{}, interface{}) error
	//Empty() (bool, error)
	//Size() (int64, error)
	Delete() error
	Type() string
	//Buckets() ([][]byte, error)
	GetString(...interface{}) string
	GetBytes(string) []byte
	GetInt(string, ...int) int
	GetBool(string, ...bool) bool
	Foreach(func([]byte, []byte) error)
	Create(interface{}) error
	First(interface{}) error
	String() string
}

type BucketJs struct {
	Get       func(key, defaultValue string) string `json:"get"`
	Set       func(key, value string)               `json:"set"`
	Keys      func() []string                       `json:"keys"`
	DeleteAll func()                                `json:"deleteAll"`
	Name      func() string                         `json:"name"`
}

func createBucket(name string) *BucketJs {
	//fmt.Println("name => ", name)
	bucket := BoltBucket(name)
	return &BucketJs{
		Get: func(key, defaultValue string) string {
			v := bucket.GetString(key)
			if len(v) == 0 {
				if len(defaultValue) > 0 {
					return defaultValue
				}
				return ""
			}
			return v
		},
		Set: func(key, value string) {
			_ = bucket.Set(key, value)
		},
		Keys: func() []string {
			var ss []string
			bucket.Foreach(func(k, _ []byte) error {
				ss = append(ss, string(k))
				return nil
			})
			return ss
		},
		DeleteAll: func() {
			var ss []string
			bucket.Foreach(func(k, _ []byte) error {
				ss = append(ss, string(k))
				return nil
			})

			for _, s := range ss {
				_ = bucket.Set(s, "")
			}

		},
		Name: func() string {
			return string(bucket)
		},
	}
}
