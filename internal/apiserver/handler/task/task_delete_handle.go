package task

import (
	"financial_statement/internal/apiserver/logic/task"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	taskType "financial_statement/internal/apiserver/types/task"

	"github.com/gin-gonic/gin"
)

// TaskDeleteHandle 删除财报任务
func TaskDeleteHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req taskType.GetTaskReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := task.NewTaskDeleteLogic(c, svcCtx)
		err := logic.TaskDelete(&req)
		response.HandleResponse(c, nil, err)
	}
}
