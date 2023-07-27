package api

import (
	"context"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/api"

	"github.com/gin-gonic/gin"
)

type TaskCreateApiLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskCreateApiLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskCreateApiLogic {
	return TaskCreateApiLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskCreateApi 创建财报识别任务
func (l *TaskCreateApiLogic) TaskCreateApi(req *api.CreateTaskReq) (resp api.CreateTaskResp, err error) {
	// todo: add your logic here and delete this line

	return
}
