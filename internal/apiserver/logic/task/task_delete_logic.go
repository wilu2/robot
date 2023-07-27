package task

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TaskDeleteLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskDeleteLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskDeleteLogic {
	return TaskDeleteLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskDelete 删除财报任务
func (l *TaskDeleteLogic) TaskDelete(req *task.GetTaskReq) (err error) {
	var (
		u     = l.ginCtx.Keys["user"].(*model.User)
		tTask = query.Use(l.svcCtx.Db).Task
	)

	query := tTask.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID))
	if l.svcCtx.Config.ServerOptions.IsSaas {
		query = query.Where(tTask.CreaterUserID.Eq(u.ID))
	}
	_, err = query.Update(tTask.Status, consts.TaskStatusDeleted)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	return
}
