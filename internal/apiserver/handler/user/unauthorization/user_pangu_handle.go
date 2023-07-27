package unauthorization

import (
	"financial_statement/internal/apiserver/logic/user/unauthorization"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	unauthorizationType "financial_statement/internal/apiserver/types/user/unauthorization"

	"github.com/gin-gonic/gin"
)

// UserPanguHandle 系统初始化创建系统管理员
func UserPanguHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req unauthorizationType.UserPangu
		if err := c.ShouldBind(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := unauthorization.NewUserPanguLogic(c, svcCtx)
		err := logic.UserPangu(&req)
		response.HandleResponse(c, nil, err)
	}
}
