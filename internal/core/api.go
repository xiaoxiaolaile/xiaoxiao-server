package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var defaultUserName string
var defaultPassword string

var loginUUid string

func initApi() {
	sillyGirl := BoltBucket("sillyGirl")
	defaultUserName = sillyGirl.GetString("name", "小小")
	defaultPassword = sillyGirl.GetString("password", GenUUID())
	loginUUid = GenUUID()

	logs.Printf("可视化面板临时账号密码：%s %s", defaultUserName, defaultPassword)

	g := server.Group("/api")
	g.POST("/login/account", userLogin)

	g.Use(userInterceptor())
	g.GET("/currentUser", currentUser)
	g.GET("/bucket", getBucketList)
}

func getBucketList(c *gin.Context) {
	list := getBucKetList()
	successRespond(c, "", list)
}

func successRespond(c *gin.Context, message string, data interface{}) {
	d := Respond{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, d)
}

type Respond struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func currentUser(c *gin.Context) {
	m := make(map[string]interface{})
	sillyGirl := BoltBucket("sillyGirl")
	defaultUserName = sillyGirl.GetString("name", "小小")
	_ = json.Unmarshal([]byte(`{
    "data": {
        "access": "admin",
        "address": "西湖区工专路 77 号",
        "avatar": "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
        "country": "China",
        "email": "cdle@apple.com",
        "group": "蚂蚁金服－某某某事业群－某某平台部－某某技术部－UED",
        "name": "`+defaultUserName+`",
        "notifyCount": 0,
        "phone": "0752-268888888",
        "plugins": [],
        "signature": "人道，损不足以奉有余",
        "title": "交互专家",
        "unreadCount": 0,
        "userid": "string"
    },
    "success": true
}`), &m)
	c.JSON(200, m)
}

func userInterceptor() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, e := context.Cookie("token")
		if token == loginUUid {
			context.Next()
		}
		if e == nil {
			context.Next()
		} else {
			context.Abort()
			context.HTML(http.StatusUnauthorized, "401.tmpl", nil)
		}

	}
}

func userLogin(c *gin.Context) {
	//{"username":"小小","password":"123456","autoLogin":false,"type":"account"}
	j := make(map[string]string) //注意该结构接受的内容
	_ = c.ShouldBind(&j)
	log.Printf("%v", &j)
	username := j["username"]
	password := j["password"]

	logs.Printf("username = %s, password = %s", username, password)

	if defaultUserName == username && defaultPassword == password {
		c.SetCookie("token", loginUUid, 24*60*60, "/", "localhost", false, true)
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
