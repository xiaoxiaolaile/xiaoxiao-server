package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	str := "{\"code\":142,\"jmp_url\":\"https://wxapplogin.m.jd.com/h5/risk/select?token=EcZ1ybVPS9tXan8NpHwveOwVdysevHK9&client_type=wxapp&guid=c1d4d36123beb6b24214f1d96b416dbbdb58be52f4191cdeeabb73478ae971df&returnurl=https%3A%2F%2Fwxapplogin.m.jd.com%2Fcgi-bin%2Fml%2Fwxapp_redirect%3Freturnurl%3D%252Fpages%252Flogin%252Fweb-view%252Fweb-view\",\"msg\":\"\\u8bf7\\u70b9\\u51fb https://wxapplogin.m.jd.com/h5/risk/select?token=EcZ1ybVPS9tXan8NpHwveOwVdysevHK9&client_type=wxapp&guid=c1d4d36123beb6b24214f1d96b416dbbdb58be52f4191cdeeabb73478ae971df&returnurl=https%3A%2F%2Fwxapplogin.m.jd.com%2Fcgi-bin%2Fml%2Fwxapp_redirect%3Freturnurl%3D%252Fpages%252Flogin%252Fweb-view%252Fweb-view \\u9a8c\\u8bc1\\u540e\\u91cd\\u65b0\\u767b\\u9646\"}"
	fmt.Println(str)
	s := unicode2utf8(str)
	fmt.Println(s)
}

func unicode2utf8(source string) string {
	var res = []string{""}
	sUnicode := strings.Split(source, "\\u")
	var context = ""
	for _, v := range sUnicode {
		var additional = ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}
