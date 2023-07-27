package standard

import (
	"financial_statement/internal/apiserver/logic/standard"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	standardType "financial_statement/internal/apiserver/types/standard"

	"github.com/gin-gonic/gin"
)

// StandardCopyHandle 复制准则
func StandardCopyHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req standardType.CopyStandardReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := standard.NewStandardCopyLogic(c, svcCtx)
		err := logic.StandardCopy(&req)
		response.HandleResponseWithStatusOk(c, nil, err)
	}
}
