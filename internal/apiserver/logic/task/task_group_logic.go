package task

import (
	"context"
	"encoding/json"
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

type TaskGroupLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskGroupLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskGroupLogic {
	return TaskGroupLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskGroup 调整财报分组
func (l *TaskGroupLogic) TaskGroup(req *task.TaskGroupReq) (err error) {
	var (
		u     = l.ginCtx.Keys["user"].(*model.User)
		tTask = query.Use(l.svcCtx.Db).Task
		// tPage = query.Use(l.svcCtx.Db).Page
	)

	query := tTask.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID))
	if l.svcCtx.Config.ServerOptions.IsSaas {
		query = query.Where(tTask.CreaterUserID.Eq(u.ID))
	}
	dbTask, err := query.First()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	if len(req.Group) > 0 {
		//解析分组数据
		var externalInfo consts.ExternalInfo
		if err = json.Unmarshal([]byte(dbTask.ExternalInfo), &externalInfo); err != nil {
			err = errors.WithCodeMsg(code.Internal, err.Error())
			return err
		}
		externalInfo.Groups = make([]consts.ExternalInfoGroup, 0)
		for _, group := range req.Group {
			item := consts.ExternalInfoGroup{
				GroupName: group.GroupName,
				GroupId:   group.GroupId,
				Files: []struct {
					ImageSrc    string "json:\"img_src\""
					RotateAngle int    "json:\"rotate_angle\""
					FileId      uint32 "json:\"file_id\""
					Type        int    "json:\"type\""
				}{},
			}
			for _, file := range group.Files {
				item.Files = append(item.Files, struct {
					ImageSrc    string "json:\"img_src\""
					RotateAngle int    "json:\"rotate_angle\""
					FileId      uint32 "json:\"file_id\""
					Type        int    "json:\"type\""
				}{
					ImageSrc:    file.ImageSrc,
					RotateAngle: file.RotateAngle,
					FileId:      file.FileId,
					Type:        file.Type,
				})
			}
			externalInfo.Groups = append(externalInfo.Groups, item)
		}

		externalInfoJson, err := json.Marshal(externalInfo)
		if err != nil {
			err = errors.WithCodeMsg(code.Internal, err.Error())
			return err
		}
		if _, err = tTask.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID)).Updates(model.Task{
			Status:         consts.TaskStatusCreated,
			StandardResult: nil,
			ExternalInfo:   string(externalInfoJson),
		}); err != nil {
			err = errors.WithCodeMsg(code.Internal, err.Error())
			return err
		}
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
