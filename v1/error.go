package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	errBadJSON = errors.New("could not parse JSON")
)

type errorStruct struct {
	Error string `json:"error"`
}

func throwError(c *gin.Context, status int, err error) {
	c.JSON(status, errorStruct{
		Error: err.Error(),
	})
}
