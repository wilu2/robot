package standard

import (
	"financial_statement/internal/apiserver/logic/standard"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	standardType "financial_statement/internal/apiserver/types/standard"

	"github.com/gin-gonic/gin"
)

// StatementsHandle 获取准则报表信息
func StatementsHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req standardType.GetStatementsReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := standard.NewStatementsLogic(c, svcCtx)
		resp, err := logic.Statements(&req)
		response.HandleResponse(c, resp, err)
	}
}
