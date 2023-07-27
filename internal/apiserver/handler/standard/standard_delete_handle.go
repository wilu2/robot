package standard

import (
	"financial_statement/internal/apiserver/logic/standard"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	standardType "financial_statement/internal/apiserver/types/standard"

	"github.com/gin-gonic/gin"
)

// StandardDeleteHandle 启用/停用准则
func StandardDeleteHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req standardType.GetStandardReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := standard.NewStandardDeleteLogic(c, svcCtx)
		err := logic.StandardDelete(&req)
		response.HandleResponse(c, nil, err)
	}
}
