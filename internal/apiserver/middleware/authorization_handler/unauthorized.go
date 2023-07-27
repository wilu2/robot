// 不验证的授权处理模块

package authorizationhandler

import (
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/response"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type UnAuthorization struct {
}

//不验证任何权限的认证方式（取id=1的admin账户当作操作账户）
func (d *UnAuthorization) Authorization(c *gin.Context) {
	dbIns, _ := dal.GetDbFactoryOr(nil)
	user := &model.User{}
	dbIns.GetDb().Where(&model.User{ID: 1}).First(
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

func (d *UnAuthorization) unauthorized(c *gin.Context) {
	err := errors.WithCodeMsg(code.Unauthorized)
	response.HandleResponse(c, nil, err)
	c.Abort()
}
