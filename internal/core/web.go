package core

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var server *gin.Engine

func ServerRun(add ...string) {
	_ = server.Run(add...)
}

func initWeb() {

	scriptStr := `
	/**
	* @author 猫咪
	* @origin 傻妞官方
	* @create_at 2022-09-10 07:35:00
	* @description 🐮网络开发demo，基础操作演示，能懂多少看悟性。
	* @version v1.0.1
	* @title 网络开发教程
	* @on_start true
	* @icon https://www.expressjs.com.cn/images/favicon.png
	* @public false
	*/
	
	const app = require("express")
	
	//客户查询ip地址
	app.get("/myip", (req, res) => {
		console.log("hello")
	   res.json({
	       data: {
	           ip: req.ip(),
	       },
	       success: true,
	   })
	});
	`

	scriptStr2 := `
	/**
	* @author 猫咪
	* @origin 测试下
	* @create_at 2022-09-10 07:35:00
	* @description 🐮网络开发demo，基础操作演示，能懂多少看悟性。
	* @version v1.0.1
	* @title 网络开发教程
	* @on_start true
	* @icon https://www.expressjs.com.cn/images/favicon.png
	* @public false
	*/
	
	const app = require("express")
	
	//客户查询ip地址
	app.get("/hello", (req, res) => {
	   res.json({
	       data: {
	           ip: "hello => " + req.ip(),
	       },
	       success: true,
	   })
	});
	`

	createPlugins()

	keyMap := initServerPlugin(scriptStr, scriptStr2)

	server.GET("/list", func(c *gin.Context) {
		c.JSON(200, getPlugins())
	})

	server.GET("/rule", func(c *gin.Context) {
		c.JSON(200, getRules())
	})
	server.GET("/module", func(c *gin.Context) {
		c.JSON(200, getModules())
	})
	server.GET("/cron", func(c *gin.Context) {
		c.JSON(200, getCron())
	})
	server.GET("/server", func(c *gin.Context) {
		c.JSON(200, getServers())
	})

	server.NoRoute(func(c *gin.Context) {
		//patchPostForm(c)
		path := c.Request.URL.Path
		method := strings.ToLower(c.Request.Method)
		//fmt.Println(method, path)
		key := method + "-" + path
		if s, ok := keyMap[key]; ok {
			s.Handle(c)
		} else {
			c.JSON(http.StatusOK, "页面不存在")
		}
	})
}

/*
*
初始化server插件
*/
func initServerPlugin(scripts ...string) map[string]*WebService {
	keyMap := make(map[string]*WebService)
	for _, script := range scripts {
		vm := newVm()
		require.RegisterNativeModule("express", func(runtime *goja.Runtime, module *goja.Object) {
			o := module.Get("exports").(*goja.Object)
			for _, m := range []string{"get", "post", "delete", "put"} {
				mm := m
				_ = o.Set(mm, func(relativePath string, handle func(*goja.Object, *goja.Object)) {
					key := mm + "-" + relativePath
					keyMap[key] = newWebService(vm, handle)
				})
			}
		})
		_, err := vm.RunString(script)
		if err != nil {
			//c.String(http.StatusBadGateway, err.Error())
			fmt.Println(err)
			//return
		}
	}
	return keyMap
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

func newVm() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(myFieldNameMapper{})
	registry := new(require.Registry) // this can be shared by multiple runtimes
	registry.Enable(vm)
	console.Enable(vm)
	return vm
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	server = gin.New()
	server.GET("/name", func(ctx *gin.Context) {
		ctx.String(200, "sillyGirl")
	})

}

type WebService struct {
	status     int
	content    string
	isJson     bool
	isRedirect bool
	vm         *goja.Runtime
	handle     func(*goja.Object, *goja.Object)
}

func newWebService(vm *goja.Runtime, handle func(*goja.Object, *goja.Object)) *WebService {
	return &WebService{
		status: http.StatusOK,
		vm:     vm,
		handle: handle,
	}
}

type Response struct {
	Send       func(goja.Value)                     `json:"send"`
	SendStatus func(int)                            `json:"sendStatus"`
	Json       func(...interface{})                 `json:"json"`
	Header     func(string, string)                 `json:"header"`
	Render     func(string, map[string]interface{}) `json:"render"`
	Redirect   func(...interface{})                 `json:"redirect"`
	Status     func(int) goja.Value                 `json:"status"`
	GetStatus  func() int                           `json:"getStatus"`
	IsComplete func() bool                          `json:"isComplete"`
	SetCookie  func(string, string, ...interface{}) `json:"setCookie"`
}

type Request struct {
	Body        func() string              `json:"body"`
	Json        func() interface{}         `json:"json"`
	IP          func() string              `json:"ip"`
	OriginalUrl func() string              `json:"originalUrl"`
	Query       func(string) string        `json:"query"`
	Querys      func() map[string][]string `json:"querys"`
	PostForm    func(string) string        `json:"postForm"`
	PostForms   func() map[string][]string `json:"postForms"`
	Path        func() string              `json:"path"`
	Header      func(string) string        `json:"header"`
	Headers     func() map[string][]string `json:"headers"`
	Method      func() string              `json:"method"`
	Cookie      func(string) string        `json:"cookie"`
}

func (s *WebService) getRes(c *gin.Context) *goja.Object {
	var res *goja.Object
	Render := func(path string, obj map[string]interface{}) {
		c.HTML(http.StatusOK, path, obj)
	}
	res = s.vm.ToValue(&Response{
		Send: func(gv goja.Value) {
			gve := gv.Export()
			switch gve := gve.(type) {
			case string:
				s.content += gve
			default:
				d, err := json.Marshal(gve)
				if err == nil {
					s.content += string(d)
					s.isJson = true
				} else {
					s.content += fmt.Sprint(gve)
				}
			}
		},
		SendStatus: func(st int) {
			s.status = st
		},
		Json: func(ps ...interface{}) {
			if len(ps) == 1 {
				d, err := json.Marshal(ps[0])
				if err == nil {
					s.content += string(d)
					s.isJson = true
				} else {
					s.content += fmt.Sprint(ps[0])
				}
			}
			s.isJson = true
		},
		Header: func(str, value string) {
			c.Header(str, value)
		},
		Render: Render,
		Redirect: func(is ...interface{}) {
			a := 302
			b := ""
			for _, i := range is {
				switch i := i.(type) {
				case string:
					b = i
				default:
					a = Int(i)
				}
			}
			c.Redirect(a, b)
			s.isRedirect = true
		},
		Status: func(i int) goja.Value {
			s.status = i
			return res
		},
		SetCookie: func(name, value string, i ...interface{}) {
			c.SetCookie(name, value, 1000*60, "/", "", false, true)
		},
		IsComplete: func() bool {
			return s.isRedirect || len(s.content) > 0
		},
		GetStatus: func() int {
			return s.status
		},
	}).(*goja.Object)
	return res
}

func (s *WebService) getReq(c *gin.Context) *goja.Object {
	var bodyData, _ = io.ReadAll(c.Request.Body)
	query := c.Request.URL.Query()
	req := s.vm.ToValue(&Request{
		Body: func() string {
			return string(bodyData)
		},
		Json: func() interface{} {
			var i interface{}
			if json.Unmarshal(bodyData, &i) != nil {
				return nil
			}
			return i
		},
		IP:          c.ClientIP,
		OriginalUrl: c.Request.URL.String,
		Query:       c.Query,
		Querys: func() map[string][]string {
			return query
		},
		PostForm: func(s string) string {
			return c.PostForm(s)
		},
		PostForms: func() map[string][]string {
			return c.Request.PostForm
		},
		Path: func() string {
			return c.Request.URL.Path
		},
		Header: c.GetHeader,
		Headers: func() map[string][]string {
			return c.Request.Header
		},
		Method: func() string {
			return c.Request.Method
		},
		Cookie: func(s string) string {
			var cookie, _ = c.Cookie(s)
			return cookie
		},
	}).(*goja.Object)
	return req
}

func (s *WebService) Handle(c *gin.Context) {
	//重置一些基础参数

	s.status = http.StatusOK
	s.content = ""
	s.isJson = true
	s.isRedirect = false

	s.handle(s.getReq(c), s.getRes(c))
	if s.isRedirect {
		return
	}
	if s.isJson {
		c.Header("Content-Type", "application/json")
	}

	c.String(s.status, s.content)
}
