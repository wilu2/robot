package middleware

import (
	"context"
	code2 "financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	user "financial_statement/internal/apiserver/logic/user/unauthorization"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	userType "financial_statement/internal/apiserver/types/user/unauthorization"
	"financial_statement/pkg/errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"net/http"
	"strconv"
	"time"
)

var (
	GinJwt         *jwt.GinJWTMiddleware
	AuthMiddleware gin.HandlerFunc
	identityKey    = "uid" // 通过 ginCtx.Keys["uid"].(float64) 直接获取 uid
	authErrorKey   = "authErrorKey"
)

func InitJWTAuth() {
	// ApiServerIss = viper.GetString("jwt.iss")
	GinJwt, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"), // 标识
		SigningAlgorithm: "HS256",                      // 加密算法
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          viper.GetDuration("jwt.timeout"),
		MaxRefresh:       viper.GetDuration("jwt.max-refresh"), // 刷新最大延长时间
		Authenticator:    authenticator(),                      // 登陆验证逻辑
		Authorizator:     authorizator(),                       // 登陆后权限校验逻辑
		PayloadFunc:      payloadFunc(),                        // 定义 Jwt 返回的 payload 数据
		Unauthorized:     unauthorized(),                       // middleware 校验错误时响应
		LoginResponse:    loginResponse(),                      // 登陆成功返回内容
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
			})
		}, // 退出登陆返回内容
		RefreshResponse: refreshResponse(), // 刷新 token 返回内容
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims["uid"]
		}, // 传递 authorizator 的参数
		IdentityKey:   identityKey,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		SendCookie:    false, // 是否插入 Cookie
	})
	AuthMiddleware = GinJwt.MiddlewareFunc()
}

// 定义 jwt 中 payload 返回的数据
func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{
			// "iss": ApiServerIss,
		}
		if uid, ok := data.(uint32); ok {
			claims["uid"] = uid
		}

		return claims
	}
}

// 登陆校验逻辑
func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (uid interface{}, err error) {
		var req userType.PwdLoginReq
		if err = c.ShouldBind(&req); err != nil {
			c.Set(authErrorKey, true)
			response.HandlerParamsResponse(c, err)
			return
		}
		svcCtx := svc.NewServiceContext(nil, nil)
		logic := user.NewUserLoginLogic(c, svcCtx)
		resp, err := logic.UserLogin(&req)
		if err != nil {
			c.Set(authErrorKey, true)
			response.HandleResponse(c, nil, err)
			return
		}

		return resp.ID, err
	}
}

// 登陆后权限校验逻辑，通过 redis 校验 jwt 签证时间
func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		var (
			uid uint
			ok  bool
		)
		if uid, ok = data.(uint); !ok {
			return false
		}
		redisIns, _ := dal.GetRedisFactoryOr(nil)
		lastJwtTime := redisIns.GetRedisDb().HGet(context.Background(), "jwt",
			strconv.FormatInt(int64(uid), 10)).Val() // 修改密码，退出登陆的最后时间戳
		lastJwtStampInt, _ := strconv.ParseInt(lastJwtTime, 10, 64)
		signatureTime := c.Keys["JWT_PAYLOAD"].(jwt.MapClaims)["orig_iat"].(float64) // jwt 签证时间
		if int64(signatureTime) < lastJwtStampInt {                                  // 之前签证的 JWT 都失效
			return false
		} else {
			userObj := &model.User{}
			dbIns, _ := dal.GetDbFactoryOr(nil)
			dbIns.GetDb().Where(&model.User{ID: uint32(uid)}).First(
				userObj)
			if userObj.ID == 0 { // 存在用户删除，jwt 还存在的情况
				return false
			}
			c.Set("user", userObj)
			return true
		}
	}
}

// 错误响应逻辑
func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		if _, ok := c.Get(authErrorKey); !ok {
			err := errors.WithCodeMsg(code2.Unauthorized, "")
			response.HandleResponse(c, nil, err)
		}
	}
}

// 登陆成功返回内容
func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": map[string]string{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			},
		})
	}
}

// 刷新 token 时长
func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": map[string]string{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			},
		})
	}
}
