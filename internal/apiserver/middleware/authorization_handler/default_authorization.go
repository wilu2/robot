package authorizationhandler

import (
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/response"
	"financial_statement/pkg/errors"
	"time"

	"github.com/gin-gonic/gin"
)

type DefaultAuthorization struct {
}

//默认的认证方式
func (d *DefaultAuthorization) Authorization(c *gin.Context) {
	token := c.GetHeader("x-token")
	if len(token) == 0 {
		d.unauthorized(c)
		return
	}
	dbIns, _ := dal.GetDbFactoryOr(nil)
	loginToken := &model.LoginToken{}
	dbIns.GetDb().Where(&model.LoginToken{Token: token}).First(loginToken)
	if loginToken.UserID == 0 || loginToken.Expiry < time.Now().Unix() {
		//没有对应的token记录,或者token已过期
		d.unauthorized(c)
		return
	}
	user := &model.User{}
	dbIns.GetDb().Where(&model.User{ID: loginToken.UserID}).First(
		user)
	if user.ID == 0 {
		//没有该用户
		d.unauthorized(c)
		return
	}
	c.Set("user", user)
	c.Set("is_sso", true)
	c.Next()
}

func (d *DefaultAuthorization) unauthorized(c *gin.Context) {
	err := errors.WithCodeMsg(code.Unauthorized)
	response.HandleResponse(c, nil, err)
	c.Abort()
}
