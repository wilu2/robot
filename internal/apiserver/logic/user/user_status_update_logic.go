package user

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/pkg/errors"
	"github.com/gin-gonic/gin"
	"time"
)

type UserStatusUpdateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserStatusUpdateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserStatusUpdateLogic {
	return UserStatusUpdateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserStatusUpdate 修改用户状态，有效期
func (l *UserStatusUpdateLogic) UserStatusUpdate(req *user.UpdateUserStatusReq) (err error) {
	var (
		tUser = query.Use(l.svcCtx.Db).User
		u = l.ginCtx.Keys["user"].(*model.User)
	)
	//验证管理员权限
	
	if u.Account != "admin"{
		err = errors.WithCodeMsg(code.BadRequest, "非管理员无法更改")
		return
	}

	//更新用户状态
	if req.Status == 1{
		_, err = tUser.WithContext(l.ctx).Where(tUser.ID.Eq(uint32(req.ID))).
			UpdateSimple(tUser.Status.Value(req.Status), tUser.ExpiryTime.Value(req.ExpiryTime), tUser.LastLoginTime.Value(time.Now().Unix()))
	}else{
		_, err = tUser.WithContext(l.ctx).Where(tUser.ID.Eq(uint32(req.ID))).
			UpdateSimple(tUser.Status.Value(req.Status), tUser.ExpiryTime.Value(req.ExpiryTime))
	}
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	return
}
