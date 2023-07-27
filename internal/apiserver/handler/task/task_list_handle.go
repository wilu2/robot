package task

import (
	"financial_statement/internal/apiserver/logic/task"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	taskType "financial_statement/internal/apiserver/types/task"

	"github.com/gin-gonic/gin"
)

// TaskListHandle 列出财报任务
func TaskListHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req taskType.ListTaskReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := task.NewTaskListLogic(c, svcCtx)
		resp, err := logic.TaskList(&req)
		response.HandleResponse(c, resp, err)
	}
}
