package unauthorization

import (
	"context"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

type UserLogoutLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogoutLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserLogoutLogic {
	return UserLogoutLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserLogout 用户登出
func (l *UserLogoutLogic) UserLogout() (err error) {
	var (
		tLoginToken = query.Use(l.svcCtx.Db).LoginToken
	)

	token := l.ginCtx.GetHeader("x-token")
	if len(token) == 0 {
		return
	}
	_, err = tLoginToken.WithContext(l.ctx).Where(tLoginToken.Token.Eq(token)).Delete()
	return
}
