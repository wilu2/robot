package unauthorization

import (
	"financial_statement/internal/apiserver/logic/user/unauthorization"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	unauthorizationType "financial_statement/internal/apiserver/types/user/unauthorization"

	"github.com/gin-gonic/gin"
)

// SsoTicketVerificationHandle sso登录ticket验证，成功后返回一个用户token
func SsoTicketVerificationHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req unauthorizationType.SsoTicketVerificationReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := unauthorization.NewSsoTicketVerificationLogic(c, svcCtx)
		resp, err := logic.SsoTicketVerification(&req)
		response.HandleResponse(c, resp, err)
	}
}
