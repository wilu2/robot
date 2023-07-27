package task

import (
	"context"
	"encoding/json"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	taskhelper "financial_statement/internal/pkg/task_helper"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TaskInfoLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskInfoLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskInfoLogic {
	return TaskInfoLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TaskInfo 获取财报信息
func (l *TaskInfoLogic) TaskInfo(req *task.GetTaskReq) (resp task.GetTaskResp, err error) {
	var (
		u     = l.ginCtx.Keys["user"].(*model.User)
		tTask = query.Use(l.svcCtx.Db).Task
		tPage = query.Use(l.svcCtx.Db).Page
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

	pages, err := tPage.WithContext(l.ctx).Where(tPage.TaskID.Eq(req.TaskID)).Find()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	for _, page := range pages {
		var angle int
		if len(*page.OcrResult) > 0 {
			var tempOcr consts.PageOcrResult
			json.Unmarshal([]byte(*page.OcrResult), &tempOcr)
			angle = tempOcr.Result.Angle
		}
		resp.Images = append(resp.Images, task.Image{
			ImageSrc:    page.FilePath,
			RotateAngle: angle,
			FileId:      page.ID,
		})
	}
	resp.CreatedAt = int(dbTask.CreatedAt)
	// resp.External = task.ExternalInfo
	var externalInfo consts.ExternalInfo
	if err = json.Unmarshal([]byte(dbTask.ExternalInfo), &externalInfo); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}

	resp.ExtraInfo = externalInfo.ExtraInfo

	for _, file := range externalInfo.FileInfo {
		if len(file) > 0 {
			resp.External = append(resp.External, file)
		}
	}

	for _, group := range externalInfo.Groups {
		g := task.TaskGroupsResp{
			GroupName: group.GroupName,
			GroupId:   group.GroupId,
			Files:     []task.TaskGroupPageInfo{},
		}
		for _, file := range group.Files {
			g.Files = append(g.Files, task.TaskGroupPageInfo{
				FileId:      file.FileId,
				ImageSrc:    file.ImageSrc,
				RotateAngle: file.RotateAngle,
				Type:        file.Type,
			})
		}
		resp.Groups = append(resp.Groups, g)
	}
	resp.FileFormat = dbTask.FileFormat
	resp.TaskName = dbTask.TaskName
	resp.Standard = task.Standard{
		Id:         dbTask.StandardID,
		ExternalId: "",
	}
	resp.Status = int(dbTask.Status)
	resp.TaskId = dbTask.ID
	resp.UpdateAt = int(dbTask.UpdatedAt)

	resp.OperType = "B"

	if resp.Status == consts.TaskStatusOcrSuccess {
		if resp.FinancialStatement, err = taskhelper.GetTaskResult(dbTask); err != nil {
			err = errors.WithCodeMsg(code.Internal, err.Error())
			return
		}
	}

	return
}
