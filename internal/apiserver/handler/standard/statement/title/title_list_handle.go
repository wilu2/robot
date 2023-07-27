package title

import (
	"financial_statement/internal/apiserver/logic/standard/statement/title"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	titleType "financial_statement/internal/apiserver/types/standard/statement/title"

	"github.com/gin-gonic/gin"
)

// TitleListHandle 获取所有科目
func TitleListHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req titleType.TitleListReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := title.NewTitleListLogic(c, svcCtx)
		resp, err := logic.TitleList(&req)
		response.HandleResponse(c, resp, err)
	}
}
