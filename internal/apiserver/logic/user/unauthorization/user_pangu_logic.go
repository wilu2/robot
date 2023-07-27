package unauthorization

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user/unauthorization"
	"financial_statement/internal/pkg/verify"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"financial_statement/pkg/stringx"

	"github.com/gin-gonic/gin"
)

type UserPanguLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserPanguLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserPanguLogic {
	return UserPanguLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserPangu 系统初始化创建系统管理员
func (l *UserPanguLogic) UserPangu(req *unauthorization.UserPangu) (err error) {
	var (
		tUser = query.Use(l.svcCtx.Db).User
	)
	userCount, err := tUser.WithContext(l.ctx).Count()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	if userCount > 0 {
		err = errors.WithCodeMsg(code.Forbidden, "not first user")
		return
	}

	var account string
	if req.Account == "" {
		account = "admin"
	} else {
		account = req.Account
	}

	salt := stringx.RandString(8)
	newUser := model.User{
		Name:     req.Name,
		Account:  account,
		Salt:     salt,
		Password: verify.CalcPassword(req.Password, salt),
	}
	err = tUser.WithContext(l.ctx).Create(&newUser)
	if err != nil {
		log.Errorf("create user error: %s", err)
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	//TODO: add user role
	return
}
