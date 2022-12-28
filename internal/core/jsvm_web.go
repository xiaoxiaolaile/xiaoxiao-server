package core

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type WebService struct {
	status     int
	content    string
	isJson     bool
	isRedirect bool
	vm         *goja.Runtime
	handle     func(*goja.Object, *goja.Object)
}

func NewWebService(vm *goja.Runtime, handle func(*goja.Object, *goja.Object)) *WebService {
	return &WebService{
		status: http.StatusOK,
		vm:     vm,
		handle: handle,
	}
}

type Response struct {
	S *WebService
	C *gin.Context
}

func (r *Response) Send(gv goja.Value) {
	gve := gv.Export()
	switch gve := gve.(type) {
	case string:
		r.S.content += gve
	default:
		d, err := json.Marshal(gve)
		if err == nil {
			r.S.content += string(d)
			r.S.isJson = true
		} else {
			r.S.content += fmt.Sprint(gve)
		}
	}
}

func (r *Response) SendStatus(st int) {
	r.S.status = st
}
func (r *Response) Json(ps ...interface{}) {
	if len(ps) == 1 {
		d, err := json.Marshal(ps[0])
		if err == nil {
			r.S.content += string(d)
			r.S.isJson = true
		} else {
			r.S.content += fmt.Sprint(ps[0])
		}
	}
	r.S.isJson = true
}
func (r *Response) Header(str, value string) {
	r.C.Header(str, value)
}
func (r *Response) Render(path string, obj map[string]interface{}) {
	r.C.HTML(http.StatusOK, path, obj)
}
func (r *Response) Redirect(is ...interface{}) {
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
	r.C.Redirect(a, b)
	r.S.isRedirect = true
}
func (r *Response) Status(i int) goja.Value {
	r.S.status = i
	return r.S.vm.ToValue(r)
}
func (r *Response) SetCookie(name, value string) {
	r.C.SetCookie(name, value, 1000*60, "/", "", false, true)
}
func (r *Response) IsComplete() bool {
	return r.S.isRedirect || len(r.S.content) > 0
}
func (r *Response) GetStatus() int {
	return r.S.status
}

type Request struct {
	S        *WebService
	C        *gin.Context
	BodyData string
}

func (r *Request) Body() string {
	return r.BodyData
}
func (r *Request) Json() interface{} {
	var i interface{}
	if json.Unmarshal([]byte(r.BodyData), &i) != nil {
		return nil
	}
	return i
}
func (r *Request) Ip() string {
	return r.C.ClientIP()
}
func (r *Request) OriginalUrl() string {
	return r.C.Request.URL.String()
}
func (r *Request) Query(key string) string {
	return r.C.Query(key)
}
func (r *Request) Querys() map[string][]string {
	return r.C.Request.URL.Query()
}
func (r *Request) PostForm(s string) string {
	return r.C.PostForm(s)
}
func (r *Request) PostForms() map[string][]string {
	return r.C.Request.PostForm
}
func (r *Request) Path() string {
	return r.C.Request.URL.Path
}
func (r *Request) Header(key string) string {
	return r.C.GetHeader(key)
}
func (r *Request) Headers() map[string][]string {
	return r.C.Request.Header
}
func (r *Request) Method() string {
	return r.C.Request.Method
}
func (r *Request) Cookie(s string) string {
	var cookie, _ = r.C.Cookie(s)
	return cookie
}

func (s *WebService) getRes(c *gin.Context) *goja.Object {
	return s.vm.ToValue(&Response{S: s, C: c}).(*goja.Object)
}

func (s *WebService) getReq(c *gin.Context) *goja.Object {
	var bodyData, _ = io.ReadAll(c.Request.Body)
	req := s.vm.ToValue(&Request{
		S:        s,
		C:        c,
		BodyData: string(bodyData),
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
