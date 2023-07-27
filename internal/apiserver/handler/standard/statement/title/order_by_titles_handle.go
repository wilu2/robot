package title

import (
	"financial_statement/internal/apiserver/logic/standard/statement/title"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	titleType "financial_statement/internal/apiserver/types/standard/statement/title"

	"github.com/gin-gonic/gin"
)

// OrderByTitlesHandle 科目排序
func OrderByTitlesHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req titleType.OrderByTitlesReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := title.NewOrderByTitlesLogic(c, svcCtx)
		err := logic.OrderByTitles(&req)
		response.HandleResponse(c, nil, err)
	}
}
