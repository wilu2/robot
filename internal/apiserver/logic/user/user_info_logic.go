package user

import (
	"context"
	errors2 "errors"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type UserInfoLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserInfoLogic {
	return UserInfoLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserInfo 查看用户
func (l *UserInfoLogic) UserInfo(req *user.GetUserReq) (resp user.GetUserResp, err error) {
	var (
		tUser = query.Use(l.svcCtx.Db).User
	)

	user, err := tUser.WithContext(l.ctx).Where(tUser.ID.Eq(uint32(req.ID))).First()
	if err != nil {
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			err = errors.WithCodeMsg(code.NotFound)
			return
		} else {
			err = errors.WithCodeMsg(code.Internal)
			return
		}
	}

	copier.Copy(&resp.UserInfo, &user)
	return
}
