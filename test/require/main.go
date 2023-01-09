package main

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	logs "github.com/sirupsen/logrus"
)

func mapFileSystemSourceLoader(files map[string]string) require.SourceLoader {
	return func(path string) ([]byte, error) {
		s, ok := files[path]
		if !ok {
			return nil, require.ModuleFileDoesNotExistError
		}
		return []byte(s), nil
	}
}

func TestStrictModule() {
	const SCRIPT = `
	var m = require("m");
	m.test();
	`

	const MODULE = `
	"use strict";

	function test() {
		var a = "passed1";
		eval("var a = 'not passed'");
		return a;
	}

	exports.test = test;
	`

	vm := goja.New()

	registry := require.NewRegistry(require.WithGlobalFolders("."), require.WithLoader(func(name string) ([]byte, error) {
		if name == "m" {
			return []byte(MODULE), nil
		}
		return nil, errors.New("Module does not exist")
	}))
	registry.Enable(vm)

	v, err := vm.RunString(SCRIPT)
	if err != nil {
		logs.Error(err)
	}

	if !v.StrictEquals(vm.ToValue("passed1")) {
		logs.Printf("Unexpected result: %v", v)
	}
}

func main() {
	vm := goja.New()
	r := require.NewRegistry(require.WithGlobalFolders("."), require.WithLoader(mapFileSystemSourceLoader(map[string]string{
		"CQ码": `
		module.exports = {  
			hello:hello,
		}

		function hello (){
			return 1
		}
`,
		"b": `exports.done = 2;`,
	})))
	r.Enable(vm)

	str := `
	const a = require("CQ码");
	const b = require('b');
	a.hello() + b.done;

`

	//require\(['"](.*)['"]\)
	//reStr := `require\(['"](.*)['"]\)`
	//re := regexp.MustCompile(reStr)
	//str = re.ReplaceAllStringFunc(str, func(s string) string {
	//	s = strings.ReplaceAll(s, "\"", "'")
	//	if !strings.Contains(s, "./") {
	//		s = s[:9] + "./" + s[9:]
	//	}
	//	return s
	//})
	res, _ := vm.RunString(str)

	v := res.Export()
	fmt.Println(v)

	//re := regexp.MustCompile("a(x*)b")
	//fmt.Println(re.ReplaceAllLiteralString("-ab-axxb-", "T"))    //-T-T-
	//fmt.Println(re.ReplaceAllLiteralString("-ab-axxb-", "$1"))   // -$1-$1-
	//fmt.Println(re.ReplaceAllLiteralString("-ab-axxb-", "${1}")) // -${1}-${1}-
	//
	////这里$1表示的是每一个匹配的第一个分组匹配结果
	////这里第一个匹配的第一个分组匹配为空,即将匹配的ab换为空值；
	////第二个匹配的第一个分组匹配为xx,即将匹配的axxb换为xx
	//fmt.Println(re.ReplaceAllString("-ab-axxb-", "$1"))    //--xx-
	//fmt.Println(re.ReplaceAllString("-ab-axxb-", "${1}w")) //-w-xxw

	//TestStrictModule()
}
