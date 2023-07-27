package title

import (
	"financial_statement/internal/apiserver/logic/standard/statement/title"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	titleType "financial_statement/internal/apiserver/types/standard/statement/title"

	"github.com/gin-gonic/gin"
)

// PreDisableTitleHandle 预删除
func PreDisableTitleHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req titleType.DisableTitleReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := title.NewPreDisableTitleLogic(c, svcCtx)
		resp, err := logic.PreDisableTitle(&req)
		response.HandleResponse(c, resp, err)
	}
}
