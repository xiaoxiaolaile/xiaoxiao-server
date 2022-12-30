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
}

func currentUser(c *gin.Context) {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(`{
    "data": {
        "access": "admin",
        "address": "西湖区工专路 77 号",
        "avatar": "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
        "country": "China",
        "email": "cdle@apple.com",
        "group": "蚂蚁金服－某某某事业群－某某平台部－某某技术部－UED",
        "name": "小小",
        "notifyCount": 0,
        "phone": "0752-268888888",
        "plugins": [{
            "path": "/script/abb09ae0-3019-11ed-8899-52540066b468",
            "name": "插件开发",
            "component": "./Script",
            "create_at": "2034-09-12 19:14:23"
        }, {
            "path": "/script/7642f5de-3300-11ed-8a79-52540066b468",
            "name": "老版命令 💫",
            "component": "./Script",
            "create_at": "2033-09-12 19:14:24"
        }, {
            "path": "/script/1247ec8a-80fe-11ed-94dc-f44d3060890a",
            "name": "测试listen",
            "component": "./Script",
            "create_at": "2022-12-27 19:30:07"
        }, {
            "path": "/script/cf0831b0-7ab1-11ed-8d47-52540066b468",
            "name": "阿克登录",
            "component": "./Script",
            "create_at": "2022-12-13 17:52:42"
        }, {
            "path": "/script/5eea689e-7054-11ed-80c2-52540066b468",
            "name": "Telegram Bot 💫",
            "component": "./Script",
            "create_at": "2022-11-30 22:47:06"
        }, {
            "path": "/script/b050abc7-7885-11ed-867f-f44d3060890a",
            "name": "广汽传祺-取token",
            "component": "./Script",
            "create_at": "2022-11-28 23:03:34"
        }, {
            "path": "/script/85c9a37f-6e5d-11ed-9170-52540066b468",
            "name": "微信订阅号 💫",
            "component": "./Script",
            "create_at": "2022-11-28 22:41:56"
        }, {
            "path": "/script/fcb5c563-6a24-11ed-a5b3-52540066b468",
            "name": "千寻 💫",
            "component": "./Script",
            "create_at": "2022-11-22 14:21:00"
        }, {
            "path": "/script/7949f358-6b96-11ed-904d-52540066b468",
            "name": "可爱猫 💫",
            "component": "./Script",
            "create_at": "2022-11-22 14:21:00"
        }, {
            "path": "/script/3f4b19ce-64f7-11ed-ab1b-52540066b468",
            "name": "聊天机器人接入 💫",
            "component": "./Script",
            "create_at": "2022-11-22 10:41:00"
        }, {
            "path": "/script/1bee075c-3d80-11ed-8fe7-52540066b468",
            "name": "CryptoJS 🔧",
            "component": "./Script",
            "create_at": "2022-10-01 16:30:33"
        }, {
            "path": "/script/d5a55c13-37ec-11ed-b91b-aaaa00117a5c",
            "name": "定时抽奖",
            "component": "./Script",
            "create_at": "2022-09-19 22:53:36"
        }, {
            "path": "/script/78b15932-334f-11ed-8b59-aaaa00117a5c",
            "name": "比价文案 🔧",
            "component": "./Script",
            "create_at": "2022-09-14 10:59:23"
        }, {
            "path": "/script/fb3607cf-3344-11ed-9e54-52540066b468",
            "name": "运行脚本",
            "component": "./Script",
            "create_at": "2022-09-13 18:08:16"
        }, {
            "path": "/script/c12617cf-325f-11ed-9f3e-aaaa00117a5c",
            "name": "时间处理",
            "component": "./Script",
            "create_at": "2022-09-12 17:36:36"
        }, {
            "path": "/script/ccc4bbc0-3187-11ed-b60c-aaaa00117a5c",
            "name": "加密脚本 🔒",
            "component": "./Script",
            "create_at": "2022-09-12 14:56:04"
        }, {
            "path": "/script/35a4a388-3046-11ed-974c-52540066b468",
            "name": "网络开发教程 💫",
            "component": "./Script",
            "create_at": "2022-09-10 07:35:00"
        }, {
            "path": "/script/cfd18d5c-2f7d-11ed-ab94-fc7c02eb3d87",
            "name": "qinglong",
            "component": "./Script",
            "create_at": "2022-09-09 16:30:33"
        }, {
            "path": "/script/84ca21f8-32ed-11ed-99a1-fc7c02eb3d87",
            "name": "something",
            "component": "./Script",
            "create_at": "2022-09-09 16:30:33"
        }, {
            "path": "/script/0a32b49c-3107-11ed-aac0-fc7c02eb3d87",
            "name": "无名脚本",
            "component": "./Script",
            "create_at": "2022-09-08 15:06:22"
        }, {
            "path": "/script/f9087971-2aa4-11ed-b8cd-52540066b468",
            "name": "摸鱼",
            "component": "./Script",
            "create_at": "2022-09-07 13:34:03"
        }, {
            "path": "/script/d97acc0d-6a2d-11ed-b47b-52540066b468",
            "name": "CQ码 🔧",
            "component": "./Script",
            "create_at": "2021-11-22 16:12:01"
        }, {
            "path": "/script/d71cf56f-85d9-11ed-aac1-f44d3060890a",
            "name": "+新增脚本",
            "component": "./Script",
            "create_at": ""
        }],
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
	json := make(map[string]string) //注意该结构接受的内容
	_ = c.ShouldBind(&json)
	log.Printf("%v", &json)
	username := json["username"]
	password := json["password"]

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
