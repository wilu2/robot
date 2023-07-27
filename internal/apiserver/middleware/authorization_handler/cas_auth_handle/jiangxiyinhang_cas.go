package casauthhandler

import (
	"context"
	"encoding/xml"
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
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/spf13/viper"
)

type JiangXiYinHangCasAuthHandler struct {
	once                sync.Once
	casAuthQueryKey     string
	casAuthServer       string
	casAuthTicketPath   string
	casAuthClientServer string
}

type Envelope struct {
	Body Body `xml:"Body"`
}

type Body struct {
	Response Response `xml:"ssoTicketValidateResponse"`
}

type Response struct {
	Return string `xml:"ssoTicketValidateReturn"`
}

type Root struct {
	Failure *Failure `xml:"failure"`
	Success *Success `xml:"success"`
}

type Success struct {
	Username string `xml:"username,attr"`
}

type Failure struct{}

func (l *JiangXiYinHangCasAuthHandler) GetTicket(c *gin.Context) (ticket string, err error) {
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

func (l *JiangXiYinHangCasAuthHandler) CheckTicket(ticket string) (ticketBody []byte, err error) {
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`<soapenv:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:int="InterfaceService">
	<soapenv:Header/>
	<soapenv:Body>
	   <int:ssoTicketValidate soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
		  <xml xsi:type="xsd:string">&lt;?xml version=&#34;1.0&#34; encoding=&#34;utf-8&#34;?&gt;&#xA;&#x9;&#x9;&lt;root&gt;&#xA;&#x9;&#x9;&#x9;&lt;token&gt;%s&lt;/token&gt;&#xA;&#x9;&#x9;&lt;/root&gt;</xml>
	   </int:ssoTicketValidate>
	</soapenv:Body>
 </soapenv:Envelope>`, ticket))
	req, err := http.NewRequest(method, l.casAuthServer, payload)
	if err != nil {
		return nil, fmt.Errorf("初始化httprequest %s 失败：%s", l.casAuthServer, err.Error())
	}
	client := &http.Client{}
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", "ssoTicketValidate")
	req.Header.Add("Authorization", "Basic cmFtc3Rlc3Q6MTIzNDQzMjE=") //ramstest:12344321
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("post %s 失败：%s", l.casAuthServer, err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post %s HttpStatus异常：%d", l.casAuthServer, res.StatusCode)
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func (l *JiangXiYinHangCasAuthHandler) GetUserUid(ticketBody []byte) (uid string, err error) {
	if len(ticketBody) == 0 {
		return "", fmt.Errorf("解析ticket验证body失败,长度为空")
	}
	xmlData := string(ticketBody)
	// 解析 XML
	var envelope Envelope
	err = xml.Unmarshal([]byte(xmlData), &envelope)
	if err != nil {
		return "", fmt.Errorf("XML:%s 解析出错：%s", xmlData, err.Error())
	}

	// 获取 ssoTicketValidateReturn 属性值
	result := envelope.Body.Response.Return
	var root Root
	err = xml.Unmarshal([]byte(result), &root)
	if err != nil {
		return "", fmt.Errorf("XML:%s 解析出错：%s", result, err.Error())
	}
	if root.Success != nil {
		uid = root.Success.Username
		return uid, nil
	} else {
		return "", fmt.Errorf("token验证不通过：%s %s", xmlData, result)
	}
}

func (l *JiangXiYinHangCasAuthHandler) CreateToken(uid string, c *gin.Context) (user *model.User, expiry int64, tokenStr string, err error) {
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
