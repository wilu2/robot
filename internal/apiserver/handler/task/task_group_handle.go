package task

import (
	"financial_statement/internal/apiserver/logic/task"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	taskType "financial_statement/internal/apiserver/types/task"

	"github.com/gin-gonic/gin"
)

// TaskGroupHandle 调整财报分组
func TaskGroupHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req taskType.TaskGroupReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := task.NewTaskGroupLogic(c, svcCtx)
		err := logic.TaskGroup(&req)
		response.HandleResponse(c, nil, err)
	}
}
