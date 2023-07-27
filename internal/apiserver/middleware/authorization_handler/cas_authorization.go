package authorizationhandler

import (
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	casauthhandler "financial_statement/internal/apiserver/middleware/authorization_handler/cas_auth_handle"
	"financial_statement/internal/apiserver/response"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type CASAuthorization struct {
	once           sync.Once
	casAuthHandler string
}

const (
	CasAuthorizationHandlerTypeLiuzhou        = "liuzhou"
	CasAuthorizationHandlerTypeIntSig         = "intsig"
	CasAuthorizationHandlerTypeBeibuwan       = "beibuwan"
	CasAuthorizationHandlerTypeJiangXiYinHang = "jiangxiyinhang"
)

type ICASAuthorization interface {
	GetTicket(c *gin.Context) (ticket string, err error)
	CheckTicket(ticket string) (ticketBody []byte, err error)
	GetUserUid(ticketBody []byte) (uid string, err error)
	CreateToken(uid string, c *gin.Context) (user *model.User, expiry int64, token string, err error)
}

func (u *CASAuthorization) Authorization(c *gin.Context) {
	u.once.Do(func() {
		u.casAuthHandler = viper.GetString("auth-cas.handler")
	})

	token := c.GetHeader("x-token")
	if len(token) != 0 {
		dbIns, _ := dal.GetDbFactoryOr(nil)
		loginToken := &model.LoginToken{}
		dbIns.GetDb().Where(&model.LoginToken{Token: token}).First(loginToken)
		if loginToken.UserID == 0 || loginToken.Expiry < time.Now().Unix() {
			//没有对应的token记录,或者token已过期
			u.unauthorized(c)
			return
		}
		user := &model.User{}
		dbIns.GetDb().Where(&model.User{ID: loginToken.UserID}).First(
			user)
		if user.ID == 0 {
			//没有该用户
			u.unauthorized(c)
			return
		}
		c.Set("user", user)
		c.Set("is_sso", true)
		c.Next()
		return
	}

	var auth ICASAuthorization
	switch u.casAuthHandler {
	case CasAuthorizationHandlerTypeLiuzhou:
		auth = &casauthhandler.LiuzhouCasAuthHandler{}
	case CasAuthorizationHandlerTypeBeibuwan:
		auth = &casauthhandler.BeibuwanCasAuthHandler{}
	case CasAuthorizationHandlerTypeJiangXiYinHang:
		auth = &casauthhandler.JiangXiYinHangCasAuthHandler{}
	default:
		auth = &casauthhandler.IntSigCasAuthHandler{}
	}
	ticket, err := auth.GetTicket(c)
	if err != nil {
		log.Errorf("获取ticket出错%s", err.Error())
		u.unauthorized(c)
		return
	}
	ticketBody, err := auth.CheckTicket(ticket)
	if err != nil {
		log.Errorf("向cas服务器检查ticket出错%s", err.Error())
		u.unauthorized(c)
		return
	}
	uid, err := auth.GetUserUid(ticketBody)
	if err != nil {
		log.Errorf("获取用户id失败：%s", err.Error())
		u.unauthorized(c)
		return
	}
	user, _, _, err := auth.CreateToken(uid, c)
	if err != nil {
		log.Errorf("CreateToken失败：%s", err.Error())
		u.unauthorized(c)
		return
	}
	c.Set("user", user)
	c.Next()
}

func (u *CASAuthorization) unauthorized(c *gin.Context) {
	err := errors.WithCodeMsg(code.Unauthorized)
	response.HandleResponse(c, nil, err)
	c.Abort()
}
