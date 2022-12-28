package main

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	logs "github.com/sirupsen/logrus"
	"time"
)

type Time struct {
}

func (t *Time) Now() time.Time {
	return time.Now()
}
func (t *Time) Sleep(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}
func (t *Time) Unix(usec int64) time.Time {
	return time.UnixMicro(usec)
}
func (t *Time) UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec)
}

func (t *Time) Parse(value, layout, name string) time.Time {
	if len(name) == 0 {
		name = "Asia/Shanghai"
	}
	loc, _ := time.LoadLocation(name)
	_t, _ := time.ParseInLocation(layout, value, loc)
	return _t
}

func main() {
	//now := time.Now()
	//logs.Info(now.String())
	//logs.Info(time.Parse("2006/01/02", "2024/12/12"))

	vm := goja.New()
	registry := new(require.Registry) // this can be shared by multiple runtimes
	registry.Enable(vm)
	console.Enable(vm)
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	_ = vm.Set("time", Time{})
	_, err := vm.RunScript("hello.js", `

var now = time.now()

console.log(now.string())
console.log(now.unix())
console.log(now.before(time.now().add(time.day)))
console.log(now.after(time.now().add(time.second)))
console.log(now.unixMilli())
console.log(now.format("2006-01-02"))
console.log(now.format("2006-01-02 15:04:05"))

time.sleep(1000)

console.log(time.unix(1662979463))
console.log(time.unixMilli(1662979500379))
console.log(time.parse("2024/12/12", "2006/01/02"))
console.log(time.parse("2024/12/12 10:11:12", "2006/01/02 15:04:05"))
console.log(time.parse("2024/12/12 10:11:12", "2006/01/02 15:04:05", "America/Los_Angeles"))
console.log(time.parse("2024/12/12 10:11:12", "2006/01/02 15:04:05", "Asia/Shanghai"))
`)

	if err != nil {
		logs.Error(err)
	}

	// 打点器和定时器的机制有点相似：一个通道用来发送数据。
	// 这里我们在这个通道上使用内置的 `range` 来迭代值每隔
	// 500ms 发送一次的值。
	ticker := time.NewTicker(time.Millisecond * 1000)
	for t := range ticker.C {
		fmt.Println("Tick at", t)
	}

}
