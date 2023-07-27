package formulas

import (
	"financial_statement/internal/apiserver/logic/standard/statement/formulas"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	formulasType "financial_statement/internal/apiserver/types/standard/statement/formulas"

	"github.com/gin-gonic/gin"
)

// FormulaCreateHandle 新增公式
func FormulaCreateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req formulasType.CreateFormulaReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := formulas.NewFormulaCreateLogic(c, svcCtx)
		resp, err := logic.FormulaCreate(&req)
		response.HandleResponse(c, resp, err)
	}
}
