package task

import (
	"financial_statement/internal/apiserver/logic/task"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	taskType "financial_statement/internal/apiserver/types/task"

	"github.com/gin-gonic/gin"
)

// TaskCreateHandle 创建财报识别任务
func TaskCreateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req taskType.CreateTaskReq
		if err := c.ShouldBind(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := task.NewTaskCreateLogic(c, svcCtx)
		resp, err := logic.TaskCreate(&req)
		response.HandleResponse(c, resp, err)
	}
}
