package task

import (
	"financial_statement/internal/apiserver/logic/task"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	taskType "financial_statement/internal/apiserver/types/task"

	"github.com/gin-gonic/gin"
)

// TaskReIdentifyHandle 重新识别财报任务
func TaskReIdentifyHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req taskType.TaskReIdentifyReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := task.NewTaskReIdentifyLogic(c, svcCtx)
		err := logic.TaskReIdentify(&req)
		response.HandleResponse(c, nil, err)
	}
}
