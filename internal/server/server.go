package server

import (
	"net/http"
	"os"
	"status-checker/internal/api"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(addr string) error {
	router := newRouter()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// router.GET("/check", api.GetChecks)
	router.GET("/history", api.GetHistory)
	router.GET("/history/:name", api.GetHistoryByCheck)

	router.LoadHTMLFiles("ui/index.html")
	router.GET("/ui", indexPage)
	router.Static("/ui/assets", "./ui/assets")
	return router.Run(addr)
}

func indexPage(c *gin.Context) {
	hostname := ""
	if osHostname, err := os.Hostname(); err == nil {
		hostname = osHostname
	} else {
		hostname = err.Error()
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"v":        strconv.FormatInt(time.Now().Unix(), 10),
		"hostname": hostname,
	})
}

func newRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies([]string{})
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	return router
}
