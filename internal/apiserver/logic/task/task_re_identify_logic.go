package task

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	jobTask "financial_statement/internal/apiserver/job/task"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TaskReIdentifyLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskReIdentifyLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskReIdentifyLogic {
	return TaskReIdentifyLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskReIdentify 重新识别财报任务
func (l *TaskReIdentifyLogic) TaskReIdentify(req *task.TaskReIdentifyReq) (err error) {
	var (
		tTask = query.Use(l.svcCtx.Db).Task
	)
	if _, err = tTask.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID)).Updates(model.Task{
		Status:         consts.TaskStatusCreated,
		StandardResult: nil,
		StandardID:     req.StandardID,
	}); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	if _, err = jobTask.NewRecognizeTask(&jobTask.RecoginzeTaskOption{
		TaskId:   req.TaskID,
		TaskType: jobTask.TaskTypeReIdentify,
	}); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	return
}
