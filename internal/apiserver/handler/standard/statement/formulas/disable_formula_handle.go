package formulas

import (
	"financial_statement/internal/apiserver/logic/standard/statement/formulas"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	formulasType "financial_statement/internal/apiserver/types/standard/statement/formulas"

	"github.com/gin-gonic/gin"
)

// DisableFormulaHandle 禁用/启用公式
func DisableFormulaHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req formulasType.DisableFormulaReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := formulas.NewDisableFormulaLogic(c, svcCtx)
		err := logic.DisableFormula(&req)
		response.HandleResponse(c, nil, err)
	}
}
