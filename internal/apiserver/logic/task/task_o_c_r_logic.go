package task

import (
	"context"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	"github.com/gin-gonic/gin"
)

type TaskOCRLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskOCRLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskOCRLogic {
	return TaskOCRLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskOCR 财报识别
func (l *TaskOCRLogic) TaskOCR(req *task.GetTaskReq) (err error) {
	// todo: add your logic here and delete this line

	return
}
