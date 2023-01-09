package core

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// 运行js脚本
func runDefaultScript(s Sender, str string) (goja.Value, error) {
	return runScript(newVm(), s, str)
}

// 运行js脚本
func runScript(vm *goja.Runtime, s Sender, str string) (value goja.Value, err error) {
	_ = vm.Set("s", s)
	_ = vm.Set("sender", s)
	_ = vm.Set("image", func(url string) interface{} {
		return `[CQ:image,file=` + url + `]`
	})
	_ = vm.Set("request", JsRequest)

	//reStr := `require\(['"](.*)['"]\)`
	//re := regexp.MustCompile(reStr)
	//str = re.ReplaceAllStringFunc(str, func(s string) string {
	//	s = strings.ReplaceAll(s, "\"", "'")
	//	if !strings.Contains(s, "./") {
	//		s = s[:9] + "./" + s[9:]
	//	}
	//	return s
	//})
	value, err = vm.RunString(str)
	return
}

// 创建一个js虚拟机
func newVm() *goja.Runtime {
	vm := goja.New()
	//vm.SetFieldNameMapper(myFieldNameMapper{})
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	loadConsole(vm)
	loadModules(vm)
	loadBucket(vm)
	loadSender(vm)
	loadTime(vm)
	loadFmt(vm)
	loadSillyGirl(vm)
	loadStrings(vm)
	return vm
}

func mapFileSystemSourceLoader(files map[string]string) require.SourceLoader {
	return func(path string) ([]byte, error) {
		s, ok := files[path]
		if !ok {
			return nil, require.ModuleFileDoesNotExistError
		}
		return []byte(s), nil
	}
}

// 加载存储数据
func loadBucket(vm *goja.Runtime) {
	_ = vm.Set("Bucket", func(call goja.ConstructorCall) *goja.Object {
		name := call.Argument(0).ToString().String()
		//fmt.Println("test =>", name)
		return vm.ToValue(BucketJs{
			Bucket: BoltBucket(name),
		}).(*goja.Object)
	})
}

// 加载Sender
func loadSender(vm *goja.Runtime) {
	_ = vm.Set("running", func() bool { return true })
	_ = vm.Set("Sender", func(call goja.ConstructorCall) *goja.Object {
		name := call.Argument(0).ToString().String()
		//fmt.Println("test =>", name)
		return vm.ToValue(NewSenderJs(name)).(*goja.Object)
	})
}

// 加载模块
func loadModules(vm *goja.Runtime) {
	arr := getModules()
	m := make(map[string]string)
	for _, function := range arr {
		m[function.Title] = function.Content
	}
	r := require.NewRegistry(require.WithGlobalFolders("."), require.WithLoader(mapFileSystemSourceLoader(m)))
	r.Enable(vm)

}

// 加载时间方法
func loadTime(vm *goja.Runtime) {
	_ = vm.Set("time", Time{})
}

func loadConsole(vm *goja.Runtime) {
	_ = vm.Set("console", Console{})
}
func loadFmt(vm *goja.Runtime) {
	_ = vm.Set("fmt", Fmt{})
}
func loadSillyGirl(vm *goja.Runtime) {
	_ = vm.Set("SillyGirl", func(call goja.ConstructorCall) *goja.Object {
		//name := call.Argument(0).ToString().String()
		//fmt.Println("test =>", name)
		return vm.ToValue(NewSillyGirl()).(*goja.Object)
	})
}
func loadStrings(vm *goja.Runtime) {
	_ = vm.Set("strings", &Strings{})
}
