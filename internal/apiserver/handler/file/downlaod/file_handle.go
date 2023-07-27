package downlaod

import (
	"financial_statement/internal/apiserver/logic/file/downlaod"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	downlaodType "financial_statement/internal/apiserver/types/file/downlaod"

	"github.com/gin-gonic/gin"
)

// FileHandle 文件下载
func FileHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req downlaodType.FileReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := downlaod.NewFileLogic(c, svcCtx)
		resp, err := logic.File(&req)
		if err != nil {
			response.HandleResponse(c, resp, err)
		} else {
			c.Header("Content-Type", "application/octet-stream")
			c.Header("Content-Transfer-Encoding", "binary")
			c.Header("Content-Disposition", "attachment;filename="+req.FileName)
			c.Data(200, "application/octet-stream", resp.File)
		}
	}
}
