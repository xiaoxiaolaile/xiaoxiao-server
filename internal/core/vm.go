package core

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"regexp"
	"strings"
	runtime2 "xiaoxiao/internal/runtime"
)

// 运行js脚本
func runScript(str string) (goja.Value, error) {
	vm := newVm()

	reStr := `require\(['"](.*)['"]\)`
	re := regexp.MustCompile(reStr)
	str = re.ReplaceAllStringFunc(str, func(s string) string {
		s = strings.ReplaceAll(s, "\"", "'")
		if !strings.Contains(s, "./") {
			s = s[:9] + "./" + s[9:]
		}
		return s
	})
	return vm.RunString(str)
}

// 创建一个js虚拟机
func newVm() *goja.Runtime {
	vm := goja.New()
	//vm.SetFieldNameMapper(myFieldNameMapper{})
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	loadModules(vm)
	loadBucket(vm)
	loadTime(vm)
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
		return vm.ToValue(runtime2.BucketJs{
			Bucket: runtime2.BoltBucket(name),
		}).(*goja.Object)
	})
}

// 加载模块
func loadModules(vm *goja.Runtime) {
	arr := getModules()
	m := make(map[string]string)
	for _, function := range arr {
		m[function.Title] = function.Content
	}
	r := require.NewRegistry(require.WithLoader(mapFileSystemSourceLoader(m)))
	r.Enable(vm)
	console.Enable(vm)
}

// 加载时间方法
func loadTime(vm *goja.Runtime) {
	_ = vm.Set("time", runtime2.Time{})
}
