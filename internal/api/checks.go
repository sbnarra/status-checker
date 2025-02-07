package api

import (
	"net/http"
	"status-checker/internal/checker"

	"github.com/gin-gonic/gin"
)

func GetChecks(c *gin.Context) {
	if checks, err := checker.ReadConfig(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	} else {
		c.IndentedJSON(http.StatusOK, checks)
	}
}
