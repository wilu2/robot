package formulas

import (
	"financial_statement/internal/apiserver/logic/standard/statement/formulas"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	formulasType "financial_statement/internal/apiserver/types/standard/statement/formulas"

	"github.com/gin-gonic/gin"
)

// FormulaListHandle 获取所有公式
func FormulaListHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req formulasType.FormulaListReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := formulas.NewFormulaListLogic(c, svcCtx)
		resp, err := logic.FormulaList(&req)
		response.HandleResponse(c, resp, err)
	}
}
