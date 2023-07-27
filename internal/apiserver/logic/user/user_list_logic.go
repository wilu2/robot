package user

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user"
	userType "financial_statement/internal/apiserver/types/user"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserListLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserListLogic {
	return UserListLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserList 列出用户
func (l *UserListLogic) UserList(req *userType.UserListReq) (resp user.ListUserResp, err error) {
	var (
		tUser = query.Use(l.svcCtx.Db).User
	)

	userDo := tUser.WithContext(l.ctx).
		Select(tUser.ID, tUser.Name, tUser.Account, tUser.Email, tUser.Mobile, tUser.CreatedAt, tUser.UpdatedAt,tUser.ExpiryTime,tUser.Status)
	switch req.OrderBy {
	case "update_time":
		if req.OrderByType == "desc" {
			userDo = userDo.Order(tUser.UpdatedAt.Desc())
		} else {
			userDo = userDo.Order(tUser.UpdatedAt)
		}
	case "create_time":
		if req.OrderByType == "desc" {
			userDo = userDo.Order(tUser.CreatedAt.Desc())
		} else {
			userDo = userDo.Order(tUser.CreatedAt)
		}
	}

	resp.Count, err = userDo.Count()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	users, err := userDo.Offset((int(req.Page) - 1) * int(req.PerPage)).Limit(int(req.PerPage)).Find()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	copier.Copy(&resp.Users, &users)
	return
}
