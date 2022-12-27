package main

import (
	"github.com/dop251/goja"

	logs "github.com/sirupsen/logrus"
)

func init() {
	logs.SetFormatter(&logs.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// Only log the warning severity or above.
	logs.SetLevel(logs.TraceLevel)
}
func main() {

	s := testGoReflectMethod_O{
		Test: "passed",
	}
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	_ = vm.Set("console", _console)
	_ = vm.Set("s", s)
	_, _ = vm.RunScript("hello.js", `

	function test(a,  b) {
		console.info(a+b)
	}
	test(1,2)
	test(2,2)
	console.log(s.test)
`)
	v, err := vm.RunString("2 + 2")
	if err != nil {
		panic(err)
	}
	if num := v.Export().(int64); num != 4 {
		panic(num)
	}

}

type testGoReflectMethod_O struct {
	field string
	Test  string
}

func (o testGoReflectMethod_O) Method(s string) string {
	return o.field + s
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
