package convert

import (
	"financial_statement/internal/apiserver/logic/file/convert"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	convertType "financial_statement/internal/apiserver/types/file/convert"

	"github.com/gin-gonic/gin"
)

// Pdf2ImageHandle pdf转图片
func Pdf2ImageHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req convertType.FileConvertReq

		if err := c.ShouldBind(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		logic := convert.NewPdf2ImageLogic(c, svcCtx)
		resp, err := logic.Pdf2Image(&req)
		response.HandleResponseWithStatusOk(c, resp, err)
	}
}
