package standard

import (
	"financial_statement/internal/apiserver/logic/standard"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	standardType "financial_statement/internal/apiserver/types/standard"

	"github.com/gin-gonic/gin"
)

// StandardSetDefaultHandle 设为默认准则
func StandardSetDefaultHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req standardType.SetDefaultStandardReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := standard.NewStandardSetDefaultLogic(c, svcCtx)
		err := logic.StandardSetDefault(&req)
		response.HandleResponse(c, nil, err)
	}
}
