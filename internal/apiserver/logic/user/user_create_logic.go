package user

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user"
	"financial_statement/internal/pkg/database/orm"
	"financial_statement/internal/pkg/verify"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"financial_statement/pkg/stringx"
	"financial_statement/internal/pkg/settings"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserCreateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserCreateLogic {
	return UserCreateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserCreate 新建用户
func (l *UserCreateLogic) UserCreate(req *user.CreateUserReq) (err error) {
	salt := stringx.RandString(8)
	password := verify.CalcPassword(req.Password, salt)
	setting_map := settings.GetSetting()
	result := l.svcCtx.Db.Create(&model.User{
		Account:           req.Account,
		Password:          password,
		Salt:              salt,
		Name:              req.Name,
		Email:             req.Email,
		Mobile:            req.Mobile,
		ExpiryTime:        time.Now().Unix() + int64(setting_map["valid_time"]*24*60*60),
		LastLoginTime:     time.Now().Unix(),
		Status:            0,
		LastLoginFailTime: 0,
	})
	if result.Error != nil {
		if orm.IsUniqueConstraintFailed(result) {
			err = errors.WithCodeMsg(code.BadRequest, fmt.Sprintf("账号:%s已存在", req.Account))
			return
		} else {
			log.Errorf("create user error: %s", result.Error)
			err = errors.WithCodeMsg(code.Internal)
			return
		}
	}

	return
}
