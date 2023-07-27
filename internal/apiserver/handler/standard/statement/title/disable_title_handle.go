package title

import (
	"financial_statement/internal/apiserver/logic/standard/statement/title"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	titleType "financial_statement/internal/apiserver/types/standard/statement/title"

	"github.com/gin-gonic/gin"
)

// DisableTitleHandle 停用/启用科目
func DisableTitleHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req titleType.DisableTitleReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := title.NewDisableTitleLogic(c, svcCtx)
		err := logic.DisableTitle(&req)
		response.HandleResponseWithStatusOk(c, nil, err)
	}
}
