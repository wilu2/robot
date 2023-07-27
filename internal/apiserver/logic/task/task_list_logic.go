package task

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	taskhelper "financial_statement/internal/pkg/task_helper"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gen/field"
)

type TaskListLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskListLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskListLogic {
	return TaskListLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskList 列出财报任务
func (l *TaskListLogic) TaskList(req *task.ListTaskReq) (resp task.ListTaskResp, err error) {
	var (
		u     = l.ginCtx.Keys["user"].(*model.User)
		tTask = query.Use(l.svcCtx.Db).Task
	)
	taskQuery := tTask.WithContext(l.ctx).Where(tTask.Status.Neq(consts.TaskStatusDeleted))
	if l.svcCtx.Config.ServerOptions.IsSaas {
		taskQuery = taskQuery.Where(tTask.CreaterUserID.Eq(u.ID))
	}

	if len(req.KeyWords) > 0 {
		// id 与 任务名的 模糊搜索
		idField := field.NewString(tTask.TableName(), tTask.ID.ColumnName().String())
		taskQuery = taskQuery.Where(tTask.WithContext(l.ctx).Where(idField.Like("%" + req.KeyWords + "%")).Or(tTask.TaskName.Like("%" + req.KeyWords + "%")))
	}

	count, err := taskQuery.Count()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	if len(req.OrderBy) > 0 { // order by 排序逻辑
		switch req.OrderBy {
		case "create_time":
			if req.OrderByType == "desc" {
				taskQuery = taskQuery.Order(tTask.ID.Desc())
			} else {
				taskQuery = taskQuery.Order(tTask.ID)
			}
		}
	} else {
		taskQuery = taskQuery.Order(tTask.ID.Desc())
	}
	tasks, err := taskQuery.Offset((int(req.Page) - 1) * int(req.PerPage)).Limit(int(req.PerPage)).Find()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	for _, dbTask := range tasks {
		var statements []task.TaskFinancialStatement
		statements, err = taskhelper.GetTaskResult(dbTask)
		if err != nil {
			err = errors.WithCodeMsg(code.Internal, err.Error())
			return
		}
		item := task.TaskListInfo{
			ID:             dbTask.ID,
			FileFormat:     dbTask.FileFormat,
			TaskName:       dbTask.TaskName,
			StatementTypes: [][]int{},
			StandardID:     dbTask.StandardID,
			Async:          int(dbTask.Async),
			Status:         int(dbTask.Status),
			Error:          *dbTask.Error,
			CreatedAt:      int(dbTask.CreatedAt),
			UpdatedAt:      int(dbTask.UpdatedAt),
		}
		for _, statement := range statements {
			tempStatementTypes := make([]int, 0)
			if statement.BalanceSheet.Count > 0 {
				tempStatementTypes = append(tempStatementTypes, consts.StatementTypeBalanceSheet)
			}
			if statement.CashFlowStatement.Count > 0 {
				tempStatementTypes = append(tempStatementTypes, consts.StatementTypeCashFlowStatement)
			}
			if statement.IncomeStatement.Count > 0 {
				tempStatementTypes = append(tempStatementTypes, consts.StatementTypeIncomeStatement)
			}
			item.StatementTypes = append(item.StatementTypes, tempStatementTypes)
		}

		resp.Tasks = append(resp.Tasks, item)
	}

	// copier.Copy(&resp.Tasks, &tasks)
	resp.Total = count
	return
}
