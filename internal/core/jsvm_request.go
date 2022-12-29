package core

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/go-resty/resty/v2"
	"strings"
	"time"
)

var client *resty.Client

func init() {
	client = resty.New() // 创建一个restry客户端
	client.RemoveProxy()

}

func JsRequest(wt interface{}, handles ...func(error, map[string]interface{}, interface{}) interface{}) interface{} {
	//默认超时一分钟
	client.SetTimeout(60 * time.Second)
	defer client.RemoveProxy()
	var method = resty.MethodGet
	var url = ""
	var isJson bool
	request := client.R()
	switch wt.(type) {
	case string:
		url = wt.(string)
	default:
		props := wt.(map[string]interface{})
		for i := range props {
			switch strings.ToLower(i) {
			case "timeout":
				timeout := time.Duration(Int64(props[i]) * 1000 * 1000)
				client.SetTimeout(timeout)
			case "headers":
				headers := make(map[string]string)
				hds := props[i].(map[string]interface{})
				for s, v := range hds {
					headers[s] = fmt.Sprintf("%v", v)
				}
				request.SetHeaders(headers)
			case "method":
				method = strings.ToUpper(props[i].(string))
			case "url":
				switch props[i].(type) {
				case string:
					url = props[i].(string)
				case func(call goja.FunctionCall) goja.Value:
					f := props[i].(func(call goja.FunctionCall) goja.Value)
					call := goja.FunctionCall{}
					url = f(call).String()
					//if fn, ok := goja.AssertFunction(f); ok {
					//	v, _ := fn(nil)
					//	url = v.String()
					//}
				}
			case "json":
				isJson = props[i].(bool)
			case "datatype":
				switch props[i].(type) {
				case string:
					switch strings.ToLower(props[i].(string)) {
					case "json":
						isJson = true
						//case "location":
						//	location = true
					}
				}
			case "body":
				body := ""
				if v, ok := props[i].(string); !ok {
					d, _ := json.Marshal(props[i])
					body = string(d)
					request.SetHeader("Content-Type", "application/json")
				} else {
					body = v
				}
				request.SetBody(body)
			case "formdata":
				data := make(map[string]string)
				formData := props[i].(map[string]interface{})
				for s, v := range formData {
					data[s] = fmt.Sprintf("%v", v)
				}
				request.SetFormData(data)
			case "proxyurl":
				proxyUrl := props[i].(string)
				client.SetProxy(proxyUrl)
			}
		}
	}

	//if location {
	//	req.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
	//		return http.ErrUseLastResponse
	//	})
	//	rsp, err := req.Response()
	//	if err == nil && (rsp.StatusCode == 301 || rsp.StatusCode == 302) {
	//		return rsp.Header.Get("Location")
	//	} else
	//	//非重定向,允许用户自定义判断
	//	if len(handles) == 0 {
	//		return err
	//	}
	//}
	//if useproxy && Transport != nil {
	//	req.SetTransport(Transport)
	//}

	client.SetProxy("http://127.0.0.1:7890")
	rsp, err := request.Execute(method, url)
	rspObj := map[string]interface{}{}
	var bd interface{}
	if err == nil {
		rspObj["status"] = rsp.StatusCode()
		rspObj["statusCode"] = rsp.StatusCode()
		data := rsp.Body()
		if isJson {
			var v interface{}
			_ = json.Unmarshal(data, &v)
			bd = v
		} else {
			bd = string(data)
		}
		rspObj["body"] = bd
		rspObj["headers"] = rsp.Header()
	}
	if len(handles) > 0 {
		return handles[0](err, rspObj, bd)
	} else {
		return rspObj
	}
}
