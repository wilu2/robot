package standard

import (
	"financial_statement/internal/apiserver/logic/standard"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	standardType "financial_statement/internal/apiserver/types/standard"

	"github.com/gin-gonic/gin"
)

// StandardListHandle 列出准则
func StandardListHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req standardType.ListStandardReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := standard.NewStandardListLogic(c, svcCtx)
		resp, err := logic.StandardList(&req)
		response.HandleResponse(c, resp, err)
	}
}
