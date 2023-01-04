package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var server *gin.Engine

func ServerRun(add ...string) {
	_ = server.Run(add...)
}

func initWeb() {

	keyMap := initServerPlugin(getServers()...)
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

	server.GET("/api/plugin/:type", func(c *gin.Context) {
		t := c.Param("type")
		name := c.Query("name")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

		db := BoltBucket("plugins")
		var functions Functions
		db.Foreach(func(k, v []byte) error {
			functions = append(functions, createPlugin(string(v)))
			return nil
		})
		f := func(d Function) bool {
			r := true
			switch t {
			case "rule":
				r = d.Rules != nil
			case "module":
				r = d.Module
			case "server":
				r = len(d.Cron) > 0
			case "cron":
				r = d.OnStart
			}
			if len(name) > 0 {
				return r && strings.Contains(d.Title, name)
			}
			return r
		}
		sort.Sort(functions)
		var list []*Function
		for _, function := range functions {
			if f(*function) {
				list = append(list, function)
			}
		}

		successList(c, "", len(list), list[(page-1)*pageSize:page*pageSize])
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
