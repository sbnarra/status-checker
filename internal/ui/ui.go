package ui

import (
	_ "embed"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var html string

//go:embed index.js
var js string

//go:embed styles.css
var css string

func Register(router *gin.Engine) {
	router.GET("/ui", func(c *gin.Context) {
		now := strconv.FormatInt(time.Now().Unix(), 10)
		c.Data(http.StatusOK, "text/html", []byte(strings.ReplaceAll(html, "t={epoch}", "t="+now)))
	})
	router.GET("/ui/index.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/javascript", []byte(js))
	})
	router.GET("/ui/styles.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css", []byte(css))
	})
}
