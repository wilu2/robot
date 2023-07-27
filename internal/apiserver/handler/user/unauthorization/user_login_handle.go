package unauthorization

import (
	"financial_statement/internal/apiserver/logic/user/unauthorization"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	unauthorizationType "financial_statement/internal/apiserver/types/user/unauthorization"

	"github.com/gin-gonic/gin"
)

// UserLoginHandle 用户登录
func UserLoginHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req unauthorizationType.PwdLoginReq
		if err := c.ShouldBind(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := unauthorization.NewUserLoginLogic(c, svcCtx)
		resp, err := logic.UserLogin(&req)
		response.HandleResponseWithStatusOk(c, resp, err)
	}
}
