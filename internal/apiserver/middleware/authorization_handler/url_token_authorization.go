package authorizationhandler

import (
	"bytes"
	"encoding/json"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/response"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// url token 认证时，检查并验证url中的token，通过后的所有操作使用admin账户进行操作
type UrlTokenAuthorization struct {
	once            sync.Once
	check_token_api string
}

var ()

func (u *UrlTokenAuthorization) Authorization(c *gin.Context) {
	userName, ok := u.getUserName(c)
	dbIns, _ := dal.GetDbFactoryOr(nil)
	user := &model.User{}
	if ok {
		//
		dbIns.GetDb().Where(&model.User{Name: "admin"}).First(
			user)
		if user.ID == 0 {
			//没有该用户
			u.unauthorized(c)
			return
		}
	} else {
		token := c.GetHeader("x-token")
		if len(token) == 0 {
			u.unauthorized(c)
			return
		}
		loginToken := &model.LoginToken{}
		dbIns.GetDb().Where(&model.LoginToken{Token: token}).First(loginToken)
		if loginToken.UserID == 0 || loginToken.Expiry < time.Now().Unix() {
			//没有对应的token记录,或者token已过期
			u.unauthorized(c)
			return
		}
		dbIns.GetDb().Where(&model.User{ID: loginToken.UserID}).First(
			user)
		if user.ID == 0 {
			//没有该用户
			u.unauthorized(c)
			return
		}
	}
	if len(userName) > 0 {
		user.Account = userName
	}
	c.Set("user", user)
	c.Set("is_sso", false)
	c.Set("user_name", userName)
	c.Next()
}

func (u *UrlTokenAuthorization) getUserName(c *gin.Context) (string, bool) {
	userName := ""
	token, ok := c.GetQuery("token")
	if len(token) == 0 || !ok {
		log.Infof("URL TOKEN 认证，获取url中token参数失败!尝试获取user_id")
		userName, ok = c.GetQuery("user_id")
		if len(userName) == 0 || !ok {
			log.Errorf("URL TOKEN 认证，获取url中user_id参数失败!,返回认证失败！")
			// u.unauthorized(c)
			return "", false
		}
		log.Infof("URL TOKEN 认证，获取url中user_id参数成功，user_id:%s", userName)
	} else {
		ok, userName = u.checkToken(token)
		if !ok {
			log.Errorf("URL TOKEN 认证，获取url中token成功，但4A认证失败!，返回认证失败！")
			// u.unauthorized(c)
			return "", false
		}
		log.Infof("URL TOKEN 认证，获取url中token成功，4A认证成功!user_name:%s", userName)
	}
	return userName, true
}
func (u *UrlTokenAuthorization) unauthorized(c *gin.Context) {
	err := errors.WithCodeMsg(code.Unauthorized)
	response.HandleResponse(c, nil, err)
	c.Abort()
}

//验证通过后返回用户工号
func (u *UrlTokenAuthorization) checkToken(token string) (bool, string) {
	u.once.Do(func() {
		u.check_token_api = viper.GetString("auth.check-token-api")
	})
	postData := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	postDataStr, _ := json.Marshal(postData)
	body := bytes.NewReader([]byte(postDataStr))
	httpReq, err := http.NewRequest("POST", u.check_token_api, body)
	httpReq.Header.Add("Content-type", "application/json")

	if err != nil {
		log.Errorf("Url Token 认证请求失败：%s", err.Error())
		return false, ""
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Errorf("Url Token 认证请求失败httpstatus:%s", err.Error())
		return false, ""
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Errorf("Url Token 认证请求失败httpstatus:%d", httpResp.StatusCode)
		return false, ""
	}
	defer httpResp.Body.Close()
	resBody, _ := ioutil.ReadAll(httpResp.Body)

	result := struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			UserName string `json:"user_name"`
		} `json:"data"`
	}{}

	if err = json.Unmarshal(resBody, &result); err != nil {
		log.Errorf("Url Token 认证请求返回的json解析失败：%s", err.Error())
		log.Errorf("Url Token 认证请求返回的Body：%s", resBody)
		return false, ""
	}
	if result.Code != "000000" {
		log.Errorf("Url Token 认证请求返回非000000的Code码：%s", resBody)
		return false, ""
	}
	return true, result.Data.UserName
}
