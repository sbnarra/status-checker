package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	message string
}

func onError(c *gin.Context, err error) {
	c.IndentedJSON(http.StatusInternalServerError, Error{
		message: err.Error(),
	})
}
