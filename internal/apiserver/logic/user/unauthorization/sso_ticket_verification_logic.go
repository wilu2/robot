package unauthorization

import (
	"context"
	"financial_statement/internal/apiserver/code"
	authorizationhandler "financial_statement/internal/apiserver/middleware/authorization_handler"
	casauthhandler "financial_statement/internal/apiserver/middleware/authorization_handler/cas_auth_handle"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/user/unauthorization"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type SsoTicketVerificationLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewSsoTicketVerificationLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) SsoTicketVerificationLogic {
	return SsoTicketVerificationLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

var (
	once           sync.Once
	casAuthHandler string
)

// SsoTicketVerification sso登录ticket验证，成功后返回一个用户token
func (l *SsoTicketVerificationLogic) SsoTicketVerification(req *unauthorization.SsoTicketVerificationReq) (resp unauthorization.PwdLoginResp, err error) {
	var (
		ticket     = ""
		ticketBody = make([]byte, 0)
		uid        = ""
	)
	once.Do(func() {
		casAuthHandler = viper.GetString("auth-cas.handler")
	})
	var auth authorizationhandler.ICASAuthorization
	switch casAuthHandler {
	case authorizationhandler.CasAuthorizationHandlerTypeLiuzhou:
		auth = &casauthhandler.LiuzhouCasAuthHandler{}
	case authorizationhandler.CasAuthorizationHandlerTypeBeibuwan:
		auth = &casauthhandler.BeibuwanCasAuthHandler{}
	case authorizationhandler.CasAuthorizationHandlerTypeJiangXiYinHang:
		auth = &casauthhandler.JiangXiYinHangCasAuthHandler{}
	default:
		auth = &casauthhandler.IntSigCasAuthHandler{}
	}
	if ticket, err = auth.GetTicket(l.ginCtx); err != nil {
		log.Errorf("获取ticket出错%s", err.Error())
		err = errors.WithCodeMsg(code.BadRequest, err.Error())
		return
	}
	if ticketBody, err = auth.CheckTicket(ticket); err != nil {
		log.Errorf("向cas服务器检查ticket出错%s", err.Error())
		err = errors.WithCodeMsg(code.BadRequest, err.Error())
		return
	}
	if uid, err = auth.GetUserUid(ticketBody); err != nil {
		log.Errorf("获取用户id失败：%s", err.Error())
		err = errors.WithCodeMsg(code.BadRequest, err.Error())
		return
	}
	user, expiry, tokenStr, err := auth.CreateToken(uid, l.ginCtx)
	if err != nil {
		log.Errorf("Create User Token 失败：%s", err.Error())
		err = errors.WithCodeMsg(code.BadRequest, err.Error())
		return
	}
	resp.Expiry = expiry
	resp.ID = user.ID
	resp.Token = tokenStr
	resp.Name = user.Name
	resp.Email = user.Email
	resp.Mobile = user.Mobile
	resp.IsAdmin = false
	return
}
