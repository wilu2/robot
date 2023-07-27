package task

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TaskReNameLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskReNameLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskReNameLogic {
	return TaskReNameLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskReName 重命名财报任务
func (l *TaskReNameLogic) TaskReName(req *task.TaskReNameReq) (err error) {
	var (
		u     = l.ginCtx.Keys["user"].(*model.User)
		tTask = query.Use(l.svcCtx.Db).Task
	)

	taskQuery := tTask.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID))
	if l.svcCtx.Config.ServerOptions.IsSaas {
		taskQuery = taskQuery.Where(tTask.CreaterUserID.Eq(u.ID))
	}
	if _, err = taskQuery.Update(tTask.TaskName, req.TaskName); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	return
}
