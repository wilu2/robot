package casauthhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/pkg/database/orm"
	"financial_statement/internal/pkg/verify"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"financial_statement/pkg/stringx"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/segmentio/ksuid"
	"github.com/spf13/viper"
)

type LiuzhouCasAuthHandler struct {
	once                sync.Once
	casAuthQueryKey     string
	casAuthServer       string
	casAuthTicketPath   string
	casAuthClientServer string
}

func (l *LiuzhouCasAuthHandler) GetTicket(c *gin.Context) (ticket string, err error) {
	l.once.Do(func() {
		l.casAuthQueryKey = viper.GetString("auth-cas.query-key")
		l.casAuthServer = viper.GetString("auth-cas.server")
		l.casAuthTicketPath = viper.GetString("auth-cas.uid-path")
		l.casAuthClientServer = viper.GetString("auth-cas.client-server")
	})
	key, ok := c.GetQuery(l.casAuthQueryKey)
	if len(key) == 0 || !ok {
		return "", fmt.Errorf("CasAuthQueryKey获取失败！")
	}
	return key, nil
}

func (l *LiuzhouCasAuthHandler) CheckTicket(ticket string) (ticketBody []byte, err error) {
	globalNo := fmt.Sprintf("11%s%s", time.Now().UTC().Format(http.TimeFormat), stringx.RandString(1))
	postData := struct {
		Ticket             string `json:"ticket"`
		Service            string `json:"service"`
		RequestGlobalJnlNo string `json:"requestGlobalJnlNo"`
		RequestJnlNo       string `json:"requestJnlNo"`
		RequestChannelCode string `json:"requestChannelCode"`
		RequestChannelID   string `json:"requestChannelId"`
		RequestDate        string `json:"requestDate"`
		BankId             string `json:"bankId"`
	}{
		Ticket:             ticket,
		Service:            l.casAuthClientServer,
		RequestGlobalJnlNo: globalNo,
		RequestJnlNo:       globalNo,
		RequestChannelCode: "OC",
		RequestChannelID:   "OC",
		RequestDate:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		BankId:             "70100",
	}
	postDataStr, _ := json.Marshal(postData)
	body := bytes.NewReader([]byte(postDataStr))

	httpReq, err := http.NewRequest("POST", l.casAuthServer, body)
	httpReq.Header.Add("Content-type", "application/json")
	if err != nil {
		return nil, fmt.Errorf("初始化httprequest %s 失败：%s", l.casAuthServer, err.Error())
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("post %s 失败：%s", l.casAuthServer, err.Error())
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post %s HttpStatus异常：%d", l.casAuthServer, httpResp.StatusCode)
	}
	defer httpResp.Body.Close()
	return ioutil.ReadAll(httpResp.Body)
}

func (l *LiuzhouCasAuthHandler) GetUserUid(ticketBody []byte) (uid string, err error) {
	str := (*string)(unsafe.Pointer(&ticketBody))
	obj, err := oj.ParseString(*str)
	if err != nil {
		return "", err
	}
	ticket, err := jp.ParseString(l.casAuthTicketPath)
	if err != nil {
		return "", err
	}
	value := ticket.Get(obj)
	if value != nil {
		return value[0].(string), nil
	}
	return "", fmt.Errorf("解析ticket path: %s失败，body:%s", l.casAuthTicketPath, ticketBody)
}

func (l *LiuzhouCasAuthHandler) CreateToken(uid string, c *gin.Context) (user *model.User, expiry int64, tokenStr string, err error) {
	dbIns, _ := dal.GetDbFactoryOr(nil)
	db := dbIns.GetDb()
	var (
		tUser      = query.Use(db).User
		tUserToken = query.Use(db).LoginToken
	)
	count, err := tUser.WithContext(context.Background()).
		Select(tUser.ID, tUser.Name, tUser.Password, tUser.Salt, tUser.Email, tUser.Mobile).
		Where(tUser.Account.Eq(uid)).Count()
	if err != nil {
		return nil, 0, "", err
	}
	if count == 0 {
		salt := stringx.RandString(8)
		password := verify.CalcPassword("INTSIG@intsig2022", salt)
		user = &model.User{
			Account:  uid,
			Password: password,
			Salt:     salt,
			Name:     uid,
			Email:    "abc@fr.net",
			Mobile:   "18888888888"}
		result := db.Create(user)
		if result.Error != nil {
			if orm.IsUniqueConstraintFailed(result) {
				err = errors.WithCodeMsg(code.BadRequest, "User already exists")
				return nil, 0, "", err
			} else {
				log.Errorf("create user error: %s", result.Error)
				err = errors.WithCodeMsg(code.Internal)
				return nil, 0, "", err
			}
		}
	} else {
		user, err = tUser.WithContext(context.Background()).Select(tUser.ID, tUser.Name, tUser.Password, tUser.Salt, tUser.Email, tUser.Mobile).
			Where(tUser.Account.Eq(uid)).First()
		if err != nil {
			return nil, 0, "", err
		}
	}

	now := time.Now()
	expiry = now.Add(consts.TokenExpiry).Unix()
	tokenStr = ksuid.New().String()
	newUserToken := model.LoginToken{
		Token:     tokenStr,
		UserID:    user.ID,
		CreatedAt: now.Unix(),
		Expiry:    expiry,
	}

	// create user token info
	err = tUserToken.WithContext(context.Background()).Create(&newUserToken)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return nil, 0, "", err
	}
	c.SetCookie("token", tokenStr, int(consts.TokenExpiry/time.Second), "/", "", false, false)
	tUserToken.WithContext(context.Background()).Where(tUserToken.Expiry.Lt(now.Unix())).Delete()
	return

}
