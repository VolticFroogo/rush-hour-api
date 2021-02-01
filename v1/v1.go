package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(r *gin.Engine) {
	v1 := r.Group("/v1")

	v1.Handle(http.MethodPost, "/solve", solve)
}
