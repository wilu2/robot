package unauthorization

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user/unauthorization"
	"financial_statement/internal/pkg/verify"
	"financial_statement/pkg/errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type UserLoginLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UserLoginLogic {
	return UserLoginLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UserLogin 用户登录
func (l *UserLoginLogic) UserLogin(req *unauthorization.PwdLoginReq) (resp unauthorization.PwdLoginResp, err error) {
	var (
		tUser      = query.Use(l.svcCtx.Db).User
		tUserToken = query.Use(l.svcCtx.Db).LoginToken
	)
	user, err := tUser.WithContext(l.ctx).
		Select(tUser.ID, tUser.Name, tUser.Password, tUser.Salt, tUser.Email, tUser.Mobile).
		Where(tUser.Account.Eq(req.Account)).
		First()
	if err != nil {
		err = errors.WithCodeMsg(code.BadRequest, "用户名或密码错误")
		return
	}
	if !verify.VerifyPassword(user, req.Password) {
		err = errors.WithCodeMsg(code.BadRequest, "用户名或密码错误")
		return
	}
	now := time.Now()
	expiry := now.Add(consts.TokenExpiry).Unix()
	tokenStr := ksuid.New().String()
	newUserToken := model.LoginToken{
		Token:     tokenStr,
		UserID:    user.ID,
		CreatedAt: now.Unix(),
		Expiry:    expiry,
	}

	// create user token info
	err = tUserToken.WithContext(l.ctx).Create(&newUserToken)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	l.ginCtx.SetCookie("token", tokenStr, int(consts.TokenExpiry/time.Second), "/", "", false, false)
	resp.Expiry = expiry
	resp.ID = user.ID
	resp.Token = tokenStr
	resp.Name = user.Name
	resp.Email = user.Email
	resp.Mobile = user.Mobile
	resp.IsAdmin = false
	tUserToken.WithContext(l.ctx).Where(tUserToken.Expiry.Lt(now.Unix())).Delete()
	return
}
