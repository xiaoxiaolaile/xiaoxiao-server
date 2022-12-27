package main

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	logs "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"xiaoxiao/internal/runtime"
)

func init() {
	logs.SetFormatter(&logs.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// Only log the warning severity or above.
	logs.SetLevel(logs.TraceLevel)
}

type myFieldNameMapper struct{}

func (tfm myFieldNameMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	tag := f.Tag.Get(`json`)
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	if parser.IsIdentifier(tag) {
		return tag
	}
	return f.Name //uncapitalize()
}

func (tfm myFieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return m.Name //uncapitalize(m.Name)
}

func main() {
	runtime.InitStore()
	vm := goja.New()
	vm.SetFieldNameMapper(myFieldNameMapper{})
	_ = vm.Set("console", _console)
	//_ = vm.Set("Bucket", createBucket)
	_ = vm.Set("Bucket", func(call goja.ConstructorCall) *goja.Object {
		name := call.Argument(0).ToString().String()
		//fmt.Println("test =>", name)
		return vm.ToValue(createBucket(name)).(*goja.Object)
	})
	_, _ = vm.RunScript("hello.js", `
	console.log("hello")
	const bucket = new Bucket("sillyGirl")
	console.log(bucket.get("name", "傻妞是谁"))
	console.log(bucket.keys())
	console.log(bucket.name())
	
`)

}

var _console = map[string]func(...interface{}){
	"info": func(v ...interface{}) {
		logs.Info(v...)
	},
	"debug": func(v ...interface{}) {
		logs.Debug(v...)
	},
	"warn": func(v ...interface{}) {
		logs.Warn(v...)
	},
	"error": func(v ...interface{}) {
		logs.Error(v...)
	},
	"log": func(v ...interface{}) {
		logs.Info(v...)
	},
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
	bucket := runtime.BoltBucket(name)
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
