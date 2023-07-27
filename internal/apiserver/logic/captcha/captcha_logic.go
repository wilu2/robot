package captcha

import (
	"context"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/captcha"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/code"
	"financial_statement/pkg/errors"
	"financial_statement/internal/apiserver/dal/model"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"financial_statement/pkg/log"
	"time"
)

var store = base64Captcha.DefaultMemStore

type CaptchaLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) CaptchaLogic {
	return CaptchaLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

//验证码生成配置
func NewDriver() *base64Captcha.DriverString {
    driver := new(base64Captcha.DriverString)
    driver.Height = 44
    driver.Width = 120
    driver.NoiseCount = 4
    // driver.ShowLineOptions = base64Captcha.OptionShowSineLine | base64Captcha.OptionShowSlimeLine | base64Captcha.OptionShowHollowLine
    driver.Length = 4
    driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
    driver.Fonts = []string{"wqy-microhei.ttc"}
    return driver
}

// Captcha 验证码获取
func (l *CaptchaLogic) Captcha() (resp captcha.CaptchaResp, err error) {
	var (
		tCaptcha = query.Use(l.svcCtx.Db).Captcha
	)
	_, e := tCaptcha.WithContext(l.ctx).Where(tCaptcha.Expiry.Lt(time.Now().Unix())).Delete()
	if e != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	driver := NewDriver().ConvertFonts()
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, e := cp.Generate()
	if e != nil{
		err = errors.WithCodeMsg(code.Internal, e.Error())
	} 
	mm, _ := time.ParseDuration("1m")
	e = tCaptcha.WithContext(l.ctx).Create(&model.Captcha{
		CaptchaID : id,
		Captcha : store.Get(id,true),
		CreatedAt: time.Now().Unix(),
		Expiry: time.Now().Add(mm).Unix(),
	})
	if e != nil{
		log.Errorf("create captcha error: %s", e)
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	resp.CaptchaId = id
	resp.CaptchaImg = b64s
	return
}
