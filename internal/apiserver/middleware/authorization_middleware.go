package middleware

import (
	"financial_statement/internal/apiserver/code"
	// authMiddleware "financial_statement/internal/apiserver/middleware/authorization_handler"
	authorizationhandler "financial_statement/internal/apiserver/middleware/authorization_handler"
	"financial_statement/internal/apiserver/response"
	"financial_statement/pkg/errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	once     sync.Once
	authType string
)

func Unauthorized(c *gin.Context) {
	err := errors.WithCodeMsg(code.Unauthorized)
	response.HandleResponse(c, nil, err)
	c.Abort()
}

func AuthorizationMiddleware(c *gin.Context) {
	once.Do(func() {
		authType = viper.GetString("auth.auth-type")
	})
	var auth authorizationhandler.IAuthorization
	switch authType {
	case authorizationhandler.UrlToken:
		auth = &authorizationhandler.UrlTokenAuthorization{}
	case authorizationhandler.CAS:
		auth = &authorizationhandler.CASAuthorization{}
	case authorizationhandler.Unauthorize:
		auth = &authorizationhandler.UnAuthorization{}
	default:
		auth = &authorizationhandler.DefaultAuthorization{}
	}
	auth.Authorization(c)
}
