package unauthorization

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user/unauthorization"
	"financial_statement/internal/pkg/verify"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/stringx"

	"github.com/gin-gonic/gin"
)

type UserUpdatePwdLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdatePwdLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserUpdatePwdLogic {
	return UserUpdatePwdLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserUpdatePwd 不使用token修改用户密码，需要旧密码
func (l *UserUpdatePwdLogic) UserUpdatePwd(req *unauthorization.UpdatePwdReq) (err error) {
	var (
		tUser      = query.Use(l.svcCtx.Db).User
	)
	user, err := tUser.WithContext(l.ctx).Select(tUser.ID,tUser.Password,tUser.Salt).Where(tUser.Account.Eq(req.Account)).First()
	if err != nil {
		err = errors.WithCodeMsg(code.BadRequest, "用户名或密码错误")
		return
	}
	if !verify.VerifyPassword(user, req.OldPassword) {
		err = errors.WithCodeMsg(code.BadRequest, "用户名或密码错误")
		return
	}
	salt := stringx.RandString(8)
	newPsw := verify.CalcPassword(req.NewPassword, salt)
	_, err = tUser.WithContext(l.ctx).Where(tUser.ID.Eq(uint32(user.ID))).
		UpdateSimple(tUser.Password.Value(newPsw), tUser.Salt.Value(salt))
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	return
}
