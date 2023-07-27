package captcha

import (
	"financial_statement/internal/apiserver/logic/captcha"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

// CaptchaHandle 验证码获取
func CaptchaHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		logic := captcha.NewCaptchaLogic(c, svcCtx)
		resp, err := logic.Captcha()
		response.HandleResponse(c, resp, err)
	}
}
