package formulas

import (
	"financial_statement/internal/apiserver/logic/standard/statement/formulas"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	formulasType "financial_statement/internal/apiserver/types/standard/statement/formulas"

	"github.com/gin-gonic/gin"
)

// UpdateFormulaHandle 更新公式
func UpdateFormulaHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req formulasType.UpdateFormulaReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := formulas.NewUpdateFormulaLogic(c, svcCtx)
		err := logic.UpdateFormula(&req)
		response.HandleResponse(c, nil, err)
	}
}
