package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var server *gin.Engine

func ServerRun(add ...string) {
	_ = server.Run(add...)
}

var keyMap map[string]*WebService

func initWeb() {

	keyMap = initServerPlugin(getServers()...)

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
	//server.Static("/", "")
}

func init() {
	//gin.SetMode(gin.ReleaseMode)
	server = gin.New()
	server.Use(gin.Logger())
	server.GET("/name", func(ctx *gin.Context) {
		ctx.String(200, "sillyGirl")
	})

}
