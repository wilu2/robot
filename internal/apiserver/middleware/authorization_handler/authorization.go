package authorizationhandler

import "github.com/gin-gonic/gin"

const (
	Default     = "default"
	UrlToken    = "url-token"
	CAS         = "cas"
	Unauthorize = "unauthorize"
)

type IAuthorization interface {
	Authorization(c *gin.Context)
	unauthorized(c *gin.Context)
}
