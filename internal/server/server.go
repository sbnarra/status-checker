package server

import (
	"status-checker/internal/api"
	"status-checker/internal/ui"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(addr string) error {
	router := newRouter()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// router.GET("/check", api.GetChecks)
	router.GET("/history", api.GetHistory)
	router.GET("/history/:check", api.GetHistoryByCheck)

	ui.Register(router)
	return router.Run(addr)
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	return router
}
