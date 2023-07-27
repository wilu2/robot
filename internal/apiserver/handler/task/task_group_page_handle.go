package task

import (
	"financial_statement/internal/apiserver/logic/task"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	taskType "financial_statement/internal/apiserver/types/task"

	"github.com/gin-gonic/gin"
)

// TaskGroupPageHandle 调整分组时的上传文件接口
func TaskGroupPageHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req taskType.TaskGroupPageReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := task.NewTaskGroupPageLogic(c, svcCtx)
		resp, err := logic.TaskGroupPage(&req)
		response.HandleResponse(c, resp, err)
	}
}
