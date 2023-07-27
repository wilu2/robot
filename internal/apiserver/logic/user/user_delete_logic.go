package user

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user"
	"financial_statement/pkg/errors"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDeleteLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserDeleteLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserDeleteLogic {
	return UserDeleteLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserDelete 删除用户
func (l *UserDeleteLogic) UserDelete(req *user.DeleteUserReq) (err error) {
	var (
		tUser       = query.Use(l.svcCtx.Db).User
		u           = l.ginCtx.Keys["user"].(*model.User)
		tLoginToken = query.Use(l.svcCtx.Db).LoginToken
		q           = query.Use(l.svcCtx.Db)
	)
	if u.ID == uint32(req.ID) {
		err = errors.WithCodeMsg(code.BadRequest, "Can not delete yourself")
		return
	}

	// update user delete tag and set user token's expired
	err = q.Transaction(func(tx *query.Query) error {
		_, err := tUser.WithContext(l.ctx).Where(tUser.ID.Eq(uint32(req.ID))).Delete()
		if err != nil {
			return errors.WithCodeMsg(code.Internal)
		}
		_, err = tLoginToken.WithContext(l.ctx).Where(tLoginToken.UserID.Eq(uint32(req.ID))).Update(tLoginToken.Expiry, time.Now().Unix())
		if err != nil {
			return errors.WithCodeMsg(code.Internal)
		}
		return nil
	})

	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	return
}
