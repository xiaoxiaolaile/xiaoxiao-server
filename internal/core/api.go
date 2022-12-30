package core

import (
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"log"
)

var defaultUserName string
var defaultPassword string

func initApi() {
	sillyGirl := BoltBucket("sillyGirl")
	defaultUserName = sillyGirl.GetString("name", "小小")
	defaultPassword = sillyGirl.GetString("password", GenUUID())

	logs.Printf("可视化面板临时账号密码：%s %s", defaultUserName, defaultPassword)

	g := server.Group("/api")
	g.POST("/login/account", userLogin)

	g.Use(userInterceptor())
}

func userInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		//// 签名校验逻辑，逻辑调用自己的签名逻辑即可
		//// checkSign := signature.CheckSign(c)
		//// if !checkSign {
		//// 	c.Abort()
		//// 	return
		//// }
		//request := c.Request
		//_ = request.ParseForm()
		//// 获取参数值
		//req, ok := getReqParam(c, request)
		//if ok {
		//	return
		//}
		//// 获取响应结果
		//blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		//c.Request = request
		//c.Writer = blw
		//// 执行下一个
		//c.Next()
		//// 获取 响应
		//realIp := util.GetRealIp(request)
		//log.Info("请求地址：", c.FullPath(), "，请求ip：", realIp, "，Referer：", request.Referer(), "，请求参数：", req, "，响应参数：", blw.body.String())
	}
}

func userLogin(c *gin.Context) {
	//{"username":"小小","password":"123456","autoLogin":false,"type":"account"}
	json := make(map[string]string) //注意该结构接受的内容
	_ = c.ShouldBind(&json)
	log.Printf("%v", &json)
	username := json["username"]
	password := json["password"]

	logs.Printf("username = %s, password = %s", username, password)

	if defaultUserName == username && defaultPassword == password {
		c.JSON(200, LoginData{
			CurrentAuthority: "admin",
			Status:           "ok",
			Type:             "account",
		})
	} else {
		c.JSON(200, LoginData{
			CurrentAuthority: "guest",
			Status:           "error",
			Type:             "account",
		})
	}

}

type LoginData struct {
	CurrentAuthority string `json:"currentAuthority"`
	Status           string `json:"status"`
	Type             string `json:"type"`
}
