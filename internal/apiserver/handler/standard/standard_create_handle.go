package standard

import (
	"financial_statement/internal/apiserver/logic/standard"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	standardType "financial_statement/internal/apiserver/types/standard"

	"github.com/gin-gonic/gin"
)

// StandardCreateHandle 创建准则
func StandardCreateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req standardType.CreateStandardReq
		if err := c.ShouldBind(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := standard.NewStandardCreateLogic(c, svcCtx)
		resp, err := logic.StandardCreate(&req)
		response.HandleResponseWithStatusOk(c, resp, err)
	}
}
