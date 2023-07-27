package user

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user"
	"financial_statement/internal/pkg/verify"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/stringx"

	"github.com/gin-gonic/gin"
)

type UserUpdateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserUpdateLogic {
	return UserUpdateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserUpdate 修改用户信息
func (l *UserUpdateLogic) UserUpdate(req *user.UpdateUserReq) (err error) {
	var (
		tUser = query.Use(l.svcCtx.Db).User
	)
	salt := stringx.RandString(8)
	newPsw := verify.CalcPassword(*req.Password, salt)
	_, err = tUser.WithContext(l.ctx).Where(tUser.ID.Eq(uint32(req.ID))).
		UpdateSimple(tUser.Password.Value(newPsw), tUser.Salt.Value(salt))
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	return
}
