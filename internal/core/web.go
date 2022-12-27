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

func initWeb() {

	var scripts []string

	servers := getServers()
	for _, function := range servers {
		scripts = append(scripts, function.Content)
	}
	keyMap := initServerPlugin(scripts...)

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

func init() {
	gin.SetMode(gin.ReleaseMode)
	server = gin.New()
	server.GET("/name", func(ctx *gin.Context) {
		ctx.String(200, "sillyGirl")
	})

}
